package flightaware_api

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
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

func (a AeroAPI) CalculateCancellationFee(f *FlightInfo, ticketPrice, cancelRate float32) (float32, error) {
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

	fee += float32(0.0001 * (math.Pow(float64(windSpeed), 3)))

	if strings.Contains(strings.ToLower(metarEx.MetarExResult.Metar[0].CloudFriendly), "snow") {
		fee += 7.5
	}

	return fee, nil
}

func (a AeroAPI) CalculateCancellation(flightNumber string, flightDate int64, ticketPrice float32) (float32, error) {
	flightInfo, err := a.GetFlightInfoEx(flightNumber, flightDate)
	if err != nil {
		log.Errorf("Unable to get flight info: %v", err)
		return 0, err
	}

	cancelRate, err := a.GetCancellationRate(flightInfo)
	if err != nil {
		log.Errorf("Unable to get cancellation rate: %v", err)
		return 0, err
	}

	return a.CalculateCancellationFee(flightInfo, ticketPrice, cancelRate)
}
////////////////////////////////////////////////////////////////////////////////////////
func (a AeroAPI) CalculateDelay(flightNumber string, flightDate int64, ticketPrice float32) (float32, error) {
	flightInfo, err := a.GetFlightInfoEx(flightNumber, flightDate)
	if err != nil {
		log.Errorf("Unable to get flight info: %v", err)
		return 0, err
	}

	delayRate, err := a.CalculateDelayRate(flightInfo)
	if err != nil {
		log.Errorf("Unable to get cancellation rate: %v", err)
		return 0, err
	}

	return a.CalculateDelayFee(flightInfo, ticketPrice, delayRate)
}

func (a AeroAPI) CalculateDelayFee(f *FlightInfo, ticketPrice, delayRate float32) (float32, error) {

	fee, err := a.CalculateFeeByRate(ticketPrice, delayRate)
	if err != nil{
		log.Errorf("Can't calculate fee by rate: %v", err)
		return 0, err
	}

	return fee, nil
}

func(a AeroAPI) CalculateFeeByRate(ticketPrice, delayRate float32) (float32, error){

	if delayRate < 0 || delayRate > 100{
		return 0, errors.New("invalid delayRate")
	}

	switch {
	case delayRate > 0 && delayRate <= 10:
		fee := ticketPrice * 0.005
		return fee, nil
	case delayRate > 10 && delayRate <= 20:
		fee := ticketPrice * 0.01
		return fee, nil
	case delayRate > 20 && delayRate <= 30:
		fee := ticketPrice * 0.02
		return fee, nil
	case delayRate > 30 && delayRate <= 40:
		fee := ticketPrice * 0.03
		return fee, nil
	case delayRate > 40 && delayRate <= 50:
		fee := ticketPrice * 0.04
		return fee, nil
	case delayRate > 50 && delayRate <= 60:
		fee := ticketPrice * 0.05
		return fee, nil
	case delayRate > 60 && delayRate <= 70:
		fee := ticketPrice * 0.075
		return fee, nil
	case delayRate > 70 && delayRate <= 80:
		fee := ticketPrice * 0.09 //TODO: miss task description
		return fee, nil
	case delayRate > 80 && delayRate <= 90:
		fee := ticketPrice * 0.1
		return fee, nil
	case delayRate > 90 && delayRate <= 100:
		fee := ticketPrice * 0.15
		return fee, nil

	}

	return 0, errors.New("something gone wrong, internal bug in calculate fee")
}

func (a AeroAPI) CalculateDelayRate(f *FlightInfo) (float32, error) {
	var rate float32
	metarEx, err := a.GetMetarExInfo(f)
	if err != nil{
		log.Errorf("Unable to get is it snowing: %v", err)
		return 0, err
	}
	//TODO : Add Holiday rate logic. Added write in report
	hRate, err := a.HolidayRate(f.FiledDeparturetime)
	if err != nil{
		log.Errorf("Can't calculate holiday rate: %v", err)
		return 0, err
	}
	rate += hRate


	if strings.Contains(strings.ToLower(metarEx.MetarExResult.Metar[0].CloudFriendly), "snow"){
		rate += 75 * 0.05
	}
	if metarEx.MetarExResult.Metar[0].WindSpeed >= 3{
		rate += float32(metarEx.MetarExResult.Metar[0].WindSpeed)
		//TODO: end wind speed rate logic
	}
	//TODO: Add Airline Delay logic

	//TODO : Add Airport Delay logic
	apRate, err := a.AirportDelaysRate(f)
	if err != nil{
		log.Errorf("Can't calculate airport delay rate rate: %v", err)
		return 0, err
	}
	rate += apRate
	//TODO : Add Scheduled Dep Time logic



	return rate, nil

}







func(a AeroAPI) HolidayRate(dateInt int64) (float32, error){

	date := time.Unix(dateInt, 0)

	if date.Month() == 12 && date.Day() >= 15{
		return 92.6 * 0.075, nil
	} else if date.Month() == 1 && date.Day() <=15 {
		return 92.6 * 0.075, nil
	} else if date.Month() == 11 && date.Day() >= 18 &&date.Day() <= 27 {
		return 67.2 * 0.075, nil
	} else if date.Month() == 3 && date.Day() >= 10 &&date.Day() <= 22 {
		return 58.3 * 0.075, nil
	}

	return 0, errors.New("something gone wrong, internal bug in holiday rate")
}


func (a AeroAPI) GetAirportDelays(f *FlightInfo) (*AirportDelaysStruct, error) {
	aeroApiURLStr := a.NewAirportDelaysURL(f.Destination)

	re, _ := http.NewRequest("POST", aeroApiURLStr, nil)

	client := &http.Client{}
	resp, _ := client.Do(re)

	airportDelays := new(AirportDelaysStruct)

	err := json.NewDecoder(resp.Body).Decode(&airportDelays)
	if err != nil {
		log.Errorf("Unable to decode json: %v", err)
		return nil, err
	}

	return airportDelays, nil
}

func (a AeroAPI) AirportDelaysRate(f *FlightInfo) (float32, error){
	//TODO: add airport delays and all flights
	return 0, nil
}
