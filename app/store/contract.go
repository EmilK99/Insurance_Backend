package store

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func (s *Store) CreateContract(ctx context.Context, c *Contract) error {

	var (
		check   int
		payment bool
	)

	err := s.Conn.QueryRow(ctx,
		"SELECT id, payment FROM contracts WHERE user_id = $1 AND flight_number = $2 AND flight_date = $3", c.UserID, c.FlightNumber, c.FlightDate).Scan(&check, &payment)
	if err != nil {
		if err != pgx.ErrNoRows {
			return err
		}
	}
	if check != 0 {
		if payment {
			return errors.New("Contract already exists")
		}
		_, err := s.Conn.Exec(ctx,
			"UPDATE contracts SET ticket_price=$2, fee=$3 WHERE id=$1",
			check, c.TicketPrice, c.Fee)
		if err != nil {
			log.Errorf("Unable to UPDATE: %v\n", err)
			return err
		}
		return nil
	}

	err = s.Conn.QueryRow(ctx,
		"INSERT INTO contracts (user_id, flight_number, flight_date, date, ticket_price, fee) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID",
		c.UserID, c.FlightNumber, c.FlightDate, c.Date, c.TicketPrice, c.Fee,
	).Scan(&c.ID)
	if err != nil {
		log.Errorf("Unable to INSERT: %v\n", err)
		return err
	}
	return nil
}

func (s *Store) VerifyPayment(ctx context.Context, contractId, paySystem, customerID string) error {

	id, _ := strconv.ParseInt(contractId, 10, 64)

	_, err := s.Conn.Exec(ctx,
		"UPDATE contracts SET payment=true WHERE id=$1",
		id)
	if err != nil {
		log.Errorf("Unable to UPDATE: %v\n", err)
		return err
	}

	_, err = s.Conn.Exec(ctx,
		"INSERT INTO payments (contract_id, pay_system, customer_id) VALUES ($1, $2, $3)", id, paySystem, customerID)
	if err != nil {
		log.Errorf("Unable to UPDATE: %v\n", err)
		return err
	}

	return nil
}

func (s *Store) GetContract(ctx context.Context, contractID int) (Contract, error) {
	var contract = Contract{ID: contractID}

	err := s.Conn.QueryRow(ctx,
		"SELECT fee, status, flight_date, payment FROM contracts WHERE id = $1", contract.ID).Scan(
		&contract.Fee,
		&contract.Status,
		&contract.FlightDate,
		&contract.Payment)
	if err != nil {
		if err != pgx.ErrNoRows {
			return Contract{}, err
		}
	}
	return contract, nil
}

func (s *Store) GetContractsByUser(ctx context.Context, userID string) ([]*ContractsInfo, error) {

	var contracts []*ContractsInfo

	rows, err := s.Conn.Query(ctx, "SELECT flight_number, status, ticket_price FROM contracts WHERE user_id = $1", userID)
	if err != nil {
		log.Errorf("Unable to SELECT: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		row := new(ContractsInfo)
		err := rows.Scan(&row.FlightNumber, &row.Status, &row.Reward)
		if err != nil {
			log.Errorf("Unable to scan: %v\n", err)
			return nil, err
		}

		if row.Status != "cancelled" {
			row.Reward = 0
		}

		contracts = append(contracts, row)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return contracts, nil
}

func (s *Store) GetPayouts(ctx context.Context, userID string) ([]*PayoutsInfo, error) {

	var payouts []*PayoutsInfo
	rows, err := s.Conn.Query(ctx, "SELECT c.id, pa.customer_id, c.ticket_price FROM contracts c left join payments pa on c.id = pa.contract_id WHERE c.user_id = $1 AND status = $2", userID, "cancelled")
	if err != nil {
		log.Errorf("Unable to SELECT: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		row := new(PayoutsInfo)
		err := rows.Scan(&row.ContractId, &row.UserEmail, &row.TicketPrice)
		if err != nil {
			log.Errorf("Unable to scan: %v\n", err)
			return nil, err
		}

		payouts = append(payouts, row)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return payouts, nil
}
