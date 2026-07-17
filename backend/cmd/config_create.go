package cmd

import (
	"fmt"
	"log/slog"

	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configCreateForce    bool
	configCreateFilename string
)

func NewConfigCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new configuration file with default values",
		Annotations: map[string]string{
			skipConfigValidationAnnotation: "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			const defaultFrontendURL = "https://localhost:3000"

			viper.SetDefault("app.frontend_url", defaultFrontendURL)
			viper.SetDefault("app.cors_allow_origins", []string{defaultFrontendURL})

			viper.SetDefault("mailer.dsn", "smtp://localhost:1025")
			viper.SetDefault("mailer.from", "kasseapparat@example.com")
			viper.SetDefault("mailer.subject_prefix", "[Kasseapparat]")

			viper.SetDefault("format.currency.locale", "de-DE")
			viper.SetDefault("format.currency.code", "EUR")
			viper.SetDefault("format.currency.fraction_digits_min", 0)
			viper.SetDefault("format.currency.fraction_digits_max", config.DefaultMinorUnit)

			viper.SetDefault("vat_rates", config.DefaultVatRates)
			viper.SetDefault("payment_methods", config.DefaultPaymentMethods)

			filename := configCreateFilename

			var writeErr error
			if configCreateForce {
				writeErr = viper.WriteConfigAs(filename)
			} else {
				writeErr = viper.SafeWriteConfigAs(filename)
			}

			if writeErr != nil {
				if _, ok := writeErr.(viper.ConfigFileAlreadyExistsError); ok {
					return fmt.Errorf(
						"file %s already exists or was not able to be created: %w",
						filename,
						writeErr,
					)
				}

				return fmt.Errorf("error writing the config: %w", writeErr)
			}

			slog.Info("Configuration file created successfully", "filename", filename)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&configCreateForce, "force", "f", false, "Overwrite existing config file if it already exists")
	cmd.Flags().
		StringVarP(&configCreateFilename, "output", "o", "config/config.yaml", "Filename for the generated config file")

	return cmd
}
