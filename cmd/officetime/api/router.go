package api

import (
	"net/http"
)

type Router struct {
	Id          int
	Name        string
	Description string
	Address     string
	Login       string
	Password    string
}

type Routers struct {
	Routers []Router `json:"content"`
}

func (otc *officeTimeClient) GetAllRouters() (*Routers, error) {
	request, err := http.NewRequest(http.MethodGet, otc.baseURL+"/router", nil)
	if err != nil {
		return nil, err
	}

	routers := Routers{}
	if err := otc.sendRequest(request, &routers); err != nil {
		return nil, err
	}

	return &routers, nil
}
