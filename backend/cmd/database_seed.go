package cmd

import (
	"fmt"
	"log/slog"

	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/spf13/cobra"
)

var includeTestData bool

var dbSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Fills the database with initial data",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("Running database seed...", "file", Cfg.App.DbFilename, "with_test_data", includeTestData)

		db, err := utils.ConnectToDatabase(Cfg.App.DbFilename)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		utils.SeedDatabase(db, includeTestData)

		slog.Info("Seed completed successfully!")

		return nil
	},
}

func init() {
	databaseCmd.AddCommand(dbSeedCmd)

	dbSeedCmd.Flags().BoolVarP(&includeTestData, "test-data", "t", false, "Additional test data (dummy products, etc.) to generate")
}
