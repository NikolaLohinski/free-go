package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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
	APIVersion(context.Context) (types.APIVersion, error)
	// authentication
	Authorize(context.Context, types.AuthorizationRequest) (types.PrivateToken, error)
	Login(context.Context) (types.Permissions, error)
	Logout(context.Context) error
	// port forwarding
	ListPortForwardingRules(context.Context) ([]types.PortForwardingRule, error)
	GetPortForwardingRule(ctx context.Context, identifier int64) (types.PortForwardingRule, error)
	CreatePortForwardingRule(ctx context.Context, payload types.PortForwardingRulePayload) (types.PortForwardingRule, error)
	UpdatePortForwardingRule(ctx context.Context, identifier int64, payload types.PortForwardingRulePayload) (types.PortForwardingRule, error)
	DeletePortForwardingRule(ctx context.Context, identifier int64) error
	// lan browser
	ListLanInterfaceInfo(context.Context) ([]types.LanInfo, error)
	GetLanInterface(ctx context.Context, name string) (result []types.LanInterfaceHost, err error)
	GetLanInterfaceHost(ctx context.Context, interfaceName, identifier string) (result types.LanInterfaceHost, err error)
	// virtual machines
	GetVirtualMachineInfo(context.Context) (result types.VirtualMachinesInfo, err error)
	GetVirtualMachineDistributions(context.Context) (result []types.VirtualMachineDistribution, err error)
	ListVirtualMachines(context.Context) (result []types.VirtualMachine, err error)
	CreateVirtualMachine(ctx context.Context, payload types.VirtualMachinePayload) (result types.VirtualMachine, err error)
	GetVirtualMachine(ctx context.Context, identifier int64) (result types.VirtualMachine, err error)
	UpdateVirtualMachine(ctx context.Context, identifier int64, payload types.VirtualMachinePayload) (result types.VirtualMachine, err error)
	DeleteVirtualMachine(ctx context.Context, identifier int64) error
	StartVirtualMachine(ctx context.Context, identifier int64) error
	KillVirtualMachine(ctx context.Context, identifier int64) error
	StopVirtualMachine(ctx context.Context, identifier int64) error
	// websocket
	ListenEvents(ctx context.Context, events []types.EventDescription) (chan types.Event, error)
	// filesystem
	GetFileInfo(ctx context.Context, path string) (types.FileInfo, error)
	RemoveFiles(ctx context.Context, paths []string) (types.FileSystemTask, error)
	GetFileSystemTask(ctx context.Context, identifier int64) (types.FileSystemTask, error)
	CreateDirectory(ctx context.Context, parent, name string) (path string, err error)
	AddHashFileTask(ctx context.Context, payload types.HashPayload) (task types.FileSystemTask, err error)
	GetHashResult(ctx context.Context, identifier int64) (result string, err error)
	// downloads
	ListDownloadTasks(ctx context.Context) ([]types.DownloadTask, error)
	GetDownloadTask(ctx context.Context, identifier int64) (types.DownloadTask, error)
	AddDownloadTask(ctx context.Context, request types.DownloadRequest) (identifier int64, err error)
	DeleteDownloadTask(ctx context.Context, identifier int64) error
	EraseDownloadTask(ctx context.Context, identifier int64) error
	UpdateDownloadTask(ctx context.Context, identifier int64, payload types.DownloadTaskUpdate) error
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

var matchHTTPSRegex = regexp.MustCompile("^https?://.*")

func New(endpoint, version string) (Client, error) {
	if !matchHTTPSRegex.MatchString(endpoint) {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}

	base, err := url.Parse(fmt.Sprintf("%s/api/%s", endpoint, version))
	if err != nil {
		return nil, fmt.Errorf("can not build base url from endpoint \"%s\" and version \"%s\"", endpoint, version)
	}

	return &client{
		httpClient: http.DefaultClient,
		base:       base,
	}, nil
}

type client struct {
	httpClient   HTTPClient
	privateToken *string
	appID        *string

	session *session
	base    *url.URL
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
