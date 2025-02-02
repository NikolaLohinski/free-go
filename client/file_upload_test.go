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
		server        *ghttp.Server
		sessionToken  string

		uploadContext context.Context

		returnedWriter io.WriteCloser
		returnedErr    error

		content []byte
		size    int
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		DeferCleanup(server.Close)

		freeboxClient = Must(client.New(server.Addr(), version)).
			WithAppID(appID).
			WithPrivateToken(privateToken)

		sessionToken = setupLoginFlow(server)

		content = []byte("data")
		size = len(content)

		uploadContext = context.Background()
	})

	JustBeforeEach(func() {
		ctx, cancel := context.WithTimeout(uploadContext, 10*time.Second)
		DeferCleanup(cancel)

		returnedWriter, returnedErr = freeboxClient.FileUploadStart(ctx, types.FileUploadStartActionInput{
			Dirname:  "dir",
			Filename: "the-file",
			Force:    types.FileUploadStartActionForceOverwrite,
			Size:     size,
		})
	})

	Describe("FileUploadStart", func() {
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						verifyAuth(sessionToken),
						wsHandler(func(ws *websocket.Conn) {
							ws.SetCloseHandler(func(code int, text string) error {
								defer GinkgoRecover()

								// Ignore default close handler to explicitly assert the websocket is closed correctly
								Expect(code).To(Equal(websocket.CloseNormalClosure))
								Expect(text).To(Equal(websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))

								_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

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

							Expect(writeJSON(ws, &types.FileUploadStartResponse{
								Success:   true,
								Action:    types.FileUploadStartActionNameUploadStart,
								RequestID: result.RequestID,
							})).To(Succeed())

							received, err := readChunk(ws)
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
						}),
					),
				)
			})

			It("should start a file upload", func() {
				Expect(returnedErr).To(BeNil())
				Expect(returnedWriter).ToNot(BeNil())

				Expect(returnedWriter.Write(content)).To(BeEquivalentTo(size))

				Expect(returnedWriter.Close()).To(Succeed())
			})
		})

		Context("File already exist", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						verifyAuth(sessionToken),
						wsHandler(func(ws *websocket.Conn) {
							ws.SetCloseHandler(func(code int, text string) error {
								defer GinkgoRecover()

								// Ignore default close handler to explicitly assert the websocket is closed correctly
								Expect(code).To(Equal(websocket.CloseNormalClosure))
								Expect(text).To(Equal(websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))

								_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

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

							Expect(writeJSON(ws, &types.FileUploadStartResponse{
								Success:   false,
								Action:    types.FileUploadStartActionNameUploadStart,
								RequestID: result.RequestID,
								FileSize:  size,
								ErrorCode: "conflict",
								Message:   "Le fichier existe déjà",
							})).To(Succeed())
						}),
					),
				)
			})

			It("conflicts", func() {
				Expect(returnedErr).To(MatchError(&types.WebSocketResponseError{
					ErrorCode: "conflict",
					Message:   "Le fichier existe déjà",
				}))
			})
		})

		Context("with non JSON message", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						verifyAuth(sessionToken),
						wsHandler(func(ws *websocket.Conn) {
							ws.SetCloseHandler(func(code int, text string) error {
								defer GinkgoRecover()

								// Ignore default close handler to explicitly assert the websocket is closed correctly
								Expect(code).To(Equal(websocket.CloseNormalClosure))
								Expect(text).To(Equal(websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))

								_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

								return nil
							})

							ws.NetConn().Write([]byte("invalid message"))
						}),
					),
				)
			})

			It("fails", func() {
				Expect(returnedErr).To(HaveOccurred())
			})
		})

		Context("Websocket close early", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						verifyAuth(sessionToken),
						wsHandler(func(ws *websocket.Conn) {
							// Do nothing and close the connection
						}),
					),
				)
			})

			It("propagate the error", func() {
				Expect(returnedErr).To(MatchError(Or(ContainSubstring("connection reset by peer"), ContainSubstring("unexpected EOF"))))
			})
		})

		Context("close error", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						verifyAuth(sessionToken),
						wsHandler(func(ws *websocket.Conn) {
							ws.SetCloseHandler(func(code int, text string) error {
								defer GinkgoRecover()

								// Ignore default close handler to explicitly assert the websocket is closed correctly
								Expect(code).To(Equal(websocket.CloseNormalClosure))
								Expect(text).To(Equal(websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))

								_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "internal error"))

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

							Expect(writeJSON(ws, &types.FileUploadStartResponse{
								Success:   true,
								Action:    types.FileUploadStartActionNameUploadStart,
								RequestID: result.RequestID,
							})).To(Succeed())

							received, err := readChunk(ws)
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
						}),
					),
				)
			})

			It("should start a file upload", func() {
				Expect(returnedErr).To(BeNil())
				Expect(returnedWriter).ToNot(BeNil())

				Expect(returnedWriter.Write(content)).To(BeEquivalentTo(size))

				Expect(returnedWriter.Close()).To(Succeed())
			})
		})

		Context("upload error", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						verifyAuth(sessionToken),
						wsHandler(func(ws *websocket.Conn) {
							ws.SetCloseHandler(func(code int, text string) error {
								defer GinkgoRecover()

								// Ignore default close handler to explicitly assert the websocket is closed correctly
								Expect(code).To(Equal(websocket.CloseNormalClosure))
								Expect(text).To(Equal(websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))

								_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

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

							Expect(writeJSON(ws, &types.FileUploadStartResponse{
								Success:   true,
								Action:    types.FileUploadStartActionNameUploadStart,
								RequestID: result.RequestID,
							})).To(Succeed())

							received, err := readChunk(ws)
							Expect(err).ToNot(HaveOccurred())
							Expect(string(received)).To(Equal(string(content)))

							Expect(writeJSON(ws, &types.WebSocketResponse[types.FileUploadChunkResponse]{
								Success:   false,
								Action:    types.FileUploadStartActionNameUploadData,
								RequestID: result.RequestID,
								Result: types.FileUploadChunkResponse{
									TotalLen:  len(received) - 1,
									Complete:  false,
									Cancelled: false,
								},
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
						}),
					),
				)
			})

			It("propagate the error", func() {
				Expect(returnedErr).To(Succeed())
				_, err := returnedWriter.Write(content)
				Expect(err).To(MatchError(ContainSubstring("test_error: no space left on device")))

				Expect(returnedWriter.Close()).To(Succeed())
			})
		})

		Context("ignore message with a different requestID", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						verifyAuth(sessionToken),
						wsHandler(func(ws *websocket.Conn) {
							ws.SetCloseHandler(func(code int, text string) error {
								defer GinkgoRecover()

								// Ignore default close handler to explicitly assert the websocket is closed correctly
								Expect(code).To(Equal(websocket.CloseNormalClosure))
								Expect(text).To(Equal(websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))

								_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

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

							Expect(writeJSON(ws, &types.FileUploadStartResponse{
								Success:   true,
								Action:    types.FileUploadStartActionNameUploadStart,
								RequestID: result.RequestID,
							})).To(Succeed())

							received, err := readChunk(ws)
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
						}),
					),
				)
			})

			It("ignore the message", func() {
				Expect(returnedErr).To(Succeed())
				Expect(returnedWriter).ToNot(BeNil())

				Expect(returnedWriter.Write(content)).To(BeEquivalentTo(size))

				Expect(returnedWriter.Close()).To(Succeed())
			})
		})
	})
})

func writeJSON(ws *websocket.Conn, data interface{}) error {
	w, err := ws.NextWriter(websocket.TextMessage)
	if err != nil {
		return fmt.Errorf("next writer: %w", err)
	}

	if err := json.NewEncoder(gbytes.TimeoutWriter(w, time.Second)).Encode(data); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
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

func readChunk(ws *websocket.Conn) ([]byte, error) {
	messageType, content, err := ws.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("next reader: %w", err)
	}

	if messageType != websocket.BinaryMessage {
		return nil, fmt.Errorf("unexpected message type: %d", messageType)
	}

	return content, nil
}

func wsHandler(predicate func(ws *websocket.Conn)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func(ws *websocket.Conn) {
			Expect(ws.Close()).To(Succeed())
		}(ws)

		predicate(ws)
	}
}
