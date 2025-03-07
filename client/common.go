package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	AuthHeader = "X-Fbx-App-Auth"
)

type genericResponse struct {
	UID       string          `json:"uid,omitempty"`
	Message   string          `json:"msg,omitempty"`
	ErrorCode string          `json:"error_code,omitempty"`
	Success   bool            `json:"success"`
	Result    json.RawMessage `json:"result"`
}

type HTTPOption = func(*http.Request) error

func (c *client) get(ctx context.Context, path string, options ...HTTPOption) (response *genericResponse, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.base, path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to forge new request: %w", err)
	}

	return c.do(request, options...)
}

func (c *client) delete(ctx context.Context, path string, options ...HTTPOption) (response *genericResponse, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", c.base, path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to forge new request: %w", err)
	}

	return c.do(request, options...)
}

func (c *client) put(ctx context.Context, path string, body interface{}, options ...HTTPOption) (*genericResponse, error) {
	requestBody := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(requestBody).Encode(body); err != nil {
			return nil, fmt.Errorf("failed to encode body to JSON: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/%s", c.base, path), requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to forge new request: %w", err)
	}

	options = append(options, c.withJSONContentType)

	return c.do(request, options...)
}

func (c *client) post(ctx context.Context, path string, body interface{}, options ...HTTPOption) (*genericResponse, error) {
	requestBody := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(requestBody).Encode(body); err != nil {
			return nil, fmt.Errorf("failed to encode body to JSON: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", c.base, path), requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to forge new request: %w", err)
	}

	options = append(options, c.withJSONContentType)

	return c.do(request, options...)
}

func (c *client) do(request *http.Request, options ...HTTPOption) (response *genericResponse, err error) {
	for _, option := range options {
		if err := option(request); err != nil {
			return nil, fmt.Errorf("failed to apply option to request: %w", err)
		}
	}

	httpResponse, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	defer func() {
		closeError := httpResponse.Body.Close()
		if err == nil {
			err = closeError
		} else if closeError != nil {
			err = fmt.Errorf("%s: %w", closeError.Error(), err)
		}
	}()

	return c.fromHTTPResponse(httpResponse)
}

func (c *client) fromGenericResponse(generic *genericResponse, target interface{}) error {
	if err := json.Unmarshal(generic.Result, target); err != nil {
		return fmt.Errorf("failed to decode response result to given target: %w", err)
	}

	return nil
}

func (c *client) fromHTTPResponse(httpResponse *http.Response) (*genericResponse, error) {
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if httpResponse.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("failed with status '%d': server returned '%s'", httpResponse.StatusCode, string(body))
	}

	response := new(genericResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body '%s': %w", string(body), err)
	}

	if !response.Success {
		return response, &APIError{
			Code:    response.ErrorCode,
			Message: response.Message,
		}
	}

	return response, nil
}

func (c *client) withJSONContentType(req *http.Request) error {
	req.Header.Add("Content-Type", "application/json")

	return nil
}

func (c *client) withWWWFormURLEncodedContentType(req *http.Request) error {
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return nil
}

func (c *client) withSession(ctx context.Context) func(req *http.Request) error {
	return func(req *http.Request) error {
		if c.session == nil {
			if _, err := c.Login(ctx); err != nil {
				return fmt.Errorf("failed to login before attempting request: %w", err)
			}
		}

		if time.Now().After(c.session.expires) {
			if _, err := c.Login(ctx); err != nil {
				return fmt.Errorf("failed to login again after session expired: %w", err)
			}
		}

		req.Header.Add(AuthHeader, c.session.token)

		return nil
	}
}

// APIError represents a structured Freebox API error.
type APIError struct {
	Code    string
	Message string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("failed with error code %q: %s", e.Code, e.Message)
	}

	return fmt.Sprintf("failed with error code %q", e.Code)
}

func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)

	return ok && t.Code == e.Code
}
