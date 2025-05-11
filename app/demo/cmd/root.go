package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"

	"github.com/suzi1037/pcmd/app/demo/apis/scheme"
	"github.com/suzi1037/pcmd/internal/validators"
	"github.com/suzi1037/pcmd/pkg/pcmd"
)

type cmdFunc func() *cobra.Command

// + Add your command here
var cmdFunctions = []cmdFunc{
	newCmdPrintScheme,
	newCmdVersion,
	newCmdPoC,
	newCmdFoo,
}

const (
	CMD = "demo"
)

var (
	v          = validator.New(validator.WithRequiredStructEnabled(), validator.WithPrivateFieldValidation())
	cmdFactory = pcmd.NewPhaseCmdFactory(scheme.Scheme, v)
	rootCmd    = NewRootCmd()
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          CMD,
		SilenceUsage: true,
	}

	cobra.OnInitialize(func() {
		validators.RegisteValidator(v)
	})

	registerSubCommands(rootCmd)

	return rootCmd
}

func Run() error {
	v.SetTagName("v")

	klog.InitFlags(nil)
	defer func() {
		klog.Flush()
	}()

	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	_ = pflag.Set("logtostderr", "false")
	_ = pflag.Set("alsologtostderr", "true")

	_ = pflag.Set("log_file", fmt.Sprintf("%s.log", CMD))

	markHidden()

	klog.Infoln(os.Args)

	return rootCmd.Execute()
}

func registerSubCommands(rootCmd *cobra.Command) {
	for _, fn := range cmdFunctions {
		rootCmd.AddCommand(fn())
	}
}

func markHidden() {
	hiddenFlags := [...]string{
		"log-flush-frequency",
		"alsologtostderr",
		"log-backtrace-at",
		"log-dir",
		"one-output",
		"logtostderr",
		"stderrthreshold",
		"vmodule",
		"add-dir-header",
		"log-file",
		"log-file-max-size",
		"skip-headers",
		"skip-log-headers",
		"version",
	}
	for _, f := range hiddenFlags {
		_ = pflag.CommandLine.MarkHidden(f)
	}
}
