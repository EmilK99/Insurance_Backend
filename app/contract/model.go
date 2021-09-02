package contract

import "time"

type Contract struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	FlightNumber string    `json:"flight_number"`
	FlightDate   int       `json:"flight_date"`
	Date         time.Time `json:"date"`
	TicketPrice  float32   `json:"ticket_price"`
	Fee          float32   `json:"fee"`
}

func NewContract(userID, flightNumber string, flightDate int, ticketPrice, fee float32) Contract {
	return Contract{UserID: userID,
		FlightNumber: flightNumber,
		FlightDate:   flightDate,
		Date:         time.Now(),
		TicketPrice:  ticketPrice,
		Fee:          fee}
}
