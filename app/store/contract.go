package store

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
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
			"UPDATE contracts SET ticket_price=$2, fee=$3, type=$4 WHERE id=$1",
			check, c.TicketPrice, c.Fee, c.Type)
		if err != nil {
			log.Errorf("Unable to UPDATE: %v\n", err)
			return err
		}
		c.ID = check
		return nil
	}
	err = s.Conn.QueryRow(ctx,
		"INSERT INTO contracts (user_id, type, flight_number, flight_date, date, ticket_price, fee) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ID",
		c.UserID, c.Type, c.FlightNumber, c.FlightDate, c.Date, c.TicketPrice, c.Fee,
	).Scan(&c.ID)
	if err != nil {
		log.Errorf("Unable to INSERT: %v\n", err)
		return err
	}
	return nil
}

func (s *Store) VerifyPayment(ctx context.Context, contractId int, paySystem, customerID string) error {
	_, err := s.Conn.Exec(ctx,
		"UPDATE contracts SET payment=true, status='waiting' WHERE id=$1",
		contractId)
	if err != nil {
		log.Errorf("Unable to UPDATE: %v\n", err)
		return err
	}
	_, err = s.Conn.Exec(ctx,
		"INSERT INTO payments (contract_id, pay_system, customer_id) VALUES ($1, $2, $3)", contractId, paySystem, customerID)
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

	rows, err := s.Conn.Query(ctx, "SELECT id, flight_number, status, ticket_price, fee FROM contracts WHERE user_id = $1 and status != 'pending' ORDER BY id DESC", userID)
	if err != nil {
		log.Errorf("Unable to SELECT: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		row := new(ContractsInfo)
		var fee float32
		err := rows.Scan(&row.ContractID, &row.FlightNumber, &row.Status, &row.Reward, &fee)
		if err != nil {
			log.Errorf("Unable to scan: %v\n", err)
			return nil, err
		}

		if row.Status != "cancelled" && row.Status != "paid" {
			row.Reward = fee
		}

		contracts = append(contracts, row)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return contracts, nil
}

func (s *Store) GetPayouts(ctx context.Context, userID string) ([]*PayoutsInfo, []string, error) {

	var payouts []*PayoutsInfo
	var accountKeys []string
	rows, err := s.Conn.Query(ctx, "SELECT c.id, pa.customer_id, c.ticket_price, c.flight_number, c.sc_key FROM contracts c "+
		"left join payments pa on c.id = pa.contract_id WHERE c.user_id = $1 AND status = $2 AND payment = $3", userID, "cancelled", true)
	if err != nil {
		log.Errorf("Unable to SELECT: %v\n", err)
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		row := new(PayoutsInfo)
		var contractAccountKey string
		err := rows.Scan(&row.ContractId, &row.UserEmail, &row.TicketPrice, &row.FlightNumber, &contractAccountKey)
		if err != nil {
			log.Errorf("Unable to scan: %v\n", err)
			return nil, nil, err
		}
		row.PaySystem = "Paypal"
		payouts = append(payouts, row)
		accountKeys = append(accountKeys, contractAccountKey)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	return payouts, accountKeys, nil
}

func (s *Store) UpdatePaidPayouts(ctx context.Context, payouts []*PayoutsInfo) error {
	for i := range payouts {
		_, err := s.Conn.Exec(ctx, "UPDATE contracts SET status = 'paid' WHERE id = $1", payouts[i].ContractId)
		if err != nil {
			return err
		}

		_, err = s.Conn.Exec(ctx, "INSERT INTO payouts(contract_id, pay_system, customer_id, amount) VALUES ($1, $2, $3, $4)",
			payouts[i].ContractId, payouts[i].PaySystem, payouts[i].UserEmail, payouts[i].TicketPrice)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) CheckCountContracts(userID string) (bool, error) {

	var contractCount int

	row := s.Conn.QueryRow(context.Background(), "SELECT COUNT(id) FROM contracts WHERE user_id = $1 and status = 'waiting'", userID)
	if err := row.Scan(&contractCount); err != nil {
		return false, err
	}

	if contractCount > 2 {
		return false, nil
	}

	return true, nil
}

func (s *Store) CheckCountPaypal(payerID string) (bool, error) {
	var contractCount int

	row := s.Conn.QueryRow(context.Background(), "SELECT COUNT(id) FROM contracts WHERE payer_id = $1 and status = 'waiting'", payerID)
	if err := row.Scan(&contractCount); err != nil {
		return false, err
	}

	if contractCount > 2 {
		return false, nil
	}

	return true, nil

}

func (s *Store) CheckCountAircraft(aircraftType, flightNumber string, flightDate int64) (bool, error) {
	var contractCount int

	row := s.Conn.QueryRow(context.Background(), "SELECT COUNT(id) FROM contracts WHERE flight_number = $1 and flight_date  = $2 and status = 'waiting'", flightNumber, flightDate)
	if err := row.Scan(&contractCount); err != nil {
		return false, err
	}

	aircraftCap, ok := aircraftCapacity[aircraftType]
	if !ok {
		aircraftCap = defaultAircraftCap
	}

	if float64(contractCount/aircraftCap) > 0.2 {
		return false, nil
	}

	return true, nil
}

func (s *Store) DeleteContractById(ctx context.Context,contractId int) error{
	_, err := s.Conn.Exec(ctx,
		"DELETE FROM contracts WHERE id=$1",
		contractId)
	return err
}

func (s *Store) InsertPayerIdInContract(ctx context.Context,contractId int,payerId string) error{
	_, err := s.Conn.Exec(ctx,
		"INSERT INTO contracts(payer_id) VALUES ($1) WHERE id=$2",
		payerId, contractId)
	return err
}