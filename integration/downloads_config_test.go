//go:build integration

package integration_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("downloads config", func() {
	Context("getting the download configuration", func() {
		It("should return the download configuration without error", func() {
			client := freeboxClient.WithAppID(appID).WithPrivateToken(token)
			config, err := client.GetDownloadConfiguration(context.Background())
			Expect(err).ToNot(HaveOccurred())
			Expect(config).ToNot(BeNil())
		})
	})
})
