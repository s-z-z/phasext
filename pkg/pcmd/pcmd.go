package pcmd

// phase command

import (
	"io"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"

	"github.com/suzi1037/phasext/pkg/pcmd/util"
	util2 "github.com/suzi1037/phasext/pkg/util"
)

type PhaseCmdFactory struct {
	scheme   *runtime.Scheme
	validate *validator.Validate
}

func (pf *PhaseCmdFactory) Create(use string, opts ...Option) *PhasesCmd {
	return pf.CreateWithProp(CmdProp{Use: use}, opts...)
}

func (pf *PhaseCmdFactory) CreateWithProp(prop CmdProp, opts ...Option) *PhasesCmd {
	opts = append(opts,
		WithScheme(pf.scheme),
		WithValidator(pf.validate),
	)
	return newPhasesCmd(prop, opts...)
}

func NewPhaseCmdFactory(s *runtime.Scheme, v *validator.Validate) *PhaseCmdFactory {
	return &PhaseCmdFactory{
		scheme:   s,
		validate: v,
	}
}

type PhasesCmd struct {
	cmd                    *cobra.Command
	Runner                 *workflow.Runner
	data                   WareHouse
	gvk                    schema.GroupVersionKind
	firstAppend            bool
	withConfirm            bool
	withConfig             bool
	configFlag             string
	configPath             string
	scheme                 *runtime.Scheme
	documentParser2Reader  DocumentParser2Redaer
	documentParser         *DocumentParser
	runnerDataInitializer  util.RunnerDataInitializer
	exportOverrideFlags    bool
	specExportIncludeFlags []string
	persistentExportedFlag bool
	extraFlagStructs       []any
	bindToCommand          bool
	finished               bool
	configWriteBack        bool
	v                      *validator.Validate
	shouldValidate         bool
	viper                  *viper.Viper
	viperFn                func(*viper.Viper)
	// preRunE1 load data前执行
	preRunE1 CobraRun
	// preRunE2 load data后执行
	preRunE2 CobraRun
	// postRunE1 writeback data前执行
	postRunE1 CobraRun
	// postRunE2 writeback data后执行
	postRunE2 CobraRun
}

func newPhaseCmdByProp(prop CmdProp) *PhasesCmd {
	runner := workflow.NewRunner()

	cmd := &cobra.Command{
		Use:                    prop.Use,
		Aliases:                prop.Aliases,
		SuggestFor:             prop.SuggestFor,
		Short:                  prop.Short,
		GroupID:                prop.GroupID,
		Long:                   prop.Long,
		Example:                prop.Example,
		ValidArgs:              prop.ValidArgs,
		ValidArgsFunction:      prop.ValidArgsFunction,
		Args:                   prop.Args,
		ArgAliases:             prop.ArgAliases,
		BashCompletionFunction: prop.BashCompletionFunction,
		Version:                prop.Version,
		Hidden:                 prop.Hidden,
		SilenceUsage:           prop.SilenceUsage,
		SilenceErrors:          prop.SilenceErrors,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Run(args)
		},
	}

	return &PhasesCmd{
		cmd:             cmd,
		Runner:          runner,
		configPath:      DefaultConfigPath,
		configFlag:      DefaultConfigFlag,
		configWriteBack: DefaultConfigWriteBack,
		shouldValidate:  DefaultGoValidate,
		firstAppend:     true,
		viper:           viper.New(),
	}
}

func newPhasesCmd(prop CmdProp, opts ...Option) *PhasesCmd {

	p := newPhaseCmdByProp(prop)

	for _, o := range opts {
		o(p)
	}

	if p.data == nil {
		if p.withConfig || p.exportOverrideFlags {
			klog.Fatalf("pcmd:New: If no WithData, can not set WithConfig or WithExportOverrideFlags")
		}
	} else {
		gvk, err := GetGVKByObject(p.scheme, p.data)
		if err != nil {
			klog.Fatalf("pcmd:New: %s", err)
		}
		p.gvk = gvk
	}

	if p.configWriteBack && !p.withConfig {
		klog.Fatalf("pcmd:New: If configWriteBack, WithConfig must be set")
	}

	// 控制顺序
	p.init()

	return p
}

func (p *PhasesCmd) init() {

	if p.withConfig && p.configPath == "" {
		util.AddConfigFlag(p.cmd, p.configFlag, &p.configPath)
	}

	// 注入PersistentPreRunE: 检查scheme, 解析文件, Unmarshal
	p.documentToDataPersistentPreRun()

	// 注入PostRun: 配置回写
	p.dataToDocumentPostRun()

	// export tag解析, 添加，支持flag > file（or nil）
	flagKind := util.Local
	if p.persistentExportedFlag {
		flagKind = util.Persistent
	}
	p._exportOverrideFlags(flagKind)
	p._exportExtraFlags(flagKind)

	if p.viperFn != nil {
		p.viperFn(p.viper)
	}
}

func (p *PhasesCmd) finalize() {
	// ! 必须Append Phase之后执行

	// runner数据初始化: 返回Data
	if p.runnerDataInitializer == nil {
		p.Runner.SetDataInitializer(util.OnlyArgsDataInitializer)
	}

	// 支持Phase
	if p.bindToCommand {
		p.Runner.BindToCommand(p.cmd)
	}
}

func (p *PhasesCmd) _exportOverrideFlags(flagKind util.FlagKind) {

	if !p.exportOverrideFlags {
		return
	}

	util.AddExportFlags(p.cmd, p.data, p.specExportIncludeFlags, flagKind, false)

	if flagKind == util.Local {
		_ = p.viper.BindPFlags(p.cmd.Flags())
	} else {
		_ = p.viper.BindPFlags(p.cmd.PersistentFlags())
	}
}

func (p *PhasesCmd) _exportExtraFlags(flagKind util.FlagKind) {
	for _, v := range p.extraFlagStructs {
		util.AddExportFlags(p.cmd, v, []string{}, flagKind, true)
	}
}

func (p *PhasesCmd) setDefaultDocumentParser() {
	// documentParser2Reader 默认使用UnmarshalSelf
	if p.withConfig && p.documentParser2Reader == nil {
		p.documentParser2Reader = OnlyUnmarshalSelf(p.data)
	}
}

func (p *PhasesCmd) documentToDataPersistentPreRun() {
	p.setDefaultDocumentParser()

	originPersistentPreRunE := p.cmd.PersistentPreRunE

	p.cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if p.preRunE1 != nil {
			if err := p.preRunE1(cmd, args); err != nil {
				return err
			}
		}

		if p.withConfig {
			p.configPath = util2.GetAbsolutePath(p.configPath)
			klog.V(1).Infof("read config from: %s", p.configPath)
		} else {
			p.configPath = ""
		}

		if p.data != nil {
			// 支持空解析, viper绑定默认数据结构,通过flag override
			var reader io.Reader = strings.NewReader("")
			if p.configPath != "" {
				documentParser, err := File2DocumentParser(p.configPath, p.scheme)
				if err != nil {
					return errors.Wrapf(err, "pcmd:parse:File2Reader:NewDocumentParser: %s", p.configPath)
				}
				p.documentParser = documentParser
				reader, err = p.GetReader()
				if err != nil {
					return errors.Wrapf(err, "pcmd:parse:File2Reader:Reader: %s", p.configPath)
				}
			}

			if err := ReaderFillData(p.viper, reader, p.data); err != nil {
				return errors.Wrapf(err, "pcmd:parse:Reader2Data: %s", p.configPath)
			}

			if p.preRunE2 != nil {
				if err := p.preRunE2(cmd, args); err != nil {
					return err
				}
			}
		}

		// init
		if err := p.dataInit(); err != nil {
			return err
		}

		// go validate
		if err := p.toValidate(); err != nil {
			return err
		}

		if originPersistentPreRunE != nil {
			return originPersistentPreRunE(cmd, args)
		}

		return nil
	}
}

func (p *PhasesCmd) dataInit() error {
	if p.data != nil {
		v, ok := p.data.(HasInit)
		if ok {
			return v.Init()
		}
	}
	return nil
}

func (p *PhasesCmd) toValidate() error {
	if !p.shouldValidate || p.v == nil {
		return nil
	}

	if p.data != nil {
		if err := p.data.ValidateStruct(p.v); err != nil {
			return err
		}

		v, ok := p.data.(HasValidate)
		if ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}

	for _, ef := range p.extraFlagStructs {
		if err := p.v.Struct(ef); err != nil {
			return err
		}
		ev, ok := ef.(HasValidate)
		if ok {
			if err := ev.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *PhasesCmd) GetValidator() *validator.Validate {
	return p.v
}

func (p *PhasesCmd) dataToDocumentPostRun() {
	originPostRunE := p.cmd.PostRunE
	p.cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		if originPostRunE != nil {
			if err := originPostRunE(cmd, args); err != nil {
				return err
			}
		}

		if p.postRunE1 != nil {
			if err := p.postRunE1(cmd, args); err != nil {
				return err
			}
		}

		if p.configWriteBack {
			if err := WriteBackFile(p.configPath, p.documentParser.Dp, p.codec(), p.data, p.gvk); err != nil {
				return errors.Wrapf(err, "pcmd:parse:WriteBackFile: %s", p.configPath)
			}
			klog.V(5).Info("write back success")
		}

		if p.postRunE2 != nil {
			if err := p.postRunE2(cmd, args); err != nil {
				return err
			}
		}

		return nil
	}
}

func (p *PhasesCmd) GetConfigPath() string {
	return p.configPath
}

func (p *PhasesCmd) codec() serializer.CodecFactory {
	return serializer.NewCodecFactory(p.scheme)
}

func (p *PhasesCmd) GetDataYaml() ([]byte, error) {
	return ObjectToYaml(p.codec(), p.data, p.gvk)
}

// SetPreRun
//
//	p1: load data前执行
//	p2: load data后执行
func (p *PhasesCmd) SetPreRun(p1, p2 CobraRun) {
	p.preRunE1 = p1
	p.preRunE2 = p2
}

// SetPostRun
//
//	p1: writeback data前执行
//	p2: writeback data后执行
func (p *PhasesCmd) SetPostRun(p1, p2 CobraRun) {
	p.postRunE1 = p1
	p.postRunE2 = p2
}

func (p *PhasesCmd) GetReader() (io.Reader, error) {
	r, err := p.documentParser.Reader(p.documentParser2Reader)
	return r, err
}

// AppendPhases 添加原生workflow.Phase: 需要自行断言
func (p *PhasesCmd) AppendPhases(phases ...workflow.Phase) {
	if p.finished {
		klog.Errorln("pcmd:AppendPhases: forbidden append now, skip append")
		return
	}

	if p.withConfirm && p.firstAppend {
		p.firstAppend = false
		p.Runner.AppendPhase(NewPhaseSpew(p.data).convert2workflowPhase())
		confirmBeforeRun, ok := p.data.(HasConfirmBeforeRun)
		if ok {
			p.Runner.AppendPhase(NewPhaseRawfn(confirmBeforeRun.ConfirmBeforeRun).convert2workflowPhase())
		}
		p.Runner.AppendPhase(NewPhaseConfirm().convert2workflowPhase())
	}

	for _, phase := range phases {
		p.Runner.AppendPhase(phase)
	}
}

// AppendPcmdPhases 添加pcmd.Phase: 不需要断言
func (p *PhasesCmd) AppendPcmdPhases(phases ...PhaseInterface) {
	for _, phase := range phases {
		p.AppendPhases(phase.convert2workflowPhase())
	}
}

func (p *PhasesCmd) AppendPhaseRunDataFn(use string, f func(workflow.RunData) error) {
	p.AppendPhases(workflow.Phase{
		Name: use,
		Run:  f,
	})
}

// AppendPhaseRawFn 单个phase, 带name
func (p *PhasesCmd) AppendPhaseRawFn(use string, f func() error) {
	p.AppendPhaseRunDataFn(use, func(r workflow.RunData) error {
		return f()
	})
}

// AppendPhaseRawFns 每个phase没有name
func (p *PhasesCmd) AppendPhaseRawFns(fns ...func() error) {
	nfs := make([]NameFn, 0, len(fns))
	for _, f := range fns {
		nfs = append(nfs, NameFn{
			Fn: f,
		})
	}
	p.AppendPhaseNamedRawFns(nfs...)
}

// AppendPhaseNamedRawFns 每个phase带有name
func (p *PhasesCmd) AppendPhaseNamedRawFns(nfs ...NameFn) {
	phases := make([]workflow.Phase, 0, len(nfs))
	for _, nf := range nfs {
		phases = append(phases, workflow.Phase{
			Name: nf.Name,
			Run: func(data workflow.RunData) error {
				return nf.Fn()
			},
		})
	}
	p.AppendPhases(phases...)
}

func (p *PhasesCmd) PersistentFlags() *pflag.FlagSet {
	return p.cmd.PersistentFlags()
}

func (p *PhasesCmd) Cmd() *cobra.Command {

	defer func() {
		p.finished = true
	}()

	if p.finished {
		klog.Fatalf("pcmd:Cmd: app has been exported")
	}

	p.finalize()
	return p.cmd
}
