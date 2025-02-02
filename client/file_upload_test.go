package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/ghttp"
	"github.com/onsi/gomega/gstruct"

	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("FileUpload", func() {
	var (
		freeboxClient client.Client

		server   *ghttp.Server
		endpoint = new(string)

		uploadContext context.Context

		sessionToken = new(string)

		returnedWriter io.WriteCloser
		returnedErr    = new(error)

		content []byte
		size    int

		ws *websocket.Conn
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		DeferCleanup(server.Close)
		*endpoint = server.Addr()

		freeboxClient = Must(client.New(*endpoint, version)).
			WithAppID(appID).
			WithPrivateToken(privateToken)

		*sessionToken = setupLoginFlow(server)

		content = []byte("data")
		size = len(content)

		uploadContext = context.Background()

		ws = new(websocket.Conn)
	})

	JustBeforeEach(func() {
		ctx, cancel := context.WithTimeout(uploadContext, 10*time.Second)
		DeferCleanup(cancel)

		returnedWriter, *returnedErr = freeboxClient.FileUploadStart(ctx, types.FileUploadStartActionInput{
			Dirname:  "dir",
			Filename: "the-file",
			Force:    types.FileUploadStartActionForceOverwrite,
			Size:     size,
		})
	})

	Describe("FileUploadStart", func() {
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()

					Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
					var err error
					ws, err = (&websocket.Upgrader{}).Upgrade(w, r, nil)
					Expect(err).ToNot(HaveOccurred())

					defer ws.Close()

					ws.SetCloseHandler(func(code int, text string) error {
						// Ignore default close handler to explicitly assert the websocket is closed correctly
						Expect(code).To(Equal(websocket.CloseNormalClosure))
						Expect(text).To(Equal(websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))

						Expect(ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))).To(Succeed())
						return nil
					})

					var result types.FileUploadStartAction
					Expect(readJSON(ws, &result)).To(Succeed())
					Expect(result).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Action":   BeEquivalentTo(types.FileUploadStartActionNameUploadStart),
						"Size":     BeEquivalentTo(4),
						"Dirname":  BeEquivalentTo("dir"),
						"Filename": Equal("the-file"),
						"Force":    BeEquivalentTo(types.FileUploadStartActionForceOverwrite),
					}))
					Expect(result.RequestID).To(BeNumerically(">", 0))

					Expect(writeJSON(ws, &types.WebSocketResponse[interface{}]{
						Success:   true,
						Action:    types.FileUploadStartActionNameUploadStart,
						RequestID: result.RequestID,
					})).To(Succeed())

					messageType, reader, err := ws.NextReader()
					Expect(err).ToNot(HaveOccurred())
					Expect(messageType).To(Equal(websocket.BinaryMessage))
					received, err := io.ReadAll(gbytes.TimeoutReader(reader, time.Second))
					Expect(err).ToNot(HaveOccurred())
					Expect(string(received)).To(Equal(string(content)))

					Expect(writeJSON(ws, &types.WebSocketResponse[types.FileUploadChunkResponse]{
						Success:   true,
						Action:    types.FileUploadStartActionNameUploadData,
						RequestID: result.RequestID,
						Result: types.FileUploadChunkResponse{
							TotalLen: len(received),
							Complete: false,
						},
					})).To(Succeed())

					var finalize types.FileUploadFinalize
					Expect(readJSON(ws, &finalize)).To(Succeed())
					Expect(finalize).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Action":    BeEquivalentTo(types.FileUploadStartActionNameUploadFinalize),
						"RequestID": BeEquivalentTo(result.RequestID),
					}))

					Expect(writeJSON(ws, &types.WebSocketResponse[types.FileUploadFinalizeResponse]{
						Success:   true,
						Action:    types.FileUploadStartActionNameUploadFinalize,
						RequestID: result.RequestID,
						Result: types.FileUploadFinalizeResponse{
							TotalLen:  len(received),
							Complete:  true,
							Cancelled: false,
						},
					})).To(Succeed())
				})
			})

			It("should start a file upload", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(returnedWriter).ToNot(BeNil())

				Expect(returnedWriter.Write(content)).To(BeEquivalentTo(size))

				Expect(returnedWriter.Close()).To(Succeed())
			})
		})

		Context("File already exist", func() {
			BeforeEach(func() {
				server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()

					Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
					var err error
					ws, err = (&websocket.Upgrader{}).Upgrade(w, r, nil)
					Expect(err).ToNot(HaveOccurred())

					defer ws.Close()

					ws.SetCloseHandler(func(code int, text string) error {
						// Ignore default close handler to explicitly assert the websocket is closed correctly
						Expect(code).To(Equal(websocket.CloseNormalClosure))
						Expect(text).To(Equal(websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))

						Expect(ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))).To(Succeed())
						return nil
					})

					var result types.FileUploadStartAction
					Expect(readJSON(ws, &result)).To(Succeed())
					Expect(result).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Action":   BeEquivalentTo(types.FileUploadStartActionNameUploadStart),
						"Size":     BeEquivalentTo(4),
						"Dirname":  BeEquivalentTo("dir"),
						"Filename": Equal("the-file"),
						"Force":    BeEquivalentTo(types.FileUploadStartActionForceOverwrite),
					}))
					Expect(result.RequestID).To(BeNumerically(">", 0))

					Expect(writeJSON(ws, &types.WebSocketResponse[interface{}]{
						Success:   false,
						Action:    types.FileUploadStartActionNameUploadStart,
						RequestID: result.RequestID,
						ErrorCode: "conflict",
						Message:   "Le fichier existe déjà",
					})).To(Succeed())
				})
			})

			It("conflicts", func() {
				Expect(*returnedErr).To(MatchError(ContainSubstring("conflict: Le fichier existe déjà")))
			})
		})

		Context("with non JSON message", func() {
			BeforeEach(func() {
				server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()

					Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
					var err error
					ws, err = (&websocket.Upgrader{}).Upgrade(w, r, nil)
					Expect(err).ToNot(HaveOccurred())

					defer ws.Close()

					ws.NetConn().Write([]byte("invalid message"))
				})
			})

			It("fails", func() {
				Expect(*returnedErr).To(HaveOccurred())
			})
		})

		Context("Websocket close early", func() {
			BeforeEach(func() {
				server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()

					Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
					var err error
					ws, err = (&websocket.Upgrader{}).Upgrade(w, r, nil)
					Expect(err).ToNot(HaveOccurred())

					Expect(ws.Close()).To(Succeed())
				})
			})

			It("propagate the error", func() {
				Expect(*returnedErr).To(MatchError(Or(ContainSubstring("connection reset by peer"), ContainSubstring("unexpected EOF"))))
			})
		})

		Context("upload error", func() {
			BeforeEach(func() {
				server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()

					Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
					var err error
					ws, err = (&websocket.Upgrader{}).Upgrade(w, r, nil)
					Expect(err).ToNot(HaveOccurred())

					defer ws.Close()

					ws.SetCloseHandler(func(code int, text string) error {
						// Ignore default close handler to explicitly assert the websocket is closed correctly
						Expect(code).To(Equal(websocket.CloseNormalClosure))
						Expect(text).To(Equal(websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))

						Expect(ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))).To(Succeed())
						return nil
					})

					var result types.FileUploadStartAction
					Expect(readJSON(ws, &result)).To(Succeed())
					Expect(result).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Action":   BeEquivalentTo(types.FileUploadStartActionNameUploadStart),
						"Size":     BeEquivalentTo(4),
						"Dirname":  BeEquivalentTo("dir"),
						"Filename": Equal("the-file"),
						"Force":    BeEquivalentTo(types.FileUploadStartActionForceOverwrite),
					}))
					Expect(result.RequestID).To(BeNumerically(">", 0))

					Expect(writeJSON(ws, &types.WebSocketResponse[interface{}]{
						Success:   true,
						Action:    types.FileUploadStartActionNameUploadStart,
						RequestID: result.RequestID,
					})).To(Succeed())

					messageType, reader, err := ws.NextReader()
					Expect(err).ToNot(HaveOccurred())
					Expect(messageType).To(Equal(websocket.BinaryMessage))
					received, err := io.ReadAll(gbytes.TimeoutReader(reader, time.Second))
					Expect(err).ToNot(HaveOccurred())
					Expect(string(received)).To(Equal(string(content)))

					Expect(writeJSON(ws, &types.WebSocketResponse[types.FileUploadFinalizeResponse]{
						Success:   false,
						Action:    types.FileUploadStartActionNameUploadFinalize,
						RequestID: result.RequestID,
						ErrorCode: "test_error",
						Message:   "no space left on device",
					})).To(Succeed())

					var cancel types.FileUploadCancelAction
					Expect(readJSON(ws, &cancel)).To(Succeed())
					Expect(cancel).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Action":    BeEquivalentTo(types.FileUploadStartActionNameUploadCancel),
						"RequestID": BeEquivalentTo(result.RequestID),
					}))

					Expect(writeJSON(ws, &types.WebSocketResponse[types.FileUploadCancelResponse]{
						Success:   true,
						Action:    types.FileUploadStartActionNameUploadCancel,
						RequestID: result.RequestID,
						Result: types.FileUploadCancelResponse{
							Cancelled: true,
						},
					})).To(Succeed())
				})
			})

			It("propagate the error", func() {
				Expect(*returnedErr).To(Succeed())
				_, err := returnedWriter.Write(content)
				Expect(err).To(MatchError(ContainSubstring("test_error: no space left on device")))

				Expect(returnedWriter.Close()).To(Succeed())
			})
		})

		Context("ignore message with a different requestID", func() {
			BeforeEach(func() {
				server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()

					Expect(r.Header[client.AuthHeader]).To(ContainElement(Equal(*sessionToken)))
					var err error
					ws, err = (&websocket.Upgrader{}).Upgrade(w, r, nil)
					Expect(err).ToNot(HaveOccurred())

					defer ws.Close()

					ws.SetCloseHandler(func(code int, text string) error {
						// Ignore default close handler to explicitly assert the websocket is closed correctly
						Expect(code).To(Equal(websocket.CloseNormalClosure))
						Expect(text).To(Equal(websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))

						Expect(ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))).To(Succeed())
						return nil
					})

					var result types.FileUploadStartAction
					Expect(readJSON(ws, &result)).To(Succeed())
					Expect(result).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Action":   BeEquivalentTo(types.FileUploadStartActionNameUploadStart),
						"Size":     BeEquivalentTo(4),
						"Dirname":  BeEquivalentTo("dir"),
						"Filename": Equal("the-file"),
						"Force":    BeEquivalentTo(types.FileUploadStartActionForceOverwrite),
					}))
					Expect(result.RequestID).To(BeNumerically(">", 0))

					Expect(writeJSON(ws, &types.WebSocketResponse[interface{}]{
						Success:   true,
						Action:    types.FileUploadStartActionNameUploadStart,
						RequestID: result.RequestID,
					})).To(Succeed())

					messageType, reader, err := ws.NextReader()
					Expect(err).ToNot(HaveOccurred())
					Expect(messageType).To(Equal(websocket.BinaryMessage))
					received, err := io.ReadAll(gbytes.TimeoutReader(reader, time.Second))
					Expect(err).ToNot(HaveOccurred())
					Expect(string(received)).To(Equal(string(content)))

					Expect(writeJSON(ws, &types.WebSocketResponse[types.FileUploadChunkResponse]{
						Success:   true,
						Action:    types.FileUploadStartActionNameUploadData,
						RequestID: result.RequestID,
						Result: types.FileUploadChunkResponse{
							TotalLen: len(received),
							Complete: false,
						},
					})).To(Succeed())

					var finalize types.FileUploadFinalize
					Expect(readJSON(ws, &finalize)).To(Succeed())
					Expect(finalize).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Action":    BeEquivalentTo(types.FileUploadStartActionNameUploadFinalize),
						"RequestID": BeEquivalentTo(result.RequestID),
					}))

					Expect(writeJSON(ws, &types.WebSocketResponse[types.FileUploadFinalizeResponse]{
						Success:   false,
						Action:    types.FileUploadStartActionNameUploadFinalize,
						RequestID: result.RequestID + 1,
						ErrorCode: "error_to_ignore",
						Message:   "This error is addressed to another request",
					})).To(Succeed())

					Expect(writeJSON(ws, &types.WebSocketResponse[types.FileUploadFinalizeResponse]{
						Success:   true,
						Action:    types.FileUploadStartActionNameUploadFinalize,
						RequestID: result.RequestID,
						Result: types.FileUploadFinalizeResponse{
							TotalLen:  len(received),
							Complete:  true,
							Cancelled: false,
						},
					})).To(Succeed())
				})
			})

			It("ignore the message", func() {
				Expect(*returnedErr).To(Succeed())
				Expect(returnedWriter).ToNot(BeNil())

				Expect(returnedWriter.Write(content)).To(BeEquivalentTo(size))

				Expect(returnedWriter.Close()).To(Succeed())
			})
		})
	})
})

func writeJSON(ws *websocket.Conn, data interface{}) error {
	w, err := ws.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return fmt.Errorf("next writer: %w", err)
	}

	if err := json.NewEncoder(gbytes.TimeoutWriter(w, time.Second)).Encode(data); err != nil {
		return fmt.Errorf("write json: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("close writer: %w", err)
	}

	return nil
}

func readJSON(ws *websocket.Conn, result interface{}) error {
	messageType, reader, err := ws.NextReader()
	if err != nil {
		return fmt.Errorf("next reader: %w", err)
	}
	if messageType != websocket.TextMessage {
		return fmt.Errorf("unexpected message type: %d", messageType)
	}

	if err := json.NewDecoder(gbytes.TimeoutReader(reader, time.Second)).Decode(result); err != nil {
		return fmt.Errorf("read json: %w", err)
	}

	return nil
}
