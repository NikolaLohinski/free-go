package client_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client suite")
}

func Must(r interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return r
}
