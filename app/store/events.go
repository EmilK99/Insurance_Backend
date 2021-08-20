package store

import (
	"flight_app/app/api"
	log "github.com/sirupsen/logrus"
)

func CheckStatus(flightTd string) {
	flightInfo, err := api.GetInFlightInfo(flightTd)
	if err != nil {
		log.Error(err)
	}
	log.Info(flightInfo.InFlightInfoResult.ArrivalTime)
}
