package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/onsi/gomega/gstruct"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("events", func() {
	var (
		freeboxClient client.Client

		server   *ghttp.Server
		endpoint = new(string)

		sessionToken = new(string)

		ctx           context.Context
		cancelContext func()
		events        = new([]types.EventDescription)

		returnedChannel = new(chan types.Event)
		returnedErr     = new(error)
	)
	BeforeEach(func() {
		ctx, cancelContext = context.WithCancel(context.Background())
		server = ghttp.NewServer()
		*endpoint = server.Addr()

		*returnedChannel = make(chan types.Event)

		freeboxClient = Must(client.New(*endpoint, version)).(client.Client).
			WithAppID(appID).
			WithPrivateToken(privateToken)

		*sessionToken = setupLoginFlow(server)
	})
	AfterEach(func() {
		server.Close()
	})
	JustBeforeEach(func() {
		*returnedChannel, *returnedErr = freeboxClient.ListenEvents(ctx, *events)
	})
	Context("default", func() {
		BeforeEach(func() {
			*events = []types.EventDescription{
				{
					Source: "foo",
					Name:   "bar",
				},
			}
			server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
				ws, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
				if err != nil {
					Expect(err).To(BeNil())
				}
				defer ws.Close()
				ws.SetCloseHandler(func(code int, text string) error {
					// Ignore default close handler to explicitly assert the websocket is closed correctly
					Expect(code).To(Equal(websocket.CloseNormalClosure))
					Expect(text).To(Equal(websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))

					Expect(ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))).To(BeNil())
					return nil
				})

				messageType, message, err := ws.ReadMessage()
				Expect(err).To(BeNil())
				Expect(messageType).To(Equal(websocket.TextMessage))
				Expect(string(message)).To(MatchJSON(`{
					"action": "register",
					"events": [
						"foo_bar"
					]
				}`))
				Expect(ws.WriteMessage(websocket.TextMessage, []byte(`{
					"action": "register",
					"success": true
				}`))).To(BeNil())
				Expect(ws.WriteMessage(websocket.TextMessage, []byte(`{
					"action": "notification",
					"success": true,
					"source": "foo",
					"event": "bar",
					"result": {"not":"checked"}
				}`))).To(BeNil())
			})
		})
		It("should send the expected events through the returned channel", func() {
			Expect(*returnedErr).To(BeNil())
			Expect(<-*returnedChannel).To(Equal(types.Event{
				Notification: types.EventNotification{
					Action:  "notification",
					Success: true,
					Source:  "foo",
					Event:   "bar",
					Result:  json.RawMessage("{\"not\":\"checked\"}"),
				},
				Error: nil,
			}))
			cancelContext()
		})
	})
	Context("when the context is canceled before receiving a notification", func() {
		BeforeEach(func() {
			*events = []types.EventDescription{
				{
					Source: "foo",
					Name:   "bar",
				},
			}
			server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
				ws, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
				if err != nil {
					Expect(err).To(BeNil())
				}
				defer ws.Close()

				messageType, message, err := ws.ReadMessage()
				Expect(err).To(BeNil())
				Expect(messageType).To(Equal(websocket.TextMessage))
				Expect(string(message)).To(MatchJSON(`{
					"action": "register",
					"events": [
						"foo_bar"
					]
				}`))
				Expect(ws.WriteMessage(websocket.TextMessage, []byte(`{
					"action": "register",
					"success": true
				}`))).To(BeNil())

				cancelContext()

				_, _, err = ws.ReadMessage()
				Expect(websocket.IsCloseError(err, websocket.CloseNormalClosure)).To(BeTrue(), "websocket should have been closed by client")
			})
		})
		It("should not return an error", func() {
			Expect(*returnedErr).To(BeNil())
		})
	})
	Context("when the received notification is unexpected", func() {
		BeforeEach(func() {
			*events = []types.EventDescription{
				{
					Source: "foo",
					Name:   "bar",
				},
			}
			server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
				ws, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
				if err != nil {
					Expect(err).To(BeNil())
				}
				defer ws.Close()

				messageType, message, err := ws.ReadMessage()
				Expect(err).To(BeNil())
				Expect(messageType).To(Equal(websocket.TextMessage))
				Expect(string(message)).To(MatchJSON(`{
					"action": "register",
					"events": [
						"foo_bar"
					]
				}`))
				Expect(ws.WriteMessage(websocket.TextMessage, []byte(`{
					"action": "register",
					"success": true
				}`))).To(BeNil())

				Expect(ws.WriteMessage(websocket.TextMessage, []byte(`{
					"action": "error",
					"success": false
				}`))).To(BeNil())
			})
		})
		It("should return an error", func() {
			Expect(*returnedErr).To(BeNil())
			Expect(<-*returnedChannel).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Error": Not(BeNil()),
			}))
		})
	})

	Context("when logging in fails", func() {
		BeforeEach(func() {
			server = ghttp.NewServer()
			*endpoint = server.Addr()
			freeboxClient = Must(client.New(*endpoint, version)).(client.Client).
				WithAppID(appID).
				WithPrivateToken(privateToken)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/login", version)),
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				),
			)
		})
		It("should return an error", func() {
			Expect(*returnedErr).ToNot(BeNil())
		})
	})
	Context("when dialing the websocket fails", func() {
		BeforeEach(func() {
			server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
				w.WriteHeader(http.StatusInternalServerError)
			})
		})
		It("should return an error", func() {
			Expect(*returnedErr).ToNot(BeNil())
		})
	})
	Context("when registering for notifications fails", func() {
		BeforeEach(func() {
			*events = []types.EventDescription{
				{
					Source: "foo",
					Name:   "bar",
				},
			}
			server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
				ws, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
				if err != nil {
					Expect(err).To(BeNil())
				}
				defer ws.Close()

				messageType, message, err := ws.ReadMessage()
				Expect(err).To(BeNil())
				Expect(messageType).To(Equal(websocket.TextMessage))
				Expect(string(message)).To(MatchJSON(`{
					"action": "register",
					"events": [
						"foo_bar"
					]
				}`))
				ws.Close()
			})
		})
		It("should return an error", func() {
			Expect(*returnedErr).ToNot(BeNil())
		})
	})
	Context("when registering for notifications returns an unexpected result", func() {
		BeforeEach(func() {
			*events = []types.EventDescription{
				{
					Source: "foo",
					Name:   "bar",
				},
			}
			server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
				ws, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
				if err != nil {
					Expect(err).To(BeNil())
				}
				defer ws.Close()

				messageType, message, err := ws.ReadMessage()
				Expect(err).To(BeNil())
				Expect(messageType).To(Equal(websocket.TextMessage))
				Expect(string(message)).To(MatchJSON(`{
					"action": "register",
					"events": [
						"foo_bar"
					]
				}`))
				Expect(ws.WriteMessage(websocket.TextMessage, []byte(`{
					"action": "register",
					"success": false
				}`))).To(BeNil())
			})
		})
		It("should return an error", func() {
			Expect(*returnedErr).ToNot(BeNil())
		})
	})

})
