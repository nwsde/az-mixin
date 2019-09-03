package main

import (
	"github.com/deislabs/porter-az/pkg/az"
	"github.com/spf13/cobra"
)

func buildUninstallCommand(m *az.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Execute the uninstall functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute()
		},
	}
	return cmd
}
