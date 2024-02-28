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
	WithHTTPClient(HTTPClient) Client
	// unauthenticated
	APIVersion() (types.APIVersion, error)
	// authentication
	Authorize(types.AuthorizationRequest) (types.PrivateToken, error)
	Login() (types.Permissions, error)
	Logout() error
	// port forwarding
	ListPortForwardingRules() ([]types.PortForwardingRule, error)
	GetPortForwardingRule(identifier int64) (types.PortForwardingRule, error)
	CreatePortForwardingRule(payload types.PortForwardingRulePayload) (types.PortForwardingRule, error)
	UpdatePortForwardingRule(identifier int64, payload types.PortForwardingRulePayload) (types.PortForwardingRule, error)
	DeletePortForwardingRule(identifier int64) error
	// lan browser
	ListLanInterfaceInfo() ([]types.LanInfo, error)
	GetLanInterface(name string) (result []types.LanInterfaceHost, err error)
	GetLanInterfaceHost(interfaceName, identifier string) (result types.LanInterfaceHost, err error)
	// virtual machines
	GetVirtualMachineInfo() (result types.VirtualMachinesInfo, err error)
	GetVirtualMachineDistributions() (result []types.VirtualMachineDistribution, err error)
	ListVirtualMachines() (result []types.VirtualMachine, err error)
	CreateVirtualMachine(payload types.VirtualMachinePayload) (result types.VirtualMachine, err error)
	GetVirtualMachine(identifier int64) (result types.VirtualMachine, err error)
	DeleteVirtualMachine(identifier int64) error
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

var matchHTTPSRegex = regexp.MustCompile("^https?://.*")

func New(endpoint, version string) (Client, error) {
	if !matchHTTPSRegex.MatchString(endpoint) {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}

	return &client{
		httpClient: http.DefaultClient,
		base:       fmt.Sprintf("%s/api/%s", endpoint, version),
	}, nil
}

type client struct {
	httpClient   HTTPClient
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

func (c *client) WithHTTPClient(httpClient HTTPClient) Client {
	c.httpClient = httpClient

	return c
}
