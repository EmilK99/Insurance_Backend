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

var mapWindSpeed = map[int]float64{
	3: 0.1,
	4: 0.3,
	5:	0.06,
	6:	0.11,
	7:	0.19,
	8:	0.30,
	9:	0.45,
	10:	0.65,
	11:	0.91,
	12:	1.24,
	13:	1.65,
	14:	2.15,
	15:	2.76,
	16:	3.49,
	17:	4.36,
	18:	5.38,
	19:	6.57,
	20:	7.95,
	21:	9.53,
	22:	11.34,
	23:	13.40,
	24:	15.72,
	25:	18.33,
	26:	21.26,
	27:	24.52,
	28:	28.15,
	29:	32.17,
	30:	36.60,
	31:	41.47,
	32:	46.81,
	33:	52.66,
	34:	59.04,
	35:	65.98,
	36:	73.51,
	37:	81.67,
	38:	90.49,
	39:	100.00,
}

type AeroAPI struct {
	Username string
	APIKey   string
	URL      string
	URLc	 string
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

func (a AeroAPI) GetCancellationAirlineRate(f *FlightInfo) (float32, error) {
	aeroApiURLStr := a.NewCancellationRateAirlineURL(f.Ident)

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

	cancelRate, err := a.GetCancellationAirlineRate(flightInfo)
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
		log.Errorf("Can't calculate fee by rate for: %s, Error: %v", f.Ident, err)
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
		fee := ticketPrice * 0.09
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



	if strings.Contains(strings.ToLower(metarEx.MetarExResult.Metar[0].CloudFriendly), "snow"){
		rate += 75 * 0.05
	}

	if metarEx.MetarExResult.Metar[0].WindSpeed >= 3 {
		var wsRate float32
		for k, v := range mapWindSpeed{
			if metarEx.MetarExResult.Metar[0].WindSpeed == k{
				wsRate += float32(v)
				break
			}
		}
		rate += wsRate*0.05
	}


	alRate, err := a.AirlineDelaysRate(f)
	if err != nil{
		log.Errorf("Can't calculate airport delay rate rate: %v", err)
		return 0, err
	}
	rate += alRate

	apRate, err := a.AirportDelaysRate(f)
	if err != nil{
		log.Errorf("Can't calculate airport delay rate rate: %v", err)
		return 0, err
	}
	rate += apRate

	sdtRate, err := a.ScheduleDepTimeAndHolidayRate(f)
	if err != nil{
		log.Errorf("Can't calculate schedual dep rate rate: %v", err)
		return 0, err
	}

	rate += sdtRate


	return rate, nil

}



func (a AeroAPI) AirportDelaysRate(f *FlightInfo) (float32, error){
	apd, err := a.GetCancellationAirportInfo(f)
	if err != nil{
		log.Errorf("Can't get airport delays: %v", err)
		return 0, err
	}
	numDelays := apd.FlightCancellationStatisticsResult.Matching[0].Delays

	numDep := apd.FlightCancellationStatisticsResult.Matching[0].Total

	airportDelayRate := float32(numDelays)/float32(numDep)

	score := airportDelayRate * 3 * 100

	return score * 0.35, nil
}

func (a AeroAPI) AirlineDelaysRate(f *FlightInfo) (float32, error){
	apd, err := a.GetCancellationAirlineInfo(f)
	if err != nil{
		log.Errorf("Can't get airport delays: %v", err)
		return 0, err
	}
	numDelays := apd.FlightCancellationStatisticsResult.Matching[0].Delays

	numDep := apd.FlightCancellationStatisticsResult.Matching[0].Total

	airportDelayRate := float32(numDelays)/float32(numDep)

	score := airportDelayRate * 3 * 100

	return score * 0.4, nil
}

func (a AeroAPI) GetCancellationAirlineInfo(f *FlightInfo) (*FlightCancellationStatisticsResponse, error) {
	aeroApiURLStr := a.NewCancellationRateAirlineURL(f.Ident)

	client := &http.Client{}
	re, _ := http.NewRequest("POST", aeroApiURLStr, nil)

	resp, err := client.Do(re)
	if err != nil {
		return nil, err
	}

	flightCancelRate := new(FlightCancellationStatisticsResponse)
	err = json.NewDecoder(resp.Body).Decode(&flightCancelRate)
	if err != nil {
		return nil, err
	}

	return flightCancelRate, nil


}

func (a AeroAPI) GetCancellationAirportInfo(f *FlightInfo) (*FlightCancellationStatisticsResponse, error) {
	aeroApiURLStr := a.NewCancellationRateAirportURL(f.Origin)

	client := &http.Client{}
	re, _ := http.NewRequest("POST", aeroApiURLStr, nil)

	resp, err := client.Do(re)
	if err != nil {
		return nil, err
	}

	flightCancelRate := new(FlightCancellationStatisticsResponse)
	err = json.NewDecoder(resp.Body).Decode(&flightCancelRate)
	if err != nil {
		return nil, err
	}

	return flightCancelRate, nil
}


func (a AeroAPI) ScheduleDepTimeAndHolidayRate(f *FlightInfo) (float32, error){

	var score float32
	depTime := time.Unix(f.FiledDeparturetime, 0)

	apInfo, err := a.GetAirportInfo(f)
	if err != nil{
		log.Errorf("Can't get count airport operations: %v", err)
		return 0, err
	}

	tzf := strings.TrimPrefix(apInfo.AirportInfoResult.Timezone,":")
	tz, err := time.LoadLocation(tzf)
	if err != nil{
		log.Errorf("Unknown location: %v", err)
		return 0, err
	}
	depTime = depTime.In(tz)

	if f.FiledDeparturetime == 0{
		return 0, errors.New("can't get departure time for rate calculation")
	}

	if depTime.Month() == 12 && depTime.Day() >= 15 || depTime.Month() == 1 && depTime.Day() <=15 {
		score += 92.6
	} else if depTime.Month() == 11 && depTime.Day() >= 18 &&depTime.Day() <= 27 {
		score += 67.2
	} else if depTime.Month() == 3 && depTime.Day() >= 10 &&depTime.Day() <= 22 {
		score += 58.3
	}

	switch {

		case depTime.Hour() >= 7 &&  depTime.Hour() < 9:
			return (score + 10) * 0.075, nil
		case (depTime.Hour() >= 9 &&  depTime.Hour() < 11) ||(depTime.Hour() == 22):
			return (score + 20) * 0.075, nil
		case depTime.Hour() == 11:
			return (score + 30) * 0.075, nil
		case depTime.Hour() == 12 || depTime.Hour() == 21:
			return (score + 40) * 0.075, nil
		case depTime.Hour() == 13:
			return (score + 50) * 0.075, nil
		case depTime.Hour() == 14 || depTime.Hour() == 20:
			return (score + 60) * 0.075, nil
		case depTime.Hour() == 15:
			return (score + 70) * 0.075, nil
		case depTime.Hour() == 16 || depTime.Hour() == 19:
			return (score + 80) * 0.075, nil
		case depTime.Hour() == 17:
			return (score + 90) * 0.075, nil
		case depTime.Hour() == 18:
			return (score + 100) * 0.075, nil
	default:
		return score * 0.075, nil
	}

}

func (a AeroAPI) GetAirportInfo(f *FlightInfo) (*AirportInfoResp, error) {

	aeroApiURLStr := a.NewAirportInfoURL(f.Origin)

	re, _ := http.NewRequest("POST", aeroApiURLStr, nil)

	client := &http.Client{}
	resp, err := client.Do(re)
	if err != nil{
		log.Errorf("Unable to do requst airport info: %v", err)
		return nil, err
	}

	apInfo := new(AirportInfoResp)

	err = json.NewDecoder(resp.Body).Decode(&apInfo)
	if err != nil {
		log.Errorf("Unable to decode json: %v", err)
		return nil, err
	}

	return apInfo, nil
}