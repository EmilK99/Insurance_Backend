package contract

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flight_app/api"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CreateContract(p *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
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

	//TODO: get userId from app

	conn, err := p.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		w.WriteHeader(500)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(r, err)
		}
		return
	}
	defer conn.Release()

	check, err := conn.Query(context.Background(),
		"SELECT id FROM contracts WHERE user_id = $1, flight_number = $2", req.UserID, req.FlightNumber)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Errorf("Unable to SELECT: %v\n", err)
			w.WriteHeader(500)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				log.Print(r, err)
			}
			return
		} else {
			log.Errorf("Contract already exists")
			w.WriteHeader(422)
			if err := json.NewEncoder(w).Encode(errors.New("Contract already exists")); err != nil {
				log.Print(r, err)
			}
			return
		}
	}
	defer check.Close()

	flightInfo, err := api.GetFlightInfo(req.FlightNumber)
	if err != nil {
		w.WriteHeader(500)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(r, err)
		}
	}

	req.Fee, err = flightInfo.CalculateFee()
	if err != nil {
		w.WriteHeader(500)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(r, err)
		}
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
	}

	res := response{ContractID: req.ID}

	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Print(r, err)
	}

}
