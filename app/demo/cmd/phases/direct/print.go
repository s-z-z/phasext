package direct

import (
	"github.com/suzi1037/pcmd/pkg/cprt"
	"github.com/suzi1037/pcmd/pkg/pcmd"
)

type HasPrint interface {
	Print() error
}

func NewPhasePrint(o any) pcmd.Phase {
	return pcmd.Phase{
		Name:   "print",
		Short:  "print config",
		Hidden: true,
		Run: func() error {
			if p, ok := o.(HasPrint); ok {
				return p.Print()
			}
			cprt.SpewInfo(o)
			return nil
		},
	}
}
