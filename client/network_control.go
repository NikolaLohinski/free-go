package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nikolalohinski/free-go/types"
)

const (
	codeNetworkControlNotFound = "noent"
)

func (c *client) ListNetworkControl(ctx context.Context) (result []types.NetworkControlInfo, err error) {
	response, err := c.get(ctx, "network_control/", c.withSession(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to GET network_control/ endpoint: %w", err)
	}

	if response.Result == nil {
		return
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to list network control from generic response: %w", err)
	}

	return result, nil
}

func (c *client) GetNetworkControl(ctx context.Context, identifier int64) (result types.NetworkControlInfo, err error) {
	response, err := c.get(ctx, "network_control/"+strconv.FormatInt(identifier, 10), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeNetworkControlNotFound {
			return result, ErrNetworkControlNotFound
		}

		return result, fmt.Errorf("failed to GET network_control/%d endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get network control from generic response: %w", err)
	}

	return result, nil
}

func (c *client) UpdateNetworkControl(ctx context.Context, payload types.NetworkControlPayload) (result types.NetworkControlInfo, err error) {
	response, err := c.put(ctx, "network_control/"+strconv.FormatInt(payload.ProfileID, 10), payload, c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeNetworkControlNotFound {
			return result, ErrNetworkControlNotFound
		}

		return result, fmt.Errorf("failed to PUT network_control/%d endpoint: %w", payload.ProfileID, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to update network control from generic response: %w", err)
	}

	return result, nil
}
