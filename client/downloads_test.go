package client_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("downloads", func() {
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

		freeboxClient = Must(client.New(*endpoint, version)).
			WithAppID(appID).
			WithPrivateToken(privateToken)

		*sessionToken = setupLoginFlow(server)
	})
	AfterEach(func() {
		server.Close()
	})
	Context("listing download tasks", func() {
		returnedDownloadTasks := new([]types.DownloadTask)
		JustBeforeEach(func() {
			*returnedDownloadTasks, *returnedErr = freeboxClient.ListDownloadTasks(context.Background())
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/downloads/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": [
								{
									"rx_bytes": 621340000,
									"tx_bytes": 0,
									"download_dir": "RnJlZWJveC9WTXM=",
									"archive_password": "",
									"eta": 0,
									"status": "done",
									"io_priority": "normal",
									"type": "http",
									"piece_length": 0,
									"queue_pos": 1,
									"id": 7,
									"info_hash": "",
									"created_ts": 1711656593,
									"stop_ratio": 0,
									"tx_rate": 0,
									"name": "freebox-vm-21141.qcow2",
									"tx_pct": 10000,
									"rx_pct": 10000,
									"rx_rate": 0,
									"error": "none",
									"size": 621340000
								},
								{
									"rx_bytes": 185610000,
									"tx_bytes": 0,
									"download_dir": "RnJlZWJveC9WTXM=",
									"archive_password": "",
									"eta": 919,
									"status": "downloading",
									"io_priority": "normal",
									"type": "http",
									"piece_length": 0,
									"queue_pos": 2,
									"id": 8,
									"info_hash": "",
									"created_ts": 1711656773,
									"stop_ratio": 0,
									"tx_rate": 0,
									"name": "ubuntu-22.04-server-cloudimg-arm64.img",
									"tx_pct": 10000,
									"rx_pct": 2987,
									"rx_rate": 473800,
									"error": "none",
									"size": 621340000
								}
							]
						}`),
					),
				)
			})
			It("should return the correct download tasks", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedDownloadTasks).To(Equal([]types.DownloadTask{
					{
						ID:                 7,
						Name:               "freebox-vm-21141.qcow2",
						Type:               types.DownloadTaskTypeHTTP,
						Status:             types.DownloadTaskStatusDone,
						IOPriority:         types.DownloadTaskIOPriorityNormal,
						SizeBytes:          621340000,
						QueuePosition:      1,
						TransmittedBytes:   0,
						ReceivedBytes:      621340000,
						TransmitRate:       0,
						ReceiveRate:        0,
						TransmitPercentage: 10000,
						ReceivedPercentage: 10000,
						Error:              "none",
						CreatedTimestamp:   types.Timestamp{Time: time.Unix(1711656593, 0).UTC()},
						ETASeconds:         0,
						DownloadDirectory:  "Freebox/VMs",
						StopRatio:          0,
						ArchivePassword:    "",
						InfoHash:           "",
						PieceLength:        0,
					},
					{
						ID:                 8,
						Name:               "ubuntu-22.04-server-cloudimg-arm64.img",
						Type:               types.DownloadTaskTypeHTTP,
						Status:             types.DownloadTaskStatusDownloading,
						IOPriority:         types.DownloadTaskIOPriorityNormal,
						SizeBytes:          621340000,
						QueuePosition:      2,
						TransmittedBytes:   0,
						ReceivedBytes:      185610000,
						TransmitRate:       0,
						ReceiveRate:        473800,
						TransmitPercentage: 10000,
						ReceivedPercentage: 2987,
						Error:              "none",
						CreatedTimestamp:   types.Timestamp{Time: time.Unix(1711656773, 0).UTC()},
						ETASeconds:         919,
						DownloadDirectory:  "Freebox/VMs",
						StopRatio:          0,
						ArchivePassword:    "",
						InfoHash:           "",
						PieceLength:        0,
					},
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
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/downloads/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": [
								"foo"
							]
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("get a download task", func() {
		const identifier = int64(42)
		returnedDownloadTask := new(types.DownloadTask)
		JustBeforeEach(func() {
			*returnedDownloadTask, *returnedErr = freeboxClient.GetDownloadTask(context.Background(), identifier)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/downloads/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
						    "success": true,
						    "result": {
						        "rx_bytes": 184300000,
						        "tx_bytes": 0,
						        "download_dir": "RnJlZWJveC9WTXM=",
						        "archive_password": "",
						        "eta": 906,
						        "status": "downloading",
						        "io_priority": "normal",
						        "type": "http",
						        "piece_length": 0,
						        "queue_pos": 2,
						        "id": 8,
						        "info_hash": "",
						        "created_ts": 1711656773,
						        "stop_ratio": 0,
						        "tx_rate": 0,
						        "name": "ubuntu-22.04-server-cloudimg-arm64.img",
						        "tx_pct": 10000,
						        "rx_pct": 2966,
						        "rx_rate": 481950,
						        "error": "none",
						        "size": 621340000
						    }
						}`),
					),
				)
			})
			It("should return the correct download tasks", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedDownloadTask).To(Equal(types.DownloadTask{
					ID:                 8,
					Type:               types.DownloadTaskTypeHTTP,
					Name:               "ubuntu-22.04-server-cloudimg-arm64.img",
					Status:             types.DownloadTaskStatusDownloading,
					IOPriority:         types.DownloadTaskIOPriorityNormal,
					SizeBytes:          621340000,
					QueuePosition:      2,
					TransmittedBytes:   0,
					ReceivedBytes:      184300000,
					TransmitRate:       0,
					ReceiveRate:        481950,
					TransmitPercentage: 10000,
					ReceivedPercentage: 2966,
					Error:              "none",
					CreatedTimestamp:   types.Timestamp{Time: time.Unix(1711656773, 0).UTC()},
					ETASeconds:         906,
					DownloadDirectory:  "Freebox/VMs",
					StopRatio:          0,
					ArchivePassword:    "",
					InfoHash:           "",
					PieceLength:        0,
				}))
			})
		})
		Context("when the task is not found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/downloads/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
						    "msg": "Erreur lors de la récupération des infos de téléchargement : Pas de tâche avec cet id",
						    "success": false,
						    "error_code": "task_not_found"
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrTaskNotFound))
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/downloads/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": [
								"foo"
							]
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
