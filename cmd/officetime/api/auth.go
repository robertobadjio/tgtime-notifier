package api

import (
	"cloud-time-tracker/cmd/officetime/config"
	"encoding/json"
	"net/http"
	"strings"
)

type Tokens struct {
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenExpires  uint64 `json:"access_token_expires"`
	RefreshTokenExpires uint64 `json:"refresh_token_expires"`
}

type LoginData struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

func (otc *officeTimeClient) Login() (*Tokens, error) {
	cfg := config.New()
	authData := new(LoginData)
	authData.Email = cfg.ApiMasterEmail
	authData.Password = cfg.ApiMasterPassword
	params, _ := json.Marshal(authData)
	payload := strings.NewReader(string(params))

	request, err := http.NewRequest(http.MethodPost, otc.baseURL+"/login", payload)
	if err != nil {
		return nil, err
	}

	tokens := Tokens{}
	if err := otc.sendRequest(request, &tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}

/*func Login() *Tokens {
	url := "https://demo.officetime.tech/api-service/login" // TODO: Убрать
	method := "POST"

	time := new(LoginData)
	time.Email = config.Config.ApiMasterEmail
	time.Password = config.Config.ApiMasterPassword

	params, _ := json.Marshal(time)
	payload := strings.NewReader(string(params))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tokens *Tokens
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		panic(err)
	}

	return tokens
}*/
