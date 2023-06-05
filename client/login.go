package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/nikolalohinski/free-go/types"
)

const (
	sessionTTL = time.Minute * 30 // Fixed by the server
)

type loginResponse struct {
	LoggedIn     bool   `json:"logged_in"`
	Challenge    string `json:"challenge"`
	PasswordSalt string `json:"password_salt"`
	PasswordSet  bool   `json:"password_set"`
}

type sessionsRequest struct {
	AppID    string `json:"app_id"`
	Password string `json:"password"`
}

type sessionResponse struct {
	SessionToken string            `json:"session_token,omitempty"`
	PasswordSet  bool              `json:"password_set,omitempty"`
	Permissions  types.Permissions `json:"permissions,omitempty"`
	Challenge    string            `json:"challenge"`
	PasswordSalt string            `json:"password_salt"`
}

func (c *client) Login() (permissions types.Permissions, err error) {
	loginHTTPResponse, err := c.httpClient.Get(fmt.Sprintf("%s/login", c.base))
	if err != nil {
		err = fmt.Errorf("failed to perform request: %s", err)
		return
	}

	body, err := ioutil.ReadAll(loginHTTPResponse.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response body: %s", err)
		return
	}

	if loginHTTPResponse.StatusCode >= http.StatusInternalServerError {
		err = fmt.Errorf("failed with status '%d': server returned '%s'", loginHTTPResponse.StatusCode, string(body))
		return
	}
	response := new(genericResponse)
	if err = json.Unmarshal(body, response); err != nil {
		err = fmt.Errorf("failed to unmarshal response body '%s': %s", string(body), err)
		return
	}
	if !response.Success {
		err = fmt.Errorf("failed with error code '%s': %s", response.ErrorCode, response.Message)
		return
	}

	loginResult := new(loginResponse)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  loginResult,
	})
	if err != nil {
		err = fmt.Errorf("failed to instantiate a map structure decoder: %s", err)
		return
	}

	if err = decoder.Decode(response.Result); err != nil {
		err = fmt.Errorf("failed to decode response result to a login result: %s", err)
		return
	}

	hash := hmac.New(sha1.New, []byte(c.privateToken))
	hash.Write([]byte(loginResult.Challenge))
	sessionRequest := sessionsRequest{
		AppID:    c.appID,
		Password: fmt.Sprintf("%x", hash.Sum(nil)),
	}
	sessionRequestBody := new(bytes.Buffer)
	json.NewEncoder(sessionRequestBody).Encode(sessionRequest)

	sessionHTTPResponse, err := c.httpClient.Post(fmt.Sprintf("%s/login/session", c.base), "application/json", sessionRequestBody)
	if err != nil {
		err = fmt.Errorf("failed to perform request: %s", err)
		return
	}
	body, err = ioutil.ReadAll(sessionHTTPResponse.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response body: %s", err)
		return
	}

	if sessionHTTPResponse.StatusCode >= http.StatusInternalServerError {
		err = fmt.Errorf("failed with status '%d': server returned '%s'", sessionHTTPResponse.StatusCode, string(body))
		return
	}

	response = new(genericResponse)
	if err = json.Unmarshal(body, response); err != nil {
		err = fmt.Errorf("failed to unmarshal response body '%s': %s", string(body), err)
		return
	}
	if !response.Success {
		err = fmt.Errorf("failed with error code '%s': %s", response.ErrorCode, response.Message)
		return
	}
	sessionResult := new(sessionResponse)
	decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  sessionResult,
	})
	if err != nil {
		err = fmt.Errorf("failed to instantiate a map structure decoder: %s", err)
		return
	}
	if err = decoder.Decode(response.Result); err != nil {
		err = fmt.Errorf("failed to decode response result to a login result: %s", err)
		return
	}

	c.sessionToken = sessionResult.SessionToken
	c.sessionExpires = time.Now().Add(sessionTTL)

	return sessionResult.Permissions, nil
}
