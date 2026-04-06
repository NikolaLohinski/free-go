package client_test

import (
	"context"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("system", func() {
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

	Context("getting the system info", func() {
		returnedSystemInfo := new(types.SystemConfig)
		JustBeforeEach(func() {
			*returnedSystemInfo, *returnedErr = freeboxClient.GetSystemInfo(context.Background())
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/system/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"mac": "F4:CA:E5:5C:EA:14",
								"box_flavor": "light",
								"temp_cpub": 63,
								"disk_status": "active",
								"box_authenticated": true,
								"board_name": "fbxgw1r",
								"fan_rpm": 1832,
								"temp_sw": 52,
								"uptime": "6 jours 22 heures 9 minutes 46 secondes",
								"uptime_val": 598186,
								"user_main_storage": "Disque 1",
								"temp_cpum": 62,
								"serial": "805400T144100853",
								"firmware_version": "6.6.6"
							}
						}`),
					),
				)
			})
			It("should return the correct system info", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedSystemInfo).To(Equal(types.SystemConfig{
					Mac:              "F4:CA:E5:5C:EA:14",
					BoxFlavor:        types.BoxFlavorLight,
					TempCPUB:         63,
					DiskStatus:       types.DiskStatusActive,
					BoxAuthenticated: true,
					BoardName:        "fbxgw1r",
					FanRPM:           1832,
					TempSW:           52,
					Uptime:           "6 jours 22 heures 9 minutes 46 secondes",
					UptimeVal:        598186,
					UserMainStorage:  "Disque 1",
					TempCPUM:         62,
					Serial:           "805400T144100853",
					FirmwareVersion:  "6.6.6",
				}))
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
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/system/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": "not-an-object"
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
