package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cybersiddhu/go-micro-auth/api"
)

type AuthClient struct {
	Host string
}

type respErr struct {
	Message string
}

type tokenString struct {
	Token string
}

func (ac *AuthClient) Login(email string, pass string) (string, error) {
	user := &api.UserJSON{Email: email, Password: pass}
	url := ac.Host + "/auth/login"

	b, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 400 || resp.StatusCode == 401 {
		var rerr respErr
		if err := json.Unmarshal(body, rerr); err != nil {
			return "", err
		}
		return "", fmt.Errorf("Unable to process request, error:%s\n", rerr.Message)
	}
	var ts tokenString
	if err := json.Unmarshal(body, ts); err != nil {
		return "", err
	}
	return ts.Token, nil
}
