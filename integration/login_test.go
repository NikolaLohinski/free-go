//go:build integration

package integration_test

import (
	"context"

	"github.com/nikolalohinski/free-go/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("login scenarios", func() {
	Context("logging in with the provided token", func() {
		It("should not return an error nor unexpected responses", func() {
			permissions, err := freeboxClient.WithAppID(appID).WithPrivateToken(token).Login(context.Background())
			Expect(err).To(BeNil())
			Expect(permissions).ToNot(BeNil())
			Expect(permissions).To(BeAssignableToTypeOf(types.Permissions{}))
		})
	})
})
