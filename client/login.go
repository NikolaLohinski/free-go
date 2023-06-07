package client

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/nikolalohinski/free-go/types"
)

var (
	ClientLoginSessionTTL = time.Minute * 30 // Fixed by the freebox server, but made into a variable for unit testing
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

	sessionResponse, err := c.getSession(challenge.Challenge)
	if err != nil {
		return nil, fmt.Errorf("failed to get a session: %s", err)
	}
	c.session = &session{
		token:   sessionResponse.SessionToken,
		expires: time.Now().Add(ClientLoginSessionTTL),
	}

	return sessionResponse.Permissions, nil
}

func (c *client) getLoginChallenge() (*loginChallenge, error) {
	response, err := c.genericGet("login")
	if err != nil {
		return nil, fmt.Errorf("failed to GET login endpoint: %s", err)
	}

	result := new(loginChallenge)
	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get login challenge from generic response: %s", err)
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

	response, err := c.genericPost("login/session", sessionsRequest{
		AppID:    *c.appID,
		Password: fmt.Sprintf("%x", hash.Sum(nil)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to POST login/session endpoint: %s", err)
	}

	result := new(sessionResponse)
	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get session response from generic response: %s", err)
	}

	return result, nil
}
