package pcmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/s-z-z/phasext/pkg/cprt"
)

func NewPhaseSpew(o any) PhaseInterface {
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

func NewPhaseRawfn(run func() error) PhaseInterface {
	return Phase{
		Name:   "_rawfn",
		Hidden: true,
		Run:    run,
	}
}

func NewPhaseConfirm() PhaseInterface {
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
