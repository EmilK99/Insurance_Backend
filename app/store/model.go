package store

import "time"

type Contract struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	FlightNumber string    `json:"flight_number"`
	FlightDate   int64     `json:"flight_date"`
	Date         time.Time `json:"date"`
	TicketPrice  float32   `json:"ticket_price"`
	Fee          float32   `json:"fee"`
	Payment      bool      `json:"payment"`
	Status       string    `json:"status"`
}

func NewContract(userID, flightNumber string, flightDate int64, ticketPrice, fee float32) Contract {
	return Contract{UserID: userID,
		FlightNumber: flightNumber,
		FlightDate:   flightDate,
		Date:         time.Now(),
		TicketPrice:  ticketPrice,
		Payment:      false,
		Fee:          fee}
}

type GetContractsReq struct {
	UserID string `json:"user_id"`
}

type CreateContractRequest struct {
	UserID       string  `json:"user_id"`
	FlightNumber string  `json:"flight_number"`
	TicketPrice  float32 `json:"ticket_price"`
	Cancellation bool    `json:"cancellation"`
	Delay        bool    `json:"delay"`
}

type CreateContractResponse struct {
	Fee        float32 `json:"fee"`
	ContractID int     `json:"contract_id"`
	AlertID    int     `json:"alert_id"`
}

type ContractsInfo struct {
	FlightNumber string  `json:"flight_number"`
	Status       string  `json:"status"`
	Reward       float32 `json:"reward"`
}
