package app

import (
	"encoding/json"
	"errors"
	flightaware_api2 "flight_app/app/api/flightaware_api"
	"flight_app/app/store"
	"flight_app/payments"
	"github.com/gogo/protobuf/sortkeys"
	"github.com/plutov/paypal/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (s *server) HandleGetFlights(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FlightNumber string `json:"flight_number"`
	}

	type response struct {
		FlightNumber string  `json:"flight_number"`
		Count        int     `json:"count"`
		Flights      []int64 `json:"flights"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(response{Flights: make([]int64, 0)})
		return
	}

	flights, err := s.aeroApi.GetFlights(req.FlightNumber)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(response{Flights: make([]int64, 0)})
		return
	}

	var res response

	sortkeys.Int64s(flights)

	res.Flights = flights
	res.FlightNumber = req.FlightNumber
	res.Count = len(res.Flights)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(response{Flights: make([]int64, 0)})
		return
	}
}

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

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(400), "message": err.Error(), "status": "Error"})
		return
	}

	payouts, _, err := s.store.GetPayouts(s.ctx, req.UserID)
	if err != nil {
		w.WriteHeader(422)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(422), "message": err.Error(), "status": "Error"})
		return
	}
	var res store.GetPayoutsResponse

	for i := range payouts {
		ctr := store.ContractsInfo{ContractID: payouts[i].ContractId,
			FlightNumber: payouts[i].FlightNumber,
			Status:       "cancelled",
			Reward:       payouts[i].TicketPrice,
		}
		res.Contracts = append(res.Contracts, &ctr)
		res.TotalPayout += payouts[i].TicketPrice
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

	var (
		req          store.CreateContractRequest
		premium      float32
		contractType string
	)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(400), "message": err.Error(), "status": "Error"})
		return
	}
	//TODO
	if req.TicketPrice > 1000 {
		log.Errorf("The ticket is too expensive, suspicion of fraud")
		w.WriteHeader(412)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(412), "message": "Invalid ticket price", "status": "Suspicion of fraud"})
		return
	}
	//TODO
	checkContr, err := s.store.CheckCountContracts(req.UserID)
	if err != nil {
		log.Errorf("Unable to count contracts")
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error()})
		return
	}
	if !checkContr {
		log.Errorf("Too many opened contracts, suspicion of fraud")
		w.WriteHeader(412)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(412), "message": "Too many opened contracts", "status": "Suspicion of fraud"})
		return
	}

	flightInfo, err := s.aeroApi.GetFlightInfoEx(req.FlightNumber, req.FlightDate)
	if err != nil {
		log.Errorf("Unable to get flight info: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}
	//TODO
	if flightInfo.FiledDeparturetime-3600 > time.Now().Unix() {
		log.Errorf("Late attempt to create contract, suspicion of fraud")
		w.WriteHeader(412)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(412), "message": "Late attempt to create contract", "status": "Suspicion of fraud"})
		return
	}
	//TODO

	checkCap, err := s.store.CheckCountAircraft(flightInfo.Aircrafttype, req.FlightNumber, req.FlightDate)
	if err != nil {
		log.Errorf("Unable to check apircraft capacity: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	if !checkCap {
		log.Errorf("Too many opened contracts on this flight, suspicion of fraud")
		w.WriteHeader(412)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(412), "message": "Too many opened contracts on this flight", "status": "Suspicion of fraud"})
		return
	}

	if req.Cancellation {
		premium, err = s.aeroApi.CalculateCancellation(req.FlightNumber, req.FlightDate, req.TicketPrice)
		if err != nil {
			log.Errorf("Unable to calculate fee: %v", err)
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
			return
		}
		contractType = "cancel"
	} else if req.Delay {
		premium, err = s.aeroApi.CalculateDelay(req.FlightNumber, req.FlightDate, req.TicketPrice)
		if err != nil {
			log.Errorf("Unable to calculate fee: %v", err)
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
			return
		}
		contractType = "delay"
	} else {
		log.Errorf("Unable to calculate fee: %v", err)
		w.WriteHeader(422)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(422), "message": "Choose cancellation or delay", "status": "Error"})
		return
	}
	contr := store.NewContract(req.UserID, contractType, req.FlightNumber, flightInfo.FiledDeparturetime,
		req.TicketPrice, premium)

	if time.Unix(contr.FlightDate, 0).Before(time.Now()) {
		log.Error("Operation can't be done", errors.New("123"))
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": errors.New("Flight already departured or cancelled").Error(), "status": "Error"})
		return
	}

	err = s.store.CreateContract(s.ctx, &contr)
	if err != nil {
		log.Errorf("Unable to create contract: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	_, err = s.aeroApi.SetAlerts(flightInfo.FaFlightID, contr.ID)
	if err != nil {
		log.Errorf("Unable to set alert: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	returnUrl, cancelURL := flightaware_api2.GetSuccessCancelURL(r.Host, false)

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

func (s *server) CalculateFeeHandler(w http.ResponseWriter, r *http.Request) {

	var req flightaware_api2.CalculateFeeRequest
	var premium float32

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)
		return
	}

	if req.FlightDate == 0 {
		w.WriteHeader(400)
		return
	}

	if req.TicketPrice > 1000 {
		log.Errorf("The ticket is too expensive, suspicion of fraud")
		w.WriteHeader(412)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(412), "message": "Invalid ticket price", "status": "Suspicion of fraud"})
		return
	}

	if req.Cancellation {
		premium, err = s.aeroApi.CalculateCancellation(req.FlightNumber, req.FlightDate, req.TicketPrice)
		if err != nil {
			log.Errorf("Unable to calculate fee: %v", err)
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
			return
		}
	} else if req.Delay {
		premium, err = s.aeroApi.CalculateDelay(req.FlightNumber, req.FlightDate, req.TicketPrice)
		if err != nil {
			log.Errorf("Unable to calculate fee: %v", err)
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
			return
		}
	} else {
		log.Errorf("Unable to calculate fee: %v", err)
		w.WriteHeader(422)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(422), "message": "Choose cancellation or delay", "status": "Error"})
		return
	}

	res := flightaware_api2.CalculateFeeResponse{Fee: premium}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(500)
		return
	}
}

func (s *server) HandleAlertWebhook(w http.ResponseWriter, r *http.Request) {

	var alert store.Alert

	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil { // bad request
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(400), "message": err.Error(), "status": "Error"})
		return
	}
	if alert.Eventcode == "cancelled" {
		err = s.store.UpdateContractsByAlert(s.ctx, alert.Flight.Ident, alert.Eventcode, alert.Flight.FiledDeparturetime)
		if err != nil {
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
			return
		}
	}

	err = s.aeroApi.DeleteAlerts(alert.AlertId)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	log.Println(alert.Flight.Ident, alert.Flight.FiledDeparturetime, alert.Eventcode)
	w.WriteHeader(200)
}

func (s *server) HandleRegisterAlertsEndpoint(w http.ResponseWriter, r *http.Request) {
	host, _ := url.QueryUnescape(r.Host)
	err := s.aeroApi.RegisterAlertsEndpoint(host)
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
	if contract.Status != "pending" || time.Unix(contract.FlightDate, 0).Before(time.Now()) {
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

	_, err = s.store.Conn.Exec(s.ctx,
		"INSERT INTO contracts(payer_id) VALUES ($1) WHERE id=$2",
		res.Payer.PayerID, contractID)

	check, err := s.store.CheckCountPaypal(res.Payer.PayerID)
	if err != nil {
		log.Errorf("Failed to get count contracts: %v", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}
	if !check {
		_, err = s.store.Conn.Exec(s.ctx,
			"DELETE FROM contracts WHERE id=$1",
			contractID)
		if err != nil {
			log.Errorf("Unable to UPDATE: %v\n", err)
			return
		}

		log.Errorf("Too many opened contracts for :%v", res.Payer.PayerID)
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

	err = s.store.VerifyPayment(s.ctx, contractID, "Paypal", res.Payer.EmailAddress)
	if err != nil {
		log.Errorf("Unable to verify: %v", err)
		return
	}

	account, err := s.solClient.CreateInsuranceContract(s.ctx, contractID)
	if err != nil {
		log.Errorf("Unable to create contract: %v", err)
		return
	}
	err = s.store.SaveContractAccount(s.ctx, contractID, account)
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

func (s *server) HandleWithdrawPremium(w http.ResponseWriter, r *http.Request) {

	var req struct {
		UserID    string `json:"user_id"`
		Contracts []int  `json:"contracts"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(400), "message": err.Error(), "status": "Error"})
		return
	}

	payouts, keys, err := s.store.GetPayouts(s.ctx, req.UserID)
	if err != nil {
		w.WriteHeader(422)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(422), "message": err.Error(), "status": "Error"})
		return
	}

	var newPayouts []*store.PayoutsInfo

	for _, v := range req.Contracts {
		for j := range payouts {
			if payouts[j].ContractId == v {
				newPayouts = append(newPayouts, payouts[j])
			}
		}
	}

	err = s.client.CreatePayout(s.ctx, newPayouts)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	err = s.store.UpdatePaidPayouts(s.ctx, newPayouts)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(500), "message": err.Error(), "status": "Error"})
		return
	}

	for i := range keys {
		err = s.solClient.CloseInsuranceContract(s.ctx, keys[i])
		if err != nil {
			log.Errorf("Unable to close contract: %v", err)
			return
		}
	}

	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(map[string]string{"code": strconv.Itoa(200), "message": "Withdraw requested", "status": "successful"})
	if err != nil {
		log.Error(err)
	}
}
