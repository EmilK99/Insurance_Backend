package api

const (
	InFlightInfo = "InFlightInfo"
)

type FlightInfo struct {
	//TODO: get flight info structure from doc
}

type CalculateFeeRequest struct {
	FlightNumber string `json:"flight_number"`
	TicketPrice  string `json:"ticket_price"`
	Cancellation bool   `json:"cancellation"`
	Delay        bool   `json:"delay"`
}

type CalculateFeeResponse struct {
	Fee float32 `json:"fee"`
}

type InFlightInfoResponse struct {
	InFlightInfoResult struct {
		FaFlightID        string  `json:"faFlightID"`
		Ident             string  `json:"ident"`
		Prefix            string  `json:"prefix"`
		Type              string  `json:"type"`
		Suffix            string  `json:"suffix"`
		Origin            string  `json:"origin"`
		Destination       string  `json:"destination"`
		Timeout           string  `json:"timeout"`
		Timestamp         int     `json:"timestamp"`
		DepartureTime     int     `json:"departureTime"`
		FirstPositionTime int     `json:"firstPositionTime"`
		ArrivalTime       int     `json:"arrivalTime"`
		Longitude         float64 `json:"longitude"`
		Latitude          float64 `json:"latitude"`
		LowLongitude      float64 `json:"lowLongitude"`
		LowLatitude       float64 `json:"lowLatitude"`
		HighLongitude     float64 `json:"highLongitude"`
		HighLatitude      float64 `json:"highLatitude"`
		Groundspeed       int     `json:"groundspeed"`
		Altitude          int     `json:"altitude"`
		Heading           int     `json:"heading"`
		AltitudeStatus    string  `json:"altitudeStatus"`
		UpdateType        string  `json:"updateType"`
		AltitudeChange    string  `json:"altitudeChange"`
		Waypoints         string  `json:"waypoints"`
	} `json:"InFlightInfoResult"`
}
