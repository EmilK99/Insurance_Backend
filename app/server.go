package app

import (
	"context"
	store2 "flight_app/app/store"
	"flight_app/payments"
	"github.com/gorilla/mux"
	"net/http"
)

type server struct {
	ctx    context.Context
	router *mux.Router
	store  *store2.Store
	client *payments.Client
	port   string
}

func newServer(store *store2.Store, router *mux.Router, ctx context.Context, port string) server {
	return server{ctx: ctx, store: store, router: router, client: &payments.Client{}, port: port}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s server) initHandlers() http.Handler {
	//fee calculation
	s.router.HandleFunc("/api/calculate", s.CalculateFeeHandler).Methods("POST")

	//contract create
	s.router.HandleFunc("/api/contract/create", s.HandleCreateContract).Methods("POST")

	//get contracts
	s.router.HandleFunc("/api/contracts", s.HandleGetContracts).Methods("POST")

	//get payout history
	s.router.HandleFunc("/api/payouts", s.HandleGetPayouts).Methods("POST")

	//alerts webhook
	s.router.HandleFunc("/api/alerts", s.HandleAlertWebhook).Methods("POST")

	//register alerts
	s.router.HandleFunc("/api/alerts/register", s.HandleRegisterAlertsEndpoint).Methods("GET")

	//paypal redirect
	s.router.HandleFunc("/api/paypal", s.HandleCreatePaypalOrder).Methods("GET")

	//paypal success
	s.router.HandleFunc("/api/success", s.HandlerSuccess).Methods("GET")

	//paypal cancel
	s.router.HandleFunc("/api/cancel", HandlerCancel).Methods("GET")

	//register webhook endpoint
	s.router.HandleFunc("/api/ipn", s.IPNHandler).Methods("POST")

	return s.router
}
