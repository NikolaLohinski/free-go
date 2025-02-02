package client_test

import (
	"context"
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/nikolalohinski/free-go/client"
)

type httpClientMock struct {
	request  *http.Request
	response func() (*http.Response, error)
}

func (m *httpClientMock) Do(request *http.Request) (*http.Response, error) {
	m.request = request
	return m.response()
}

var _ = Describe("client", func() {
	var (
		server   *ghttp.Server
		endpoint = new(string)

		freeboxClient client.Client
		returnedErr   = new(error)
	)
	BeforeEach(func() {
		server = ghttp.NewServer()
		DeferCleanup(server.Close)

		*endpoint = server.Addr()
	})
	JustBeforeEach(func() {
		freeboxClient, *returnedErr = client.New(*endpoint, version)
	})
	Context("default", func() {
		It("should not return an error", func() {
			Expect(*returnedErr).To(BeNil())
		})
	})
	Context("when the built endpoint is invalid", func() {
		BeforeEach(func() {
			*endpoint = ":{/@)=$£"
		})
		It("should return an error", func() {
			Expect(*returnedErr).ToNot(BeNil())
		})
	})
	Context("when overriding the default http client", func() {
		httpMock := new(httpClientMock)
		BeforeEach(func() {
			httpMock = &httpClientMock{
				response: func() (*http.Response, error) {
					return nil, errors.New("just fail")
				},
			}
		})
		JustBeforeEach(func() {
			_, *returnedErr = freeboxClient.WithHTTPClient(httpMock).APIVersion(context.Background())
		})
		It("should have called the overridden HTTP client", func() {
			Expect(*httpMock.request).ToNot(BeNil())
			Expect((*returnedErr).Error()).To(MatchRegexp("just fail"))
		})
	})
})
