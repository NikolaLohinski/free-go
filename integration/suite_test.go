//go:build integration

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
	root     string

	freeboxClient client.Client
)

func init() {
	var ok bool
	endpoint, ok = os.LookupEnv("FREEBOX_ENDPOINT")
	if !ok {
		panic("FREEBOX_ENDPOINT environment variable must be set")
	}
	version, ok = os.LookupEnv("FREEBOX_VERSION")
	if !ok {
		panic("FREEBOX_VERSION environment variable must be set")
	}
	appID, ok = os.LookupEnv("FREEBOX_APP_ID")
	if !ok {
		panic("FREEBOX_APP_ID environment variable must be set")
	}
	token, ok = os.LookupEnv("FREEBOX_TOKEN")
	if !ok {
		panic("FREEBOX_TOKEN environment variable must be set")
	}
	root, ok = os.LookupEnv("FREEBOX_ROOT")
	if !ok {
		root = "Freebox"
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
