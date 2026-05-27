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

var _ = Describe("network control", func() {
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

	// ── ListNetworkControl ──────────────────────────────────────────────────────

	Context("listing network controls", func() {
		returnedControls := new([]types.NetworkControlInfo)
		JustBeforeEach(func() {
			*returnedControls, *returnedErr = freeboxClient.ListNetworkControl(ctx)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/network_control/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": [
								{
									"profile_id": 1,
									"next_change": 0,
									"override_mode": "allowed",
									"current_mode": "allowed",
									"rule_mode": "allowed",
									"override_until": 0,
									"override": false,
									"macs": ["aa:bb:cc:dd:ee:ff"],
									"hosts": [],
									"resolution": 3600,
									"cdayranges": []
								}
							]
						}`),
					),
				)
			})
			It("should return the correct list of network controls", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedControls).To(HaveLen(1))
				Expect((*returnedControls)[0]).To(MatchFields(IgnoreExtras, Fields{
					"ProfileID":    Equal(int64(1)),
					"OverrideMode": Equal(types.RuleModeAllowed),
					"CurrentMode":  Equal(types.RuleModeAllowed),
					"Override":     Equal(false),
					"Macs":         ConsistOf("aa:bb:cc:dd:ee:ff"),
				}))
			})
		})
		Context("when the result is empty", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/network_control/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{"success": true}`),
					),
				)
			})
			It("should return an empty slice without error", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedControls).To(BeEmpty())
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

	// ── GetNetworkControl ───────────────────────────────────────────────────────

	Context("getting a network control", func() {
		const identifier = int64(42)
		returnedControl := new(types.NetworkControlInfo)
		JustBeforeEach(func() {
			*returnedControl, *returnedErr = freeboxClient.GetNetworkControl(ctx, identifier)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/network_control/42", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"profile_id": 42,
								"next_change": 0,
								"override_mode": "denied",
								"current_mode": "denied",
								"rule_mode": "allowed",
								"override_until": 3600,
								"override": true,
								"macs": ["11:22:33:44:55:66"],
								"hosts": [],
								"resolution": 1800,
								"cdayranges": [":fr_bank_holidays"]
							}
						}`),
					),
				)
			})
			It("should return the correct network control", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedControl).To(MatchFields(IgnoreExtras, Fields{
					"ProfileID":     Equal(int64(42)),
					"OverrideMode":  Equal(types.RuleModeDenied),
					"CurrentMode":   Equal(types.RuleModeDenied),
					"RuleMode":      Equal(types.RuleModeAllowed),
					"Override":      Equal(true),
					"OverrideUntil": Equal(3600),
					"Macs":          ConsistOf("11:22:33:44:55:66"),
					"CustomDayRanges": ConsistOf(types.DayRangeFrenchBankHolidays),
				}))
			})
		})
		Context("when the network control is not found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/network_control/42", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": false,
							"error_code": "noent",
							"msg": "Network control not found"
						}`),
					),
				)
			})
			It("should return ErrNetworkControlNotFound", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrNetworkControlNotFound))
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

	// ── UpdateNetworkControl ────────────────────────────────────────────────────

	Context("updating a network control", func() {
		var (
			returnedControl = new(types.NetworkControlInfo)
			payload         = new(types.NetworkControlPayload)
		)
		BeforeEach(func() {
			*payload = types.NetworkControlPayload{
				ProfileID:    42,
				OverrideMode: types.RuleModeDenied,
				Override:     true,
				OverrideUntil: 7200,
				Macs:         []string{"11:22:33:44:55:66"},
				CustomDayRanges: []types.DayRange{},
			}
		})
		JustBeforeEach(func() {
			*returnedControl, *returnedErr = freeboxClient.UpdateNetworkControl(ctx, *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/network_control/42", version)),
						ghttp.VerifyContentType("application/json"),
						verifyAuth(*sessionToken),
						ghttp.VerifyJSON(`{
							"profile_id": 42,
							"override_mode": "denied",
							"override_until": 7200,
							"override": true,
							"macs": ["11:22:33:44:55:66"],
							"cdayranges": []
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"profile_id": 42,
								"next_change": 0,
								"override_mode": "denied",
								"current_mode": "denied",
								"rule_mode": "allowed",
								"override_until": 7200,
								"override": true,
								"macs": ["11:22:33:44:55:66"],
								"hosts": [],
								"resolution": 3600,
								"cdayranges": []
							}
						}`),
					),
				)
			})
			It("should return the updated network control", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedControl).To(MatchFields(IgnoreExtras, Fields{
					"ProfileID":     Equal(int64(42)),
					"OverrideMode":  Equal(types.RuleModeDenied),
					"CurrentMode":   Equal(types.RuleModeDenied),
					"Override":      Equal(true),
					"OverrideUntil": Equal(7200),
				}))
			})
		})
		Context("when the network control is not found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/network_control/42", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": false,
							"error_code": "noent",
							"msg": "Network control not found"
						}`),
					),
				)
			})
			It("should return ErrNetworkControlNotFound", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrNetworkControlNotFound))
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
