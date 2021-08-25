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
	"os/signal"
	"strings"
	"time"
)

type server struct {
	store  *Store
	router *mux.Router
}

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

	store := NewStore(pool)

	log.Infof("Connected!")

	if !skipMigration {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Fatalf("Unable to acquire a database connection: %v", err)
		}
		conn.Release()
	}

	ctx, cancel := context.WithCancel(context.Background())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	scheduler := event.NewScheduler(pool, eventListeners)
	scheduler.CheckEventsInInterval(ctx, time.Minute)

	scheduler.Schedule("checkStatus", "DLH418", time.Now().Add(1*time.Minute))
	scheduler.Schedule("checkStatus", "JAF7DY", time.Now().Add(2*time.Minute))

	go func() {
		for range interrupt {
			log.Println("\n❌ Interrupt received closing...")
			cancel()
		}
	}()

	listenAddr := viper.GetString("listen")
	log.Infof("Starting HTTP server at %s...", listenAddr)
	router := mux.NewRouter()

	srv := newServer(store, router)
	router.Handle("/", cors.AllowAll().Handler(srv.initHandlers()))
	err = http.ListenAndServe(listenAddr, router)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}

	<-ctx.Done()
	log.Info("HTTP server terminated")
}

func newServer(store *Store, router *mux.Router) server {
	return server{store: store, router: router}
}
