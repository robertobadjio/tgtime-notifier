package api

import (
	"net/http"
)

func (otc *officeTimeClient) GetAllPeriods() (*Periods, error) {
	request, err := http.NewRequest(http.MethodGet, otc.baseURL+"/period", nil)
	if err != nil {
		return nil, err
	}

	periods := Periods{}
	if err := otc.sendRequest(request, &periods); err != nil {
		return nil, err
	}

	return &periods, nil
}
