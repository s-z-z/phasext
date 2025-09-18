package util

import (
	"github.com/spf13/cobra"

	"github.com/s-z-z/phasext/workflow"
)

type RunnerDataInitializer func(cmd *cobra.Command, args []string) (workflow.RunData, error)

func OnlyArgsDataInitializer(cmd *cobra.Command, args []string) (workflow.RunData, error) {
	return args, nil
}
