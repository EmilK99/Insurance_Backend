package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

type server struct {
	router *mux.Router
	store  *Store
}

func newServer(store *Store, router *mux.Router) server {
	return server{store: store, router: router}
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

	//alerts webhook
	s.router.HandleFunc("/api/alerts", s.HandleAlertWebhook).Methods("POST")

	//register alerts
	s.router.HandleFunc("/api/alerts/register", s.HandleRegisterAlertsEndpoint).Methods("GET")

	//register webhook endpoint
	s.router.HandleFunc("/api/ipn", s.IPNHandler).Methods("GET")

	return s.router
}
