package pcmd

import (
	"bufio"
	"bytes"
	"io"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	errorsutil "k8s.io/apimachinery/pkg/util/errors"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
)

func SplitYAMLDocuments(yamlBytes []byte) (DocumentMap, error) {
	gvkmap := DocumentMap{}
	knownKinds := map[string]bool{}
	errs := []error{}
	buf := bytes.NewBuffer(yamlBytes)
	reader := utilyaml.NewYAMLReader(bufio.NewReader(buf))
	for {
		// Read one YAML document at a time, until io.EOF is returned
		b, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if len(b) == 0 {
			break
		}
		// Deserialize the TypeMeta information of this byte slice
		gvk, err := yamlserializer.DefaultMetaFactory.Interpret(b)
		if err != nil {
			return nil, err
		}
		if len(gvk.Group) == 0 || len(gvk.Version) == 0 || len(gvk.Kind) == 0 {
			return nil, errors.Errorf("invalid configuration for GroupVersionKind %+v: kind and apiVersion is mandatory information that must be specified", gvk)
		}

		// Check whether the kind has been registered before. If it has, throw an error
		if known := knownKinds[gvk.Kind]; known {
			errs = append(errs, errors.Errorf("invalid configuration: kind %q is specified twice in YAML file", gvk.Kind))
			continue
		}
		knownKinds[gvk.Kind] = true

		// Save the mapping between the gvk and the bytes that object consists of
		gvkmap[*gvk] = b
	}
	if err := errorsutil.NewAggregate(errs); err != nil {
		return nil, err
	}
	return gvkmap, nil
}

// VerifyUnmarshalStrict takes a slice of schems, a JSON/YAML byte slice and a GroupVersionKind
// and verifies if the schema is known and if the byte slice unmarshals with strict mode.
func VerifyUnmarshalStrict(schemes []*runtime.Scheme, gvk schema.GroupVersionKind, bytes []byte) error {
	var scheme *runtime.Scheme
	for _, s := range schemes {
		if _, err := s.New(gvk); err == nil {
			scheme = s
			break
		}
	}
	if scheme == nil {
		return errors.Errorf("unknown configuration %#v", gvk)
	}

	opt := json.SerializerOptions{Yaml: true, Pretty: false, Strict: true}
	serializer := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme, scheme, opt)
	_, _, err := serializer.Decode(bytes, &gvk, nil)
	if err != nil {
		return errors.Wrapf(err, "error unmarshaling configuration %#v", gvk)
	}

	return nil
}
