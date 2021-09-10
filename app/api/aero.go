package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func CalculateFeeHandler(w http.ResponseWriter, r *http.Request) {

	var req CalculateFeeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)

		return
	}

	premium, err := Calculate(req.FlightNumber, req.TicketPrice)
	if err != nil {
		log.Errorf("Unable to calculate fee: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	res := CalculateFeeResponse{Fee: premium}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(500)
		return
	}
}

func HandleAlertWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	w.WriteHeader(200)
}

func HandleRegisterAlertsEndpoint(w http.ResponseWriter, r *http.Request) {
	err := RegisterAlertsEndpoint(r.Host)
	if err != nil {
		log.Errorf("Unable to register endpoint: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}
	w.WriteHeader(200)
}
