package client_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/onsi/gomega/gstruct"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("DHCPStatic", func() {
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

	Context("ListDHCPStaticLease", func() {
		const (
			IPAddress  = "192.168.1.10"
			MACAddress = "00:11:22:33:44:55"
		)
		returnedDHCPStatjcLeases := new([]types.DHCPStaticLeaseInfo)

		JustBeforeEach(func(ctx SpecContext) {
			*returnedDHCPStatjcLeases, *returnedErr = freeboxClient.ListDHCPStaticLease(ctx)
		})

		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", fmt.Sprintf("/api/%s/dhcp/static_lease/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWithJSONEncoded(http.StatusOK, struct {
							Success bool                        `json:"success"`
							Result  []types.DHCPStaticLeaseInfo `json:"result"`
						}{
							Success: true,
							Result: []types.DHCPStaticLeaseInfo{
								{
									ID:       "identifier",
									IP:       IPAddress,
									Mac:      MACAddress,
									Comment:  CurrentSpecReport().FullText(),
									Hostname: "test-hostname",
									Host: types.LanInterfaceHost{
										L2Ident: types.L2Ident{
											ID:   MACAddress,
											Type: types.MacAddress,
										},
										Active:     true,
										Persistent: true,
										Names: []types.HostName{
											{
												Name:   "test",
												Source: "dhcp",
											},
										},
										VendorName:        "vendor",
										Type:              "workstation",
										Interface:         "pub",
										ID:                "ether-7e:ec:37:cd:5b:6a",
										LastTimeReachable: types.Timestamp{Time: time.Unix(1682578724, 0)},
										PrimaryNameManual: true,
										L3Connectivities: []types.LanHostL3Connectivity{
											{
												Address:           IPAddress,
												Active:            true,
												Reachable:         true,
												LastActivity:      types.Timestamp{Time: time.Unix(1682578724, 0)},
												Type:              "ipv4",
												LastTimeReachable: types.Timestamp{Time: time.Unix(1682578724, 0)},
											},
										},
										DefaultName:   "testing",
										FirstActivity: types.Timestamp{Time: time.Unix(1682578724, 0)},
										Reachable:     true,
										LastActivity:  types.Timestamp{Time: time.Unix(1682578724, 0)},
										PrimaryName:   "testing",
									},
								},
							},
						}),
					),
				)
			})

			It("should return the list of DHCP static leases", func() {
				Expect(*returnedErr).ToNot(HaveOccurred())
				Expect(*returnedDHCPStatjcLeases).To(ConsistOf(
					gstruct.MatchAllFields(gstruct.Fields{
						"ID":       Equal("identifier"),
						"IP":       Equal(IPAddress),
						"Mac":      Equal(MACAddress),
						"Comment":  Equal(CurrentSpecReport().FullText()),
						"Hostname": Equal("test-hostname"),
						"Host": gstruct.MatchAllFields(gstruct.Fields{
							"L2Ident": gstruct.MatchAllFields(gstruct.Fields{
								"ID":   Equal(MACAddress),
								"Type": Equal(types.MacAddress),
							}),
							"Active":     BeTrue(),
							"Persistent": BeTrue(),
							"Names": ConsistOf(gstruct.MatchAllFields(gstruct.Fields{
								"Name":   Equal("test"),
								"Source": Equal("dhcp"),
							})),
							"VendorName": Equal("vendor"),
							"Type":       Equal("workstation"),
							"Interface":  Equal("pub"),
							"ID":         Equal("ether-7e:ec:37:cd:5b:6a"),
							"LastTimeReachable": gstruct.MatchAllFields(gstruct.Fields{
								"Time": BeTemporally("==", time.Unix(1682578724, 0)),
							}),
							"PrimaryNameManual": BeTrue(),
							"L3Connectivities": ConsistOf(gstruct.MatchAllFields(gstruct.Fields{
								"Address":   Equal(IPAddress),
								"Active":    BeTrue(),
								"Reachable": BeTrue(),
								"LastActivity": gstruct.MatchAllFields(gstruct.Fields{
									"Time": BeTemporally("==", time.Unix(1682578724, 0)),
								}),
								"Type": Equal("ipv4"),
								"LastTimeReachable": gstruct.MatchAllFields(gstruct.Fields{
									"Time": BeTemporally("==", time.Unix(1682578724, 0)),
								}),
								"Model": BeEmpty(),
							})),
							"DefaultName": Equal("testing"),
							"FirstActivity": gstruct.MatchAllFields(gstruct.Fields{
								"Time": BeTemporally("==", time.Unix(1682578724, 0)),
							}),
							"Reachable": BeTrue(),
							"LastActivity": gstruct.MatchAllFields(gstruct.Fields{
								"Time": BeTemporally("==", time.Unix(1682578724, 0)),
							}),
							"PrimaryName":    Equal("testing"),
							"NetworkControl": BeNil(),
							"Model":          BeEmpty(),
							"AccessPoint":    BeNil(),
						}),
					})),
				)
			})
		})

		Context("when the server returns an error", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", fmt.Sprintf("/api/%s/dhcp/static_lease/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWithJSONEncoded(http.StatusOK, struct {
							Success bool `json:"success"`
						}{
							Success: false,
						}),
					),
				)
			})

			It("should return an error", func() {
				Expect(*returnedErr).To(HaveOccurred())
			})
		})

		Context("when the server returns an empty result", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", fmt.Sprintf("/api/%s/dhcp/static_lease/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWithJSONEncoded(http.StatusOK, struct {
							Success bool `json:"success"`
						}{
							Success: true,
						}),
					),
				)
			})

			It("should return an empty list", func() {
				Expect(*returnedErr).ToNot(HaveOccurred())
				Expect(*returnedDHCPStatjcLeases).To(BeEmpty())
			})
		})
	})

	Context("GetDHCPStaticLease", func() {
		const (
			IPAddress  = "192.168.1.11"
			MACAddress = "00:11:22:33:44:56"
		)

		returnedDHCPStaticLease := new(types.DHCPStaticLeaseInfo)

		JustBeforeEach(func(ctx SpecContext) {
			*returnedDHCPStaticLease, *returnedErr = freeboxClient.GetDHCPStaticLease(ctx, MACAddress)
		})

		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", fmt.Sprintf("/api/%s/dhcp/static_lease/%s", version, MACAddress)),
						verifyAuth(*sessionToken),
						ghttp.RespondWithJSONEncoded(http.StatusOK, struct {
							Success bool                      `json:"success"`
							Result  types.DHCPStaticLeaseInfo `json:"result"`
						}{
							Success: true,
							Result: types.DHCPStaticLeaseInfo{
								ID:       "identifier",
								IP:       IPAddress,
								Mac:      MACAddress,
								Comment:  CurrentSpecReport().FullText(),
								Hostname: "test-hostname",
								Host: types.LanInterfaceHost{
									L2Ident: types.L2Ident{
										ID:   MACAddress,
										Type: types.MacAddress,
									},
									Active:     true,
									Persistent: true,
									Names: []types.HostName{
										{
											Name:   "test",
											Source: "dhcp",
										},
									},
									VendorName:        "vendor",
									Type:              "workstation",
									Reachable:         true,
									PrimaryNameManual: true,
									Interface:         "pub",
									ID:                "ether-7e:ec:37:cd:5b:6a",
									LastTimeReachable: types.Timestamp{Time: time.Unix(1682578724, 0)},
									FirstActivity:     types.Timestamp{Time: time.Unix(1682578724, 0)},
									LastActivity:      types.Timestamp{Time: time.Unix(1682578724, 0)},
									PrimaryName:       "testing",
									DefaultName:       "testing",
									L3Connectivities: []types.LanHostL3Connectivity{
										{
											Address:           IPAddress,
											Active:            true,
											Reachable:         true,
											LastActivity:      types.Timestamp{Time: time.Unix(1682578724, 0)},
											Type:              "ipv4",
											LastTimeReachable: types.Timestamp{Time: time.Unix(1682578724, 0)},
										},
									},
								},
							},
						}),
					),
				)
			})

			It("should return the DHCP static lease", func() {
				Expect(*returnedErr).ToNot(HaveOccurred())
				Expect(*returnedDHCPStaticLease).To(gstruct.MatchAllFields(gstruct.Fields{
					"ID":       Equal("identifier"),
					"IP":       Equal(IPAddress),
					"Mac":      Equal(MACAddress),
					"Comment":  Equal(CurrentSpecReport().FullText()),
					"Hostname": Equal("test-hostname"),
					"Host": gstruct.MatchAllFields(gstruct.Fields{
						"L2Ident": gstruct.MatchAllFields(gstruct.Fields{
							"ID":   Equal(MACAddress),
							"Type": Equal(types.MacAddress),
						}),
						"Active":     BeTrue(),
						"Persistent": BeTrue(),
						"Names": ConsistOf(gstruct.MatchAllFields(gstruct.Fields{
							"Name":   Equal("test"),
							"Source": Equal("dhcp"),
						})),
						"VendorName":        Equal("vendor"),
						"Type":              Equal("workstation"),
						"Reachable":         BeTrue(),
						"PrimaryNameManual": BeTrue(),
						"Interface":         Equal("pub"),
						"ID":                Equal("ether-7e:ec:37:cd:5b:6a"),
						"LastTimeReachable": gstruct.MatchAllFields(gstruct.Fields{
							"Time": BeTemporally("==", time.Unix(1682578724, 0)),
						}),
						"FirstActivity": gstruct.MatchAllFields(gstruct.Fields{
							"Time": BeTemporally("==", time.Unix(1682578724, 0)),
						}),
						"LastActivity": gstruct.MatchAllFields(gstruct.Fields{
							"Time": BeTemporally("==", time.Unix(1682578724, 0)),
						}),
						"PrimaryName": Equal("testing"),
						"DefaultName": Equal("testing"),
						"L3Connectivities": ConsistOf(gstruct.MatchAllFields(gstruct.Fields{
							"Address":   Equal(IPAddress),
							"Active":    BeTrue(),
							"Reachable": BeTrue(),
							"LastActivity": gstruct.MatchAllFields(gstruct.Fields{
								"Time": BeTemporally("==", time.Unix(1682578724, 0)),
							}),
							"Type": Equal("ipv4"),
							"LastTimeReachable": gstruct.MatchAllFields(gstruct.Fields{
								"Time": BeTemporally("==", time.Unix(1682578724, 0)),
							}),
							"Model": BeEmpty(),
						})),
						"Model":          BeEmpty(),
						"NetworkControl": BeNil(),
						"AccessPoint":    BeNil(),
					}),
				}))
			})
		})

		Context("when the server returns an error", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", fmt.Sprintf("/api/%s/dhcp/static_lease/%s", version, MACAddress)),
						verifyAuth(*sessionToken),
						ghttp.RespondWithJSONEncoded(http.StatusOK, struct {
							Success bool `json:"success"`
						}{
							Success: false,
						}),
					),
				)
			})

			It("should return an error", func() {
				Expect(*returnedErr).To(HaveOccurred())
			})
		})
	})

	Context("UpdateDHCPStaticLease", func() {
		const (
			IPAddress  = "192.168.1.12"
			MACAddress = "00:11:22:33:44:57"
		)

		returnedDHCPStaticLease := new(types.LanInterfaceHost)

		JustBeforeEach(func(ctx SpecContext) {
			*returnedDHCPStaticLease, *returnedErr = freeboxClient.UpdateDHCPStaticLease(ctx, MACAddress, types.DHCPStaticLeasePayload{
				Comment: CurrentSpecReport().FullText(),
			})
		})

		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PUT", fmt.Sprintf("/api/%s/dhcp/static_lease/%s", version, MACAddress)),
						verifyAuth(*sessionToken),
						ghttp.VerifyJSONRepresenting(map[string]interface{}{
							"comment": CurrentSpecReport().FullText(),
						}),
						ghttp.RespondWithJSONEncoded(http.StatusOK, struct {
							Success bool                   `json:"success"`
							Result  types.LanInterfaceHost `json:"result"`
						}{
							Success: true,
							Result: types.LanInterfaceHost{
								Active:            true,
								Persistent:        true,
								Reachable:         true,
								PrimaryNameManual: true,
								VendorName:        "vendor",
								Type:              "workstation",
								Interface:         "pub",
								ID:                "ether-7e:ec:37:cd:5b:6a",
								LastTimeReachable: types.Timestamp{Time: time.Unix(1682578724, 0)},
								FirstActivity:     types.Timestamp{Time: time.Unix(1682578724, 0)},
								LastActivity:      types.Timestamp{Time: time.Unix(1682578724, 0)},
								PrimaryName:       "testing",
								DefaultName:       "testing",
								L2Ident: types.L2Ident{
									ID:   MACAddress,
									Type: types.MacAddress,
								},
								Names: []types.HostName{
									{
										Name:   "test",
										Source: "dhcp",
									},
								},
								L3Connectivities: []types.LanHostL3Connectivity{
									{
										Address:           IPAddress,
										Active:            true,
										Reachable:         true,
										LastActivity:      types.Timestamp{Time: time.Unix(1682578724, 0)},
										Type:              "ipv4",
										LastTimeReachable: types.Timestamp{Time: time.Unix(1682578724, 0)},
									},
								},
							},
						}),
					),
				)
			})
		})

		Context("when the server returns an error", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PUT", fmt.Sprintf("/api/%s/dhcp/static_lease/%s", version, MACAddress)),
						verifyAuth(*sessionToken),
						ghttp.RespondWithJSONEncoded(http.StatusOK, struct {
							Success bool `json:"success"`
						}{
							Success: false,
						}),
					),
				)
			})

			It("should return an error", func() {
				Expect(*returnedErr).To(HaveOccurred())
			})
		})
	})

	Context("CreateDHCPStaticLease", func() {
		const (
			IPAddress  = "192.168.1.13"
			MACAddress = "00:11:22:33:44:58"
		)

		returnedDHCPStaticLease := new(types.LanInterfaceHost)

		JustBeforeEach(func(ctx SpecContext) {
			*returnedDHCPStaticLease, *returnedErr = freeboxClient.CreateDHCPStaticLease(ctx, types.DHCPStaticLeasePayload{
				Comment:  CurrentSpecReport().FullText(),
				Mac:      MACAddress,
				Hostname: "test-hostname",
				IP:       IPAddress,
			})
		})

		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", fmt.Sprintf("/api/%s/dhcp/static_lease/", version)),
						verifyAuth(*sessionToken),
						ghttp.VerifyJSONRepresenting(map[string]interface{}{
							"comment":  CurrentSpecReport().FullText(),
							"mac":      MACAddress,
							"hostname": "test-hostname",
							"ip":       IPAddress,
						}),
						ghttp.RespondWithJSONEncoded(http.StatusOK, struct {
							Success bool                   `json:"success"`
							Result  types.LanInterfaceHost `json:"result"`
						}{
							Success: true,
							Result: types.LanInterfaceHost{
								Active:            true,
								Persistent:        true,
								Reachable:         true,
								PrimaryNameManual: true,
								VendorName:        "vendor",
								Type:              "workstation",
								Interface:         "pub",
								ID:                "ether-7e:ec:37:cd:5b:6a",
								LastTimeReachable: types.Timestamp{Time: time.Unix(1682578724, 0)},
								FirstActivity:     types.Timestamp{Time: time.Unix(1682578724, 0)},
								LastActivity:      types.Timestamp{Time: time.Unix(1682578724, 0)},
								PrimaryName:       "testing",
								DefaultName:       "testing",
								L2Ident: types.L2Ident{
									ID:   MACAddress,
									Type: types.MacAddress,
								},
								Names: []types.HostName{
									{
										Name:   "test",
										Source: "dhcp",
									},
								},
								L3Connectivities: []types.LanHostL3Connectivity{
									{
										Address:           IPAddress,
										Active:            true,
										Reachable:         true,
										LastActivity:      types.Timestamp{Time: time.Unix(1682578724, 0)},
										Type:              "ipv4",
										LastTimeReachable: types.Timestamp{Time: time.Unix(1682578724, 0)},
									},
								},
							},
						}),
					),
				)
			})
		})

		Context("when the server returns an error", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", fmt.Sprintf("/api/%s/dhcp/static_lease/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWithJSONEncoded(http.StatusOK, struct {
							Success bool `json:"success"`
						}{
							Success: false,
						}),
					),
				)
			})

			It("should return an error", func() {
				Expect(*returnedErr).To(HaveOccurred())
			})
		})
	})
})
