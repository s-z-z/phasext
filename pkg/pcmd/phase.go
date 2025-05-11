package pcmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"

	"github.com/suzi1037/pcmd/pkg/cprt"
)

var (
	ErrUserAbort = errors.New("won't proceed; the user didn't answer (Y|y) in order to continue")
)

type Phase struct {
	// name of the phase.
	// Phase name should be unique among peer phases (phases belonging to
	// the same workflow or phases belonging to the same parent phase).
	Name string

	// Aliases returns the aliases for the phase.
	Aliases []string

	// Short description of the phase.
	Short string

	// Long returns the long description of the phase.
	Long string

	// Example returns the example for the phase.
	Example string

	// Hidden define if the phase should be hidden in the workflow help.
	// e.g. PrintFilesIfDryRunning phase in the kubeadm init workflow is candidate for being hidden to the users
	Hidden bool

	// RunAllSiblings allows to assign to a phase the responsibility to
	// run all the sibling phases
	// Nb. phase marked as RunAllSiblings can not have Run functions
	RunAllSiblings bool

	// Run 封装
	Run func() error

	// InheritFlags defines the list of flags that the cobra command generated for this phase should Inherit
	// from local flags defined in the parent command / or additional flags defined in the phase runner.
	// If the values is not set or empty, no flags will be assigned to the command
	// Nb. global flags are automatically inherited by nested cobra command
	InheritFlags []string

	// Dependencies is a list of phases that the specific phase depends on.
	Dependencies []string

	// DONTADD_PHASE
	DontAdd bool
}

func (p Phase) convert2workflowPhase() workflow.Phase {
	return workflow.Phase{
		Name:           p.Name,
		Aliases:        p.Aliases,
		Short:          p.Short,
		Long:           p.Long,
		Example:        p.Example,
		Hidden:         p.Hidden,
		RunAllSiblings: p.RunAllSiblings,
		Run: func(data workflow.RunData) error {
			return p.Run()
		},
		InheritFlags: p.InheritFlags,
		Dependencies: p.Dependencies,
	}
}

func NewPhaseSpew(o any) Phase {
	return Phase{
		Name:   "print",
		Short:  "print config",
		Hidden: true,
		Run: func() error {
			cprt.SpewInfo(o)
			return nil
		},
	}
}

func NewPhaseRawfn(run func() error) Phase {
	return Phase{
		Name:   "_rawfn",
		Hidden: true,
		Run:    run,
	}
}

func NewPhaseConfirm() Phase {
	return Phase{
		Name:   "_confirm",
		Hidden: true,
		Run: func() error {
			return InteractivelyConfirmAction("Are you sure you want to proceed?")
		},
	}
}

func InteractivelyConfirmAction(question string) error {
	fmt.Printf("%s [y/N]: ", question)
	r := os.Stdin
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "couldn't read from standard input")
	}
	answer := scanner.Text()
	if strings.EqualFold(answer, "y") || strings.EqualFold(answer, "yes") {
		return nil
	}

	return ErrUserAbort
}
