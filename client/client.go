package client

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/nikolalohinski/free-go/types"
)

type Client interface {
	// Configuration
	WithAppID(string) Client
	WithPrivateToken(types.PrivateToken) Client
	WithHTTPClient(*http.Client) Client
	// Methods
	APIVersion() (types.APIVersion, error)
	Authorize(types.AuthorizationRequest) (types.PrivateToken, error)
	Login() (types.Permissions, error)
	ListPortForwardingRules() ([]types.PortForwardingRule, error)
	GetPortForwardingRule(identifier int64) (types.PortForwardingRule, error)
}

type Error string

func (e Error) Error() string {
	return string(e)
}

func New(endpoint, version string) (Client, error) {
	match, err := regexp.MatchString("^https?://.*", endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to match endpoint string against regex: %w", err)
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

	session *session
	base    string
}

type session struct {
	token   string
	expires time.Time
}

func (c *client) WithAppID(appID string) Client {
	c.appID = &appID

	return c
}

func (c *client) WithPrivateToken(privateToken types.PrivateToken) Client {
	c.privateToken = &privateToken

	return c
}

func (c *client) WithHTTPClient(httpClient *http.Client) Client {
	c.httpClient = httpClient

	return c
}
