package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

func (c *client) GetSambaConfiguration(ctx context.Context) (result types.SambaConfiguration, err error) {
	response, err := c.get(ctx, "netshare/samba/", c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to GET netshare/samba/ endpoint: %w", err)
	}

	if response.Result != nil {
		if err = c.fromGenericResponse(response, &result); err != nil {
			return types.SambaConfiguration{}, fmt.Errorf("failed to get Samba configuration from generic response: %w", err)
		}
	}

	return result, nil
}

func (c *client) UpdateSambaConfiguration(ctx context.Context, payload types.SambaConfigurationPayload) (result types.SambaConfiguration, err error) {
	response, err := c.put(ctx, "netshare/samba/", payload, c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to PUT netshare/samba/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to update Samba configuration from generic response: %w", err)
	}

	return result, nil
}

func (c *client) GetAFPConfiguration(ctx context.Context) (result types.AFPConfiguration, err error) {
	response, err := c.get(ctx, "netshare/afp/", c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to GET netshare/afp/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get AFP configuration from generic response: %w", err)
	}

	return result, nil
}

func (c *client) UpdateAFPConfiguration(ctx context.Context, payload types.AFPConfigurationPayload) (result types.AFPConfiguration, err error) {
	response, err := c.put(ctx, "netshare/afp/", payload, c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to PUT netshare/afp/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to update AFP configuration from generic response: %w", err)
	}

	return result, nil
}
