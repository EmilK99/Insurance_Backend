package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (a AeroAPI) RegisterAlertsEndpoint(host string) error {
	aeroApiURLStr := a.NewRegisterAlertEndpointURL("https://" + host + "/api/alerts")

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

func (a AeroAPI) SetAlerts(faFlightID string, contractID int) (int, error) {
	aeroApiURLStr := a.NewSetAlertURL(faFlightID, contractID)

	client := &http.Client{}
	re, err := http.NewRequest("POST", aeroApiURLStr, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(re)
	if err != nil {
		return 0, err
	}

	alertId := new(SetAlertResponse)

	err = json.NewDecoder(resp.Body).Decode(&alertId)
	if err != nil {
		log.Errorf("Unable to decode json: %v", err)
		return 0, err
	}

	return alertId.SetAlertResult, nil
}

func (a AeroAPI) DeleteAlerts(id int) error {
	aeroApiURLStr := a.NewDeleteAlertURL(id)

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

//TODO: check cancellation/departure status
func (a AeroAPI) CheckStatus(flightTd string) {
	flightInfo, err := a.GetFlightInfoEx(flightTd)
	if err != nil {
		log.Error(err)
	}
	log.Info(flightInfo.FlightInfoExResult)
}
