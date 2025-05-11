package cmd

import (
	"github.com/spf13/cobra"
	v1 "github.com/suzi1037/pcmd/app/demo/apis/v1"
	"github.com/suzi1037/pcmd/app/demo/cmd/phases/direct"
	"github.com/suzi1037/pcmd/pkg/pcmd"
)

func newCmdFoo() *cobra.Command {
	c := &v1.Foo{}
	p := cmdFactory.Create(
		"foo",
		pcmd.WithData(c),
		pcmd.WithExportOverrideFlags(),
	)
	p.AppendPcmdPhases(
		direct.NewPhasePrint(c),
	)
	return p.Cmd()
}
