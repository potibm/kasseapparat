package cmd

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/initializer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version = "dev"

var Cfg config.Config

var rootCmd = &cobra.Command{
	Use:           "kasseapparat",
	Short:         "Kasseapparat ist a POS system for demoparties",
	Version:       Version,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.Unmarshal(&Cfg); err != nil {
			return fmt.Errorf("error parsing the config: %w", err)
		}

		if err := Cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}

		if !cmd.Flags().Changed("log-format") {
			if cmd.Name() != "serve" {
				Cfg.App.LogFormat = "text"
			}
		}

		setupLogger(Cfg.App.LogFormat, Cfg.App.LogLevel)

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("log-level", "info", "Log Level (debug, info, warn, error)")
	viper.BindPFlag("app.log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.PersistentFlags().String("log-format", "json", "Log Format (json, text)")
	viper.BindPFlag("app.log_format", rootCmd.PersistentFlags().Lookup("log-format"))

	rootCmd.PersistentFlags().String("db-file", "kasseapparat.db", "Dateiname der SQLite Datenbank")
	viper.BindPFlag("app.db_filename", rootCmd.PersistentFlags().Lookup("db-file"))
}

func initConfig() {
	_ = godotenv.Load()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()

	config.InitViper()
}

func setupLogger(format, level string) {
	initializer.InitLogger(format, level)
}
