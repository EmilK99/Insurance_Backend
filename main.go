package main

import (
	"context"
	"flight_app/api"
	"flight_app/contract"
	_ "flight_app/docs"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func initViper(configPath string) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("flight-app")

	viper.SetDefault("loglevel", "debug")
	viper.SetDefault("listen", "localhost:8080")
	viper.SetDefault("db_url", "host=localhost user=postgres database=users")

	if configPath != "" {
		log.Infof("Parsing config: %s", configPath)
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Unable to read config file: %s", err)
		}
	} else {
		log.Infof("Config file is not specified.")
	}
}

func run(configPath string, skipMigration bool) {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	initViper(configPath)

	logLevelString := viper.GetString("loglevel")
	logLevel, err := log.ParseLevel(logLevelString)
	if err != nil {
		log.Fatalf("Unable to parse loglevel: %s", logLevelString)
	}

	log.SetLevel(logLevel)

	dbURL := viper.GetString("db_url")
	log.Infof("Using DB URL: %s", dbURL)

	pool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer pool.Close()
	log.Infof("Connected!")

	if !skipMigration {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Fatalf("Unable to acquire a database connection: %v", err)
		}
		conn.Release()
	}

	listenAddr := viper.GetString("listen")
	log.Infof("Starting HTTP server at %s...", listenAddr)
	http.Handle("/", cors.AllowAll().Handler(InitHandlers(pool)))
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}

	log.Info("HTTP server terminated")
}

func InitHandlers(pool *pgxpool.Pool) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/calculate",
		func(w http.ResponseWriter, r *http.Request) {
			api.CalculateFeeHandler(w, r)
		}).Methods("POST")

	r.HandleFunc("/contract/create",
		func(w http.ResponseWriter, r *http.Request) {
			contract.CreateContract(pool, w, r)
		}).Methods("POST")

	// TODO: implement handlers for flight app

	return r
}

func main() {
	var configPath string
	var skipMigration bool

	rootCmd := cobra.Command{
		Use:     "flight-app",
		Version: "v1.0",
		Run: func(cmd *cobra.Command, args []string) {
			run(configPath, skipMigration)
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
