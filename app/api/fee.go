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

func GetFlightInfoEx(flightNumber string) (*FlightInfoExResponse, error) {
	username := viper.GetString("aeroapi_username")
	apiKey := viper.GetString("aeroapi_apikey")
	aeroApiURL := "https://" + username + ":" + apiKey + "@flightxml.flightaware.com/json/FlightXML2/"

	aeroApiURLStr := NewFlightInfoExURL(aeroApiURL, flightNumber)

	client := &http.Client{}
	re, _ := http.NewRequest("POST", aeroApiURLStr, nil)

	resp, err := client.Do(re)
	if err != nil {
		return nil, err
	}

	flightInfoEx := new(FlightInfoExResponse)
	err = json.NewDecoder(resp.Body).Decode(&flightInfoEx)
	if err != nil {
		return nil, err
	}

	if len(flightInfoEx.FlightInfoExResult.Flights) == 0 {
		return nil, errors.New(fmt.Sprintf("Info about this flight doesn't exist: %s", flightNumber))
	}

	if flightInfoEx.FlightInfoExResult.Flights[0].Actualarrivaltime != 0 {
		return nil, errors.New(fmt.Sprintf("Flight already arrived: %s", flightNumber))
	}

	return flightInfoEx, nil
}

func (f *FlightInfoExResponse) GetCancellationRate() (float32, error) {
	username := viper.GetString("aeroapi_username")
	apiKey := viper.GetString("aeroapi_apikey")
	aeroApiURL := "https://" + username + ":" + apiKey + "@flightxml.flightaware.com/json/FlightXML2c/"

	aeroApiURLStr := NewCancellationRateURL(aeroApiURL, f.FlightInfoExResult.Flights[0].Ident)

	client := &http.Client{}
	re, _ := http.NewRequest("POST", aeroApiURLStr, nil)

	resp, err := client.Do(re)
	if err != nil {
		return 0, err
	}

	flightCancelRate := new(FlightCancellationStatisticsResponse)
	err = json.NewDecoder(resp.Body).Decode(&flightCancelRate)
	if err != nil {
		return 0, err
	}

	if len(flightCancelRate.FlightCancellationStatisticsResult.Matching) == 0 {
		return 0.1, nil
	}

	cancelations := flightCancelRate.FlightCancellationStatisticsResult.Matching[0].Cancellations
	total := flightCancelRate.FlightCancellationStatisticsResult.Matching[0].Total

	return 100 * float32(cancelations) / float32(total), nil
}

func (f *FlightInfoExResponse) GetMetarExInfo() (*MetarExResponse, error) {
	username := viper.GetString("aeroapi_username")
	apiKey := viper.GetString("aeroapi_apikey")
	aeroApiURL := "https://" + username + ":" + apiKey + "@flightxml.flightaware.com/json/FlightXML2/"

	aeroApiURLStr := NewMetarExURL(aeroApiURL, f.FlightInfoExResult.Flights[0].Origin)

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

func (f *FlightInfoExResponse) CalculateFee(ticketPrice, cancelRate float32) (float32, error) {
	var fee float32
	//cancel rate addition
	fee += cancelRate * cancelRate / 2

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
		return 0, errors.New(fmt.Sprintf("Unable to get weather conditions in: %s", f.FlightInfoExResult.Flights[0].Origin))
	}

	windSpeed := metarEx.MetarExResult.Metar[0].WindSpeed

	fee += 0.001 * float32(windSpeed^3)
	if strings.Contains(strings.ToLower(metarEx.MetarExResult.Metar[0].CloudFriendly), "snow") {
		fee += 7.5
	}

	return fee, nil
}

func Calculate(flightNumber string, ticketPrice float32) (float32, error) {
	flightInfo, err := GetFlightInfoEx(flightNumber)
	if err != nil {
		log.Errorf("Unable to get flight info: %v", err)
		return 0, err
	}

	//fmt.Println(flightInfo.FlightInfoExResult.Flights[0].FaFlightID)

	cancelRate, err := flightInfo.GetCancellationRate()
	if err != nil {
		log.Errorf("Unable to get cancellation rate: %v", err)
		return 0, err
	}

	return flightInfo.CalculateFee(ticketPrice, cancelRate)
}
