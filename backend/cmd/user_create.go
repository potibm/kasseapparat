package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	RunE: func(cmd *cobra.Command, args []string) error {

		email, _ := cmd.Flags().GetString("email")
		fmt.Printf("Erstelle User mit Email: %s\n", email)

		// Hier rufst du z.B. userlogic.CreateUser(cfg, email) auf
		return nil
	},
}

func init() {
	userCmd.AddCommand(userCreateCmd)

	userCreateCmd.Flags().String("email", "", "E-Mail des Benutzers")
	userCreateCmd.MarkFlagRequired("email")
}
