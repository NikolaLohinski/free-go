package client_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/onsi/gomega/gstruct"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("Port forwarding", func() {
	const (
		version      = "v0"
		appID        = "test"
		privateToken = "xxxxxxxxxxxxxxxxxxxxxxxxx"
		sessionToken = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	)
	var (
		freeboxClient client.Client

		server   *ghttp.Server
		endpoint = new(string)

		returnedErr = new(error)
	)
	BeforeEach(func() {
		server = ghttp.NewServer()
		*endpoint = server.Addr()

		freeboxClient = Must(client.New(*endpoint, version)).(client.Client).WithAppID(appID).WithPrivateToken(privateToken)

		// setup login flow
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/login", version)),
				ghttp.RespondWith(http.StatusOK, `{
					"success": true,
					"result": {
						"logged_in": false,
						"challenge": "9Va31tSgQWM853j0kSCtBUyzYNhPN7IY"
					}
				}`),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/login/session", version)),
				ghttp.VerifyContentType("application/json"),
				ghttp.VerifyJSON(`{
				    "app_id": "`+appID+`",
				    "password": "aac0e1f841e689230f7e2f5815807918ced96e82"
				}`),
				ghttp.RespondWith(http.StatusOK, `{
					"result": {
						"session_token": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
						"challenge": "9Va31tSgQWM853j0kSCtBUyzYNhPN7IY",
						"permissions": {}
					},
					"success": true
				}`),
			),
		)
	})
	AfterEach(func() {
		server.Close()
	})
	Context("listing port forwarding rules", func() {
		var (
			returnedRules = new([]types.PortForwardingRule)
		)
		JustBeforeEach(func() {
			*returnedRules, *returnedErr = freeboxClient.PortForwardingRules()
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fw/redir/", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": [
								{
									"enabled": true,
									"comment": "test",
									"id": 5,
									"valid": true,
									"host": {
										"l2ident": {
											"id": "7E:EC:37:CD:5B:6A",
											"type": "mac_address"
										},
										"active": false,
										"persistent": false,
										"names": [
											{
												"name": "test",
												"source": "dhcp"
											}
										],
										"vendor_name": "",
										"host_type": "workstation",
										"interface": "pub",
										"id": "ether-7e:ec:37:cd:5b:6a",
										"last_time_reachable": 1682579132,
										"primary_name_manual": false,
										"l3connectivities": [
											{
												"addr": "192.168.1.254",
												"active": false,
												"reachable": false,
												"last_activity": 1682579111,
												"af": "ipv4",
												"last_time_reachable": 1682579111
											}
										],
										"default_name": "testing",
										"first_activity": 1682578724,
										"reachable": false,
										"last_activity": 1682579132,
										"primary_name": "testing"
									},
									"src_ip": "0.0.0.0",
									"hostname": "testing",
									"lan_port": 80,
									"wan_port_end": 12345,
									"wan_port_start": 12345,
									"lan_ip": "192.168.1.254",
									"ip_proto": "tcp"
								}
							]
						}`),
					),
				)
			})
			It("should return the private token provided by the server", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedRules).ToNot(BeEmpty())
				Expect(*returnedRules).To(HaveLen(1))
				Expect((*returnedRules)[0]).To(MatchFields(IgnoreExtras, Fields{
					"Hostname": Equal("testing"),
					"LanPort":  Equal(int64(80)),
					"ProtoIP":  Equal("tcp"),
					"Host": MatchFields(IgnoreExtras, Fields{
						"LastActivity": Equal(types.Timestamp{Time: time.Unix(1682579132, 0)}),
					}),
				}))
			})
		})

		Context("when there are no port forwarding rules", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fw/redir/", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true
						}`),
					),
				)
			})
			It("should return the private token provided by the server", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedRules).To(BeEmpty())
			})
		})

	})
})
