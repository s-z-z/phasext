package pcmd

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kubeadm/app/componentconfigs"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/config/strict"
)

type DocumentParser struct {
	Dp     kubeadm.DocumentMap
	scheme *runtime.Scheme
}

type DocumentParser2Redaer func(dp *DocumentParser) (io.Reader, error)

func NewDocumentParser(fpath string, scheme *runtime.Scheme) (*DocumentParser, error) {
	allBytes, err := os.ReadFile(fpath)
	if err != nil {
		return nil, errors.Wrapf(err, "NewDocumentParser: read file %s error", fpath)
	}

	//cprt.Debug("NewDocumentParser: %s", allBytes)

	// 分割yaml对象, 不允许有重复对象
	gvk2b, err := kubeadmutil.SplitYAMLDocuments(allBytes)
	if err != nil {
		return nil, errors.Wrap(err, "NewDocumentParser:SplitYAMLDocuments: split yaml document error")
	}

	// 校验版本和字段
	for gvk, b := range gvk2b {
		if err := strict.VerifyUnmarshalStrict(
			[]*runtime.Scheme{scheme, componentconfigs.Scheme}, gvk, b); err != nil {
			return nil, errors.Wrap(err, "NewDocumentParser:VerifyUnmarshalStrict: verify unmarshal strict error")
		}
	}
	return &DocumentParser{
		Dp:     gvk2b,
		scheme: scheme,
	}, nil
}

func (g *DocumentParser) GetBytesByGvk(gvk schema.GroupVersionKind) ([]byte, bool) {
	b, ok := g.Dp[gvk]
	return b, ok
}

func (g *DocumentParser) GetBytesByObj(o schema.ObjectKind) ([]byte, bool) {
	return g.GetBytesByGvk(o.GroupVersionKind())
}

func (g *DocumentParser) GetBytes(o WareHouse) ([]byte, error) {
	gvk, err := GetGVKByObject(g.scheme, o)
	if err != nil {
		return nil, errors.Wrap(err, "GetBytes:GetGVKByObject: get object kind error")
	}
	b, ok := g.GetBytesByGvk(gvk)
	if !ok {
		return nil, errors.Errorf("GetBytes:GetObjBytes: not found: %+v", gvk)
	}
	return b, nil
}

// Reader 自定义解析g2b
func (g *DocumentParser) Reader(f DocumentParser2Redaer) (io.Reader, error) {
	return f(g)
}

// SelfReader 只读取o自身的段
func (g *DocumentParser) SelfReader(o WareHouse) (io.Reader, error) {
	b, err := g.GetBytes(o)
	if err != nil {
		return nil, errors.Wrap(err, "GetReader:GetBytes: get bytes error")
	}

	klog.V(5).Infof("GetReader:GetBytes:\n%s\n", string(b))

	return strings.NewReader(string(b)), nil
}

func (g *DocumentParser) Fill(o WareHouse, codecs serializer.CodecFactory) error {
	b, err := g.GetBytes(o)
	if err != nil {
		return errors.Wrap(err, "Fill:GetBytes: get bytes error")
	}
	return runtime.DecodeInto(codecs.UniversalDecoder(), b, o)
}

func OnlyUnmarshalSelf(o WareHouse) DocumentParser2Redaer {
	return func(dp *DocumentParser) (io.Reader, error) {
		return dp.SelfReader(o)
	}
}

func GetGVKByObject(s *runtime.Scheme, o WareHouse) (schema.GroupVersionKind, error) {
	gvks, _, err := s.ObjectKinds(o)
	if err != nil {
		return schema.GroupVersionKind{}, errors.Wrap(err, "GetBytesByObj: get object kind error")
	}
	return gvks[0], nil
}

func File2DocumentParser(configPath string, scheme *runtime.Scheme) (*DocumentParser, error) {
	documentParser, err := NewDocumentParser(configPath, scheme)
	return documentParser, err
}

func ReaderFillData(viper *viper.Viper, reader io.Reader, o interface{}) error {
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(reader); err != nil {
		return errors.Wrap(err, "pcmd:parse:ReaderFillData:ReadConfig")
	}

	if err := viper.Unmarshal(o); err != nil {
		return errors.Wrap(err, "pcmd:parse:ReaderFillData:Unmarshal")
	}
	return nil
}

func ObjectToYaml(codec serializer.CodecFactory, o WareHouse, gvk schema.GroupVersionKind) ([]byte, error) {
	jsonData, err := ObjectToJson(codec, o, gvk)
	if err != nil {
		return nil, err
	}
	b, err := JsonToYaml(jsonData)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ObjectToJson(codec serializer.CodecFactory, o WareHouse, gvk schema.GroupVersionKind) ([]byte, error) {
	b, err := runtime.Encode(codec.LegacyCodec(gvk.GroupVersion()), o)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func JsonToYaml(jsonData []byte) ([]byte, error) {
	var yamlData map[string]interface{}
	if err := json.Unmarshal(jsonData, &yamlData); err != nil {
		return nil, err
	}
	b, err := yaml.Marshal(yamlData)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func RepalceDocument(parser kubeadm.DocumentMap, codec serializer.CodecFactory, o WareHouse, gvk schema.GroupVersionKind) (string, error) {
	var yamls []string
	for gvkI, b := range parser {
		if gvkI == gvk {
			_b, err := ObjectToYaml(codec, o, gvk)
			if err != nil {
				return "", errors.Wrap(err, "pcmd:parse:RepalceDocument:ObjectToYaml")
			}
			b = _b
		}
		yamls = append(yamls, string(b))
	}
	return strings.Join(yamls, "---\n"), nil
}

func WriteBackFile(configPath string, parser kubeadm.DocumentMap, codec serializer.CodecFactory, o WareHouse, gvk schema.GroupVersionKind) error {

	data, err := RepalceDocument(parser, codec, o, gvk)
	if err != nil {
		return errors.Wrap(err, "pcmd:parse:WriteBackFile:RepalceDocument")
	}

	klog.V(7).Infof("write back to file: %s", configPath)

	f, err := os.OpenFile(configPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrap(err, "pcmd:parse:WriteBackFile:OpenFile")
	}

	if _, err = f.WriteString(data); err != nil {
		return errors.Wrap(err, "pcmd:parse:WriteBackFile:WriteString")
	}

	if err = f.Close(); err != nil {
		return errors.Wrap(err, "pcmd:parse:WriteBackFile:CloseFile")
	}
	return nil
}
