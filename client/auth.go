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

type respMsg struct {
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
		return "", fmt.Errorf("error in marshaling user %s", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return "", fmt.Errorf("error in new request %s", err)
	}
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error in response %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error in reading response body %s", err)
	}
	if resp.StatusCode == 400 || resp.StatusCode == 401 {
		var rerr respMsg
		if err := json.Unmarshal(body, &rerr); err != nil {
			return "", fmt.Errorf("error in unmarshaling HTTP error response\nerror: %s\ncode: %d body: %s", err, resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("Unable to process request, error:%s\n", rerr.Message)
	}
	var ts tokenString
	if err := json.Unmarshal(body, &ts); err != nil {
		return "", fmt.Errorf("error in unmarshaling token response %s", err)
	}
	return ts.Token, nil
}

func (ac *AuthClient) SignUp(email string, pass string) (string, error) {
	user := &api.UserJSON{Email: email, Password: pass}
	url := ac.Host + "/auth/signup"

	b, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("error in marshaling user %s", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return "", fmt.Errorf("error in new request %s", err)
	}
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error in response %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error in reading response body %s", err)
	}

	var jsp respMsg
	if resp.StatusCode == 400 || resp.StatusCode == 401 {
		if err := json.Unmarshal(body, &jsp); err != nil {
			return "", fmt.Errorf("error in unmarshaling HTTP error response\nerror: %s\ncode: %d body: %s", err, resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("Unable to process request\nerror:%s", jsp.Message)
	}
	if err := json.Unmarshal(body, &jsp); err != nil {
		return "", fmt.Errorf("error in unmarshaling successful HTTP response\nbody: %s\n error: %s", string(body), err)
	}
	return jsp.Message, nil
}
