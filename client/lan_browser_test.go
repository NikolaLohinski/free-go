package client_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("lan browser", func() {
	var (
		freeboxClient client.Client

		server   *ghttp.Server
		endpoint = new(string)

		sessionToken = new(string)

		returnedErr = new(error)
	)
	BeforeEach(func() {
		server = ghttp.NewServer()
		*endpoint = server.Addr()

		freeboxClient = Must(client.New(*endpoint, version)).(client.Client).
			WithAppID(appID).
			WithPrivateToken(privateToken)

		*sessionToken = setupLoginFlow(server)
	})
	AfterEach(func() {
		server.Close()
	})
	Context("listing lan metadata", func() {
		returnedLanInfos := new([]types.LanInfo)
		JustBeforeEach(func() {
			*returnedLanInfos, *returnedErr = freeboxClient.ListLanInterfaceInfo()
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/browser/interfaces/", version)),
						ghttp.RespondWith(http.StatusOK, `{
						    "success": true,
						    "result": [
						        {
						            "name": "pub",
						            "host_count": 41
						        },
						        {
						            "name": "wifiguest",
						            "host_count": 0
						        }
						    ]
						}`),
					),
				)
			})
			It("should return the correct lan info", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedLanInfos).ToNot(BeEmpty())
				Expect(*returnedLanInfos).To(HaveLen(2))
				Expect(*returnedLanInfos).To(ContainElements(
					Equal(types.LanInfo{
						Name:      "pub",
						HostCount: 41,
					}),
					Equal(types.LanInfo{
						Name:      "wifiguest",
						HostCount: 0,
					}),
				))
			})
		})
		Context("when server fails to respond", func() {
			BeforeEach(func() {
				server.Close()
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/browser/interfaces/", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"foo": "bar"
							}
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
})
