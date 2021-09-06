package app

import (
	"flight_app/app/api"
	"flight_app/app/contract"
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
			contract.HandleCreateContract(s.store.pool, w, r)
		}).Methods("POST")

	//get contracts
	s.router.HandleFunc("/api/contracts",
		func(w http.ResponseWriter, r *http.Request) {
			contract.HandleGetContracts(s.store.pool, w, r)
		}).Methods("POST")

	return s.router
}
