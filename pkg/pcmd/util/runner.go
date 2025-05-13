package util

import (
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"
)

type RunnerDataInitializer func(cmd *cobra.Command, args []string) (workflow.RunData, error)

func OnlyArgsDataInitializer(cmd *cobra.Command, args []string) (workflow.RunData, error) {
	return args, nil
}
