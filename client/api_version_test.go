package client_test

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("api version", func() {
	const (
		version = "v0"
		token   = ""
	)
	var (
		server   *ghttp.Server
		endpoint = new(string)

		apiVersion  = new(types.APIVersion)
		returnedErr = new(error)
	)
	BeforeEach(func() {
		server = ghttp.NewServer()
		*endpoint = server.Addr()
	})
	JustBeforeEach(func() {
		*apiVersion, *returnedErr = Must(client.New(*endpoint, version)).(client.Client).APIVersion()
	})
	AfterEach(func() {
		server.Close()
	})
	Context("default", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/api_version", version)),
					ghttp.RespondWith(http.StatusOK, heredoc.Doc(`{
						"box_model_name": "Freebox v0",
						"api_base_url": "/api/",
						"https_port": 12345,
						"device_name": "Freebox Server",
						"https_available": true,
						"box_model": "unit/test",
						"api_domain": "test.fbxos.fr",
						"uid": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
						"api_version": "0",
						"device_type": "FreeboxServer0,0"
					}`)),
				),
			)
		})
		It("should return the correct API version", func() {
			Expect(*returnedErr).To(BeNil())
			Expect(*apiVersion).To(Equal(types.APIVersion{
				APIBaseURL:     "/api/",
				HTTPSPort:      12345,
				DeviceName:     "Freebox Server",
				HTTPSAvailable: true,
				BoxModel:       "unit/test",
				BoxModelName:   "Freebox v0",
				APIDomain:      "test.fbxos.fr",
				UID:            "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				APIVersion:     "0",
				DeviceType:     "FreeboxServer0,0",
			}))
		})
	})
	Context("when the server is unavailable", func() {
		BeforeEach(func() {
			server.Close()
		})
		It("should return an error", func() {
			Expect(*returnedErr).ToNot(BeNil())
			Expect((*returnedErr).Error()).To(MatchRegexp("failed to perform request: .*"))
		})
	})
	Context("when the server returns an unexpected status", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/api_version", version)),
					ghttp.RespondWith(http.StatusBadRequest, "test error from server"),
				),
			)
		})
		It("should return an error", func() {
			Expect(*returnedErr).ToNot(BeNil())
			Expect((*returnedErr).Error()).To(Equal("failed with status '400': server returned 'test error from server'"))
		})
	})
	Context("when the server returns an invalid JSON object", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/api_version", version)),
					ghttp.RespondWith(http.StatusOK, "{"),
				),
			)
		})
		It("should return an error", func() {
			Expect(*returnedErr).ToNot(BeNil())
			Expect((*returnedErr).Error()).To(Equal("failed to unmarshal response body '{': unexpected end of JSON input"))
		})
	})
})
