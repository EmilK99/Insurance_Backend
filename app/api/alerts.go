package api

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

func RegisterAlertsEndpoint() error {
	username := viper.GetString("aeroapi_username")
	apiKey := viper.GetString("aeroapi_apikey")
	aeroApiURL := "https://" + username + ":" + apiKey + "@flightxml.flightaware.com/json/FlightXML2/"

	aeroApiURLStr := NewRegisterAlertEndpointURL(aeroApiURL, "https://safe-beyond-32265.herokuapp.com/api/alerts")

	client := &http.Client{}
	re, err := http.NewRequest("POST", aeroApiURLStr, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(re)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func SetAlerts() error {

	//TODO: implement alerts setting
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
