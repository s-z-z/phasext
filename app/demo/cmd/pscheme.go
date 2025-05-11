package cmd

import (
	"github.com/spf13/cobra"
	"github.com/suzi1037/pcmd/app/demo/apis"
	"github.com/suzi1037/pcmd/app/demo/apis/scheme"
	"github.com/suzi1037/pcmd/pkg/cprt"
)

func newCmdPrintScheme() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pscheme",
		Hidden:  true,
		Aliases: []string{"ps"},
		Run: func(cmd *cobra.Command, args []string) {
			s := scheme.Scheme
			for gvk, tp := range s.AllKnownTypes() {

				if gvk.Group != apis.GroupName {
					continue
				}
				cprt.Ok("group: %s", gvk.Group)
				cprt.Info("version: %s", gvk.Version)
				cprt.Warning("kind: %s", gvk.Kind)
				cprt.Error("type: %s", tp)
				cprt.Debug("-----------------------------")
			}
		},
	}
	return cmd
}
