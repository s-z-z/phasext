package cmd

import (
	"github.com/spf13/cobra"

	"github.com/suzi1037/pcmd/pkg/pcmd"
)

func newCmdPoC() *cobra.Command {
	p := cmdFactory.CreateWithProp(
		pcmd.CmdProp{
			Use:    "poc",
			Hidden: true,
		},
	)

	return p.Cmd()
}
