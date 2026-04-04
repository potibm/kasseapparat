package cmd

import "github.com/spf13/cobra"

func NewUserCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "user",
		Short: "User management commands",
	}

	return cmd
}
