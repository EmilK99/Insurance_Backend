package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func GetInFlightInfo(flightNumber string) (*InFlightInfoResponse, error) {
	username := viper.GetString("aeroapi_username")
	apiKey := viper.GetString("aeroapi_apikey")
	aeroApiURL := "https://" + username + ":" + apiKey + "@flightxml.flightaware.com/json/FlightXML2/"

	aeroApiURLStr := NewInFlightInfoURL(aeroApiURL, flightNumber)

	client := &http.Client{}
	re, _ := http.NewRequest("POST", aeroApiURLStr, nil)

	resp, _ := client.Do(re)

	flightInfo := new(InFlightInfoResponse)
	err := json.NewDecoder(resp.Body).Decode(&flightInfo)
	if err != nil {
		log.Errorf("Unable to decode json: %v", err)
		return nil, err
	}
	return flightInfo, nil
}

func GetMetarExInfo(airport string) (*MetarExResponse, error) {
	username := viper.GetString("aeroapi_username")
	apiKey := viper.GetString("aeroapi_apikey")
	aeroApiURL := "https://" + username + ":" + apiKey + "@flightxml.flightaware.com/json/FlightXML2/"

	aeroApiURLStr := NewMetarExURL(aeroApiURL, airport)

	re, _ := http.NewRequest("POST", aeroApiURLStr, nil)

	client := &http.Client{}
	resp, _ := client.Do(re)

	metarEx := new(MetarExResponse)

	err := json.NewDecoder(resp.Body).Decode(&metarEx)
	if err != nil {
		log.Errorf("Unable to decode json: %v", err)
		return nil, err
	}

	return metarEx, nil
}

func (f *InFlightInfoResponse) CalculateFee() (float32, error) {

	panic("implement me")
	return 0, nil
}
