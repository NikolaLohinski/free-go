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
		DeferCleanup(server.Close)

		*endpoint = server.Addr()

		freeboxClient = Must(client.New(*endpoint, version)).
			WithAppID(appID).
			WithPrivateToken(privateToken)

		*sessionToken = setupLoginFlow(server)
	})
	Context("listing lan metadata", func() {
		returnedLanInfos := new([]types.LanInfo)
		JustBeforeEach(func() {
			*returnedLanInfos, *returnedErr = freeboxClient.ListLanInterfaceInfo(context.Background())
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/browser/interfaces/", version)),
						verifyAuth(*sessionToken),
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
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/browser/interfaces/", version)),
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
	Context("getting a lan interface", func() {
		const (
			interfaceName = "pub"
		)
		returnedLanInterfaceHosts := new([]types.LanInterfaceHost)
		JustBeforeEach(func() {
			*returnedLanInterfaceHosts, *returnedErr = freeboxClient.GetLanInterface(context.Background(), interfaceName)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/browser/%s", version, interfaceName)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": [{
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
							}]
						}`),
					),
				)
			})
			It("should return the correct lan interface", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedLanInterfaceHosts).To(Equal([]types.LanInterfaceHost{
					{
						Active:            false,
						Persistent:        false,
						Reachable:         false,
						PrimaryNameManual: false,
						VendorName:        "",
						Type:              "workstation",
						Interface:         "pub",
						ID:                "ether-7e:ec:37:cd:5b:6a",
						LastTimeReachable: types.Timestamp{
							Time: Must(time.Parse(time.RFC3339, "2023-04-27T09:05:32+02:00")).UTC(),
						},
						FirstActivity: types.Timestamp{
							Time: Must(time.Parse(time.RFC3339, "2023-04-27T08:58:44+02:00")).UTC(),
						},
						LastActivity: types.Timestamp{
							Time: Must(time.Parse(time.RFC3339, "2023-04-27T09:05:32+02:00")).UTC(),
						},
						PrimaryName: "testing",
						DefaultName: "testing",
						L2Ident: types.L2Ident{
							ID:   "7E:EC:37:CD:5B:6A",
							Type: "mac_address",
						},
						Names: []types.HostName{
							{Name: "test", Source: "dhcp"},
						},
						L3Connectivities: []types.LanHostL3Connectivity{
							{
								Address:   "192.168.1.254",
								Active:    false,
								Reachable: false,
								LastActivity: types.Timestamp{
									Time: Must(time.Parse(time.RFC3339, "2023-04-27T09:05:11+02:00")).UTC(),
								},
								LastTimeReachable: types.Timestamp{
									Time: Must(time.Parse(time.RFC3339, "2023-04-27T09:05:11+02:00")).UTC(),
								},
								Type: "ipv4",
							},
						},
					},
				},
				))
			})
		})
		Context("when the interface does not exist", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/browser/%s", version, interfaceName)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"msg": "Erreur lors de la récupération de la liste des hôtes : Interface invalide",
							"success": false,
							"error_code": "nodev"
						}`),
					),
				)
			})
			It("should return the correct error", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrInterfaceNotFound))
				Expect(client.ErrInterfaceNotFound.Error()).To(Equal("interface not found"))
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/browser/%s", version, interfaceName)),
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
	Context("getting a lan interface host", func() {
		const (
			interfaceName  = "pub"
			hostIdentifier = "ether-7e:ec:37:cd:5b:6a"
		)
		returnedLanInterfaceHost := new(types.LanInterfaceHost)
		JustBeforeEach(func() {
			*returnedLanInterfaceHost, *returnedErr = freeboxClient.GetLanInterfaceHost(context.Background(), interfaceName, hostIdentifier)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/browser/%s/%s", version, interfaceName, hostIdentifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
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
							}
						}`),
					),
				)
			})
			It("should return the correct lan interface", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedLanInterfaceHost).To(Equal(types.LanInterfaceHost{
					Active:            false,
					Persistent:        false,
					Reachable:         false,
					PrimaryNameManual: false,
					VendorName:        "",
					Type:              "workstation",
					Interface:         "pub",
					ID:                "ether-7e:ec:37:cd:5b:6a",
					LastTimeReachable: types.Timestamp{
						Time: Must(time.Parse(time.RFC3339, "2023-04-27T09:05:32+02:00")).UTC(),
					},
					FirstActivity: types.Timestamp{
						Time: Must(time.Parse(time.RFC3339, "2023-04-27T08:58:44+02:00")).UTC(),
					},
					LastActivity: types.Timestamp{
						Time: Must(time.Parse(time.RFC3339, "2023-04-27T09:05:32+02:00")).UTC(),
					},
					PrimaryName: "testing",
					DefaultName: "testing",
					L2Ident: types.L2Ident{
						ID:   "7E:EC:37:CD:5B:6A",
						Type: "mac_address",
					},
					Names: []types.HostName{
						{Name: "test", Source: "dhcp"},
					},
					L3Connectivities: []types.LanHostL3Connectivity{
						{
							Address:   "192.168.1.254",
							Active:    false,
							Reachable: false,
							LastActivity: types.Timestamp{
								Time: Must(time.Parse(time.RFC3339, "2023-04-27T09:05:11+02:00")).UTC(),
							},
							LastTimeReachable: types.Timestamp{
								Time: Must(time.Parse(time.RFC3339, "2023-04-27T09:05:11+02:00")).UTC(),
							},
							Type: "ipv4",
						},
					},
				},
				))
			})
		})
		Context("when the interface does not exist", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/browser/%s/%s", version, interfaceName, hostIdentifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"msg": "Erreur lors de la récupération de la liste des hôtes : Interface invalide",
							"success": false,
							"error_code": "nodev"
						}`),
					),
				)
			})
			It("should return the correct error", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrInterfaceNotFound))
				Expect((*returnedErr).Error()).To(Equal("interface not found"))
			})
		})
		Context("when the interface host does not exist", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/browser/%s/%s", version, interfaceName, hostIdentifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"msg": "Erreur lors de la récupération de la liste des hôtes : Pas d'hôte avec cet identifiant",
							"success": false,
							"error_code": "nohost"
						}`),
					),
				)
			})
			It("should return the correct error", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrInterfaceHostNotFound))
				Expect((*returnedErr).Error()).To(Equal("interface host not found"))
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/lan/browser/%s/%s", version, interfaceName, hostIdentifier)),
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
