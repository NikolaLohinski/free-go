package integration_test

import (
	"os"
	"testing"

	"github.com/nikolalohinski/free-go/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	version  string
	endpoint string
	token    string
	appID    string

	freeboxClient client.Client
)

func init() {
	endpoint = os.Getenv("FREEBOX_ENDPOINT")
	if endpoint == "" {
		panic("FREEBOX_ENDPOINT environment variable must be set")
	}
	version = os.Getenv("FREEBOX_VERSION")
	if version == "" {
		panic("FREEBOX_VERSION environment variable must be set")
	}
	appID = os.Getenv("FREEBOX_APP_ID")
	if appID == "" {
		panic("FREEBOX_APP_ID environment variable must be set")
	}
	token = os.Getenv("FREEBOX_TOKEN")
	if token == "" {
		panic("FREEBOX_TOKEN environment variable must be set")
	}

	freeboxClient = Must(client.New(endpoint, version))
}

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "integration")
}

func Must[T interface{}](r T, err error) T {
	if err != nil {
		panic(err)
	}
	return r
}
