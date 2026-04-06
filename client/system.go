package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

func (c *client) GetSystemInfo(ctx context.Context) (result types.SystemConfig, err error) {
	response, err := c.get(ctx, "system/", c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to GET system/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get system info from generic response: %w", err)
	}

	return result, nil
}
