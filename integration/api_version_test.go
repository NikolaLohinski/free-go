//go:build integration

package integration_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("APIVersion", func() {
	var (
		apiVersion  = new(types.APIVersion)
		returnedErr = new(error)
	)
	JustBeforeEach(func() {
		client := Must(client.New(client.Config{
			Endpoint: endpoint,
			Version:  version,
			Token:    "not used in these tests",
		})).(client.Client)

		*apiVersion, *returnedErr = client.APIVersion()
	})
	Context("default", func() {
		It("should not return an error nor unexpected response", func() {
			Expect(*returnedErr).To(BeNil())
			Expect(*apiVersion).To(MatchFields(IgnoreExtras, Fields{
				"APIBaseURL": Equal("/api/"),
				"HTTPSPort":  Not(BeZero()),
				"DeviceName": Equal("Freebox Server"),
				"APIVersion": MatchRegexp(`%s.\d+`, strings.TrimLeft(version, "v")),
			}))
		})
	})
})
