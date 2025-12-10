package client_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("lan", func() {
	var (
		freeboxClient client.Client

		server   *ghttp.Server
		endpoint = new(string)

		sessionToken = new(string)

		returnedErr = new(error)
	)
	BeforeEach(func() {
		server = ghttp.NewServer()
		DeferCleanup(server.Close)

		*endpoint = server.Addr()

		freeboxClient = Must(client.New(*endpoint, version)).
			WithAppID(appID).
			WithPrivateToken(privateToken)

		*sessionToken = setupLoginFlow(server)
	})

	Context("getting lan config", func() {
		returnedLanConfig := new(types.LanConfig)
		JustBeforeEach(func() {
			*returnedLanConfig, *returnedErr = freeboxClient.GetLanConfig(context.Background())
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/config/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"name_dns": "freebox-r0ro",
								"name_mdns": "Freebox-r0ro",
								"name": "Freebox r0ro",
								"mode": "router",
								"name_netbios": "Freebox_r0ro",
								"ip": "192.168.1.254"
							}
						}`),
					),
				)
			})
			It("should return the correct lan config", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedLanConfig).To(Equal(types.LanConfig{
					NameDNS:     "freebox-r0ro",
					NameMDNS:    "Freebox-r0ro",
					Name:        "Freebox r0ro",
					Mode:        types.LanModeRouter,
					NameNetBIOS: "Freebox_r0ro",
					IP:          "192.168.1.254",
				}))
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
		Context("when server fails to respond", func() {
			BeforeEach(func() {
				server.Close()
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when reading the body returns an error", func() {
			BeforeEach(func() {
				freeboxClient = Must(client.New(*endpoint, version)).
					WithAppID(appID).
					WithPrivateToken(privateToken).
					WithHTTPClient(mockHTTPClient{
						returnedBody: errorReader{},
					})
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when closing the body returns an error while an error was already there", func() {
			BeforeEach(func() {
				freeboxClient = Must(client.New(*endpoint, version)).
					WithAppID(appID).
					WithPrivateToken(privateToken).
					WithHTTPClient(mockHTTPClient{
						statusCode: http.StatusInternalServerError,
						returnedBody: errorCloser{
							strings.NewReader(`{}`),
						},
					})
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when the sessions has expired and trying to login again fails", func() {
			BeforeEach(func() {
				DeferCleanup(func(previous time.Duration) {
					client.LoginSessionTTL = previous
				}, client.LoginSessionTTL)
				client.LoginSessionTTL = 0

				Must(freeboxClient.Login(context.Background()))
				server.Close()
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("updating lan config", func() {
		const (
			interfaceName = "pub"
		)
		returnedLanConfig := new(types.LanConfig)
		JustBeforeEach(func() {
			*returnedLanConfig, *returnedErr = freeboxClient.UpdateLanConfig(context.Background(), *returnedLanConfig)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/lan/config/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"name_dns": "freebox-r0ro",
								"name_mdns": "Freebox-r0ro",
								"name": "Freebox r0ro",
								"mode": "router",
								"name_netbios": "Freebox_r0ro",
								"ip": "192.168.1.254"
							}
						}`),
					),
				)
			})
			It("should return the correct lan config", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedLanConfig).To(Equal(types.LanConfig{
					NameDNS:     "freebox-r0ro",
					NameMDNS:    "Freebox-r0ro",
					Name:        "Freebox r0ro",
					Mode:        types.LanModeRouter,
					NameNetBIOS: "Freebox_r0ro",
					IP:          "192.168.1.254",
				}))
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
	})
})
