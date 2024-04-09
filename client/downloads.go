package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

func (c *client) ListDownloadTasks(ctx context.Context) (result []types.DownloadTask, err error) {
	response, err := c.get(ctx, "downloads/", c.withSession(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to GET downloads/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get download tasks from generic response: %w", err)
	}

	return result, nil
}
