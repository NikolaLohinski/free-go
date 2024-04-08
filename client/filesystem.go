package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

const (
	pathNotFoundCode = "path_not_found"
)

func (c *client) GetFileInfo(ctx context.Context, path string) (types.FileInfo, error) {
	base64Path := base64.StdEncoding.EncodeToString([]byte(path))

	response, err := c.get(ctx, fmt.Sprintf("fs/info/%s", base64Path), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == pathNotFoundCode {
			return types.FileInfo{}, ErrPathNotFound
		}
		return types.FileInfo{}, fmt.Errorf("failed to GET fs/info/%s endpoint: %w", base64Path, err)
	}

	result := types.FileInfo{}
	if response.Result != nil {
		if err = c.fromGenericResponse(response, &result); err != nil {
			return types.FileInfo{}, fmt.Errorf("failed to get file info from generic response: %w", err)
		}
	}

	return result, nil
}

func (c *client) RemoveFiles(ctx context.Context, paths []string) (task types.FileSystemTask, err error) {
	files := make([]types.Base64Path, len(paths))
	for i, p := range paths {
		files[i] = types.Base64Path(p)
	}

	response, err := c.post(ctx, "fs/rm/", map[string]interface{}{"files": files}, c.withSession(ctx))
	if err != nil {
		return task, fmt.Errorf("failed to POST to fs/rm/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &task); err != nil {
		return task, fmt.Errorf("failed to get a filesystem task from a generic response: %w", err)
	}

	return task, nil
}
