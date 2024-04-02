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
})
