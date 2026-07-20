package client_test

import (
	"context"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/onsi/gomega/gstruct"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("profile", func() {
	var (
		freeboxClient client.Client

		ctx context.Context

		server   *ghttp.Server
		endpoint = new(string)

		sessionToken = new(string)

		returnedErr = new(error)
	)

	BeforeEach(func() {
		ctx = context.Background()

		server = ghttp.NewServer()
		DeferCleanup(server.Close)

		*endpoint = server.Addr()

		freeboxClient = Must(client.New(*endpoint, version)).
			WithAppID(appID).
			WithPrivateToken(privateToken)

		*sessionToken = setupLoginFlow(server)
	})

	// ── ListProfiles ─────────────────────────────────────────────────────────────

	Context("listing profiles", func() {
		returnedProfiles := new([]types.Profile)
		JustBeforeEach(func() {
			*returnedProfiles, *returnedErr = freeboxClient.ListProfiles(ctx)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/profile/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": [
								{
									"id": 1,
									"name": "Pierre",
									"icon": "/resources/images/profile/profile_02.png"
								},
								{
									"id": 2,
									"name": "Nathalie",
									"icon": "/resources/images/profile/profile_01.png"
								}
							]
						}`),
					),
				)
			})
			It("should return the correct list of profiles", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedProfiles).To(HaveLen(2))
				Expect((*returnedProfiles)[0]).To(MatchFields(IgnoreExtras, Fields{
					"ID":   Equal(int64(1)),
					"Name": Equal("Pierre"),
					"Icon": Equal("/resources/images/profile/profile_02.png"),
				}))
				Expect((*returnedProfiles)[1]).To(MatchFields(IgnoreExtras, Fields{
					"ID":   Equal(int64(2)),
					"Name": Equal("Nathalie"),
					"Icon": Equal("/resources/images/profile/profile_01.png"),
				}))
			})
		})
		Context("when the result is empty", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/profile/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{"success": true}`),
					),
				)
			})
			It("should return an empty slice without error", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedProfiles).To(BeEmpty())
			})
		})
		Context("when the server fails to respond", func() {
			BeforeEach(func() {
				server.Close()
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
})
