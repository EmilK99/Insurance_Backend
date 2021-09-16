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
