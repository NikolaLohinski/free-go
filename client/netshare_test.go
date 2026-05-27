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

var _ = Describe("netshare", func() {
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

	// ── Samba ────────────────────────────────────────────────────────────────────

	Context("getting the Samba configuration", func() {
		returnedConfig := new(types.SambaConfiguration)
		JustBeforeEach(func() {
			*returnedConfig, *returnedErr = freeboxClient.GetSambaConfiguration(ctx)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/netshare/samba/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"file_share_enabled": true,
								"print_share_enabled": false,
								"logon_enabled": true,
								"logon_user": "freebox",
								"workgroup": "WORKGROUP",
								"smbv2_enabled": true
							}
						}`),
					),
				)
			})
			It("should return the correct Samba configuration", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedConfig).To(MatchFields(IgnoreExtras, Fields{
					"FileShareEnabled":  BeTrue(),
					"PrintShareEnabled": BeFalse(),
					"LogonEnabled":      BeTrue(),
					"LogonUser":         Equal("freebox"),
					"Workgroup":         Equal("WORKGROUP"),
					"V2Enabled":         BeTrue(),
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/netshare/samba/", version)),
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

	Context("updating the Samba configuration", func() {
		var (
			returnedConfig = new(types.SambaConfiguration)
			payload        = new(types.SambaConfigurationPayload)
		)
		BeforeEach(func() {
			*payload = types.SambaConfigurationPayload{
				SambaConfiguration: types.SambaConfiguration{
					FileShareEnabled: true,
					Workgroup:        "MYGROUP",
					LogonEnabled:     true,
					LogonUser:        "freebox",
					V2Enabled:        true,
				},
				LoginPassword: "s3cr3t",
			}
		})
		JustBeforeEach(func() {
			*returnedConfig, *returnedErr = freeboxClient.UpdateSambaConfiguration(ctx, *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/netshare/samba/", version)),
						ghttp.VerifyContentType("application/json"),
						verifyAuth(*sessionToken),
						ghttp.VerifyJSON(`{
							"file_share_enabled": true,
							"print_share_enabled": false,
							"logon_enabled": true,
							"logon_user": "freebox",
							"workgroup": "MYGROUP",
							"smbv2_enabled": true,
							"login_password": "s3cr3t"
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"file_share_enabled": true,
								"print_share_enabled": false,
								"logon_enabled": true,
								"logon_user": "freebox",
								"workgroup": "MYGROUP",
								"smbv2_enabled": true
							}
						}`),
					),
				)
			})
			It("should return the updated Samba configuration", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedConfig).To(MatchFields(IgnoreExtras, Fields{
					"FileShareEnabled": BeTrue(),
					"Workgroup":        Equal("MYGROUP"),
					"LogonEnabled":     BeTrue(),
					"LogonUser":        Equal("freebox"),
					"V2Enabled":        BeTrue(),
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
		Context("when the server returns an api error", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/netshare/samba/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": false,
							"error_code": "invalid_workgroup_name",
							"msg": "Invalid workgroup name"
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/netshare/samba/", version)),
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

	// ── AFP ──────────────────────────────────────────────────────────────────────

	Context("getting the AFP configuration", func() {
		returnedConfig := new(types.AFPConfiguration)
		JustBeforeEach(func() {
			*returnedConfig, *returnedErr = freeboxClient.GetAFPConfiguration(ctx)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/netshare/afp/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"enabled": true,
								"guest_allow": "none",
								"server_type": "macpro",
								"login_name": "afpuser"
							}
						}`),
					),
				)
			})
			It("should return the correct AFP configuration", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedConfig).To(MatchFields(IgnoreExtras, Fields{
					"Enabled":    BeTrue(),
					"GuestAllow": Equal("none"),
					"ServerType": Equal(types.NetshareAFPServerTypeMacPro),
					"LoginName":  Equal("afpuser"),
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/netshare/afp/", version)),
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

	Context("updating the AFP configuration", func() {
		var (
			returnedConfig = new(types.AFPConfiguration)
			payload        = new(types.AFPConfigurationPayload)
		)
		BeforeEach(func() {
			*payload = types.AFPConfigurationPayload{
				AFPConfiguration: types.AFPConfiguration{
					Enabled:    true,
					GuestAllow: "none",
					ServerType: types.NetshareAFPServerTypeMacBook,
					LoginName:  "afpuser",
				},
				LoginPassword: "afppass",
			}
		})
		JustBeforeEach(func() {
			*returnedConfig, *returnedErr = freeboxClient.UpdateAFPConfiguration(ctx, *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/netshare/afp/", version)),
						ghttp.VerifyContentType("application/json"),
						verifyAuth(*sessionToken),
						ghttp.VerifyJSON(`{
							"enabled": true,
							"guest_allow": "none",
							"server_type": "macbook",
							"login_name": "afpuser",
							"login_password": "afppass"
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"enabled": true,
								"guest_allow": "none",
								"server_type": "macbook",
								"login_name": "afpuser"
							}
						}`),
					),
				)
			})
			It("should return the updated AFP configuration", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedConfig).To(MatchFields(IgnoreExtras, Fields{
					"Enabled":    BeTrue(),
					"GuestAllow": Equal("none"),
					"ServerType": Equal(types.NetshareAFPServerTypeMacBook),
					"LoginName":  Equal("afpuser"),
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
		Context("when the server returns an api error", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/netshare/afp/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": false,
							"error_code": "invalid_afp_login_name",
							"msg": "Invalid AFP user name"
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/netshare/afp/", version)),
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
