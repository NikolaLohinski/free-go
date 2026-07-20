package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

func (c *client) ListProfiles(ctx context.Context) (result []types.Profile, err error) {
	response, err := c.get(ctx, "profile/", c.withSession(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to GET profile/ endpoint: %w", err)
	}

	if response.Result == nil {
		return
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to list profiles from generic response: %w", err)
	}

	return result, nil
}
