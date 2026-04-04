package cmd

import (
	"fmt"

	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Hier kommt die Magie: Viper füllt dein Struct!
		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return fmt.Errorf("error while parsing config: %w", err)
		}

		//		fmt.Printf("Starting server on port: %d\n", cfg.AppConfig.Port) // Beispiel

		// Hier rufst du deine eigentliche Start-Logik auf:
		// return server.Start(cfg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// define flags
	serveCmd.Flags().Int("port", 8080, "Port that the server listens on")
	serveCmd.Flags().String("log-level", "info", "Loglevel (debug, info, warn)")

	// bind flags to viper keys
	viper.BindPFlag("AppConfig.Port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("AppConfig.LogLevel", serveCmd.Flags().Lookup("log-level"))
}
