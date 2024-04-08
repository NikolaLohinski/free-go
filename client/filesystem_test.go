package client_test

import (
	"context"
	"encoding/json"
	"fmt"
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
		*endpoint = server.Addr()

		freeboxClient = Must(client.New(*endpoint, version)).(client.Client).
			WithAppID(appID).
			WithPrivateToken(privateToken)

		*sessionToken = setupLoginFlow(server)
	})
	AfterEach(func() {
		server.Close()
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
})
