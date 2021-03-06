package app

import (
	"context"
	"flight_app/app/api"
	"flight_app/app/sc"
	store2 "flight_app/app/store"
	"flight_app/payments"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"net/http"
)

type server struct {
	ctx       context.Context
	router    http.Handler
	store     *store2.Store
	client    *payments.Client
	solClient *sc.Client
	port      string
	aeroApi   api.AeroAPI
}

func newServer(store *store2.Store, ctx context.Context, solClient *sc.Client, port string) *server {
	var aeroAPI api.AeroAPI
	aeroAPI.Username = viper.GetString("aeroapi_username")
	aeroAPI.APIKey = viper.GetString("aeroapi_apikey")
	aeroAPI.URL = "https://" + aeroAPI.Username + ":" + aeroAPI.APIKey + "@flightxml.flightaware.com/json/FlightXML2/"

	return &server{ctx: ctx, store: store, solClient: solClient, client: &payments.Client{}, port: port, aeroApi: aeroAPI}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	router := mux.NewRouter()
	//fee calculation
	router.HandleFunc("/api/flights", s.HandleGetFlights).Methods("POST")

	//fee calculation
	router.HandleFunc("/api/calculate", s.CalculateFeeHandler).Methods("POST")

	//contract create
	router.HandleFunc("/api/contract/create", s.HandleCreateContract).Methods("POST")

	//get contracts
	router.HandleFunc("/api/contracts", s.HandleGetContracts).Methods("POST")

	//get payout history
	router.HandleFunc("/api/payouts", s.HandleGetPayouts).Methods("POST")

	//alerts webhook
	router.HandleFunc("/api/alerts", s.HandleAlertWebhook).Methods("GET", "POST")

	//register alerts
	router.HandleFunc("/api/alerts/register", s.HandleRegisterAlertsEndpoint).Methods("GET")

	//paypal success
	router.HandleFunc("/api/success", s.HandlerSuccess).Methods("GET")

	//paypal cancel
	router.HandleFunc("/api/cancel", s.HandlerCancel).Methods("GET")

	//withdraw successful contract
	router.HandleFunc("/api/withdraw", s.HandleWithdrawPremium).Methods("POST")

	s.router = cors.AllowAll().Handler(router)
}
