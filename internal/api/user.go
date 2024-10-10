package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type UserData struct {
	User User `json:"user"`
}

type User struct {
	MacAddress string `json:"macAddress"`
	TelegramId int64  `json:"telegramId"`
}

func (otc *officeTimeClient) GetUserByMacAddress(macAddress string) (*UserData, error) {
	/*cfg := config.New()
	authData := new(LoginData)
	authData.Email = cfg.ApiMasterEmail
	authData.Password = cfg.ApiMasterPassword*/
	params, _ := json.Marshal(authData)
	payload := strings.NewReader(string(params))

	fmt.Println(otc.baseURL + "/user-by-mac-address/" + macAddress)
	request, err := http.NewRequest(http.MethodGet, otc.baseURL+"/user-by-mac-address/"+macAddress, payload)
	if err != nil {
		return nil, err
	}

	user := UserData{}
	if err = otc.sendRequest(request, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
