package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("filesystem", func() {
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
	Context("getting file info", func() {
		const filePath = "path/to/file"
		const filePathBase64 = "cGF0aC90by9maWxl"
		returnedFileInfo := new(types.FileInfo)
		JustBeforeEach(func() {
			*returnedFileInfo, *returnedErr = freeboxClient.GetFileInfo(context.Background(), filePath)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fs/info/%s", version, filePathBase64)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
    						    "type": "file",
    						    "index": 0,
    						    "link": false,
    						    "parent": "RnJlZWJveC9WTXM=",
    						    "modification": 1711657506,
    						    "hidden": false,
    						    "mimetype": "application/x-raw-disk-image",
    						    "name": "file",
    						    "path": "cGF0aC90by9maWxl",
    						    "size": 353773718
    						}
						}`),
					),
				)
			})
			It("should return the correct file info", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedFileInfo).To(Equal(types.FileInfo{
					Type:         "file",
					Index:        0,
					Link:         false,
					Parent:       "Freebox/VMs",
					Modification: 1711657506,
					Hidden:       false,
					MimeType:     "application/x-raw-disk-image",
					Name:         "file",
					Path:         "path/to/file",
					SizeBytes:    353773718,
				}))
			})
		})
		Context("when the file does not exist", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fs/info/%s", version, filePathBase64)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"msg": "Erreur lors de la récupération de la liste des fichiers : Le fichier n'existe pas",
							"success": false,
							"error_code": "path_not_found"
						}`),
					),
				)
			})
			It("should return the correct error", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrPathNotFound))
				Expect(client.ErrPathNotFound.Error()).To(Equal("path not found"))
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fs/info/%s", version, filePathBase64)),
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
	Context("removing files", func() {
		var (
			filePaths       = []string{"path/to/file"}
			filePathsBase64 = []string{"cGF0aC90by9maWxl"}
			returnedTask    = new(types.FileSystemTask)
		)
		JustBeforeEach(func() {
			*returnedTask, *returnedErr = freeboxClient.RemoveFiles(context.Background(), filePaths)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/rm/", version)),
						ghttp.VerifyJSON(string(Must(json.Marshal(map[string][]string{
							"files": filePathsBase64,
						})))),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"curr_bytes_done": 0,
								"total_bytes": 0,
								"nfiles_done": 0,
								"started_ts": 1711657904,
								"duration": 0,
								"done_ts": 0,
								"src": [
									"path/to/file"
								],
								"curr_bytes": 0,
								"type": "rm",
								"to": "",
								"id": 15,
								"nfiles": 0,
								"created_ts": 1711657904,
								"state": "running",
								"total_bytes_done": 0,
								"rate": 0,
								"from": "path/to/file",
								"dst": "",
								"eta": 0,
								"error": "none",
								"progress": 0
							}
						}`),
					),
				)
			})
			It("should return the task", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedTask).To(Equal(types.FileSystemTask{
					ID:                            15,
					Type:                          types.FileTaskTypeRemove,
					State:                         types.FileTaskStateRunning,
					Error:                         types.FileTaskErrorNone,
					CurrentBytesDone:              0,
					TotalBytes:                    0,
					NumberFilesDone:               0,
					StartedTimestamp:              1711657904,
					DurationSeconds:               0,
					DoneTimestamp:                 0,
					CurrentBytes:                  0,
					To:                            "",
					NumberFiles:                   0,
					CreatedTimestamp:              1711657904,
					TotalBytesDone:                0,
					From:                          "path/to/file",
					ProcessingRate:                0,
					EstimatedTimeRemainingSeconds: 0,
					ProgressPercent:               0,
					Sources:                       []string{"path/to/file"},
					Destination:                   "",
				},
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
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/rm/", version)),
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
	Context("updating filesystem task", func() {
		const identifier int64 = 42
		returnedTask := new(types.FileSystemTask)
		JustBeforeEach(func(ctx SpecContext) {
			*returnedTask, *returnedErr = freeboxClient.UpdateFileSystemTask(ctx, identifier, types.FileSytemTaskUpdate{
				State: types.FileTaskStatePaused,
			})
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/fs/tasks/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"curr_bytes_done": 0,
								"total_bytes": 0,
								"nfiles_done": 1,
								"started_ts": 1712862825,
								"duration": 0,
								"done_ts": 1712862825,
								"src": [
									"Freebox/Other/free-go.svg"
								],
								"curr_bytes": 0,
								"type": "rm",
								"to": "",
								"id": 6,
								"nfiles": 1,
								"created_ts": 1712862825,
								"state": "done",
								"total_bytes_done": 0,
								"rate": 0,
								"from": "Freebox/Other/free-go.svg",
								"dst": "",
								"eta": 0,
								"error": "none",
								"progress": 100
							}
						}`),
					),
				)
			})
			It("should return the correct file info", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedTask).To(Equal(types.FileSystemTask{
					ID:                            6,
					Type:                          types.FileTaskTypeRemove,
					State:                         types.FileTaskStateDone,
					Error:                         types.FileTaskErrorNone,
					CurrentBytesDone:              0,
					TotalBytes:                    0,
					NumberFilesDone:               1,
					StartedTimestamp:              1712862825,
					DurationSeconds:               0,
					DoneTimestamp:                 1712862825,
					CurrentBytes:                  0,
					To:                            "",
					NumberFiles:                   1,
					CreatedTimestamp:              1712862825,
					TotalBytesDone:                0,
					From:                          "Freebox/Other/free-go.svg",
					ProcessingRate:                0,
					EstimatedTimeRemainingSeconds: 0,
					ProgressPercent:               100,
					Sources: []string{
						"Freebox/Other/free-go.svg",
					},
					Destination: "",
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
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/fs/tasks/%d", version, identifier)),
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
	Context("getting filesystem task", func() {
		const identifier int64 = 42
		returnedTask := new(types.FileSystemTask)
		JustBeforeEach(func() {
			*returnedTask, *returnedErr = freeboxClient.GetFileSystemTask(context.Background(), identifier)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fs/tasks/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"curr_bytes_done": 0,
								"total_bytes": 0,
								"nfiles_done": 1,
								"started_ts": 1712862825,
								"duration": 0,
								"done_ts": 1712862825,
								"src": [
									"Freebox/Other/free-go.svg"
								],
								"curr_bytes": 0,
								"type": "rm",
								"to": "",
								"id": 6,
								"nfiles": 1,
								"created_ts": 1712862825,
								"state": "done",
								"total_bytes_done": 0,
								"rate": 0,
								"from": "Freebox/Other/free-go.svg",
								"dst": "",
								"eta": 0,
								"error": "none",
								"progress": 100
							}
						}`),
					),
				)
			})
			It("should return the correct file info", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedTask).To(Equal(types.FileSystemTask{
					ID:                            6,
					Type:                          types.FileTaskTypeRemove,
					State:                         types.FileTaskStateDone,
					Error:                         types.FileTaskErrorNone,
					CurrentBytesDone:              0,
					TotalBytes:                    0,
					NumberFilesDone:               1,
					StartedTimestamp:              1712862825,
					DurationSeconds:               0,
					DoneTimestamp:                 1712862825,
					CurrentBytes:                  0,
					To:                            "",
					NumberFiles:                   1,
					CreatedTimestamp:              1712862825,
					TotalBytesDone:                0,
					From:                          "Freebox/Other/free-go.svg",
					ProcessingRate:                0,
					EstimatedTimeRemainingSeconds: 0,
					ProgressPercent:               100,
					Sources: []string{
						"Freebox/Other/free-go.svg",
					},
					Destination: "",
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fs/tasks/%d", version, identifier)),
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
	Context("listing filesystem task", func() {
		returnedTasks := new([]types.FileSystemTask)
		JustBeforeEach(func(ctx SpecContext) {
			*returnedTasks, *returnedErr = freeboxClient.ListFileSystemTasks(ctx)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fs/tasks/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": [
								{
									"curr_bytes_done": 0,
									"total_bytes": 0,
									"nfiles_done": 0,
									"started_ts": 1355834253,
									"duration": 3,
									"done_ts": 0,
									"curr_bytes": 0,
									"type": "extract",
									"to": "oxygennosvg/128x128/mimetypes/application_x_nzb.png",
									"id": 12,
									"nfiles": 0,
									"created_ts": 1355834253,
									"state": "paused",
									"total_bytes_done": 0,
									"from": "/Disque dur/tests/oxygennosvg.tar.gz",
									"rate": 0,
									"eta": 0,
									"error": "none",
									"progress": 0,
									"src": [
										"/Disque dur/tests/oxygennosvg.tar.gz"
									],
									"dst": "/Disque dur/tests/oxygennosvg"
								},
								{
									"id": 11,
									"curr_bytes_done": 0,
									"total_bytes": 0,
									"nfiles_done": 0,
									"started_ts": 1355834187,
									"duration": 0,
									"done_ts": 1355834187,
									"curr_bytes": 0,
									"type": "rm",
									"to": "",
									"nfiles": 0,
									"created_ts": 1355834187,
									"state": "done",
									"total_bytes_done": 0,
									"from": "/Disque dur/test/testiso.1.iso",
									"rate": 0,
									"eta": 0,
									"error": "none",
									"progress": 100,
									"src": [
										"/Disque dur/test/testiso.1.iso"
									]
								}
							]
						}`),
					),
				)
			})
			It("should return the correct file info", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedTasks).To(ConsistOf(
					Equal(types.FileSystemTask{
						ID:                            12,
						Type:                          "extract",
						State:                         "paused",
						Error:                         "none",
						CurrentBytesDone:              0,
						TotalBytes:                    0,
						NumberFilesDone:               0,
						StartedTimestamp:              1355834253,
						DurationSeconds:               3,
						DoneTimestamp:                 0,
						CurrentBytes:                  0,
						To:                            "oxygennosvg/128x128/mimetypes/application_x_nzb.png",
						NumberFiles:                   0,
						CreatedTimestamp:              1355834253,
						TotalBytesDone:                0,
						From:                          "/Disque dur/tests/oxygennosvg.tar.gz",
						ProcessingRate:                0,
						EstimatedTimeRemainingSeconds: 0,
						ProgressPercent:               0,
						Sources: []string{
							"/Disque dur/tests/oxygennosvg.tar.gz",
						},
						Destination: "/Disque dur/tests/oxygennosvg",
					}),
					Equal(types.FileSystemTask{
						ID:                            11,
						Type:                          "rm",
						State:                         "done",
						Error:                         "none",
						CurrentBytesDone:              0,
						TotalBytes:                    0,
						NumberFilesDone:               0,
						StartedTimestamp:              1355834187,
						DurationSeconds:               0,
						DoneTimestamp:                 1355834187,
						CurrentBytes:                  0,
						To:                            "",
						NumberFiles:                   0,
						CreatedTimestamp:              1355834187,
						TotalBytesDone:                0,
						From:                          "/Disque dur/test/testiso.1.iso",
						ProcessingRate:                0,
						EstimatedTimeRemainingSeconds: 0,
						ProgressPercent:               100,
						Sources: []string{
							"/Disque dur/test/testiso.1.iso",
						},
						Destination: "",
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
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fs/tasks/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": "a string"
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("deleting filesystem task", func() {
		const identifier int64 = 42
		JustBeforeEach(func(ctx SpecContext) {
			*returnedErr = freeboxClient.DeleteFileSystemTask(ctx, identifier)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/fs/tasks/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true
						}`),
					),
				)
			})
			It("should return the correct file info", func() {
				Expect(*returnedErr).To(BeNil())
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
						ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/fs/tasks/%d", version, identifier)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": false
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("get a file", func() {
		var (
			path         = "path/to/file"
			returnedFile = new(types.File)
		)
		JustBeforeEach(func(ctx SpecContext) {
			*returnedFile, *returnedErr = freeboxClient.GetFile(ctx, path)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/dl/cGF0aC90by9maWxl", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `the-content`, http.Header{
							"Content-Type":        []string{"application/octet-stream"},
							"Content-Disposition": []string{`attachment; filename="file"`},
						}),
					),
				)
			})
			It("should return the task", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(returnedFile.ContentType).To(Equal("application/octet-stream"))
				Expect(returnedFile.FileName).To(Equal("file"))
				Expect(io.ReadAll(returnedFile.Content)).To(BeEquivalentTo([]byte("the-content")))
			})
		})
		Context("when the server does not mention content disposition", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/dl/cGF0aC90by9maWxl", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `content`, http.Header{
							"Content-Type": []string{"application/octet-stream"},
						}),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(returnedFile.ContentType).To(Equal("application/octet-stream"))
				Expect(returnedFile.FileName).To(BeEmpty())
			})
		})
		Context("when the server does not mention content type", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/dl/cGF0aC90by9maWxl", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `content`, http.Header{
							"Content-Type":        []string{},
							"Content-Disposition": []string{`attachment; filename="file"`},
						}),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(returnedFile.ContentType).To(BeEmpty())
				Expect(returnedFile.FileName).To(Equal("file"))
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/dl/cGF0aC90by9maWxl", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusBadRequest, `bad request`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("creating a directory", func() {
		var (
			parent       = "path/to/parent"
			name         = "folder"
			returnedPath = new(string)
		)
		JustBeforeEach(func() {
			*returnedPath, *returnedErr = freeboxClient.CreateDirectory(context.Background(), parent, name)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/mkdir/", version)),
						ghttp.VerifyJSON(`{
							"parent": "cGF0aC90by9wYXJlbnQ=",
							"dirname": "folder"
						}`),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": "cGF0aC90by9wYXJlbnQvZm9sZGVy"
						}`),
					),
				)
			})
			It("should return the task", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedPath).To(Equal("path/to/parent/folder"))
			})
		})
		Context("when the folder already exists", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/mkdir/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"msg": "Erreur lors de la création du dossier: Le fichier existe déjà",
							"success": false,
							"error_code": "destination_conflict"
						}`),
					),
				)
			})
			It("should return the correct error", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrDestinationConflict))
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
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/mkdir/", version)),
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
	Context("hash a file", func() {
		const path = "path/to/file"

		var (
			returnedTask = new(types.FileSystemTask)
			returnedErr  = new(error)
		)

		JustBeforeEach(func(ctx context.Context) {
			*returnedTask, *returnedErr = freeboxClient.AddHashFileTask(ctx, types.HashPayload{
				HashType: types.HashTypeSHA256,
				Path:     path,
			})
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/hash/", version)),
						ghttp.VerifyJSON(`{
							"hash_type": "sha256",
							"src": "cGF0aC90by9maWxl"
						}`),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"id": 1234
							}
						}`),
					),
				)
			})
			It("should return the task", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*&returnedTask.ID).To(BeEquivalentTo(1234))
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
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/hash/", version)),
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
	Context("copy files", func() {
		const (
			path1 = "path/to/file1"
			path2 = "path/to/file2"
			dest  = "path/to/dest"
		)

		var (
			returnedTask = new(types.FileSystemTask)
			returnedErr  = new(error)
		)

		JustBeforeEach(func(ctx context.Context) {
			*returnedTask, *returnedErr = freeboxClient.CopyFiles(ctx, []string{path1, path2}, dest, types.FileCopyModeSkip)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/cp/", version)),
						ghttp.VerifyJSON(`{
							"dst": "cGF0aC90by9kZXN0",
							"files": [
								"cGF0aC90by9maWxlMQ==",
								"cGF0aC90by9maWxlMg=="
							],
							"mode": "skip"
						}`),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"id": 1234
							}
						}`),
					),
				)
			})
			It("should return the task", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*&returnedTask.ID).To(BeEquivalentTo(1234))
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
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/cp/", version)),
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
	Context("move files", func() {
		const (
			path1 = "path/to/file1"
			path2 = "path/to/file2"
			dest  = "path/to/dest"
		)

		var (
			returnedTask = new(types.FileSystemTask)
			returnedErr  = new(error)
		)

		JustBeforeEach(func(ctx context.Context) {
			*returnedTask, *returnedErr = freeboxClient.MoveFiles(ctx, []string{path1, path2}, dest, types.FileMoveModeSkip)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/mv/", version)),
						ghttp.VerifyJSON(`{
							"dst": "cGF0aC90by9kZXN0",
							"files": [
								"cGF0aC90by9maWxlMQ==",
								"cGF0aC90by9maWxlMg=="
							],
							"mode": "skip"
						}`),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"id": 1234
							}
						}`),
					),
				)
			})
			It("should return the task", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*&returnedTask.ID).To(BeEquivalentTo(1234))
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
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/mv/", version)),
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
	Context("get a hash result", func() {
		const path = "path/to/file"

		var (
			hashResult  = new(string)
			returnedErr = new(error)
		)

		JustBeforeEach(func(ctx context.Context) {
			*hashResult, *returnedErr = freeboxClient.GetHashResult(ctx, 1234)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fs/tasks/1234/hash/", version)),
						verifyAuth(*sessionToken),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": "the-hash-result"
						}`),
					),
				)
			})
			It("should return the task", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*hashResult).To(BeEquivalentTo("the-hash-result"))
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/fs/tasks/1234/hash/", version)),
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
	Context("ExtractFile", func() {
		returnedTask := new(types.FileSystemTask)

		JustBeforeEach(func(ctx SpecContext) {
			*returnedTask, *returnedErr = freeboxClient.ExtractFile(ctx, types.ExtractFilePayload{
				Src:       "src",
				Dst:       "dst",
				Overwrite: true,
			})
		})

		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/extract/", version)),
						verifyAuth(*sessionToken),
						ghttp.VerifyJSON(`{
							"src": "L3NyYw==",
							"dst": "L2RzdA==",
							"password": "",
							"delete_archive": false,
							"overwrite": true
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"id": 48,
								"type": "extract",
								"state": "running",
								"error": "none",
								"curr_bytes_done": 0,
								"total_bytes": 0,
								"nfiles_done": 0,
								"started_ts": 1355842252,
								"duration": 0,
								"done_ts": 0,
								"curr_bytes": 0,
								"to": "/Disque dur/old_hdd",
								"nfiles": 0,
								"created_ts": 1355842252,
								"total_bytes_done": 0,
								"from": "/Disque dur/old_hdd/testiso.1.iso",
								"rate": 0,
								"eta": 0,
								"progress": 0
							}
						}`),
					),
				)
			})

			It("should return the task", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedTask).To(Equal(types.FileSystemTask{
					ID:                            48,
					Type:                          "extract",
					State:                         "running",
					Error:                         "none",
					CurrentBytesDone:              0,
					TotalBytes:                    0,
					NumberFilesDone:               0,
					StartedTimestamp:              1355842252,
					DurationSeconds:               0,
					DoneTimestamp:                 0,
					CurrentBytes:                  0,
					To:                            "/Disque dur/old_hdd",
					NumberFiles:                   0,
					CreatedTimestamp:              1355842252,
					TotalBytesDone:                0,
					From:                          "/Disque dur/old_hdd/testiso.1.iso",
					ProcessingRate:                0,
					EstimatedTimeRemainingSeconds: 0,
					ProgressPercent:               0,
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
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/fs/extract/", version)),
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
