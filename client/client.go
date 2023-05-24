package client

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/nikolalohinski/free-go/types"
)

type Client interface {
	APIVersion() (types.APIVersion, error)
}

type Config struct {
	Endpoint string
	Version  string
	Token    string
}

func New(config Config, httpClient ...*http.Client) (Client, error) {
	match, err := regexp.MatchString("^https?://.*", config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to match endpoint string against regex: %s", err)
	}
	if !match {
		config.Endpoint = fmt.Sprintf("http://%s", config.Endpoint)
	}
	if httpClient == nil {
		httpClient = []*http.Client{http.DefaultClient}
	}
	if len(httpClient) > 1 {
		return nil, errors.New("only one http client can be provided")
	}
	return &client{
		httpClient: httpClient[0],
		token:      config.Token,
		base:       fmt.Sprintf("%s/api/%s", config.Endpoint, config.Version),
	}, nil
}

type client struct {
	httpClient *http.Client
	token      string
	base       string
}
