package flightaware_api

const (
	FlightInformation      = "FlightInfo"
	GetFlightID            = "GetFlightID"
	FlightInfoEx           = "FlightInfoEx"
	MetarEx                = "MetarEx"
	CancellationStat       = "FlightCancellationStatistics"
	RegisterAlertEndpoint  = "RegisterAlertEndpoint"
	SetAlert               = "SetAlert"
	GetAlert               = "GetAlert"
	DeleteAlert            = "DeleteAlert"
	AirportInfo 		   = "AirportInfo"
)

type WS map[int]float64



type CalculateFeeRequest struct {
	FlightNumber string  `json:"flight_number"`
	TicketPrice  float32 `json:"ticket_price"`
	Cancellation bool    `json:"cancellation"`
	Delay        bool    `json:"delay"`
	FlightDate   int64   `json:"flight_date"`
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

type FlightInfoResponse struct {
	FlightInfoResult struct {
		NextOffset int `json:"next_offset"`
		Flights    []struct {
			Ident                string `json:"ident"`
			Aircrafttype         string `json:"aircrafttype"`
			FiledEte             string `json:"filed_ete"`
			FiledTime            int    `json:"filed_time"`
			FiledDeparturetime   int    `json:"filed_departuretime"`
			FiledAirspeedKts     int    `json:"filed_airspeed_kts"`
			FiledAirspeedMach    string `json:"filed_airspeed_mach"`
			FiledAltitude        int    `json:"filed_altitude"`
			Route                string `json:"route"`
			Actualdeparturetime  int    `json:"actualdeparturetime"`
			Estimatedarrivaltime int    `json:"estimatedarrivaltime"`
			Actualarrivaltime    int    `json:"actualarrivaltime"`
			Diverted             string `json:"diverted"`
			Origin               string `json:"origin"`
			Destination          string `json:"destination"`
			OriginName           string `json:"originName"`
			OriginCity           string `json:"originCity"`
			DestinationName      string `json:"destinationName"`
			DestinationCity      string `json:"destinationCity"`
		} `json:"flights"`
	} `json:"FlightInfoResult"`
}

type FlightInfoExResponse struct {
	FlightInfoExResult struct {
		NextOffset int          `json:"next_offset"`
		Flights    []FlightInfo `json:"flights"`
	} `json:"FlightInfoExResult"`
}

type FlightCancellationStatisticsResponse struct {
	FlightCancellationStatisticsResult struct {
		TotalCancellationsWorldwide int    `json:"total_cancellations_worldwide"`
		TotalCancellationsNational  int    `json:"total_cancellations_national"`
		TotalDelaysWorldwide        int    `json:"total_delays_worldwide"`
		TypeMatching                string `json:"type_matching"`
		NumMatching                 int    `json:"num_matching"`
		Matching                    []struct {
			Ident         string `json:"ident"`
			Description   string `json:"description"`
			Cancellations int    `json:"cancellations"`
			Delays        int    `json:"delays"`
			Total         int    `json:"total"`
		} `json:"matching"`
	} `json:"FlightCancellationStatisticsResult"`
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
			Pressure      float64 `json:"pressure"`
			TempAir       int     `json:"temp_air"`
			TempDewpoint  int     `json:"temp_dewpoint"`
			TempRelhum    int     `json:"temp_relhum"`
			Visibility    float64 `json:"visibility"`
			WindFriendly  string  `json:"wind_friendly"`
			WindDirection int     `json:"wind_direction"`
			WindSpeed     int     `json:"wind_speed"`
			WindSpeedGust int     `json:"wind_speed_gust"`
			RawData       string  `json:"raw_data"`
		} `json:"metar"`
	} `json:"MetarExResult"`
}

type SetAlertResponse struct {
	SetAlertResult int `json:"SetAlertResult"`
}

type FlightInfo struct {
	FaFlightID           string `json:"faFlightID"`
	Ident                string `json:"ident"`
	Aircrafttype         string `json:"aircrafttype"`
	FiledEte             string `json:"filed_ete"`
	FiledTime            int64  `json:"filed_time"`
	FiledDeparturetime   int64  `json:"filed_departuretime"`
	FiledAirspeedKts     int    `json:"filed_airspeed_kts"`
	FiledAirspeedMach    string `json:"filed_airspeed_mach"`
	FiledAltitude        int    `json:"filed_altitude"`
	Route                string `json:"route"`
	Actualdeparturetime  int64  `json:"actualdeparturetime"`
	Estimatedarrivaltime int64  `json:"estimatedarrivaltime"`
	Actualarrivaltime    int64  `json:"actualarrivaltime"`
	Diverted             string `json:"diverted"`
	Origin               string `json:"origin"`
	Destination          string `json:"destination"`
	OriginName           string `json:"originName"`
	OriginCity           string `json:"originCity"`
	DestinationName      string `json:"destinationName"`
	DestinationCity      string `json:"destinationCity"`
}


type AirportInfoReq struct {
	AirportCode string `json:"airport_code"`
}

type AirportInfoResp struct {
	AirportInfoResult struct {
		Name      string  `json:"name"`
		Location  string  `json:"location"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		Timezone  string  `json:"timezone"`
	} `json:"AirportInfoResult"`
}


