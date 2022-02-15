package store

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

var aircraftCapacity = map[string]int{
	"A320": 180,
	"B738": 189,
	"A321": 220,
	"A319": 116,
	"A20N": 236,
	"B737": 149,
	"B77W": 408,
	"B789": 290,
	"E75L": 86,
	"E190": 100,
	"C172": 2,
	"A21N": 206,
	"B739": 189,
	"B38M": 178,
	"A333": 335,
	"CRJ9": 90,
	"A359": 366,
	"B763": 269,
	"B788": 242,
	"737":  160,
	"P28A": 3,
	"B744": 524,
	"B77L": 400,
	"B772": 400,
	"A332": 293,
	"CRJ2": 50,
	"CRJ7": 68,
	"AT72": 74,
	"AT43": 50,
	"E145": 50,
	"B752": 201,
	"B748": 410,
	"PC12": 10,
	"SR22": 2,
	"B407": 6,
	"BE20": 9,
	"C208": 13,
	"E55P": 7,
	"DH8B": 39,
	"E170": 78,
	"DH8D": 90,
	"A388": 555,
	"B39M": 220,
	"B735": 132,
	"BCS3": 130,
	"B06":  4,
	"S22T": 3,
	"C56X": 7,
	"A35K": 412,
	"B712": 117,
	"C182": 3,
	"A330": 335,
	"B350": 9,
	"B78X": 330,
	"BE36": 6,
	"CL30": 9,
	"SU95": 98,
	"787":  294,
	"C25B": 8,
	"E45X": 50,
	"A339": 300,
	"BCS1": 108,
	"DA40": 2,
	"A343": 335,
	"B773": 479,
	"C402": 9,
	"C68A": 9,
	"AJ27": 98,
	"C25A": 5,
	"GLEX": 19,
	"GLF4": 19,
	"LJ35": 8,
	"B733": 128,
	"EC35": 8,
	"B753": 243,
	"E135": 37,
	"AS50": 6,
	"B762": 224,
	"B764": 304,
	"C510": 5,
	"CRJX": 100,
	"E295": 100,
	"SF34": 37,
	"767":  304,
	"B734": 146,
	"E35L": 13,
	"A139": 15,
	"EC30": 7,
	"A318": 107,
	"B190": 19,
	"B732": 102,
	"C525": 6,
	"EC45": 9,
	"PA31": 9,
	"A30B": 345,
	"AC50": 6,
	"BE9L": 9,
	"CL60": 19,
	"DH8C": 56,
	"777":  550,
	"A109": 7,
	"BE58": 5,
	"D228": 19,
	"E290": 106,
	"F2TH": 19,
	"F70":  79,
	"P180": 9,
	"747":  524,
	//"74F": 581,
	"B429": 7,
	"B736": 108,
	"B742": 423,
	"C441": 10,
	"C55B": 7,
	"C560": 8,
	"C680": 12,
	"CN1":  2,
	"DA42": 3,
	"DHC6": 20,
	"E120": 30,
	"E50P": 7,
	"FA8X": 19,
	"GL7T": 19,
	"GLF6": 18,
	"LJ60": 9,
	"M20P": 3,
	"P210": 5,
	"PA44": 3,
	"PC24": 10,
	"R44":  3,
	"SW4":  19,
	"340":  440,
	"350":  325,
	"A124": 438,
	"A346": 440,
	"AA5":  1,
	"AN12": 90,
	"AN24": 50,
	"AS65": 12,
	"B462": 100,
	"B463": 112,
	"BE10": 15,
	"BE40": 9,
	"BE55": 3,
	"C25C": 9,
	"C25M": 6,
	"C310": 4,
	"C340": 5,
	"C404": 11,
	"CL35": 10,
	"CRJ1": 50,
	"DA62": 6,
	"E550": 12,
	"FA20": 14,
	"FA7X": 19,
	"G280": 10,
	"GALX": 12,
	"GL5T": 16,
	"GLF5": 19,
	"H500": 5,
	"IL62": 195,
	//"IL76": 10,
	"JS32": 19,
	"L410": 17,
	"LJ25": 8,
	"LJ45": 9,
	"M7":   4,
	"MA60": 60,
	"MD88": 172,
	"P32R": 6,
	"PA38": 1,
	"RV10": 3,
	"SH36": 39,
	"T210": 5,
	"TBM8": 6,
}

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
		aircraftCap = 50 //CONST DefaultAircraftCap
	}

	for k, v := range aircraftCapacity {
		if aircraftType == k {
			aircraftCap = v
		}
	}

	if float64(contractCount/aircraftCap) > 0.2 { //CONST MaxContractsPerFlight
		return false, nil
	}

	return true, nil
}
