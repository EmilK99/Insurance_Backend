package app

import (
	"bytes"
	"encoding/json"
	"flight_app/app/api"
	"flight_app/app/contract"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (s *server) HandleGetContracts(w http.ResponseWriter, r *http.Request) {

	var req contract.GetContractsReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(400), "message": err.Error(), "status": "Error"})
		return
	}

	contracts, err := contract.GetContracts(s.store.pool, req.UserID)
	if err != nil {
		w.WriteHeader(422)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(422), "message": err.Error(), "status": "Error"})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(contracts)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(500)
		return
	}
}

func (s *server) HandleCreateContract(w http.ResponseWriter, r *http.Request) {

	var req contract.CreateContractRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(400), "message": err.Error(), "status": "Error"})
		return
	}

	flightInfo, err := api.GetFlightInfoEx(req.FlightNumber)
	if err != nil {
		log.Errorf("Unable to get flight info: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	premium, err := api.Calculate(req.FlightNumber, req.TicketPrice)
	if err != nil {
		log.Errorf("Unable to calculate fee: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	contr := contract.NewContract(req.UserID, req.FlightNumber, int(flightInfo.FlightInfoExResult.Flights[0].FiledDeparturetime),
		req.TicketPrice, premium)

	err = contr.CreateContract(s.store.pool)
	if err != nil {
		log.Errorf("Unable to create contract: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	alertID, err := api.SetAlerts(flightInfo.FlightInfoExResult.Flights[0].FaFlightID, contr.ID)
	if err != nil {
		log.Errorf("Unable to set alert: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	res := contract.CreateContractResponse{Fee: premium, ContractID: contr.ID, AlertID: alertID}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(500)
		return
	}
}

func (s *server) IPNHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	// Switch for production and live
	isProduction := true

	urlSimulator := "https://www.sandbox.paypal.com/cgi-bin/webscr"
	urlLive := "https://www.paypal.com/cgi-bin/webscr"
	paypalURL := urlSimulator

	if isProduction {
		paypalURL = urlLive
	}

	// Verify that the POST HTTP Request method was used.
	// A more sophisticated router would have handled this before calling this handler.
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("No route for %v", r.Method), http.StatusNotFound)
		return
	}

	log.Printf("Write Status 200")
	w.WriteHeader(http.StatusOK)

	// Get Content-Type of request to be parroted back to paypal
	contentType := r.Header.Get("Content-Type")
	// Read the raw POST body
	body, _ := ioutil.ReadAll(r.Body)
	// Prepend POST body with required field
	body = append([]byte("cmd=_notify-validate&"), body...)
	// Make POST request to paypal
	resp, _ := http.Post(paypalURL, contentType, bytes.NewBuffer(body))

	verifyStatus, _ := ioutil.ReadAll(resp.Body)

	if string(verifyStatus) != "VERIFIED" {
		log.Printf("Response: %v", string(verifyStatus))
		log.Println("This indicates that an attempt was made to spoof this interface, or we have a bug.")
		return
	}
	// We can now assume that the POSTed information in `body` is VERIFIED to be from Paypal.
	log.Printf("Response: %v", string(verifyStatus))

	values, err := url.ParseQuery(string(body))
	if err != nil {
		log.Println("Error in parsing url", err)
	}

	switch values["txn_type"][0] {
	case "payment_finished":
		//TODO handle success webhook
	case "payment_cancelled":
		//TODO handle fail webhook
	}

}

func (s *server) CalculateFeeHandler(w http.ResponseWriter, r *http.Request) {

	var req api.CalculateFeeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)

		return
	}

	premium, err := api.Calculate(req.FlightNumber, req.TicketPrice)
	if err != nil {
		log.Errorf("Unable to calculate fee: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	res := api.CalculateFeeResponse{Fee: premium}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(500)
		return
	}
}

func (s *server) HandleAlertWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	w.WriteHeader(200)
}

func (s *server) HandleRegisterAlertsEndpoint(w http.ResponseWriter, r *http.Request) {
	err := api.RegisterAlertsEndpoint(r.Host)
	if err != nil {
		log.Errorf("Unable to register endpoint: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}
	w.WriteHeader(200)
}
