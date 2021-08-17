package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"net/http"
	"net/url"
)

func CalculateFeeHandler(w http.ResponseWriter, r *http.Request) {

	var req CalculateFeeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)

		return
	}

	var res = CalculateFeeResponse{Fee: rand.Float32() * 5}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(500)
		return
	}

	username := viper.GetString("aeroapi_username")
	apiKey := viper.GetString("aeroapi_apikey")
	aeroApiURL := "https://" + username + ":" + apiKey + "@flightxml.flightaware.com/json/FlightXML2/"

	data := url.Values{}
	data.Set("ident", "RYR5004")

	u, _ := url.ParseRequestURI(aeroApiURL + InFlightInfo)
	u.RawQuery = data.Encode()
	aeroApiURLStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	re, _ := http.NewRequest("POST", aeroApiURLStr, nil)

	resp, _ := client.Do(re)

	var flightInfo InFlightInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&flightInfo)
	if err != nil {
		log.Errorf("Unable to decode json: %v", err)
		w.WriteHeader(500)
		return
	}

	fmt.Printf("%+v\n", flightInfo)

	//TODO: connect AeroAPI and get flight info
}
