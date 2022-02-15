package store

import "time"

type Contract struct {
	ID           int       `json:"id"`
	Type         string    `json:"type"`
	UserID       string    `json:"user_id"`
	FlightNumber string    `json:"flight_number"`
	FlightDate   int64     `json:"flight_date"`
	Date         time.Time `json:"date"`
	TicketPrice  float32   `json:"ticket_price"`
	Fee          float32   `json:"fee"`
	Payment      bool      `json:"payment"`
	Status       string    `json:"status"`
}

func NewContract(userID, contractType, flightNumber string, flightDate int64, ticketPrice, fee float32) Contract {
	return Contract{UserID: userID,
		FlightNumber: flightNumber,
		Type:         contractType,
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
	FlightDate   int64   `json:"flight_date"`
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
	ContractID   int     `json:"contract_id"`
	FlightNumber string  `json:"flight_number"`
	Status       string  `json:"status"`
	Reward       float32 `json:"reward"`
}

type PayoutsInfo struct {
	ContractId   int     `db:"contract_id"`
	UserEmail    string  `db:"customer_id"`
	FlightNumber string  `json:"flight_number"`
	TicketPrice  float32 `db:"amount"`
	PaySystem    string  `db:"pay_system"`
}

type GetPayoutsResponse struct {
	Contracts   []*ContractsInfo `json:"contracts"`
	TotalPayout float32          `json:"total_payout"`
}

const defaultAircraftCap int = 50

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