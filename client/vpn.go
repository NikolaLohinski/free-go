package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

const (
	codeVPNUserNotFound = "noent"

	// The real Freebox OS API serves OpenVPN client configs per VPN server
	// (openvpn_routed or openvpn_bridge), not per user in isolation. The
	// home-automation setup only ever configures the routed server.
	vpnDownloadConfigServerName = "openvpn_routed"
)

// GetVPNServerConfig returns the current configuration of the given VPN server.
func (c *client) GetVPNServerConfig(ctx context.Context, id types.VPNServerID) (config types.VPNServerConfig, err error) {
	endpoint := fmt.Sprintf("vpn/%s/config/", id)

	response, err := c.get(ctx, endpoint, c.withSession(ctx))
	if err != nil {
		return config, fmt.Errorf("failed to GET %s endpoint: %w", endpoint, err)
	}

	if err = c.fromGenericResponse(response, &config); err != nil {
		return config, fmt.Errorf("failed to get VPN server config from generic response: %w", err)
	}

	return config, nil
}

// UpdateVPNServerConfig updates the configuration of the given VPN server.
func (c *client) UpdateVPNServerConfig(ctx context.Context, id types.VPNServerID, payload types.VPNServerConfig) (config types.VPNServerConfig, err error) {
	endpoint := fmt.Sprintf("vpn/%s/config/", id)

	response, err := c.put(ctx, endpoint, payload, c.withSession(ctx))
	if err != nil {
		return config, fmt.Errorf("failed to PUT %s endpoint: %w", endpoint, err)
	}

	if err = c.fromGenericResponse(response, &config); err != nil {
		return config, fmt.Errorf("failed to get updated VPN server config from generic response: %w", err)
	}

	return config, nil
}

// ListVPNUsers returns all configured VPN user accounts.
func (c *client) ListVPNUsers(ctx context.Context) ([]types.VPNUser, error) {
	response, err := c.get(ctx, "vpn/user/", c.withSession(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to GET vpn/user/ endpoint: %w", err)
	}

	result := make([]types.VPNUser, 0)
	if response.Result != nil {
		if err = c.fromGenericResponse(response, &result); err != nil {
			return nil, fmt.Errorf("failed to get VPN users from generic response: %w", err)
		}
	}

	return result, nil
}

// GetVPNUser returns the VPN user account with the given login.
func (c *client) GetVPNUser(ctx context.Context, login string) (user types.VPNUser, err error) {
	response, err := c.get(ctx, fmt.Sprintf("vpn/user/%s", login), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeVPNUserNotFound {
			return user, ErrVPNUserNotFound
		}

		return user, fmt.Errorf("failed to GET vpn/user/%s endpoint: %w", login, err)
	}

	if err = c.fromGenericResponse(response, &user); err != nil {
		return user, fmt.Errorf("failed to get VPN user from generic response: %w", err)
	}

	return user, nil
}

// CreateVPNUser creates a new VPN user account.
func (c *client) CreateVPNUser(ctx context.Context, payload types.VPNUserPayload) (user types.VPNUser, err error) {
	response, err := c.post(ctx, "vpn/user/", payload, c.withSession(ctx))
	if err != nil {
		return user, fmt.Errorf("failed to POST vpn/user/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &user); err != nil {
		return user, fmt.Errorf("failed to get created VPN user from generic response: %w", err)
	}

	return user, nil
}

// UpdateVPNUser updates an existing VPN user account.
func (c *client) UpdateVPNUser(ctx context.Context, login string, payload types.VPNUserPayload) (user types.VPNUser, err error) {
	response, err := c.put(ctx, fmt.Sprintf("vpn/user/%s", login), payload, c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeVPNUserNotFound {
			return user, ErrVPNUserNotFound
		}

		return user, fmt.Errorf("failed to PUT vpn/user/%s endpoint: %w", login, err)
	}

	if err = c.fromGenericResponse(response, &user); err != nil {
		return user, fmt.Errorf("failed to get updated VPN user from generic response: %w", err)
	}

	return user, nil
}

// DeleteVPNUser deletes a VPN user account.
func (c *client) DeleteVPNUser(ctx context.Context, login string) error {
	response, err := c.delete(ctx, fmt.Sprintf("vpn/user/%s", login), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeVPNUserNotFound {
			return ErrVPNUserNotFound
		}

		return fmt.Errorf("failed to DELETE vpn/user/%s endpoint: %w", login, err)
	}

	return nil
}

// GetVPNUserClientConfig returns the OpenVPN client configuration (.ovpn content) for the given user.
// The real Freebox OS API serves this as a raw file download at
// vpn/download_config/{server_name}/{login} (Content-Type: application/x-openvpn-profile),
// not as a JSON-wrapped result under vpn/user/{login}/config/openvpn (which 404s).
func (c *client) GetVPNUserClientConfig(ctx context.Context, login string) (string, error) {
	endpoint := fmt.Sprintf("vpn/download_config/%s/%s", vpnDownloadConfigServerName, login)

	body, err := c.getRaw(ctx, endpoint, c.withSession(ctx))
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Code == codeVPNUserNotFound {
			return "", ErrVPNUserNotFound
		}

		return "", fmt.Errorf("failed to GET %s endpoint: %w", endpoint, err)
	}

	return string(body), nil
}
