package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
)

func CalculateFeeHandler(w http.ResponseWriter, r *http.Request) {

	var req CalculateFeeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil { // bad request
		w.WriteHeader(400)

		return
	}

	res := CalculateFeeResponse{Fee: rand.Float32() * 5}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(500)
		return
	}

	flightInfo, err := GetInFlightInfo("CCA612")
	if err != nil {
		log.Errorf("Unable to get flight info: %v", err)
		w.WriteHeader(500)
		return
	}

	metarEx, err := GetMetarExInfo(flightInfo.InFlightInfoResult.Origin)

	fmt.Printf("%+v\n", metarEx.MetarExResult.Metar[0].WindSpeed)

	//TODO: connect AeroAPI and get flight info
}
