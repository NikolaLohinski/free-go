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

var _ = Describe("vpn", func() {
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

	// ── OpenVPN server config ───────────────────────────────────────────────────

	Context("getting the OpenVPN server config", func() {
		returnedConfig := new(types.OpenVPNServerConfig)
		JustBeforeEach(func() {
			*returnedConfig, *returnedErr = freeboxClient.GetOpenVPNServerConfig(ctx)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vpn/openvpn/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"enabled": true,
								"server_port": 1194,
								"server_ip": "10.8.0.0",
								"server_mask": "255.255.255.0",
								"push_default_gw": false,
								"push_dhcp": true,
								"ca": "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----\n"
							}
						}`),
					),
				)
			})
			It("should return the correct config", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedConfig).To(MatchFields(IgnoreExtras, Fields{
					"Enabled":    Equal(true),
					"ServerPort": Equal(int64(1194)),
					"ServerIP":   Equal("10.8.0.0"),
					"ServerMask": Equal("255.255.255.0"),
					"PushDHCP":   Equal(true),
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
		Context("when the context is nil", func() {
			BeforeEach(func() {
				ctx = nil
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})

	Context("updating the OpenVPN server config", func() {
		var (
			returnedConfig = new(types.OpenVPNServerConfig)
			payload        = new(types.OpenVPNServerConfig)
		)
		BeforeEach(func() {
			*payload = types.OpenVPNServerConfig{
				Enabled:    true,
				ServerPort: 1194,
				ServerIP:   "10.8.0.0",
				ServerMask: "255.255.255.0",
			}
		})
		JustBeforeEach(func() {
			*returnedConfig, *returnedErr = freeboxClient.UpdateOpenVPNServerConfig(ctx, *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/vpn/openvpn/", version)),
						ghttp.VerifyContentType("application/json"),
						verifyAuth(*sessionToken),
						ghttp.VerifyJSON(`{
							"enabled": true,
							"server_port": 1194,
							"server_ip": "10.8.0.0",
							"server_mask": "255.255.255.0",
							"push_default_gw": false,
							"push_dhcp": false
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"enabled": true,
								"server_port": 1194,
								"server_ip": "10.8.0.0",
								"server_mask": "255.255.255.0",
								"push_default_gw": false,
								"push_dhcp": false
							}
						}`),
					),
				)
			})
			It("should return the updated config", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedConfig).To(MatchFields(IgnoreExtras, Fields{
					"Enabled":    Equal(true),
					"ServerPort": Equal(int64(1194)),
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
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/vpn/openvpn/", version)),
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

	// ── VPN users ───────────────────────────────────────────────────────────────

	Context("listing VPN users", func() {
		returnedUsers := new([]types.VPNUser)
		JustBeforeEach(func() {
			*returnedUsers, *returnedErr = freeboxClient.ListVPNUsers(ctx)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vpn/user/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": [
								{"login": "perreux", "password": "", "description": "Perreux NAS"}
							]
						}`),
					),
				)
			})
			It("should return the correct users", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedUsers).To(HaveLen(1))
				Expect((*returnedUsers)[0]).To(MatchFields(IgnoreExtras, Fields{
					"VPNUserPayload": MatchFields(IgnoreExtras, Fields{
						"Login":       Equal("perreux"),
						"Description": Equal("Perreux NAS"),
					}),
				}))
			})
		})
		Context("when there are no VPN users", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vpn/user/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{"success": true}`),
					),
				)
			})
			It("should return an empty list", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedUsers).To(BeEmpty())
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

	Context("getting a VPN user", func() {
		const login = "perreux"
		returnedUser := new(types.VPNUser)
		JustBeforeEach(func() {
			*returnedUser, *returnedErr = freeboxClient.GetVPNUser(ctx, login)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vpn/user/%s", version, login)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {"login": "perreux", "password": "", "description": "Perreux NAS"}
						}`),
					),
				)
			})
			It("should return the correct user", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedUser).To(MatchFields(IgnoreExtras, Fields{
					"VPNUserPayload": MatchFields(IgnoreExtras, Fields{
						"Login": Equal("perreux"),
					}),
				}))
			})
		})
		Context("when the VPN user is not found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vpn/user/%s", version, login)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": false,
							"error_code": "noent",
							"msg": "User not found"
						}`),
					),
				)
			})
			It("should return ErrVPNUserNotFound", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrVPNUserNotFound))
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

	Context("creating a VPN user", func() {
		var (
			returnedUser = new(types.VPNUser)
			payload      = new(types.VPNUserPayload)
		)
		BeforeEach(func() {
			*payload = types.VPNUserPayload{
				Login:       "perreux",
				Password:    "secret",
				Description: "Perreux NAS",
			}
		})
		JustBeforeEach(func() {
			*returnedUser, *returnedErr = freeboxClient.CreateVPNUser(ctx, *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vpn/user/", version)),
						ghttp.VerifyContentType("application/json"),
						verifyAuth(*sessionToken),
						ghttp.VerifyJSON(`{
							"login": "perreux",
							"password": "secret",
							"description": "Perreux NAS"
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {"login": "perreux", "password": "", "description": "Perreux NAS"}
						}`),
					),
				)
			})
			It("should return the created user", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedUser).To(MatchFields(IgnoreExtras, Fields{
					"VPNUserPayload": MatchFields(IgnoreExtras, Fields{
						"Login": Equal("perreux"),
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
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vpn/user/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{"success": true, "result": []}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})

	Context("updating a VPN user", func() {
		const login = "perreux"
		var (
			returnedUser = new(types.VPNUser)
			payload      = new(types.VPNUserPayload)
		)
		BeforeEach(func() {
			*payload = types.VPNUserPayload{
				Login:    "perreux",
				Password: "new-secret",
			}
		})
		JustBeforeEach(func() {
			*returnedUser, *returnedErr = freeboxClient.UpdateVPNUser(ctx, login, *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/vpn/user/%s", version, login)),
						ghttp.VerifyContentType("application/json"),
						verifyAuth(*sessionToken),
						ghttp.VerifyJSON(`{"login": "perreux", "password": "new-secret"}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {"login": "perreux", "password": ""}
						}`),
					),
				)
			})
			It("should return the updated user", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedUser).To(MatchFields(IgnoreExtras, Fields{
					"VPNUserPayload": MatchFields(IgnoreExtras, Fields{
						"Login": Equal("perreux"),
					}),
				}))
			})
		})
		Context("when the VPN user is not found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/vpn/user/%s", version, login)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": false,
							"error_code": "noent",
							"msg": "User not found"
						}`),
					),
				)
			})
			It("should return ErrVPNUserNotFound", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrVPNUserNotFound))
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

	Context("deleting a VPN user", func() {
		const login = "perreux"
		JustBeforeEach(func() {
			*returnedErr = freeboxClient.DeleteVPNUser(ctx, login)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/vpn/user/%s", version, login)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{"success": true}`),
					),
				)
			})
			It("should not return an error", func() {
				Expect(*returnedErr).To(BeNil())
			})
		})
		Context("when the VPN user is not found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/vpn/user/%s", version, login)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": false,
							"error_code": "noent",
							"msg": "User not found"
						}`),
					),
				)
			})
			It("should return ErrVPNUserNotFound", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrVPNUserNotFound))
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
		Context("when the context is nil", func() {
			BeforeEach(func() {
				ctx = nil
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})

	// ── VPN client config ───────────────────────────────────────────────────────

	Context("getting the VPN client config", func() {
		const login = "perreux"
		returnedConfig := new(string)
		JustBeforeEach(func() {
			*returnedConfig, *returnedErr = freeboxClient.GetVPNUserClientConfig(ctx, login)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vpn/user/%s/config/openvpn", version, login)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": "client\nproto udp\nremote freebox.example.com 1194\n"
						}`),
					),
				)
			})
			It("should return the .ovpn content", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedConfig).To(ContainSubstring("client"))
				Expect(*returnedConfig).To(ContainSubstring("proto udp"))
			})
		})
		Context("when the VPN user is not found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vpn/user/%s/config/openvpn", version, login)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": false,
							"error_code": "noent",
							"msg": "User not found"
						}`),
					),
				)
			})
			It("should return ErrVPNUserNotFound", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrVPNUserNotFound))
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
