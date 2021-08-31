package app

import (
	"flight_app/app/api"
	"flight_app/payments"
	"net/http"
)

func (s server) initHandlers() http.Handler {

	s.router.HandleFunc("/api/calculate",
		func(w http.ResponseWriter, r *http.Request) {
			api.CalculateFeeHandler(w, r)
		}).Methods("POST")

	s.router.HandleFunc("/api/contract/payment",
		func(w http.ResponseWriter, r *http.Request) {
			payments.HandleCreatePaymentIntent(s.store.pool, w, r)
		}).Methods("POST")

	s.router.HandleFunc("/webhook",
		func(w http.ResponseWriter, r *http.Request) {
			payments.HandleCreatePaymentIntent(s.store.pool, w, r)
		}).Methods("POST")

	//TODO: add query endpoint "/contracts"

	return s.router
}
