package types_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "type-tests")
}

func Must(r interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return r
}
