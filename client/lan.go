package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

func (c *client) GetLanConfig(ctx context.Context) (result types.LanConfig, err error) {
	response, err := c.get(ctx, "lan/config/", c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to GET lan/config/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get lan config from generic response: %w", err)
	}

	return result, nil
}

func (c *client) UpdateLanConfig(ctx context.Context, payload types.LanConfig) (result types.LanConfig, err error) {
	response, err := c.put(ctx, "lan/config/", payload, c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to PUT lan/config/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to update lan config from generic response: %w", err)
	}

	return result, nil
}
