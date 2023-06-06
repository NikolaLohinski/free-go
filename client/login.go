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

type loginChallenge struct {
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

func (c *client) Login() (types.Permissions, error) {
	challenge, err := c.getLoginChallenge()
	if err != nil {
		return nil, fmt.Errorf("failed to get login challenge: %s", err)
	}

	session, err := c.getSession(challenge.Challenge)
	if err != nil {
		return nil, fmt.Errorf("failed to get a session: %s", err)
	}

	c.sessionToken = session.SessionToken
	c.sessionExpires = time.Now().Add(sessionTTL)

	return session.Permissions, nil
}

func (c *client) getLoginChallenge() (*loginChallenge, error) {
	httpResponse, err := c.httpClient.Get(fmt.Sprintf("%s/login", c.base))
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %s", err)
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	if httpResponse.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("failed with status '%d': server returned '%s'", httpResponse.StatusCode, string(body))
	}
	response := new(genericResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body '%s': %s", string(body), err)
	}
	if !response.Success {
		return nil, fmt.Errorf("failed with error code '%s': %s", response.ErrorCode, response.Message)
	}

	result := new(loginChallenge)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  result,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a map structure decoder: %s", err)
	}

	if err = decoder.Decode(response.Result); err != nil {
		return nil, fmt.Errorf("failed to decode response result to a login result: %s", err)
	}

	return result, nil
}

func (c *client) getSession(challenge string) (*sessionResponse, error) {
	hash := hmac.New(sha1.New, []byte(c.privateToken))
	hash.Write([]byte(challenge))
	sessionRequest := sessionsRequest{
		AppID:    c.appID,
		Password: fmt.Sprintf("%x", hash.Sum(nil)),
	}
	sessionRequestBody := new(bytes.Buffer)
	json.NewEncoder(sessionRequestBody).Encode(sessionRequest)

	httpResponse, err := c.httpClient.Post(fmt.Sprintf("%s/login/session", c.base), "application/json", sessionRequestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %s", err)
	}
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	if httpResponse.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("failed with status '%d': server returned '%s'", httpResponse.StatusCode, string(body))
	}

	response := new(genericResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body '%s': %s", string(body), err)
	}
	if !response.Success {
		return nil, fmt.Errorf("failed with error code '%s': %s", response.ErrorCode, response.Message)
	}

	result := new(sessionResponse)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  result,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a map structure decoder: %s", err)
	}
	if err = decoder.Decode(response.Result); err != nil {
		return nil, fmt.Errorf("failed to decode response result to a login result: %s", err)
	}

	return result, nil
}
