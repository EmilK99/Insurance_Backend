package main

import (
	"flight_app/app"
	_ "flight_app/docs"
	_ "github.com/jackc/pgx/v4"
	"github.com/spf13/cobra"
)

func main() {
	var configPath string
	var skipMigration bool

	rootCmd := cobra.Command{
		Use:     "flight-app",
		Version: "v1.0",
		Run: func(cmd *cobra.Command, args []string) {
			app.Run(configPath, skipMigration)
		},
	}

	rootCmd.Flags().StringVarP(&configPath, "config", "c", "config/flight_app.toml", "Config file path")
	rootCmd.Flags().BoolVarP(&skipMigration, "skip-migration", "s", false, "Skip migration")
	err := rootCmd.Execute()
	if err != nil {
		// Required arguments are missing, etc
		return
	}
}
