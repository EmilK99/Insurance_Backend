package contract

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
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
	fmt.Printf("%+v\n", c)

	err = conn.QueryRow(context.Background(),
		"SELECT id FROM contracts WHERE user_id = $1 AND flight_number = $2 AND flight_date = $3", c.UserID, c.FlightNumber, c.FlightDate).Scan(&check)
	if err != nil {
		if err != pgx.ErrNoRows {
			return err
		}
	}
	if check != 0 {
		return errors.New("Contract already exists")
	}

	err = conn.QueryRow(context.Background(),
		"INSERT INTO contracts (user_id, flight_number, flight_date, date, ticket_price, fee) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID",
		c.UserID, c.FlightNumber, c.FlightDate, c.Date, c.TicketPrice, c.Fee,
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
