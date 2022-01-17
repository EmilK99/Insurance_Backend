package flightaware_api

import (
	"fmt"
	"net/url"
	"regexp"
)

func GetSuccessCancelURL(host string, tls bool) (string, string) {
	var url string
	if tls {
		url += "https://"
	} else {
		url += "http://"
	}
	url += host
	return url + "/flightaware_api/success", url + "/flightaware_api/cancel"
}

func (a AeroAPI) NewFlightInfoURL(ident string) string {
	data := url.Values{}
	data.Set("ident", ident)
	data.Set("howMany", "10")

	u, _ := url.ParseRequestURI(a.URL + FlightInformation)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func (a AeroAPI) NewFlightInfoExURL(ident string) string {
	data := url.Values{}
	data.Set("ident", ident)

	u, _ := url.ParseRequestURI(a.URL + FlightInfoEx)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func (a AeroAPI) NewGetFlightIDURL(ident string, departureTime int) string {
	data := url.Values{}
	data.Set("ident", ident)
	data.Set("departureTime", fmt.Sprint(departureTime))

	u, _ := url.ParseRequestURI(a.URL + GetFlightID)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func (a AeroAPI) NewCancellationRateAirlineURL(ident string) string {
	re := regexp.MustCompile("[A-Z]+")

	data := url.Values{}
	data.Set("time_period", "today")
	data.Set("type_matching", "airline")
	data.Set("ident_filter", re.FindAllString(ident, -1)[0])

	u, _ := url.ParseRequestURI(a.URLc + CancellationStat)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func (a AeroAPI) NewCancellationRateAirportURL(ident string) string {


	data := url.Values{}
	data.Set("time_period", "today")
	data.Set("type_matching", "origin")
	data.Set("ident_filter", ident)

	u, _ := url.ParseRequestURI(a.URLc + CancellationStat)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func (a AeroAPI) NewMetarExURL(airport string) string {
	data := url.Values{}
	data.Set("airport", airport)
	data.Add("startTime", "0")
	data.Add("howMany", "1")
	data.Add("offset", "0")

	u, _ := url.ParseRequestURI(a.URL + MetarEx)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func (a AeroAPI) NewRegisterAlertEndpointURL(endpoint string) string {
	data := url.Values{}
	data.Set("address", endpoint)
	data.Add("format_type", "json/post")

	u, _ := url.ParseRequestURI(a.URL + RegisterAlertEndpoint)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func (a AeroAPI) NewSetAlertURL(faFlightId string, contractID int) string {
	data := url.Values{}
	data.Set("alert_id", "0")
	data.Add("ident", faFlightId)
	data.Add("channels", "{16 e_departure e_arrival e_cancelled}")
	data.Add("max_weekly", "1000")

	u, _ := url.ParseRequestURI(a.URL + SetAlert)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func (a AeroAPI) NewDeleteAlertURL(id int) string {
	data := url.Values{}
	data.Set("alert_id", fmt.Sprint(id))

	u, _ := url.ParseRequestURI(a.URL + DeleteAlert)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}



func (a AeroAPI) NewAirportDelaysURL(airport string) string{
	data := url.Values{}

	data.Set("airport", airport)

	u, _ := url.ParseRequestURI(a.URLc + AirportDelays)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)


}

func (a AeroAPI) NewCountAirportOperationsURL(airport string) string{
	data := url.Values{}

	data.Set("airport", airport)

	u, _ := url.ParseRequestURI(a.URL + CountAirportOperations)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func (a AeroAPI) NewAirportInfoURL(airport string) string {
	data := url.Values{}

	data.Set("airportCode", airport)

	u, _ := url.ParseRequestURI(a.URL + AirportInfo)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}