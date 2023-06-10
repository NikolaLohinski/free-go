package client_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("authorize", func() {
	const (
		returnedTrackID      = "123"
		returnedPrivateToken = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	)
	var (
		freeboxClient        client.Client
		authorizationRequest = types.AuthorizationRequest{
			Name:    "name",
			Version: "0.0.0",
			Device:  "device",
		}

		server   *ghttp.Server
		endpoint = new(string)

		privateToken = new(string)
		returnedErr  = new(error)
	)
	BeforeEach(func() {
		server = ghttp.NewServer()
		*endpoint = server.Addr()

		freeboxClient = Must(client.New(*endpoint, version)).(client.Client).WithAppID(appID)
	})
	JustBeforeEach(func() {
		*privateToken, *returnedErr = freeboxClient.Authorize(authorizationRequest)
	})
	AfterEach(func() {
		server.Close()
	})
	Context("when the authorization is approved after some time", func() {
		BeforeEach(func() {
			client.ClientAuthorizeRetryDelay = time.Millisecond * 50
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/login/authorize", version)),
					ghttp.VerifyContentType("application/json"),
					ghttp.VerifyJSON(`{
						"app_id": "`+appID+`",
						"app_name": "`+authorizationRequest.Name+`",
						"app_version": "`+authorizationRequest.Version+`",
						"device_name": "`+authorizationRequest.Device+`"
					}`),
					ghttp.RespondWith(http.StatusOK, `{
					    "success": true,
					    "result": {
					        "app_token": "`+returnedPrivateToken+`",
					        "track_id": `+returnedTrackID+`
					    }
					}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/login/authorize/%s", version, returnedTrackID)),
					ghttp.RespondWith(http.StatusOK, `{
					    "success": true,
					    "result": {
					        "status": "pending",
					        "challenge": "KWmElA9q9R49DsZUzjVpe0D/3aze2sBf",
					        "password_salt": "PJpG867vNjvbYY2z67Yy4164kEmmfrOC"
					    }
					}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/login/authorize/%s", version, returnedTrackID)),
					ghttp.RespondWith(http.StatusOK, `{
					    "success": true,
					    "result": {
					        "status": "granted",
					        "challenge": "KWmElA9q9R49DsZUzjVpe0D/3aze2sBf",
					        "password_salt": "PJpG867vNjvbYY2z67Yy4164kEmmfrOC"
					    }
					}`),
				),
			)
		})
		It("should return the private token provided by the server", func() {
			Expect(*returnedErr).To(BeNil())
			Expect(*privateToken).To(Equal(returnedPrivateToken))
		})
	})
	Context("when the authorization ends in an unexpected status", func() {
		BeforeEach(func() {
			client.ClientAuthorizeRetryDelay = time.Millisecond * 50
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/login/authorize", version)),
					ghttp.VerifyContentType("application/json"),
					ghttp.VerifyJSON(`{
						"app_id": "`+appID+`",
						"app_name": "`+authorizationRequest.Name+`",
						"app_version": "`+authorizationRequest.Version+`",
						"device_name": "`+authorizationRequest.Device+`"
					}`),
					ghttp.RespondWith(http.StatusOK, `{
					    "success": true,
					    "result": {
					        "app_token": "`+returnedPrivateToken+`",
					        "track_id": `+returnedTrackID+`
					    }
					}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/login/authorize/%s", version, returnedTrackID)),
					ghttp.RespondWith(http.StatusOK, `{
					    "success": true,
					    "result": {
					        "status": "denied"
					    }
					}`),
				),
			)
		})
		It("should return an explicit error", func() {
			Expect(*returnedErr).ToNot(BeNil())
			Expect((*returnedErr).Error()).To(MatchRegexp(".* received unexpected track status: denied"))
		})
	})
	Context("when the authorization times out on client side", func() {
		BeforeEach(func() {
			client.ClientAuthorizeRetryDelay = time.Millisecond * 50
			client.ClientAuthorizeGrantingTimeout = time.Millisecond * 1
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/login/authorize", version)),
					ghttp.VerifyContentType("application/json"),
					ghttp.VerifyJSON(`{
						"app_id": "`+appID+`",
						"app_name": "`+authorizationRequest.Name+`",
						"app_version": "`+authorizationRequest.Version+`",
						"device_name": "`+authorizationRequest.Device+`"
					}`),
					ghttp.RespondWith(http.StatusOK, `{
					    "success": true,
					    "result": {
					        "app_token": "`+returnedPrivateToken+`",
					        "track_id": `+returnedTrackID+`
					    }
					}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/login/authorize/%s", version, returnedTrackID)),
					ghttp.RespondWith(http.StatusOK, `{
					    "success": true,
					    "result": {
					        "status": "pending",
					        "challenge": "KWmElA9q9R49DsZUzjVpe0D/3aze2sBf",
					        "password_salt": "PJpG867vNjvbYY2z67Yy4164kEmmfrOC"
					    }
					}`),
				),
			)
		})
		It("should return the private token provided by the server", func() {
			Expect(*returnedErr).ToNot(BeNil())
			Expect((*returnedErr).Error()).To(MatchRegexp(".* reached hard timeout after .* waiting for token approval"))
		})
	})

	Context("when appID is not set", func() {
		BeforeEach(func() {
			freeboxClient = Must(client.New(*endpoint, version)).(client.Client)
		})
		It("should return the private token provided by the server", func() {
			Expect(*returnedErr).ToNot(BeNil())
			Expect((*returnedErr).Error()).To(MatchRegexp(".* app ID is not set"))
		})
	})
})
