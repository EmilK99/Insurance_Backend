package contract

import (
	"context"
	"database/sql"
	"errors"
	"flight_app/app/api"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func (c *Contract) CreateContract(pool *pgxpool.Pool) error {

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		return err
	}
	defer conn.Release()

	var check int

	err = conn.QueryRow(context.Background(),
		"SELECT id FROM contracts WHERE user_id = $1, flight_number = $2, flight_date = $3", c.UserID, c.FlightNumber, c.FlightDate).Scan(&check)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("Unable to SELECT: %v\n", err)
		return err
	}
	if check != 0 {
		log.Errorf("Contract already exists")
		return errors.New("Contract already exists")
	}

	flightInfo, err := api.GetFlightInfoEx(c.FlightNumber)
	if err != nil {
		log.Errorf("%s", err)
		return err
	}

	cancelRate, err := flightInfo.GetCancellationRate()
	if err != nil {
		log.Errorf("%s", err)
		return err
	}

	c.Fee, err = flightInfo.CalculateFee(c.TicketPrice, cancelRate)
	if err != nil {
		log.Errorf("%s", err)
		return err
	}

	err = conn.QueryRow(context.Background(),
		"INSERT INTO contract (user_id, flight_number, date, ticket_price, fee) VALUES ($1, $2, $3, $4, $5) RETURNING ID",
		c.UserID, c.FlightNumber, c.Date, c.TicketPrice, c.Fee,
	).Scan(&c.ID)
	if err != nil {
		log.Errorf("Unable to INSERT: %v\n", err)
		return err
	}
	return nil
}

func VerifyPayment(pool *pgxpool.Pool, contractId string) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		return err
	}
	defer conn.Release()

	id, _ := strconv.ParseInt(contractId, 10, 0)

	_, err = conn.Exec(context.Background(),
		"UPDATE contract SET payment=true WHERE id=$1",
		id)
	if err != nil {
		log.Errorf("Unable to INSERT: %v\n", err)
		return err
	}
	return nil
}
