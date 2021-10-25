package api

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type AeroAPI struct {
	Username string
	APIKey   string
	URL      string
}

func (a AeroAPI) GetFlightInfoEx(flightNumber string, flightDate int64) (*FlightInfo, error) {
	aeroApiURLStr := a.NewFlightInfoExURL(flightNumber)

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

	flights := flightInfoEx.FlightInfoExResult.Flights

	if len(flights) == 0 {
		return nil, errors.New(fmt.Sprintf("Info about this flight doesn't exist: %s", flightNumber))
	}

	if flights[0].Actualdeparturetime != 0 || time.Now().After(time.Unix(flights[0].FiledDeparturetime, 0)) {
		return nil, errors.New(fmt.Sprintf("Flight already departured: %s", flightNumber))
	}

	flight := new(FlightInfo)
	for i := range flights {
		if flights[i].FiledDeparturetime == flightDate && flights[i].Actualdeparturetime == 0 {
			flight = &flights[i]
			break
		}
	}

	if flight.Ident == "" {
		return nil, errors.New(fmt.Sprintf("Info about this flight doesn't exist: %s", flightNumber))
	}

	return flight, nil
}

func (a AeroAPI) GetFlights(flightNumber string) ([]int64, error) {
	aeroApiURLStr := a.NewFlightInfoExURL(flightNumber)

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

	flights := flightInfoEx.FlightInfoExResult.Flights

	if len(flights) == 0 {
		return nil, errors.New(fmt.Sprintf("Info about this flight doesn't exist: %s", flightNumber))
	}

	if flights[0].Actualdeparturetime != 0 || time.Now().After(time.Unix(flights[0].FiledDeparturetime, 0)) {
		return nil, errors.New(fmt.Sprintf("Flight already departured: %s", flightNumber))
	}

	flight := make([]int64, 0)
	if len(flights) > 1 {
		for i := 0; i < len(flights); i++ {
			if time.Now().Before(time.Unix(flights[i].FiledDeparturetime, 0)) {
				flight = append(flight, flights[i].FiledDeparturetime)
			}
		}
	}

	return flight, nil
}

func (a AeroAPI) GetCancellationRate(f *FlightInfo) (float32, error) {
	aeroApiURLStr := a.NewCancellationRateURL(f.Ident)

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

func (a AeroAPI) GetMetarExInfo(f *FlightInfo) (*MetarExResponse, error) {
	aeroApiURLStr := a.NewMetarExURL(f.Origin)

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

func (a AeroAPI) CalculateFee(f *FlightInfo, ticketPrice, cancelRate float32) (float32, error) {
	var fee float32
	//cancel rate addition
	fee += cancelRate * cancelRate / 2

	//ticket price premium addition
	if ticketPrice > 100 {
		fee += 0.025 * (ticketPrice - 100)
	}

	//wind and snow premium addition
	metarEx, err := a.GetMetarExInfo(f)
	if err != nil {
		return 0, err
	}
	if len(metarEx.MetarExResult.Metar) == 0 {
		return 0, errors.New(fmt.Sprintf("Unable to get weather conditions in: %s", f.Origin))
	}

	windSpeed := metarEx.MetarExResult.Metar[0].WindSpeed

	fee += 0.001 * float32(windSpeed^3)
	if strings.Contains(strings.ToLower(metarEx.MetarExResult.Metar[0].CloudFriendly), "snow") {
		fee += 7.5
	}

	return fee, nil
}

func (a AeroAPI) Calculate(flightNumber string, flightDate int64, ticketPrice float32) (float32, error) {
	flightInfo, err := a.GetFlightInfoEx(flightNumber, flightDate)
	if err != nil {
		log.Errorf("Unable to get flight info: %v", err)
		return 0, err
	}

	//fmt.Println(flightInfo.FlightInfoExResult.Flights[0].FaFlightID)

	cancelRate, err := a.GetCancellationRate(flightInfo)
	if err != nil {
		log.Errorf("Unable to get cancellation rate: %v", err)
		return 0, err
	}

	return a.CalculateFee(flightInfo, ticketPrice, cancelRate)
}
