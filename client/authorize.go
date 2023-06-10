package client

import (
	"fmt"
	"time"

	"github.com/nikolalohinski/free-go/types"
)

var (
	ClientAuthorizeGrantingTimeout = time.Minute * 5
	ClientAuthorizeRetryDelay      = time.Second * 5
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
	if c.appID == nil {
		return nil, fmt.Errorf("app ID is not set")
	}

	response, err := c.Post("login/authorize", authorizationRequest{
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
	expiration := time.Now().Add(ClientAuthorizeGrantingTimeout)

	for {
		if time.Now().After(expiration) {
			return fmt.Errorf("reached hard timeout after %s waiting for token approval", ClientAuthorizeGrantingTimeout)
		}

		response, err := c.Get(fmt.Sprintf("login/authorize/%d", trackID))
		if err != nil {
			return fmt.Errorf("failed to GET login/authorize/%d endpoint: %w", trackID, err)
		}

		result := new(trackResponse)
		if err = c.fromGenericResponse(response, &result); err != nil {
			return fmt.Errorf("failed to get track response from generic response: %w", err)
		}

		if result.Status == "pending" {
			time.Sleep(ClientAuthorizeRetryDelay)

			continue
		}

		if result.Status != "granted" {
			return fmt.Errorf("received unexpected track status: %s", result.Status)
		}

		return nil
	}
}
