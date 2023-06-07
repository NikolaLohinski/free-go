package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/nikolalohinski/free-go/types"
)

type genericResponse struct {
	UID       string      `json:"uid,omitempty"`
	Message   string      `json:"msg,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
	Success   bool        `json:"success"`
	Result    interface{} `json:"result"`
}

type HTTPOption = func(*http.Request) error

func (c *client) genericGet(path string, options ...HTTPOption) (*genericResponse, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.base, path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to forge new request: %s", err)
	}
	for _, option := range options {
		if err := option(request); err != nil {
			return nil, fmt.Errorf("failed to apply option to request: %s", err)
		}
	}

	httpResponse, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET %s request: %s", path, err)
	}
	return c.fromHTTPResponse(httpResponse)
}

func (c *client) genericPost(path string, body interface{}, options ...HTTPOption) (*genericResponse, error) {
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(body)

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.base, path), requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to forge new request: %s", err)
	}
	for _, option := range append(options, c.withJSONContentType) {
		if err := option(request); err != nil {
			return nil, fmt.Errorf("failed to apply option to request: %s", err)
		}
	}

	httpResponse, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET %s request: %s", path, err)
	}
	return c.fromHTTPResponse(httpResponse)
}

func (c *client) fromGenericResponse(generic *genericResponse, target interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:    "json",
		Result:     target,
		DecodeHook: types.Float64ToTimeHookFunc(),
	})
	if err != nil {
		return fmt.Errorf("failed to instantiate a map structure decoder: %s", err)
	}

	if err = decoder.Decode(generic.Result); err != nil {
		return fmt.Errorf("failed to decode response result to given target: %s", err)
	}

	return nil
}

func (c *client) fromHTTPResponse(httpResponse *http.Response) (*genericResponse, error) {
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}
	if httpResponse.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("failed with status '%d': server returned '%s'", httpResponse.StatusCode, string(body))
	}
	response := new(genericResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body '%s': %s", string(body), err)
	}
	if !response.Success {
		return nil, fmt.Errorf("failed with error code '%s': %s", response.ErrorCode, response.Message)
	}
	return response, nil
}

func (c *client) withJSONContentType(req *http.Request) error {
	req.Header.Add("Content-Type", "application/json")
	return nil
}

func (c *client) withSession(req *http.Request) error {
	if c.session == nil {
		if _, err := c.Login(); err != nil {
			return fmt.Errorf("failed to login before attempting request: %s", err)
		}
	}
	if time.Now().After(c.session.expires) {
		if _, err := c.Login(); err != nil {
			return fmt.Errorf("failed to login again after session expired: %s", err)
		}
	}
	req.Header.Add("X-Fbx-App-Auth", c.session.token)
	return nil
}
