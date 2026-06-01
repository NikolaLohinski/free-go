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

var _ = Describe("downloads config", func() {
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

	Context("getting the download configuration", func() {
		returnedConfig := new(types.DownloadConfiguration)
		JustBeforeEach(func() {
			*returnedConfig, *returnedErr = freeboxClient.GetDownloadConfiguration(context.Background())
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/downloads/config/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"max_downloading_tasks": 5,
								"download_dir": "L0Rpc3F1ZSBkdXIvVMOpbMOpY2hhcmdlbWVudHMv",
								"watch_dir": "L0Rpc3F1ZSBkdXIvLnF1ZXVl",
								"use_watch_dir": true,
								"dns1": "",
								"dns2": "",
								"throttling": {
									"normal": {"rx_rate": 0, "tx_rate": 0},
									"slow": {"rx_rate": 512, "tx_rate": 42},
									"schedule": ["slow", "normal", "normal"],
									"mode": "normal"
								},
								"news": {
									"server": "news.free.fr",
									"port": 119,
									"ssl": false,
									"user": "",
									"nthreads": 1,
									"auto_repair": true,
									"lazy_par2": true,
									"auto_extract": true,
									"erase_tmp": true
								},
								"bt": {
									"max_peers": 50,
									"stop_ratio": 150,
									"crypto_support": "allowed",
									"enable_dht": false,
									"enable_pex": false,
									"announce_timeout": 0,
									"main_port": 0,
									"dht_port": 0
								},
								"feed": {
									"fetch_interval": 60,
									"max_items": 0
								},
								"blocklist": {
									"sources": []
								}
							}
						}`),
					),
				)
			})
			It("should return the correct download configuration", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedConfig).To(Equal(types.DownloadConfiguration{
					MaxDownloadingTasks: 5,
					DownloadDir:         "/Disque dur/Téléchargements/",
					WatchDir:            "/Disque dur/.queue",
					UseWatchDir:         true,
					DNS1:                "",
					DNS2:                "",
					Throttling: types.DlThrottlingConfig{
						Normal:   types.DlRate{RxRate: 0, TxRate: 0},
						Slow:     types.DlRate{RxRate: 512, TxRate: 42},
						Schedule: []types.DlThrottlingMode{"slow", "normal", "normal"},
						Mode:     types.DlThrottlingModeNormal,
					},
					News: types.DlNewsConfig{
						Server:      "news.free.fr",
						Port:        119,
						SSL:         false,
						User:        "",
						NThreads:    1,
						AutoRepair:  true,
						LazyPar2:    true,
						AutoExtract: true,
						EraseTmp:    true,
					},
					Bt: types.DlBtConfig{
						MaxPeers:      50,
						StopRatio:     150,
						CryptoSupport: types.DlBtCryptoSupportAllowed,
					},
					Feed: types.DlFeedConfig{
						FetchInterval: 60,
						MaxItems:      0,
					},
					BlockList: types.DlBlockListConfig{
						Sources: []interface{}{},
					},
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/downloads/config/", version)),
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

	Context("updating the download configuration", func() {
		var (
			payload = new(types.DownloadConfiguration)
		)
		returnedConfig := new(types.DownloadConfiguration)
		BeforeEach(func() {
			*payload = types.DownloadConfiguration{
				MaxDownloadingTasks: 3,
				DownloadDir:         "/Disque dur/Téléchargements/",
				UseWatchDir:         false,
			}
		})
		JustBeforeEach(func() {
			*returnedConfig, *returnedErr = freeboxClient.UpdateDownloadConfiguration(context.Background(), *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/downloads/config/", version)),
						ghttp.VerifyMimeType("application/json"),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"max_downloading_tasks": 3,
								"download_dir": "L0Rpc3F1ZSBkdXIvVMOpbMOpY2hhcmdlbWVudHMv",
								"watch_dir": "",
								"use_watch_dir": false,
								"dns1": "",
								"dns2": "",
								"throttling": {
									"normal": {"rx_rate": 0, "tx_rate": 0},
									"slow": {"rx_rate": 0, "tx_rate": 0},
									"schedule": [],
									"mode": "normal"
								},
								"news": {
									"server": "",
									"port": 0,
									"ssl": false,
									"user": "",
									"nthreads": 0,
									"auto_repair": false,
									"lazy_par2": false,
									"auto_extract": false,
									"erase_tmp": false
								},
								"bt": {
									"max_peers": 0,
									"stop_ratio": 0,
									"crypto_support": ""
								},
								"feed": {
									"fetch_interval": 0,
									"max_items": 0
								},
								"blocklist": {
									"sources": null
								}
							}
						}`),
					),
				)
			})
			It("should return the updated download configuration", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(returnedConfig.MaxDownloadingTasks).To(Equal(3))
				Expect(string(returnedConfig.DownloadDir)).To(Equal("/Disque dur/Téléchargements/"))
				Expect(returnedConfig.UseWatchDir).To(BeFalse())
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
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/downloads/config/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": false,
							"error_code": "auth_required",
							"msg": "not logged in"
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
