package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Shows the current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			configJson, err := json.MarshalIndent(Cfg, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}

			fmt.Println(string(configJson))

			return nil
		},
	}
}
