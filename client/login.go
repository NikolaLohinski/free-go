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

	"github.com/nikolalohinski/free-go/types"
)

const (
	sessionTTL = time.Minute * 30 // Fixed by the server
)

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
	loginResponse := types.LoginResponse{}
	if err = json.Unmarshal(body, &loginResponse); err != nil {
		err = fmt.Errorf("failed to unmarshal response body '%s': %s", string(body), err)
		return
	}
	if !loginResponse.Success {
		err = fmt.Errorf("failed with error code '%s': %s", loginResponse.ErrorCode, loginResponse.Message)
		return
	}

	hash := hmac.New(sha1.New, []byte(c.apiKey))
	hash.Write([]byte(loginResponse.Result.Challenge))
	sessionRequest := types.SessionsRequest{
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

	sessionResponse := types.SessionResponse{}
	if err = json.Unmarshal(body, &sessionResponse); err != nil {
		err = fmt.Errorf("failed to unmarshal response body '%s': %s", string(body), err)
		return
	}
	if !sessionResponse.Success {
		err = fmt.Errorf("failed with error code '%s': %s", sessionResponse.ErrorCode, sessionResponse.Message)
		return
	}

	c.session_token = sessionResponse.Result.SessionToken
	c.session_expires = time.Now().Add(sessionTTL)

	return sessionResponse.Result.Permissions, nil
}
