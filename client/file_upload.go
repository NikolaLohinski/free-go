package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nikolalohinski/free-go/types"
)

const codeUploadTaskNotFound = "noent"

func (c *client) FileUploadStart(ctx context.Context, input types.FileUploadStartActionInput) (io.WriteCloser, error) {
	header := http.Header{}
	if err := c.withSession(ctx)(&http.Request{
		Header: header,
	}); err != nil {
		return nil, fmt.Errorf("get a session: %w", err)
	}

	url := *c.base
	url.Scheme = "ws"

	if strings.ToLower(c.base.Scheme) == "https" {
		url.Scheme = "wss"
	}

	url.Path = url.Path + "/ws/upload"

	ws, dialResponse, err := websocket.DefaultDialer.DialContext(ctx, url.String(), header)
	if err != nil {
		return nil, fmt.Errorf("dialing websocket returned a status %s: %w", dialResponse.Status, err)
	}

	ctx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		<-ctx.Done()

		_ = ws.Close()
	}(ctx)

	requestID := newRequestID()

	if err := ws.WriteJSON(&types.FileUploadStartAction{
		Action:    types.FileUploadStartActionNameUploadStart,
		Size:      input.Size,
		RequestID: requestID,
		Dirname:   input.Dirname,
		Filename:  input.Filename,
		Force:     input.Force,
	}); err != nil {
		cancel()

		return nil, fmt.Errorf("upload start action: %w", err)
	}

	response := &types.FileUploadStartResponse{
		Success: true,
	}
	for {
		if err := waitJSONResponse(ctx, ws, response); err != nil {
			cancel()

			return nil, fmt.Errorf("waiting start confirmation: %w", err)
		}

		if response.RequestID != requestID || response.Action != types.FileUploadStartActionNameUploadStart {
			continue
		}

		if err := response.GetError(); err != nil {
			cancel()

			return nil, fmt.Errorf("start confirmation: %w", err)
		}

		break
	}

	// caller should close the writer
	return &ChunkWriter{
		Conn:      ws,
		RequestID: requestID,
		written:   0,
		expected:  input.Size,
		cancel:    cancel,
	}, nil
}

// ListUploadTasks returns a list of upload tasks.
func (c *client) ListUploadTasks(ctx context.Context) ([]types.UploadTask, error) {
	response, err := c.get(ctx, "upload/", c.withSession(ctx))
	if err != nil {
		return nil, fmt.Errorf("GET upload/ endpoint: %w", err)
	}

	if response.Result == nil {
		return nil, nil
	}

	var result []types.UploadTask
	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("upload tasks from generic response: %w", err)
	}

	return result, nil
}

// GetUploadTask returns a upload task by its identifier.
func (c *client) GetUploadTask(ctx context.Context, identifier int64) (*types.UploadTask, error) {
	response, err := c.get(ctx, fmt.Sprintf("upload/%d", identifier), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeTaskNotFound {
			return nil, ErrTaskNotFound
		}

		return nil, fmt.Errorf("GET upload/%d endpoint: %w", identifier, err)
	}

	var result types.UploadTask
	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("upload task from generic response: %w", err)
	}

	return &result, nil
}

// CancelUploadTask cancels a upload task by its identifier.
func (c *client) CancelUploadTask(ctx context.Context, identifier int64) error {
	response, err := c.delete(ctx, fmt.Sprintf("upload/%d/cancel", identifier), nil, c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeTaskNotFound {
			return ErrTaskNotFound
		}

		return fmt.Errorf("DELETE upload/%d/cancel endpoint: %w", identifier, err)
	}

	return nil
}

// DeleteUploadTask deletes a upload task by its identifier.
func (c *client) DeleteUploadTask(ctx context.Context, identifier int64) error {
	response, err := c.delete(ctx, fmt.Sprintf("upload/%d", identifier), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeUploadTaskNotFound {
			return ErrTaskNotFound
		}

		return fmt.Errorf("DELETE upload/%d endpoint: %w", identifier, err)
	}

	return nil
}

// CleanUploadTasks deletes all the FileUpload not in_progress.
func (c *client) CleanUploadTasks(ctx context.Context) error {
	_, err := c.delete(ctx, "upload/clean", c.withSession(ctx))
	if err != nil {
		return fmt.Errorf("DELETE upload/clean endpoint: %w", err)
	}

	return nil
}

func newRequestID() types.UploadRequestID {
	return types.UploadRequestID(time.Now().Nanosecond())
}

type ChunkWriter struct {
	*websocket.Conn

	cancel context.CancelFunc

	RequestID         types.UploadRequestID
	lock              sync.Mutex
	expected, written int
}

var _ io.WriteCloser = (*ChunkWriter)(nil)

func (w *ChunkWriter) Write(data []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if err := w.Conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return 0, fmt.Errorf("write chunk: %w", err)
	}

	responseData, err := waitWebSocketResponse(context.Background(), w.Conn, &types.WebSocketResponse[types.FileUploadChunkResponse]{}, w.RequestID, types.FileUploadStartActionNameUploadData)
	if err != nil {
		return 0, fmt.Errorf("chunk upload confirmation: %w", err)
	}

	written := responseData.TotalLen - w.written
	w.written = responseData.TotalLen

	if responseData.Cancelled {
		return written, errors.New("upload cancelled")
	}

	if written != len(data) {
		return written, io.ErrShortWrite
	}

	return written, nil
}

func (w *ChunkWriter) Close() (finalErr error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	ctx := context.Background()

	defer func(ctx context.Context, conn *websocket.Conn) {
		var errs []error

		defer w.cancel() // cancel the context to stop the goroutine

		errs = append(errs, finalErr)

		if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
			if !errors.Is(err, net.ErrClosed) {
				errs = append(errs, fmt.Errorf("send close message: %w", err))
			}
		}

		_ = conn.Close()

		finalErr = errors.Join(errs...)
	}(ctx, w.Conn)

	switch w.written {
	case w.expected:
		if err := w.Conn.WriteJSON(&types.FileUploadFinalize{
			Action:    types.FileUploadStartActionNameUploadFinalize,
			RequestID: w.RequestID,
		}); err != nil {
			finalErr = fmt.Errorf("finalize upload: %w", err)
			return finalErr
		}

		if _, err := waitWebSocketResponse(ctx, w.Conn, &types.WebSocketResponse[types.FileUploadFinalizeResponse]{}, w.RequestID, types.FileUploadStartActionNameUploadFinalize); err != nil {
			finalErr = fmt.Errorf("finalize upload confirmation: %w", err)
			return finalErr
		}
	default:
		if err := w.Conn.WriteJSON(&types.FileUploadFinalize{
			Action:    types.FileUploadStartActionNameUploadCancel,
			RequestID: w.RequestID,
		}); err != nil {
			finalErr = fmt.Errorf("cancel upload: %w", err)
			return finalErr
		}

		if _, err := waitWebSocketResponse(ctx, w.Conn, &types.WebSocketResponse[types.FileUploadCancelAction]{}, w.RequestID, types.FileUploadStartActionNameUploadCancel); err != nil {
			finalErr = fmt.Errorf("cancel upload confirmation: %w", err)
			return finalErr
		}
	}

	return nil
}

func waitWebSocketResponse[R interface{}](ctx context.Context, ws *websocket.Conn, message *types.WebSocketResponse[R], requestID types.UploadRequestID, action types.WebSocketAction) (*R, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("cancelled: %w", ctx.Err())
		default:
			if err := waitJSONResponse(ctx, ws, message); err != nil {
				return nil, err
			}

			if message.RequestID == requestID && message.Action == action {
				if err := message.GetError(); err != nil {
					return nil, fmt.Errorf("response: %w", err)
				}

				return &message.Result, nil
			}
		}
	}
}

func waitJSONResponse(ctx context.Context, ws *websocket.Conn, message interface{}) error {
	content, err := waitResponse(ctx, ws, websocket.TextMessage)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, message); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}

func waitResponse(ctx context.Context, ws *websocket.Conn, expectedType int) ([]byte, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("cancelled: %w", ctx.Err())
		default:
			messageType, data, err := ws.ReadMessage()
			if err != nil {
				if expectedType == websocket.CloseMessage && isCloseError(err) {
					return data, nil
				}

				return nil, fmt.Errorf("read websocket message: %w", err)
			}

			if messageType == expectedType {
				return data, nil
			}
			if messageType == websocket.CloseMessage {
				return data, fmt.Errorf("websocket closed: %w", errors.New(string(data)))
			}
		}
	}
}

func isCloseError(err error) bool {
	return websocket.IsUnexpectedCloseError(err) ||
		strings.HasSuffix(err.Error(), net.ErrClosed.Error()) || // gorilla package does not return net.ErrClosed
		errors.Is(err, net.ErrClosed) ||
		errors.Is(err, websocket.ErrCloseSent)
}
