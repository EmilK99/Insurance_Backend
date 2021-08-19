package contract

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flight_app/api"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
)

func CreateContract(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	var req Contract

	type response struct {
		ContractID int `json:"contract_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(r, err)
		}
		return
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		w.WriteHeader(500)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(r, err)
		}
		return
	}
	defer conn.Release()

	var check int

	err = conn.QueryRow(context.Background(),
		"SELECT id FROM contracts WHERE user_id = $1, flight_number = $2", req.UserID, req.FlightNumber).Scan(&check)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("Unable to SELECT: %v\n", err)
		w.WriteHeader(500)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(r, err)
		}
		return

	} else if err == sql.ErrNoRows {

	}

	if true {

	} else {
		log.Errorf("Contract already exists")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(errors.New("Contract already exists")); err != nil {
			log.Print(r, err)
		}
		return
	}

	flightInfo, err := api.GetInFlightInfo(req.FlightNumber)
	if err != nil {
		w.WriteHeader(500)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(r, err)
		}
		return
	}

	req.Fee, err = flightInfo.CalculateFee(req.TicketPrice)
	if err != nil {
		w.WriteHeader(500)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(r, err)
		}
		return
	}

	err = conn.QueryRow(context.Background(),
		"INSERT INTO contract (user_id, flight_number, date, ticket_price, fee) VALUES ($1, $2, $3, $4, $5) RETURNING ID",
		req.UserID, req.FlightNumber, req.Date, req.TicketPrice, req.Fee,
	).Scan(&req.ID)
	if err != nil {
		w.WriteHeader(500)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(r, err)
		}
		return
	}

	res := response{ContractID: rand.Intn(100000)}

	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Print(r, err)
		return
	}
}
