//go:build integration

package integration_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	version  string
	endpoint string
	token    string
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
	token = os.Getenv("FREEBOX_TOKEN")
	if token == "" {
		panic("FREEBOX_TOKEN environment variable must be set")
	}
}

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "integration")
}

func Must(r interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return r
}
