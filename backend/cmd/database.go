package cmd

import "github.com/spf13/cobra"

func NewDatabaseCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "database",
		Short:   "Database management commands",
		Aliases: []string{"db"},
	}

	return cmd
}
