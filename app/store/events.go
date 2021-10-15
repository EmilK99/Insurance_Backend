package store

type Alert struct {
	LongDesc  string `json:"long_desc"`
	ShortDesc string `json:"short_desc"`
	Summary   string `json:"summary"`
	Eventcode string `json:"eventcode"`
	AlertId   int    `json:"alert_id"`
	Flight    struct {
		FaFlightID            string `json:"faFlightID"`
		Ident                 string `json:"ident"`
		Reg                   string `json:"reg"`
		Aircrafttype          string `json:"aircrafttype"`
		Origin                string `json:"origin"`
		Destination           string `json:"destination"`
		Route                 string `json:"route"`
		FiledEte              string `json:"filed_ete"`
		FiledAltitude         int    `json:"filed_altitude"`
		FiledAirspeedKts      int    `json:"filed_airspeed_kts"`
		FiledTime             int64  `json:"filed_time"`
		FiledBlockoutTime     int64  `json:"filed_blockout_time"`
		EstimatedBlockoutTime int64  `json:"estimated_blockout_time"`
		ActualBlockoutTime    int64  `json:"actual_blockout_time"`
		FiledDeparturetime    int64  `json:"filed_departuretime"`
		Actualdeparturetime   int64  `json:"actualdeparturetime"`
		FiledArrivaltime      int64  `json:"filed_arrivaltime"`
		Estimatedarrivaltime  int64  `json:"estimatedarrivaltime"`
		Actualarrivaltime     int64  `json:"actualarrivaltime"`
		FiledBlockinTime      int64  `json:"filed_blockin_time"`
		EstimatedBlockinTime  int64  `json:"estimated_blockin_time"`
		ActualBlockinTime     int64  `json:"actual_blockin_time"`
	} `json:"flight"`
}
