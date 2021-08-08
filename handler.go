package main

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

func InitHandlers(pool *pgxpool.Pool) http.Handler {
	r := mux.NewRouter()

	// TODO: implement handlers for flight app

	return r
}
