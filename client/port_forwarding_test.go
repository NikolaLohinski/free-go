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

var _ = Describe("port forwarding", func() {
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
	Context("listing port forwarding rules", func() {
		returnedRules := new([]types.PortForwardingRule)
		JustBeforeEach(func() {
			*returnedRules, *returnedErr = freeboxClient.ListPortForwardingRules()
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fw/redir/", version)),
						verifyAuth(*sessionToken),
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
			It("should return the correct forwarding rules", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedRules).ToNot(BeEmpty())
				Expect(*returnedRules).To(HaveLen(1))
				Expect((*returnedRules)[0]).To(MatchFields(IgnoreExtras, Fields{
					"Hostname": Equal("testing"),
					"Host": MatchFields(IgnoreExtras, Fields{
						"LastActivity": Equal(types.Timestamp{Time: time.Unix(1682579132, 0)}),
					}),
					"PortForwardingRulePayload": MatchFields(IgnoreExtras, Fields{
						"LanPort":    Equal(int64(80)),
						"IPProtocol": Equal(types.TCP),
					}),
				}))
			})
		})
		Context("when there are no port forwarding rules", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fw/redir/", version)),
						verifyAuth(*sessionToken),
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fw/redir/", version)),
						verifyAuth(*sessionToken),
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
	Context("getting a port forwarding rule", func() {
		const (
			identifier = int64(5)
		)
		returnedRule := new(types.PortForwardingRule)
		JustBeforeEach(func() {
			*returnedRule, *returnedErr = freeboxClient.GetPortForwardingRule(identifier)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fw/redir/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
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
						}`),
					),
				)
			})

			It("should return the correct forwarding rule", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedRule).To(MatchFields(IgnoreExtras, Fields{
					"Hostname": Equal("testing"),
					"Host": MatchFields(IgnoreExtras, Fields{
						"LastActivity": Equal(types.Timestamp{Time: time.Unix(1682579132, 0)}),
					}),
					"ID": Equal(int64(5)),
					"PortForwardingRulePayload": MatchFields(IgnoreExtras, Fields{
						"LanPort":    Equal(int64(80)),
						"IPProtocol": Equal(types.TCP),
					}),
				}))
			})
		})
		Context("when the port forwarding rule is not found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fw/redir/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"msg": "Impossible de récupérer la redirection : Entrée non trouvée",
							"success": false,
							"error_code": "noent"
						}`),
					),
				)
			})
			It("should return the correct error", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrPortForwardingRuleNotFound))
				Expect(client.ErrPortForwardingRuleNotFound.Error()).To(Equal("port forwarding rule not found"))
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fw/redir/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": []
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("creating a port forwarding rule", func() {
		var (
			returnedRule = new(types.PortForwardingRule)
			payload      = new(types.PortForwardingRulePayload)
		)
		BeforeEach(func() {
			enabled := true
			*payload = types.PortForwardingRulePayload{
				Enabled:      &enabled,
				Comment:      "test",
				SourceIP:     "0.0.0.0",
				LanPort:      80,
				WanPortStart: 12345,
				WanPortEnd:   12345,
				LanIP:        "192.168.1.254",
				IPProtocol:   types.TCP,
			}
		})
		JustBeforeEach(func() {
			*returnedRule, *returnedErr = freeboxClient.CreatePortForwardingRule(*payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fw/redir/", version)),
						ghttp.VerifyContentType("application/json"),
						verifyAuth(*sessionToken),
						ghttp.VerifyJSON(`{
							"enabled": true,
							"comment": "test",
							"src_ip": "0.0.0.0",
							"lan_port": 80,
							"wan_port_end": 12345,
							"wan_port_start": 12345,
							"lan_ip": "192.168.1.254",
							"ip_proto": "tcp"
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
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
						}`),
					),
				)
			})

			It("should return newly created forwarding rule", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedRule).To(MatchFields(IgnoreExtras, Fields{
					"Hostname": Equal("testing"),
					"Host": MatchFields(IgnoreExtras, Fields{
						"LastActivity": Equal(types.Timestamp{Time: time.Unix(1682579132, 0)}),
					}),
					"ID":    Equal(int64(5)),
					"Valid": Equal(true),
					"PortForwardingRulePayload": MatchFields(IgnoreExtras, Fields{
						"LanPort":    Equal(int64(80)),
						"IPProtocol": Equal(types.TCP),
					}),
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
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fw/redir/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": []
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("deleting a port forwarding rule", func() {
		const (
			identifier = int64(5)
		)
		JustBeforeEach(func() {
			*returnedErr = freeboxClient.DeletePortForwardingRule(identifier)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/fw/redir/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true
						}`),
					),
				)
			})

			It("should not return an error", func() {
				Expect(*returnedErr).To(BeNil())
			})
		})
		Context("when the port forwarding rule is not found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/fw/redir/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"msg": "Impossible de récupérer la redirection : Entrée non trouvée",
							"success": false,
							"error_code": "noent"
						}`),
					),
				)
			})
			It("should return the correct error", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrPortForwardingRuleNotFound))
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
	Context("updating a port forwarding rule", func() {
		const (
			identifier = int64(5)
		)
		var (
			returnedRule = new(types.PortForwardingRule)
			payload      = new(types.PortForwardingRulePayload)
		)
		BeforeEach(func() {
			*payload = types.PortForwardingRulePayload{
				Enabled: new(bool),
			}
		})
		JustBeforeEach(func() {
			*returnedRule, *returnedErr = freeboxClient.UpdatePortForwardingRule(identifier, *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/fw/redir/%d", version, identifier)),
						ghttp.VerifyContentType("application/json"),
						verifyAuth(*sessionToken),
						ghttp.VerifyJSON(`{
							"enabled": false
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"enabled": false,
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
						}`),
					),
				)
			})

			It("should return the updated forwarding rule", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedRule).To(MatchFields(IgnoreExtras, Fields{
					"Hostname": Equal("testing"),
					"Host": MatchFields(IgnoreExtras, Fields{
						"LastActivity": Equal(types.Timestamp{Time: time.Unix(1682579132, 0)}),
					}),
					"ID":    Equal(int64(5)),
					"Valid": Equal(true),
					"PortForwardingRulePayload": MatchFields(IgnoreExtras, Fields{
						"Enabled":    PointTo(Equal(false)),
						"LanPort":    Equal(int64(80)),
						"IPProtocol": Equal(types.TCP),
					}),
				}))
			})
		})
		Context("when the port forwarding rule is not found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/fw/redir/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"msg": "Impossible de récupérer la redirection : Entrée non trouvée",
							"success": false,
							"error_code": "noent"
						}`),
					),
				)
			})
			It("should return the correct error", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrPortForwardingRuleNotFound))
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
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/fw/redir/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": []
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
