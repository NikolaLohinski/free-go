package client

import (
	"crypto/hmac"
	"crypto/sha1" //nolint:gosec
	"fmt"
	"time"

	"github.com/nikolalohinski/free-go/types"
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

func (c *client) Login() (permissions types.Permissions, err error) {
	if c.appID == nil {
		return permissions, ErrAppIDIsNotSet
	}

	if c.privateToken == nil {
		return permissions, ErrPrivateTokenIsNotSet
	}

	challenge, err := c.getLoginChallenge()
	if err != nil {
		return permissions, fmt.Errorf("failed to get login challenge: %w", err)
	}

	sessionResponse, err := c.getSession(challenge.Challenge)
	if err != nil {
		return permissions, fmt.Errorf("failed to get a session: %w", err)
	}

	c.session = &session{
		token:   sessionResponse.SessionToken,
		expires: time.Now().Add(LoginSessionTTL),
	}

	return sessionResponse.Permissions, nil
}

func (c *client) getLoginChallenge() (*loginChallenge, error) {
	response, err := c.get("login")
	if err != nil {
		return nil, fmt.Errorf("failed to GET login endpoint: %w", err)
	}

	result := new(loginChallenge)
	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get login challenge from generic response: %w", err)
	}

	return result, nil
}

func (c *client) getSession(challenge string) (*sessionResponse, error) {
	hash := hmac.New(sha1.New, []byte(*c.privateToken))

	hash.Write([]byte(challenge))

	response, err := c.post("login/session", sessionsRequest{
		AppID:    *c.appID,
		Password: fmt.Sprintf("%x", hash.Sum(nil)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to POST login/session endpoint: %w", err)
	}

	result := new(sessionResponse)
	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get session response from generic response: %w", err)
	}

	return result, nil
}
