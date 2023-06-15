package client_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nikolalohinski/free-go/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

const (
	version      = "v0"
	appID        = "test"
	privateToken = "xXXyyX9999wwwwwwwwxxx99999XXYYYYYYWWW000000000999999XXXXX9999Yx"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "unit-tests")
}

func Must(r interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return r
}

func verifyAuth(sessionToken string) http.HandlerFunc {
	return ghttp.VerifyHeaderKV(client.AuthHeader, sessionToken)
}

func setupLoginFlow(server *ghttp.Server) string {
	sessionToken := "EfETzVibY7K5vZVsq+MjtD6pDJoAaYQiqyXwS5kFvooTczPMk7Tz+6//aTe9zZNy"

	server.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/login", version)),
			ghttp.RespondWith(http.StatusOK, `{
					"success": true,
					"result": {
						"logged_in": false,
						"challenge": "9Va31tSgQWM853j0kSCtBUyzYNhPN7IY"
					}
				}`),
		),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/login/session", version)),
			ghttp.VerifyContentType("application/json"),
			ghttp.VerifyJSON(`{
			    "app_id": "`+appID+`",
			    "password": "c3464d210c1be4f1ef6f34c578d463fc28d40a61"
			}`),
			ghttp.RespondWith(http.StatusOK, `{
				"result": {
					"session_token": "`+sessionToken+`",
					"challenge": "9Va31tSgQWM853j0kSCtBUyzYNhPN7IY",
					"permissions": {}
				},
				"success": true
			}`),
		),
	)

	return sessionToken
}
