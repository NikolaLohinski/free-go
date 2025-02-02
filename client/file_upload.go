package client

import (
	"context"
	"errors"
	"fmt"
	"io"
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

	go func() {
		<-ctx.Done()

		_ = ws.Close()
	}()

	requestID := newRequestID()

	if err := ws.WriteJSON(&types.FileUploadStartAction{
		Action:    types.FileUploadStartActionNameUploadStart,
		Size:      input.Size,
		RequestID: requestID,
		Dirname:   input.Dirname,
		Filename:  input.Filename,
		Force:     input.Force,
	}); err != nil {
		return nil, fmt.Errorf("upload start action: %w", err)
	}

	if _, err := waitJSONResponse(ctx, ws, &types.WebSocketResponse[interface{}]{}, requestID, types.FileUploadStartActionNameUploadStart); err != nil {
		return nil, fmt.Errorf("upload start confirmation: %w", err)
	}

	// caller should close the writer
	return &ChunkWriter{
		Conn:      ws,
		RequestID: requestID,
		written:   0,
		expected:  input.Size,
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

	responseData, err := waitJSONResponse(context.Background(), w.Conn, &types.WebSocketResponse[types.FileUploadChunkResponse]{}, w.RequestID, types.FileUploadStartActionNameUploadData)
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

func (w *ChunkWriter) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	switch w.written {
	case w.expected:
		if err := w.Conn.WriteJSON(&types.FileUploadFinalize{
			Action:    types.FileUploadStartActionNameUploadFinalize,
			RequestID: w.RequestID,
		}); err != nil {
			return fmt.Errorf("finalize upload: %w", err)
		}

	default:
		if err := w.Conn.WriteJSON(&types.FileUploadFinalize{
			Action:    types.FileUploadStartActionNameUploadCancel,
			RequestID: w.RequestID,
		}); err != nil {
			return fmt.Errorf("cancel upload: %w", err)
		}

		if _, err := waitJSONResponse(context.Background(), w.Conn, &types.WebSocketResponse[types.FileUploadCancelAction]{}, w.RequestID, types.FileUploadStartActionNameUploadCancel); err != nil {
			return fmt.Errorf("cancel upload confirmation: %w", err)
		}
	}

	if err := w.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		return fmt.Errorf("send close message: %w", err)
	}

	if err := w.Conn.Close(); err != nil {
		return fmt.Errorf("close websocket: %w", err)
	}

	return nil
}

func waitJSONResponse[R interface{}](ctx context.Context, ws *websocket.Conn, message *types.WebSocketResponse[R], requestID types.UploadRequestID, action types.WebSocketAction) (*R, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("cancelled: %w", ctx.Err())
		default:
			if err := ws.ReadJSON(message); err != nil {
				return nil, fmt.Errorf("read websocket response: %w", err)
			}

			if err := message.GetError(); err != nil {
				return nil, fmt.Errorf("upload start: %w", err)
			}

			if message.RequestID == types.WebSocketRequestID(requestID) && message.Action == action {
				return &message.Result, nil
			}
		}
	}
}
