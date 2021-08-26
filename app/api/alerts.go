package api

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

func GetRegisterAlertEndpoint() error {
	username := viper.GetString("aeroapi_username")
	apiKey := viper.GetString("aeroapi_apikey")
	aeroApiURL := "https://" + username + ":" + apiKey + "@flightxml.flightaware.com/json/FlightXML2/"

	//TODO: add endpoint listener
	aeroApiURLStr := NewRegisterAlertEndpointURL(aeroApiURL, "http://my_endpoint")

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
