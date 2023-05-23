package client

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/nikolalohinski/free-go/types"
)

type Client interface {
	APIVersion() (types.APIVersion, error)
}

func New(endpoint string, version string, token string) Client {
	if match, err := regexp.MatchString("^https?://.*", endpoint); !match || err != nil {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}
	return &client{
		httpClient: &http.Client{},
		token:      token,
		base:       fmt.Sprintf("%s/api/%s", endpoint, version),
	}
}

type client struct {
	httpClient *http.Client
	token      string
	base       string
}
