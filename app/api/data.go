package api

import (
	"fmt"
	"net/url"
	"regexp"
)

func NewFlightInfoURL(aeroApiURL, ident string) string {
	data := url.Values{}
	data.Set("ident", ident)
	data.Set("howMany", "1")

	u, _ := url.ParseRequestURI(aeroApiURL + FlightInfo)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func NewFlightInfoExURL(aeroApiURL, ident string) string {
	data := url.Values{}
	data.Set("ident", ident)

	u, _ := url.ParseRequestURI(aeroApiURL + FlightInfoEx)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func NewGetFlightIDURL(aeroApiURL, ident string, departureTime int) string {
	data := url.Values{}
	data.Set("ident", ident)
	data.Set("departureTime", fmt.Sprint(departureTime))

	u, _ := url.ParseRequestURI(aeroApiURL + GetFlightID)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func NewCancellationRateURL(aeroApiURL, ident string) string {
	re := regexp.MustCompile("[A-Z]+")

	data := url.Values{}
	data.Set("time_period", "today")
	data.Set("type_matching", "airline")
	data.Set("ident_filter", re.FindAllString(ident, -1)[0])

	u, _ := url.ParseRequestURI(aeroApiURL + CancellationStat)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func NewMetarExURL(aeroApiURL, airport string) string {
	data := url.Values{}
	data.Set("airport", airport)
	data.Add("startTime", "0")
	data.Add("howMany", "1")
	data.Add("offset", "0")

	u, _ := url.ParseRequestURI(aeroApiURL + MetarEx)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func NewRegisterAlertEndpointURL(aeroApiURL, endpoint string) string {
	data := url.Values{}
	data.Set("address", endpoint)
	data.Add("format_type", "json/post")

	u, _ := url.ParseRequestURI(aeroApiURL + RegisterAlertEndpoint)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func NewSetAlertURL(aeroApiURL, faFlightId string, contractID int) string {
	data := url.Values{}
	data.Set("alert_id", "0")
	data.Add("ident", faFlightId)
	data.Add("channels", "{16 e_departure e_cancelled}")
	data.Add("max_weekly", "1000")

	u, _ := url.ParseRequestURI(aeroApiURL + SetAlert)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}
