package client

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/nikolalohinski/free-go/types"
)

//nolint:interfacebloat
type Client interface {
	// configuration
	WithAppID(string) Client
	WithPrivateToken(types.PrivateToken) Client
	WithHTTPClient(*http.Client) Client
	// unauthenticated
	APIVersion() (types.APIVersion, error)
	Authorize(types.AuthorizationRequest) (types.PrivateToken, error)
	Login() (types.Permissions, error)
	// port forwarding
	ListPortForwardingRules() ([]types.PortForwardingRule, error)
	GetPortForwardingRule(identifier int64) (types.PortForwardingRule, error)
	CreatePortForwardingRule(payload types.PortForwardingRulePayload) (types.PortForwardingRule, error)
	UpdatePortForwardingRule(identifier int64, payload types.PortForwardingRulePayload) (types.PortForwardingRule, error)
	DeletePortForwardingRule(identifier int64) error
	// lan browser
	ListLanInterfaceInfo() ([]types.LanInfo, error)
}

type Error string

func (e Error) Error() string {
	return string(e)
}

var matchHTTPS = regexp.MustCompile("^https?://.*")

func New(endpoint, version string) (Client, error) {
	if !matchHTTPS.MatchString(endpoint) {
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
