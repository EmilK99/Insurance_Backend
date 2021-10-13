package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (a AeroAPI) RegisterAlertsEndpoint(host string) error {
	aeroApiURLStr := a.NewRegisterAlertEndpointURL(fmt.Sprintf("https://%s/api/alerts", host))

	client := &http.Client{}
	re, err := http.NewRequest("POST", aeroApiURLStr, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(re)
	if err != nil {
		return err
	}

	log.Println(res)
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
func (a AeroAPI) CheckStatus(flightId string) {
	flightInfo, err := a.GetFlightInfoEx(flightId)
	if err != nil {
		log.Error(err)
	}
	log.Info(flightInfo)
}
