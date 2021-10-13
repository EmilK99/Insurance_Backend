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
		FiledTime             int    `json:"filed_time"`
		FiledBlockoutTime     int    `json:"filed_blockout_time"`
		EstimatedBlockoutTime int    `json:"estimated_blockout_time"`
		ActualBlockoutTime    int    `json:"actual_blockout_time"`
		FiledDeparturetime    int    `json:"filed_departuretime"`
		Actualdeparturetime   int    `json:"actualdeparturetime"`
		FiledArrivaltime      int    `json:"filed_arrivaltime"`
		Estimatedarrivaltime  int    `json:"estimatedarrivaltime"`
		Actualarrivaltime     int    `json:"actualarrivaltime"`
		FiledBlockinTime      int    `json:"filed_blockin_time"`
		EstimatedBlockinTime  int    `json:"estimated_blockin_time"`
		ActualBlockinTime     int    `json:"actual_blockin_time"`
	} `json:"flight"`
}
