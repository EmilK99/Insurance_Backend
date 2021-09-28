package app

import (
	"context"
	store2 "flight_app/app/store"
	"flight_app/payments"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

type server struct {
	ctx    context.Context
	router http.Handler
	store  *store2.Store
	client *payments.Client
	port   string
}

func newServer(store *store2.Store, ctx context.Context, port string) *server {
	return &server{ctx: ctx, store: store, client: &payments.Client{}, port: port}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	router := mux.NewRouter()
	//fee calculation
	router.HandleFunc("/api/calculate", s.CalculateFeeHandler).Methods("POST")

	//contract create
	router.HandleFunc("/api/contract/create", s.HandleCreateContract).Methods("POST")

	//get contracts
	router.HandleFunc("/api/contracts", s.HandleGetContracts).Methods("POST")

	//get payout history
	router.HandleFunc("/api/payouts", s.HandleGetPayouts).Methods("POST")

	//alerts webhook
	router.HandleFunc("/api/alerts", s.HandleAlertWebhook).Methods("POST")

	//register alerts
	router.HandleFunc("/api/alerts/register", s.HandleRegisterAlertsEndpoint).Methods("GET")

	//paypal redirect
	router.HandleFunc("/api/paypal", s.HandleCreatePaypalOrder).Methods("GET")

	//paypal success
	router.HandleFunc("/api/success", s.HandlerSuccess).Methods("GET")

	//paypal cancel
	router.HandleFunc("/api/cancel", HandlerCancel).Methods("GET")

	//register webhook endpoint
	router.HandleFunc("/api/ipn", s.IPNHandler).Methods("POST")

	s.router = cors.AllowAll().Handler(router)
}
