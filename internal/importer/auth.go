package importer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AuthResponse struct {
	Token string `json:"token"`
}

// Authenticate выполняет авторизацию и возвращает токен
func Authenticate(apiBaseURL, login, password string) (string, error) {
	authURL := fmt.Sprintf("%s/auth/signin?login=%s&password=%s", apiBaseURL, login, password)
	resp, err := http.Get(authURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to authenticate: " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var authResponse AuthResponse
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return "", err
	}

	if authResponse.Token == "" {
		return "", errors.New("no token received")
	}

	return authResponse.Token, nil
}
