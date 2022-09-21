package main

import (
	"get.porter.sh/mixin/az/pkg/az"
	"github.com/spf13/cobra"
)

func buildBuildCommand(m *az.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Generate Dockerfile lines for the bundle invocation image",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Build(cmd.Context())
		},
	}
	return cmd
}
