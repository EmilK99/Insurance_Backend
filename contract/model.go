package contract

import "time"

type Contract struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	FlightNumber string    `json:"flight_number"`
	Date         time.Time `json:"date"`
	TicketPrice  float32   `json:"ticket_price"`
	Fee          float32   `json:"fee"`
}
