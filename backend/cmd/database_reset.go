package cmd

import (
	"fmt"
	"log/slog"

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
			if !resetForce && !confirmReset(Cfg.App.DbFilename) {
				slog.Info("Reset aborted by user")

				return nil
			}

			return performDatabaseReset(Cfg.App.DbFilename, resetSeed || resetTestData, resetTestData)
		},
	}

	cmd.Flags().BoolVarP(&resetSeed, "seed", "s", false, "Runs a seed after the reset")
	cmd.Flags().BoolVarP(&resetTestData, "test-data", "t", false, "Adds test data (implies --seed)")
	cmd.Flags().BoolVarP(&resetForce, "force", "f", false, "Skips the confirmation prompt (for CI/CD)")

	return cmd
}

func confirmReset(dbName string) bool {
	return confirm(fmt.Sprintf("WARNING: The database '%s' will be COMPLETELY deleted!", dbName))
}

func performDatabaseReset(dbName string, shouldSeed, withTestData bool) error {
	slog.Warn("Starting database reset...", "file", dbName)

	db, err := utils.ConnectToDatabase(dbName)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	defer func() {
		if closeErr := utils.CloseDatabase(db); closeErr != nil {
			slog.Error("Failed to close database connection", "error", closeErr)
		}
	}()

	// Atomic steps of the reset
	steps := []struct {
		msg string
		fn  func() error
	}{
		{"Deleting old tables...", func() error { return utils.PurgeDatabase(db) }},
		{"Rebuilding table structure...", func() error { return utils.MigrateDatabase(db) }},
	}

	for _, step := range steps {
		slog.Info(step.msg)

		if err := step.fn(); err != nil {
			return fmt.Errorf("reset failed during step '%s': %w", step.msg, err)
		}
	}

	if shouldSeed {
		slog.Info("Running seeding...", "with_test_data", withTestData)
		utils.SeedDatabase(db, withTestData)
	}

	slog.Info("Database reset completed successfully!")

	return nil
}
