package app

import (
	"context"
	event "flight_app/app/store"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strings"
)

var eventListeners = event.Listeners{
	"checkStatus": event.CheckStatus,
}

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

	dbURL := viper.GetString("db_url")
	log.Infof("Using DB URL: %s", dbURL)

	pool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer pool.Close()

	store := event.NewStore(pool)

	log.Infof("Connected!")

	store.Conn, err = store.Pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
	}
	defer store.Conn.Release()

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)

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
	router := mux.NewRouter()

	srv := newServer(store, router, port)
	err = srv.client.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize paypal client: %v", err)
	}

	router.Handle("/", cors.AllowAll().Handler(srv.initHandlers()))
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}

	//<-ctx.Done()
	log.Info("HTTP server terminated")
}
