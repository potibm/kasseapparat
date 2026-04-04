package cmd

import "github.com/spf13/cobra"

var databaseCmd = &cobra.Command{
	Use:     "database",
	Short:   "Database management commands",
	Aliases: []string{"db"},
}

func init() {
	rootCmd.AddCommand(databaseCmd)
}
