package api

import (
	"net/http"
)

type User struct {
	Id         int
	Name       string
	Email      string
	MacAddress string
	TelegramId int64
}

type Users struct {
	Users []User `json:"users"`
}

func (otc *officeTimeClient) GetAllUsers() (*Users, error) {
	request, err := http.NewRequest(http.MethodGet, otc.baseURL+"/user", nil)
	if err != nil {
		return nil, err
	}

	users := Users{}
	if err := otc.sendRequest(request, &users); err != nil {
		return nil, err
	}

	return &users, nil
}
