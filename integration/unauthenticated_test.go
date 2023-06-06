//go:build integration

package integration_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("unauthenticated scenarios", func() {
	Context("getting API version information", func() {
		It("should not return an error nor unexpected responses", func() {
			apiVersion, err := freeboxClient.APIVersion()
			Expect(err).To(BeNil())
			Expect(apiVersion).To(MatchFields(IgnoreExtras, Fields{
				"APIBaseURL": Equal("/api/"),
				"HTTPSPort":  Not(BeZero()),
				"DeviceName": Equal("Freebox Server"),
				"APIVersion": MatchRegexp(`%s.\d+`, strings.TrimLeft(version, "v")),
			}))
		})
	})
})
