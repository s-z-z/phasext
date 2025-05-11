package util

import (
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"
)

type RunnerDataInitializer func(cmd *cobra.Command, args []string) (workflow.RunData, error)

func SimpleDataInitializer(c workflow.RunData) RunnerDataInitializer {
	return func(cmd *cobra.Command, args []string) (workflow.RunData, error) {
		return c, nil
	}
}

func SetSimpleDataInitializer(r *workflow.Runner, c workflow.RunData) {
	r.SetDataInitializer(SimpleDataInitializer(c))
}
