package cmd

import (
	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	"github.com/potibm/kasseapparat/internal/app/service/user"
	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/spf13/cobra"
)

func NewUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "User management commands",
	}

	return cmd
}

func setupUserService() (*user.UserService, func(), error) {
	db, err := utils.ConnectToDatabase(Cfg.App.DbFilename)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() { _ = utils.CloseDatabase(db) }

	repo := sqlite.NewRepository(db, Cfg.Format.Currency.FractionDigitsMax)

	mailerClient, err := mailer.NewMailer(Cfg.Mailer.DSN)
	if err != nil {
		cleanup()

		return nil, nil, err
	}

	return user.NewUserService(repo, mailerClient), cleanup, nil
}
