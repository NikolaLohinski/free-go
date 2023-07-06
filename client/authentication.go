package client

import (
	"crypto/hmac"
	"crypto/sha1" //nolint:gosec
	"fmt"
	"time"

	"github.com/nikolalohinski/free-go/types"
)

type authorizationRequest struct {
	AppID      string `json:"app_id"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	DeviceName string `json:"device_name"`
}

type authorizationResponse struct {
	PrivateToken string `json:"app_token"`
	TrackID      int64  `json:"track_id"`
}

type trackResponse struct {
	Status string `json:"status"`
}

func (c *client) Authorize(request types.AuthorizationRequest) (types.PrivateToken, error) {
	if c.appID == nil {
		return "", ErrAppIDIsNotSet
	}

	authorization, err := c.requestToken(request)
	if err != nil {
		return "", fmt.Errorf("failed to request a private token: %w", err)
	}

	if err := c.waitForTokenApproval(authorization.TrackID); err != nil {
		return "", fmt.Errorf("failed to wait for the token to be approved: %w", err)
	}

	return authorization.PrivateToken, nil
}

func (c *client) requestToken(request types.AuthorizationRequest) (*authorizationResponse, error) {
	response, err := c.post("login/authorize", authorizationRequest{
		AppID:      *c.appID,
		AppName:    request.Name,
		AppVersion: request.Version,
		DeviceName: request.Device,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to POST to login/authorize endpoint: %w", err)
	}

	result := new(authorizationResponse)
	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get authorization response from generic response: %w", err)
	}

	return result, nil
}

func (c *client) waitForTokenApproval(trackID int64) error {
	expiration := time.Now().Add(AuthorizeGrantingTimeout)

	for {
		if time.Now().After(expiration) {
			return fmt.Errorf("reached hard timeout after %s waiting for token approval", AuthorizeGrantingTimeout)
		}

		response, err := c.get(fmt.Sprintf("login/authorize/%d", trackID))
		if err != nil {
			return fmt.Errorf("failed to GET login/authorize/%d endpoint: %w", trackID, err)
		}

		result := new(trackResponse)
		if err = c.fromGenericResponse(response, &result); err != nil {
			return fmt.Errorf("failed to get track response from generic response: %w", err)
		}

		if result.Status == "pending" {
			time.Sleep(AuthorizeRetryDelay)

			continue
		}

		if result.Status != "granted" {
			return fmt.Errorf("received unexpected track status: %s", result.Status)
		}

		return nil
	}
}

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

func (c *client) Logout() error {
	_, err := c.post("login/logout/", nil, c.withSession)
	if err != nil {
		return fmt.Errorf("failed to POST to login/logout/ endpoint: %w", err)
	}

	return nil
}
