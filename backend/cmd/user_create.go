package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewUserCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			email, _ := cmd.Flags().GetString("email")
			fmt.Printf("Create an user with email: %s\n", email)

			// @TODO: Implement user creation logic here
			return nil
		},
	}

	cmd.Flags().String("email", "", "The users email address")
	_ = cmd.MarkFlagRequired("email")

	return cmd
}
