package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

func (c *client) GetDownloadConfiguration(ctx context.Context) (result types.DownloadConfiguration, err error) {
	response, err := c.get(ctx, "downloads/config/", c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to GET downloads/config/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get download configuration from generic response: %w", err)
	}

	return result, nil
}

func (c *client) UpdateDownloadConfiguration(ctx context.Context, payload types.DownloadConfiguration) (result types.DownloadConfiguration, err error) {
	response, err := c.put(ctx, "downloads/config/", payload, c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to PUT downloads/config/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get download configuration from generic response: %w", err)
	}

	return result, nil
}
