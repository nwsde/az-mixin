package main

import (
	"github.com/deislabs/porter-az/pkg/az"
	"github.com/spf13/cobra"
)

func buildUpgradeCommand(m *az.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Execute the invoke functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute()
		},
	}
	return cmd
}
