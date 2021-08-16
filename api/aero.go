package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
)

func CalculateFeeHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		FlightNumber string `json:"flight_number"`
		TicketPrice  string `json:"ticket_price"`
		Cancellation bool   `json:"cancellation"`
		Delay        bool   `json:"delay"`
	}

	type response struct {
		Fee float32 `json:"fee"`
	}

	var req request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)

		return
	}

	var res = response{Fee: rand.Float32() * 5}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(500)
		return
	}

	//TODO: connect AeroAPI and get flight info
}
