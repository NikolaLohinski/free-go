package client

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

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

	if response.Result == nil {
		return
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

// “application/x-www-form-urlencoded” instead of “application/json”.
func (c *client) AddDownloadTask(ctx context.Context, downloadRequest types.DownloadRequest) (int64, error) {
	form := url.Values{}

	if len(downloadRequest.DownloadURLs) == 1 {
		form.Set("download_url", downloadRequest.DownloadURLs[0])
	} else {
		form.Set("download_url_list", strings.Join(downloadRequest.DownloadURLs, "\n"))
	}

	if downloadRequest.DownloadDirectory != "" {
		form.Set("download_dir", base64.StdEncoding.EncodeToString([]byte(downloadRequest.DownloadDirectory)))
	}

	if downloadRequest.Filename != "" {
		if len(downloadRequest.DownloadURLs) > 1 {
			return 0, errors.New("can not set filename with more than one download URL")
		}

		if downloadRequest.Recursive {
			return 0, errors.New("can not set filename with a recursive download")
		}

		form.Set("filename", downloadRequest.Filename)
	}

	if downloadRequest.Hash != "" {
		if len(downloadRequest.DownloadURLs) > 1 {
			return 0, errors.New("can not set hash with more than one download URL")
		}

		if downloadRequest.Recursive {
			return 0, errors.New("can not set hash with a recursive download")
		}

		form.Set("hash", downloadRequest.Hash)
	}

	if downloadRequest.Recursive {
		form.Set("recursive", "true")
	}

	if downloadRequest.Username != "" {
		form.Set("username", downloadRequest.Username)
	}

	if downloadRequest.Password != "" {
		form.Set("password", downloadRequest.Password)
	}

	if downloadRequest.ArchivePassword != "" {
		form.Set("archive_password", downloadRequest.ArchivePassword)
	}

	if len(downloadRequest.Cookies) > 0 {
		arguments := []string{}
		for name, value := range downloadRequest.Cookies {
			arguments = append(arguments, fmt.Sprintf("%s=%s", name, value))
		}

		sort.Strings(arguments)
		form.Set("cookies", strings.Join(arguments, "; "))
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/downloads/add", c.base), strings.NewReader(form.Encode()))
	if err != nil {
		return 0, fmt.Errorf("failed to forge new request: %w", err)
	}

	response, err := c.do(request, c.withSession(ctx), c.withWWWFormURLEncodedContentType)
	if err != nil {
		return 0, fmt.Errorf("failed to POST downloads/add endpoint: %w", err)
	}

	var responseBody struct {
		ID int64 `json:"id"`
	}

	if err = c.fromGenericResponse(response, &responseBody); err != nil {
		return 0, fmt.Errorf("failed to get an ID from generic response: %w", err)
	}

	return responseBody.ID, nil
}

// DeleteDownloadTask deletes a download task by its identifier.
func (c *client) DeleteDownloadTask(ctx context.Context, identifier int64) error {
	response, err := c.delete(ctx, fmt.Sprintf("downloads/%d", identifier), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeTaskNotFound {
			return ErrTaskNotFound
		}

		return fmt.Errorf("failed to DELETE downloads/%d endpoint: %w", identifier, err)
	}

	return nil
}

// EraseDownloadTask erases a download task and the downloaded files.
func (c *client) EraseDownloadTask(ctx context.Context, identifier int64) error {
	response, err := c.delete(ctx, fmt.Sprintf("downloads/%d/erase", identifier), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeTaskNotFound {
			return ErrTaskNotFound
		}

		return fmt.Errorf("failed to DELETE downloads/%d endpoint: %w", identifier, err)
	}

	return nil
}

// UpdateDownloadTask updates a download task by its identifier.
func (c *client) UpdateDownloadTask(ctx context.Context, identifier int64, downloadRequest types.DownloadTaskUpdate) error {
	resp, err := c.put(ctx, fmt.Sprintf("downloads/%d", identifier), downloadRequest, c.withSession(ctx))
	if err != nil {
		if resp != nil && resp.ErrorCode == codeTaskNotFound {
			return ErrTaskNotFound
		}

		return fmt.Errorf("failed to PUT downloads/%d endpoint: %w", identifier, err)
	}

	return nil
}
