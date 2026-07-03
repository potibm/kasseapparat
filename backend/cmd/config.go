package cmd

import (
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Config commands",
		Annotations: map[string]string{
			skipConfigValidationAnnotation: "true",
		},
	}
}
