package main

import (
	"flight_app/api"
	"flight_app/contract"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

func InitHandlers(pool *pgxpool.Pool) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/calculate",
		func(w http.ResponseWriter, r *http.Request) {
			api.CalculateFeeHandler(w, r)
		}).Methods("POST")

	r.HandleFunc("/contract/create",
		func(w http.ResponseWriter, r *http.Request) {
			contract.CreateContract(pool, w, r)
		}).Methods("POST")

	// TODO: implement handlers for flight app

	return r
}
