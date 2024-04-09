package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

const (
	codeTaskNotFound = "task_not_found"
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

func (c *client) GetDownloadTask(ctx context.Context, identifier int64) (result types.DownloadTask, err error) {
	response, err := c.get(ctx, fmt.Sprintf("downloads/%d", identifier), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeTaskNotFound {
			return result, ErrTaskNotFound
		}
		return result, fmt.Errorf("failed to GET downloads/%d endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get a download task from generic response: %w", err)
	}

	return result, nil
}
