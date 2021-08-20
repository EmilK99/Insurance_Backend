package api

const (
	InFlightInfo = "InFlightInfo"
	MetarEx      = "MetarEx"
)

type FlightInfo struct {
	//TODO: get flight info structure from doc
}

type CalculateFeeRequest struct {
	FlightNumber string  `json:"flight_number"`
	TicketPrice  float32 `json:"ticket_price"`
	Cancellation bool    `json:"cancellation"`
	Delay        bool    `json:"delay"`
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

type MetarExResponse struct {
	MetarExResult struct {
		NextOffset int `json:"next_offset"`
		Metar      []struct {
			Airport       string  `json:"airport"`
			Time          int     `json:"time"`
			CloudFriendly string  `json:"cloud_friendly"`
			CloudAltitude int     `json:"cloud_altitude"`
			CloudType     string  `json:"cloud_type"`
			Conditions    string  `json:"conditions"`
			Pressure      int     `json:"pressure"`
			TempAir       int     `json:"temp_air"`
			TempDewpoint  int     `json:"temp_dewpoint"`
			TempRelhum    int     `json:"temp_relhum"`
			Visibility    float32 `json:"visibility"`
			WindFriendly  string  `json:"wind_friendly"`
			WindDirection int     `json:"wind_direction"`
			WindSpeed     int     `json:"wind_speed"`
			WindSpeedGust int     `json:"wind_speed_gust"`
			RawData       string  `json:"raw_data"`
		} `json:"metar"`
	} `json:"MetarExResult"`
}
