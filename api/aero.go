package api

import (
	"encoding/json"
	"net/http"
)

func CalculateFeeHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		FlightNumber string `json:"flight_number"`
		TicketPrice  string `json:"ticket_price"`
		Cancellation bool   `json:"cancellation"`
		Delay        bool   `json:"delay"`
	}

	var req request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)

		return
	}

	//TODO: connect AeroAPI and get flight info
}
