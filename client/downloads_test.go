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
		DeferCleanup(server.Close)

		*endpoint = server.Addr()

		freeboxClient = Must(client.New(*endpoint, version)).
			WithAppID(appID).
			WithPrivateToken(privateToken)

		*sessionToken = setupLoginFlow(server)
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
	Context("add a download task", func() {
		var (
			ctx             = new(context.Context)
			downloadRequest = new(types.DownloadRequest)

			returnedID = new(int64)

			setupServer = func(handlers ...http.HandlerFunc) {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						append(handlers,
							ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/downloads/add", version)),
							ghttp.VerifyMimeType("application/x-www-form-urlencoded"),
							verifyAuth(*sessionToken),
							ghttp.RespondWith(http.StatusOK, `{
						    "success": true,
						    "result": {
						        "id": 1234
						    }
						}`),
						)...,
					),
				)
			}
		)
		BeforeEach(func() {
			*ctx = context.Background()
			*downloadRequest = types.DownloadRequest{
				DownloadURLs: []string{"http://localhost:8080/file"},
			}
		})

		JustBeforeEach(func() {
			*returnedID, *returnedErr = freeboxClient.AddDownloadTask(*ctx, *downloadRequest)
		})
		Context("when a single URL is passed", func() {
			BeforeEach(func() {
				setupServer(ghttp.VerifyFormKV("download_url", "http://localhost:8080/file"))
			})
			It("should return the correct task ID", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedID).To(Equal(int64(1234)))
			})
		})
		Context("when there are multiple URLs", func() {
			BeforeEach(func() {
				downloadRequest.DownloadURLs = append(downloadRequest.DownloadURLs, "ftp://localhost/other-file")
				setupServer(ghttp.VerifyFormKV("download_url_list", "http://localhost:8080/file\nftp://localhost/other-file"))
			})
			It("should return the correct task ID", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedID).To(Equal(int64(1234)))
			})
		})
		Context("when a directory is set", func() {
			BeforeEach(func() {
				downloadRequest.DownloadDirectory = "/Freebox/Téléchargements"
				setupServer(ghttp.VerifyFormKV("download_dir", "L0ZyZWVib3gvVMOpbMOpY2hhcmdlbWVudHM="))
			})
			It("should return the correct task ID", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedID).To(Equal(int64(1234)))
			})
		})
		Context("when a filename is set", func() {
			BeforeEach(func() {
				downloadRequest.Filename = "foo.bar"
				setupServer(ghttp.VerifyFormKV("filename", "foo.bar"))
			})
			It("should return the correct task ID", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedID).To(Equal(int64(1234)))
			})
			Context("when recursive is set", func() {
				BeforeEach(func() {
					downloadRequest.Recursive = true
				})
				It("should return an error", func() {
					Expect(*returnedErr).ToNot(BeNil())
				})
			})
			Context("when multiple URLs are passed", func() {
				BeforeEach(func() {
					downloadRequest.DownloadURLs = append(downloadRequest.DownloadURLs, "foo")
				})
				It("should return an error", func() {
					Expect(*returnedErr).ToNot(BeNil())
				})
			})
		})
		Context("when a hash is set", func() {
			BeforeEach(func() {
				downloadRequest.Hash = "sha256:f2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fd2"
				setupServer(ghttp.VerifyFormKV("hash", "sha256:f2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fd2"))
			})
			It("should return the correct task ID", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedID).To(Equal(int64(1234)))
			})
			Context("when recursive is set", func() {
				BeforeEach(func() {
					downloadRequest.Recursive = true
				})
				It("should return an error", func() {
					Expect(*returnedErr).ToNot(BeNil())
				})
			})
			Context("when multiple URLs are passed", func() {
				BeforeEach(func() {
					downloadRequest.DownloadURLs = append(downloadRequest.DownloadURLs, "foo")
				})
				It("should return an error", func() {
					Expect(*returnedErr).ToNot(BeNil())
				})
			})
		})
		Context("when recursive is set", func() {
			BeforeEach(func() {
				downloadRequest.Recursive = true
				setupServer(ghttp.VerifyFormKV("recursive", "true"))
			})
			It("should return the correct task ID", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedID).To(Equal(int64(1234)))
			})
		})
		Context("when username is set", func() {
			BeforeEach(func() {
				downloadRequest.Username = "foo"
				setupServer(ghttp.VerifyFormKV("username", "foo"))
			})
			It("should return the correct task ID", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedID).To(Equal(int64(1234)))
			})
		})
		Context("when password is set", func() {
			BeforeEach(func() {
				downloadRequest.Password = "bar"
				setupServer(ghttp.VerifyFormKV("password", "bar"))
			})
			It("should return the correct task ID", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedID).To(Equal(int64(1234)))
			})
		})
		Context("when archive_password is set", func() {
			BeforeEach(func() {
				downloadRequest.ArchivePassword = "zip"
				setupServer(ghttp.VerifyFormKV("archive_password", "zip"))
			})
			It("should return the correct task ID", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedID).To(Equal(int64(1234)))
			})
		})
		Context("when cookies is set", func() {
			BeforeEach(func() {
				downloadRequest.Cookies = map[string]string{
					"cookie1": "foo",
					"cookie2": "bar",
				}
				setupServer(ghttp.VerifyFormKV("cookies", "cookie1=foo; cookie2=bar"))
			})
			It("should return the correct task ID", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedID).To(Equal(int64(1234)))
			})
		})
		Context("when the context is nil", func() {
			BeforeEach(func() {
				*ctx = nil
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
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
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/downloads/add", version)),
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
	Context("delete a download task", func() {
		var (
			taskID = int64(1234)

			setupServer = func(handlers ...http.HandlerFunc) {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						append(handlers,
							ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/downloads/%d", version, taskID)),
							verifyAuth(*sessionToken),
							ghttp.RespondWith(http.StatusOK, `{"success": true}`),
						)...,
					),
				)
			}
		)
		BeforeEach(func() {
			setupServer()
		})

		It("Should work", func(ctx SpecContext) {
			Expect(freeboxClient.DeleteDownloadTask(ctx, taskID)).To(Succeed())
		})
	})
	Context("erase a download task", func() {
		var (
			taskID = int64(1234)

			setupServer = func(handlers ...http.HandlerFunc) {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						append(handlers,
							ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/downloads/%d/erase", version, taskID)),
							verifyAuth(*sessionToken),
							ghttp.RespondWith(http.StatusOK, `{"success": true}`),
						)...,
					),
				)
			}
		)
		BeforeEach(func() {
			setupServer()
		})

		It("Should work", func(ctx SpecContext) {
			Expect(freeboxClient.EraseDownloadTask(ctx, taskID)).To(Succeed())
		})
	})
	Context("update a download task", func() {
		var (
			payload = new(types.DownloadTaskUpdate)

			taskID = int64(1234)

			setupServer = func(handlers ...http.HandlerFunc) {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						append(handlers,
							ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/downloads/%d", version, taskID)),
							ghttp.VerifyMimeType("application/json"),
							verifyAuth(*sessionToken),
							ghttp.RespondWith(http.StatusOK, `{
						    "success": true,
						    "result": {
						        "id": 1234
						    }
						}`),
						)...,
					),
				)
			}
		)
		BeforeEach(func() {
			*payload = types.DownloadTaskUpdate{}
		})

		Context("when fields are not set", func() {
			BeforeEach(func() {
				setupServer(ghttp.VerifyJSON("{}"))
			})
			It("Should omit the fields", func(ctx SpecContext) {
				Expect(freeboxClient.UpdateDownloadTask(ctx, taskID, *payload)).To(Succeed())
			})
		})
		Context("when status is set", func() {
			BeforeEach(func() {
				payload.Status = types.DownloadTaskStatusStopped
				setupServer(ghttp.VerifyJSON(`{"status":"stopped"}`))
			})
			It("should work", func(ctx SpecContext) {
				Expect(freeboxClient.UpdateDownloadTask(ctx, taskID, *payload)).To(Succeed())
			})
		})
		Context("when priority is low", func() {
			BeforeEach(func() {
				payload.IOPriority = types.DownloadTaskIOPriorityLow
				setupServer(ghttp.VerifyJSON(`{"io_priority":"low"}`))
			})
			It("should work", func(ctx SpecContext) {
				Expect(freeboxClient.UpdateDownloadTask(ctx, taskID, *payload)).To(Succeed())
			})
		})
		Context("when priority is normal", func() {
			BeforeEach(func() {
				payload.IOPriority = types.DownloadTaskIOPriorityNormal
				setupServer(ghttp.VerifyJSON(`{"io_priority":"normal"}`))
			})
			It("should work", func(ctx SpecContext) {
				Expect(freeboxClient.UpdateDownloadTask(ctx, taskID, *payload)).To(Succeed())
			})
		})
	})
})
