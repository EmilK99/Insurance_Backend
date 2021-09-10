package api

import (
	"github.com/spf13/viper"
	"net/http"
)

func RegisterAlertsEndpoint(host string) error {
	username := viper.GetString("aeroapi_username")
	apiKey := viper.GetString("aeroapi_apikey")
	aeroApiURL := "https://" + username + ":" + apiKey + "@flightxml.flightaware.com/json/FlightXML2/"

	aeroApiURLStr := NewRegisterAlertEndpointURL(aeroApiURL, "https://"+host+"/api/alerts")

	client := &http.Client{}
	re, err := http.NewRequest("POST", aeroApiURLStr, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(re)
	if err != nil {
		return err
	}
	return nil
}

func SetAlerts(faFlightID string, contractID int) error {
	username := viper.GetString("aeroapi_username")
	apiKey := viper.GetString("aeroapi_apikey")
	aeroApiURL := "https://" + username + ":" + apiKey + "@flightxml.flightaware.com/json/FlightXML2/"

	aeroApiURLStr := NewSetAlertURL(aeroApiURL, faFlightID, contractID)

	client := &http.Client{}
	re, err := http.NewRequest("POST", aeroApiURLStr, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(re)
	if err != nil {
		return err
	}

	return nil
}

func GetAlerts() error {
	//TODO: implement alerts getting
	return nil
}

func DeleteAlerts() error {
	//TODO: implement alerts deletting
	return nil
}
