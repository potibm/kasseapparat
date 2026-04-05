package cmd

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/spf13/cobra"
)

var (
	resetSeed     bool
	resetTestData bool
	resetForce    bool
)

func NewDbResetCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "reset",
		Short: "Deletes the database, runs migrations and (optionally) seeds it with initial data",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !resetForce {
				fmt.Printf("WARNING: The database '%s' will be COMPLETELY deleted!\n", Cfg.App.DbFilename)
				fmt.Print("Are you sure? [y/N]: ")

				reader := bufio.NewReader(os.Stdin)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))

				if response != "y" && response != "yes" {
					slog.Info("Reset aborted by user")

					return nil
				}
			}

			slog.Warn("Starting database reset...", "file", Cfg.App.DbFilename)

			db, err := utils.ConnectToDatabase(Cfg.App.DbFilename)
			if err != nil {
				return fmt.Errorf("could not connect to database: %w", err)
			}

			defer func() {
				if closeErr := utils.CloseDatabase(db); closeErr != nil {
					slog.Error("Failed to close database connection", "error", closeErr)
				}
			}()

			slog.Info("Deleting old tables...")

			if err := utils.PurgeDatabase(db); err != nil {
				return fmt.Errorf("error while deleting tables: %w", err)
			}

			slog.Info("Rebuilding table structure...")

			if err := utils.MigrateDatabase(db); err != nil {
				return fmt.Errorf("error while rebuilding table structure: %w", err)
			}

			if resetSeed || resetTestData {
				slog.Info("Running seeding...", "with_test_data", resetTestData)
				utils.SeedDatabase(db, resetTestData)
			}

			slog.Info("Database reset completed successfully!")

			return nil
		},
	}

	cmd.Flags().BoolVarP(&resetSeed, "seed", "s", false, "Runs a seed after the reset")
	cmd.Flags().BoolVarP(&resetTestData, "test-data", "t", false, "Adds test data (implies --seed)")
	cmd.Flags().BoolVarP(&resetForce, "force", "f", false, "Skips the confirmation prompt (for CI/CD)")

	return cmd
}
