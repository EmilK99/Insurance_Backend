package app

import (
	"context"
	"flight_app/app/api/flightaware_api"
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
	port    string
	aeroApi flightaware_api.AeroAPI
}

func newServer(store *store2.Store, ctx context.Context, solClient *sc.Client, port string) *server {
	var aeroAPI flightaware_api.AeroAPI
	aeroAPI.Username = viper.GetString("aeroapi_username")
	aeroAPI.APIKey = viper.GetString("aeroapi_apikey")
	aeroAPI.URL = "http://" + aeroAPI.Username + ":" + aeroAPI.APIKey + "@flightxml.flightaware.com/json/FlightXML2/"
	aeroAPI.URLc = "http://" + aeroAPI.Username + ":" + aeroAPI.APIKey + "@flightxml.flightaware.com/json/FlightXML2c/"
	return &server{ctx: ctx, store: store, solClient: solClient, client: &payments.Client{}, port: port, aeroApi: aeroAPI}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	router := mux.NewRouter()
	//fee calculation
	router.HandleFunc("/flightaware_api/flights", s.HandleGetFlights).Methods("POST")

	//fee calculation
	router.HandleFunc("/flightaware_api/calculate", s.CalculateFeeHandler).Methods("POST")

	//contract create
	router.HandleFunc("/flightaware_api/contract/create", s.HandleCreateContract).Methods("POST")

	//get contracts
	router.HandleFunc("/flightaware_api/contracts", s.HandleGetContracts).Methods("POST")

	//get payout history
	router.HandleFunc("/flightaware_api/payouts", s.HandleGetPayouts).Methods("POST")

	//alerts webhook
	router.HandleFunc("/flightaware_api/alerts", s.HandleAlertWebhook).Methods("GET", "POST")

	//register alerts
	router.HandleFunc("/flightaware_api/alerts/register", s.HandleRegisterAlertsEndpoint).Methods("GET")

	//paypal success
	router.HandleFunc("/flightaware_api/success", s.HandlerSuccess).Methods("GET")

	//paypal cancel
	router.HandleFunc("/flightaware_api/cancel", s.HandlerCancel).Methods("GET")

	//withdraw successful contract
	router.HandleFunc("/flightaware_api/withdraw", s.HandleWithdrawPremium).Methods("POST")

	s.router = cors.AllowAll().Handler(router)
}
