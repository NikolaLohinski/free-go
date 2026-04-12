package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

const (
	codeVPNUserNotFound = "noent"
)

// GetOpenVPNServerConfig returns the current OpenVPN server configuration.
func (c *client) GetOpenVPNServerConfig(ctx context.Context) (config types.OpenVPNServerConfig, err error) {
	response, err := c.get(ctx, "vpn/openvpn/", c.withSession(ctx))
	if err != nil {
		return config, fmt.Errorf("failed to GET vpn/openvpn/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &config); err != nil {
		return config, fmt.Errorf("failed to get OpenVPN server config from generic response: %w", err)
	}

	return config, nil
}

// UpdateOpenVPNServerConfig updates the OpenVPN server configuration.
func (c *client) UpdateOpenVPNServerConfig(ctx context.Context, payload types.OpenVPNServerConfig) (config types.OpenVPNServerConfig, err error) {
	response, err := c.put(ctx, "vpn/openvpn/", payload, c.withSession(ctx))
	if err != nil {
		return config, fmt.Errorf("failed to PUT vpn/openvpn/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &config); err != nil {
		return config, fmt.Errorf("failed to get updated OpenVPN server config from generic response: %w", err)
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
func (c *client) GetVPNUserClientConfig(ctx context.Context, login string) (string, error) {
	response, err := c.get(ctx, fmt.Sprintf("vpn/user/%s/config/openvpn", login), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeVPNUserNotFound {
			return "", ErrVPNUserNotFound
		}

		return "", fmt.Errorf("failed to GET vpn/user/%s/config/openvpn endpoint: %w", login, err)
	}

	var config string
	if err = json.Unmarshal(response.Result, &config); err != nil {
		return "", fmt.Errorf("failed to decode VPN client config from response: %w", err)
	}

	return config, nil
}
