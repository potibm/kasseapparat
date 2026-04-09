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
			safeCfg := Cfg.RedactConfigForDisplay()

			configJSON, err := json.MarshalIndent(safeCfg, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(configJSON))

			return nil
		},
	}
}
