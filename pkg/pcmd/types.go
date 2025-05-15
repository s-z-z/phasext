package pcmd

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"
)

const (
	DefaultConfigPath      = "./config.yaml"
	DefaultConfigFlag      = "config"
	DefaultGoValidate      = true
	DefaultConfigWriteBack = false
)

type CobraRun func(cmd *cobra.Command, args []string) error

type WareHouse interface {
	schema.ObjectKind
	runtime.Object
	ValidateStruct(v *validator.Validate) error
}

type HasInit interface {
	Init() error
}

type HasValidate interface {
	Validate() error
}

type HasConfirmBeforeRun interface {
	ConfirmBeforeRun() error
}

type NameFn struct {
	Name string
	Fn   func() error
}

type Generator func() workflow.Phase

type Option func(cmd *PhasesCmd)

type Factory func(string, ...Option) *PhasesCmd

type CmdProp struct {
	Use                    string
	Aliases                []string
	SuggestFor             []string
	Short                  string
	GroupID                string
	Long                   string
	Example                string
	ValidArgs              []string
	ValidArgsFunction      func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)
	Args                   cobra.PositionalArgs
	ArgAliases             []string
	BashCompletionFunction string
	Version                string
	Hidden                 bool
	SilenceErrors          bool
	SilenceUsage           bool
}
