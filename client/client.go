package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/nikolalohinski/free-go/types"
)

type Client interface {
	WithAppID(string) Client
	WithPrivateToken(types.PrivateToken) Client
	WithHttpClient(*http.Client) Client
	APIVersion() (types.APIVersion, error)
	Authorize(types.AuthorizationRequest) (types.PrivateToken, error)
	Login() (types.Permissions, error)
}

func New(endpoint, version string) (Client, error) {
	match, err := regexp.MatchString("^https?://.*", endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to match endpoint string against regex: %s", err)
	}
	if !match {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}
	return &client{
		httpClient: http.DefaultClient,
		base:       fmt.Sprintf("%s/api/%s", endpoint, version),
	}, nil
}

type client struct {
	httpClient   *http.Client
	privateToken *string
	appID        *string

	sessionToken   string
	sessionExpires time.Time
	base           string
}

func (c *client) WithAppID(appID string) Client {
	c.appID = &appID
	return c
}

func (c *client) WithPrivateToken(privateToken types.PrivateToken) Client {
	c.privateToken = &privateToken
	return c
}

func (c *client) WithHttpClient(httpClient *http.Client) Client {
	c.httpClient = httpClient
	return c
}

type genericResponse struct {
	UID       string      `json:"uid,omitempty"`
	Message   string      `json:"msg,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
	Success   bool        `json:"success"`
	Result    interface{} `json:"result"`
}

func (c *client) genericGET(path string) (*genericResponse, error) {
	httpResponse, err := c.httpClient.Get(fmt.Sprintf("%s/%s", c.base, path))
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET %s request: %s", path, err)
	}
	return c.fromHTTPResponse(httpResponse)
}

func (c *client) genericPOST(path string, body interface{}) (*genericResponse, error) {
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(body)

	httpResponse, err := c.httpClient.Post(fmt.Sprintf("%s/%s", c.base, path), "application/json", requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET %s request: %s", path, err)
	}
	return c.fromHTTPResponse(httpResponse)
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
