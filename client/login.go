package client

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
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
	response, err := c.genericGET("login")
	if err != nil {
		return nil, fmt.Errorf("failed to GET login endpoint: %s", err)
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
	if c.privateToken == nil {
		return nil, fmt.Errorf("private token is not set")
	}
	hash := hmac.New(sha1.New, []byte(*c.privateToken))
	hash.Write([]byte(challenge))

	if c.appID == nil {
		return nil, fmt.Errorf("app ID is not set")
	}

	response, err := c.genericPOST("login/session", sessionsRequest{
		AppID:    *c.appID,
		Password: fmt.Sprintf("%x", hash.Sum(nil)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to POST login/session endpoint: %s", err)
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
