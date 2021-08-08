package contract

import "time"

type Contract struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	FligthNumber string    `json:"fligth_number"`
	Date         time.Time `json:"date"`
	TicketPrice  float32   `json:"ticket_price"`
	Fee          float32   `json:"fee"`
	CreateTx     string    `json:"create_tx"`
}
