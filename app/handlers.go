package app

import (
	"flight_app/app/api"
	"flight_app/payments"
	"net/http"
)

func (s server) initHandlers() http.Handler {
	//fee calculation
	s.router.HandleFunc("/api/calculate",
		func(w http.ResponseWriter, r *http.Request) {
			api.CalculateFeeHandler(w, r)
		}).Methods("POST")

	//contract create
	s.router.HandleFunc("/api/contract/create",
		func(w http.ResponseWriter, r *http.Request) {
			api.HandleCreateContract(s.store.pool, w, r)
		}).Methods("POST")

	//get contracts

	s.router.HandleFunc("/webhook",
		func(w http.ResponseWriter, r *http.Request) {
			payments.HandleStripeWebhook(s.store.pool, w, r)
		}).Methods("POST")

	//TODO: add query endpoint "/contracts"

	return s.router
}
