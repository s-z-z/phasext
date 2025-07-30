package pcmd

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"

	"github.com/s-z-z/phasext/pkg/pcmd/util"
)

// WithData 绑定数据
func WithData(data WareHouse) Option {
	return func(p *PhasesCmd) {
		if p.data != nil {
			klog.Fatalf("WithData: can only be called once")
		}
		p.data = data
	}
}

func WithRunE(runE func(cmd *cobra.Command, args []string) error) Option {
	return func(p *PhasesCmd) {
		p.cmd.RunE = func(cmd *cobra.Command, args []string) error {
			if runE != nil {
				if err := runE(cmd, args); err != nil {
					return err
				}
			}
			return p.Runner.Run(args)
		}
	}
}

// WithConfig 支持配置文件参数，默认参数名是config
func WithConfig() Option {
	return WithConfigSpecFlag("")
}

// WithSpecConfigPath 支持配置文件参数，指定configSpecPath
func WithSpecConfigPath(configPath string) Option {
	return func(p *PhasesCmd) {
		if p.withConfig {
			klog.Fatalf("WithConfig: can only be called once")
		}
		p.withConfig = true
		p.configPath = configPath
	}
}

func WithConfirm() Option {
	return func(p *PhasesCmd) {
		p.withConfirm = true
	}
}

// WithConfigSpecFlag 支持配置文件参数，指定参数名
func WithConfigSpecFlag(cFlag string) Option {
	return func(p *PhasesCmd) {

		if p.withConfig {
			klog.Fatalf("WithConfig: can only be called once")
		}

		if cFlag != "" {
			p.configFlag = cFlag
		}
		p.withConfig = true
	}
}

// WithExportOverrideFlags 导出export tag参数至命令行，flag > file
// specIncludes不指定, 默认导出所有export=true字段
// specIncludes指定, 仅导出specIncludes中的字段, 值为struct字段名
func WithExportOverrideFlags(specIncludes ...string) Option {
	return func(p *PhasesCmd) {
		p.exportOverrideFlags = true
		p.specExportIncludeFlags = specIncludes
	}
}

func WithViperFn(vfn func(v *viper.Viper)) Option {
	return func(p *PhasesCmd) {
		p.viperFn = vfn
	}
}

// WithPhaseBind 支持phase单步，跳步执行
func WithPhaseBind() Option {
	return func(p *PhasesCmd) {
		p.bindToCommand = true
	}
}

// WithPersistentExportedFlag 导出的命令绑定到 PersistentFlags，意味着所有的phase子命令将自动继承
func WithPersistentExportedFlag() Option {
	return func(p *PhasesCmd) {
		p.persistentExportedFlag = true
	}
}

// WithExtraFlagStruct 添加额外的flag
//
//	name: flag name
//	usage: flag usage
//	v: 传入指针
func WithExtraFlagStruct(v any, specIncludes ...string) Option {
	return func(p *PhasesCmd) {
		p.extraFlagStructs = append(p.extraFlagStructs, v)
		if len(specIncludes) > 0 {
			p.specExtraExportIncludeFlags = append(p.specExtraExportIncludeFlags, specIncludes...)
		}
	}
}

// WithConfigWriteBack 配置回写到文件
func WithConfigWriteBack() Option {
	return func(p *PhasesCmd) {
		p.configWriteBack = true
	}
}

// WithDocumentParser 自定义解析器：高级用法
func WithDocumentParser(fn DocumentParser2Redaer) Option {
	return func(p *PhasesCmd) {
		p.documentParser2Reader = fn
	}
}

// WithRunnerDataInitializer 自定义runner数据初始化函数：高级用法
func WithRunnerDataInitializer(fn util.RunnerDataInitializer) Option {
	return func(p *PhasesCmd) {
		p.runnerDataInitializer = fn
	}
}

func WithScheme(s *runtime.Scheme) Option {
	return func(p *PhasesCmd) {
		if p.scheme != nil {
			klog.Fatalf("scheme can only be set once")
		}
		p.scheme = s
	}
}

func WithValidator(v *validator.Validate) Option {
	return func(p *PhasesCmd) {
		p.v = v
	}
}

func WithoutValidate() Option {
	return func(p *PhasesCmd) {
		p.shouldValidate = false
	}
}

func WithPreRun(p1, p2, p3 CobraRun) Option {
	return func(p *PhasesCmd) {
		p.preRunE1 = p1
		p.preRunE2 = p2
		p.preRunE3 = p3
	}
}

func WithPostRun(p1, p2 CobraRun) Option {
	return func(p *PhasesCmd) {
		p.postRunE1 = p1
		p.postRunE2 = p2
	}
}
