package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"flight_app/app/api"
	"flight_app/app/store"
	"flight_app/payments"
	"fmt"
	"github.com/plutov/paypal/v4"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (s *server) HandleGetContracts(w http.ResponseWriter, r *http.Request) {

	var req store.GetContractsReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(400), "message": err.Error(), "status": "Error"})
		return
	}

	contracts, err := s.store.GetContractsByUser(s.ctx, req.UserID)
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

func (s *server) HandleGetPayouts(w http.ResponseWriter, r *http.Request) {

	var req store.GetContractsReq

	type response struct {
		Contracts   []*store.ContractsInfo `json:"contracts"`
		TotalPayout float32                `json:"total_payout"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(400), "message": err.Error(), "status": "Error"})
		return
	}

	contracts, err := s.store.GetPayouts(s.ctx, req.UserID)
	if err != nil {
		w.WriteHeader(422)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(422), "message": err.Error(), "status": "Error"})
		return
	}
	var res response
	res.Contracts = contracts

	for i := range contracts {
		res.TotalPayout += contracts[i].Reward
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(500)
		return
	}
}

func (s *server) HandleCreateContract(w http.ResponseWriter, r *http.Request) {

	var req store.CreateContractRequest

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

	contr := store.NewContract(req.UserID, req.FlightNumber, int64(flightInfo.FlightInfoExResult.Flights[0].FiledDeparturetime),
		req.TicketPrice, premium)

	err = s.store.CreateContract(s.ctx, &contr)
	if err != nil {
		log.Errorf("Unable to create contract: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	_, err = api.SetAlerts(flightInfo.FlightInfoExResult.Flights[0].FaFlightID, contr.ID)
	if err != nil {
		log.Errorf("Unable to set alert: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	if time.Unix(contr.FlightDate, 0).Before(time.Now()) {
		log.Error("Operation can't be done", errors.New("123"))
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": errors.New("Flight already departured or cancelled").Error(), "status": "Error"})
		return
	}
	returnUrl, cancelURL := api.GetSuccessCancelURL(r.Host, false)

	href, err := s.client.CreateOrder(s.ctx, contr, returnUrl, cancelURL)
	if err != nil {
		log.Errorf("Failed to get contract: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(map[string]string{"url": href})
	if err != nil {
		log.Errorf("Failed to encode: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
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
	case "invoice_paid":
		err := s.store.VerifyPayment(s.ctx, values["custom"][0], "Paypal", values["payer_email"][0])
		if err != nil {
			log.Println("Failed to verify", err)
		}

	case "invoice_cancelled", "invoice_refunded":
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

func (s *server) HandlerSuccess(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error(err)
	}
	token := r.Form.Get("token")

	res, err := s.client.Client.GetOrder(s.ctx, token)
	if err != nil {
		log.Error(err)
	}

	contractID, err := strconv.Atoi(res.PurchaseUnits[0].ReferenceID)
	if err != nil {
		log.Errorf("Unable to parse contractID: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	contract, err := s.store.GetContract(s.ctx, contractID)
	if err != nil {
		log.Errorf("Failed to get contract: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	if contract.Status != "waiting" || time.Unix(contract.FlightDate, 0).Before(time.Now()) {
		_, err = s.store.Conn.Exec(s.ctx,
			"DELETE FROM contracts WHERE id=$1",
			contractID)
		if err != nil {
			log.Errorf("Unable to UPDATE: %v\n", err)
			return
		}

		log.Errorf("Flight info already changed: flight date is %v", contract.FlightDate)
		return
	}

	req, err := s.client.Client.NewRequest(s.ctx, http.MethodPost, "https://api.sandbox.paypal.com/v2/checkout/orders/"+token+"/capture", nil)
	if err != nil {
		log.Error(err)
	}

	resp := paypal.CaptureOrderResponse{}
	err = s.client.Client.SendWithAuth(req, &resp)
	if err != nil {
		log.Error(err)
	}

	err = s.store.VerifyPayment(s.ctx, res.PurchaseUnits[0].ReferenceID, "Paypal", res.Payer.EmailAddress)
	if err != nil {
		log.Error(err)
	}

	err = payments.SuccessTemplate.Execute(w, nil)
	if err != nil {
		log.Error(err)
	}
}

func (s *server) HandlerCancel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error(err)
	}
	token := r.Form.Get("token")

	res, err := s.client.Client.GetOrder(s.ctx, token)
	if err != nil {
		log.Error(err)
	}

	contractID, err := strconv.Atoi(res.PurchaseUnits[0].ReferenceID)
	if err != nil {
		log.Errorf("Unable to parse contractID: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	_, err = s.store.Conn.Exec(s.ctx,
		"DELETE FROM contracts WHERE id=$1",
		contractID)
	if err != nil {
		log.Errorf("Unable to UPDATE: %v\n", err)
		return
	}

	err = payments.CancelTemplate.Execute(w, nil)
	if err != nil {
		log.Error(err)
	}
}
