package cmd

import (
	"bufio"
	"fmt"
	"net/mail"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func NewUserCreateCmd() *cobra.Command {
	var (
		username string
		email    string
		isAdmin  bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			username, email, isAdmin, err := collectUserData(cmd, username, email, isAdmin)
			if err != nil {
				return err
			}

			userService, cleanup, err := setupUserService()
			if err != nil {
				return err
			}
			defer cleanup()

			fmt.Printf("\nCreating user '%s' (%s). Admin: %v\n", username, email, isAdmin)

			if err := userService.CreateUser(username, email, isAdmin); err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			fmt.Println("✅ User created successfully!")

			return nil
		},
	}

	cmd.Flags().StringVar(&username, "username", "", "username (optional, otherwise interactive)")
	cmd.Flags().StringVar(&email, "email", "", "email address (optional, otherwise interactive)")
	cmd.Flags().BoolVar(&isAdmin, "admin", false, "Set admin rights (optional, otherwise interactive)")

	return cmd
}

func collectUserData(cmd *cobra.Command, u, e string, admin bool) (user, email string, adm bool, err error) {
	reader := bufio.NewReader(os.Stdin)

	// Username Prompt
	if u == "" {
		fmt.Print("Username: ")

		input, _ := reader.ReadString('\n')
		u = strings.TrimSpace(input)
	}

	if u == "" {
		return "", "", false, fmt.Errorf("abort: username cannot be empty")
	}

	// Email Prompt & Validation
	if e == "" {
		fmt.Print("Email Address: ")

		input, _ := reader.ReadString('\n')
		e = strings.TrimSpace(input)
	}

	if e == "" {
		return "", "", false, fmt.Errorf("abort: email cannot be empty")
	}

	if _, err := mail.ParseAddress(e); err != nil {
		return "", "", false, fmt.Errorf("abort: invalid email address format")
	}

	// Admin Prompt
	if !cmd.Flags().Changed("admin") {
		fmt.Print("Should the user have admin rights? (y/N): ")

		input, _ := reader.ReadString('\n')

		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			admin = true
		}
	}

	return u, e, admin, nil
}
