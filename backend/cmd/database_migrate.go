package cmd

import (
	"fmt"
	"log/slog"

	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/spf13/cobra"
)

func NewDbMigrateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "migrate",
		Short: "Runs database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("Running database migrations...", "file", Cfg.App.DbFilename)

			db, err := utils.ConnectToDatabase(Cfg.App.DbFilename)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}

			err = utils.MigrateDatabase(db)
			if err != nil {
				return fmt.Errorf("failed to migrate database: %w", err)
			}

			slog.Info("Database migrations completed successfully")

			return nil
		},
	}

	return cmd
}
