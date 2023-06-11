package client_test

import (
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/nikolalohinski/free-go/client"
)

type transportMock struct {
	request  *http.Request
	response func() (*http.Response, error)
}

func (m *transportMock) RoundTrip(request *http.Request) (*http.Response, error) {
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
		*endpoint = server.Addr()
	})
	AfterEach(func() {
		server.Close()
	})
	JustBeforeEach(func() {
		freeboxClient, *returnedErr = client.New(*endpoint, version)
	})
	Context("default", func() {
		It("should not return an error", func() {
			Expect(*returnedErr).To(BeNil())
		})
	})
	Context("when overriding the default http client", func() {
		mock := new(transportMock)
		BeforeEach(func() {
			mock = &transportMock{
				response: func() (*http.Response, error) {
					return nil, errors.New("just fail")
				},
			}
		})
		JustBeforeEach(func() {
			_, *returnedErr = freeboxClient.WithHTTPClient(&http.Client{
				Transport: mock,
			}).APIVersion()
		})
		It("should have called the overridden HTTP client", func() {
			Expect(*mock.request).ToNot(BeNil())
			Expect((*returnedErr).Error()).To(MatchRegexp("just fail"))
		})
	})
})
