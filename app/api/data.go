package api

import (
	"fmt"
	"net/url"
)

func NewFlightInfoURL(aeroApiURL, ident string) string {
	data := url.Values{}
	data.Set("ident", ident)
	data.Set("howMany", "1")

	u, _ := url.ParseRequestURI(aeroApiURL + FlightInfo)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}

func NewFlightInfoExURL(aeroApiURL, ident string, departureTime int) string {
	data := url.Values{}
	data.Set("ident", ident)
	data.Set("departureTime", fmt.Sprint(departureTime))

	u, _ := url.ParseRequestURI(aeroApiURL + FlightInfo)
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

	u, _ := url.ParseRequestURI(aeroApiURL + MetarEx)
	u.RawQuery = data.Encode()

	return fmt.Sprintf("%v", u)
}
