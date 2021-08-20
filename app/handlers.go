package app

import (
	"flight_app/app/api"
	"flight_app/app/contract"
	"net/http"
)

func (s server) initHandlers() http.Handler {

	s.router.HandleFunc("/api/calculate",
		func(w http.ResponseWriter, r *http.Request) {
			api.CalculateFeeHandler(w, r)
		}).Methods("POST")

	s.router.HandleFunc("/contract/create",
		func(w http.ResponseWriter, r *http.Request) {
			contract.CreateContract(s.store.pool, w, r)
		}).Methods("POST")

	// TODO: paypal integration

	return s.router
}
