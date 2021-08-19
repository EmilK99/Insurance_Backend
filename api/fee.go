package api

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"strings"
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
		return nil, err
	}
	if flightInfo.InFlightInfoResult.Origin == "" {
		return nil, errors.New(fmt.Sprintf("Empty API response: %s", flightNumber))
	}
	return flightInfo, nil
}

func (f *InFlightInfoResponse) GetMetarExInfo() (*MetarExResponse, error) {
	username := viper.GetString("aeroapi_username")
	apiKey := viper.GetString("aeroapi_apikey")
	aeroApiURL := "https://" + username + ":" + apiKey + "@flightxml.flightaware.com/json/FlightXML2/"

	aeroApiURLStr := NewMetarExURL(aeroApiURL, f.InFlightInfoResult.Origin)

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

func (f *InFlightInfoResponse) CalculateFee(ticketPrice float32) (float32, error) {
	var fee float32
	//ticket price premium addition
	if ticketPrice > 100 {
		fee += 0.025 * (ticketPrice - 100)
	}

	//wind and snow premium addition
	metarEx, err := f.GetMetarExInfo()
	if err != nil {
		return 0, err
	}
	if len(metarEx.MetarExResult.Metar) == 0 {
		return 0, errors.New(fmt.Sprintf("Unable to get weather conditions in: %s", f.InFlightInfoResult.Origin))
	}

	windSpeed := metarEx.MetarExResult.Metar[0].WindSpeed

	fee += 0.001 * float32(windSpeed^3)
	if strings.Contains(strings.ToLower(metarEx.MetarExResult.Metar[0].CloudFriendly), "snow") {
		fee += 7.5
	}

	return fee, nil
}
