package app

import (
	"context"
	"flight_app/app/sc"
	event "flight_app/app/store"
	"github.com/jackc/pgx/v4"
	"github.com/portto/solana-go-sdk/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
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

func Run(configPath string, skipMigration bool) {
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbURL := viper.GetString("db_url")
	log.Infof("Using DB URL: %s", dbURL)

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer conn.Close(ctx)

	store := event.NewStore(conn)

	log.Infof("Connected!")

	//scheduler := event.NewScheduler(pool, eventListeners)
	//scheduler.CheckEventsInInterval(ctx, 10*time.Second)
	//
	//scheduler.Schedule("checkStatus", "BAW920", time.Now().Add(10*time.Second))
	//scheduler.Schedule("checkStatus", "CCA680", time.Now().Add(20*time.Second))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	listenAddr := viper.GetString("listen") + ":" + port
	log.Infof("Starting HTTP server at %s...", listenAddr)

	programID := viper.GetString("program_id")
	solAccount, _ := sc.NewClient(common.PublicKeyFromString(programID))

	solAccount.CreateSmartContract(ctx, 0)

	srv := newServer(store, ctx, port)
	srv.client.ClientID = viper.GetString("client_id")
	srv.client.SecretID = viper.GetString("secret_id")
	err = srv.client.Initialize(ctx)
	if err != nil {
		log.Fatalf("Unable to initialize paypal client: %v", err)
	}

	srv.configureRouter()

	err = http.ListenAndServe(":"+port, srv)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}

	//<-ctx.Done()
	log.Info("HTTP server terminated")
}
