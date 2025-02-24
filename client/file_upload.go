package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nikolalohinski/free-go/types"
)

const codeUploadTaskNotFound = "noent"

func (c *client) FileUploadStart(ctx context.Context, input types.FileUploadStartActionInput) (io.WriteCloser, types.UploadRequestID, error) {
	ws, err := c.webSocket(ctx, "/ws/upload")
	if err != nil {
		return nil, 0, fmt.Errorf("websocket connection: %w", err)
	}

	requestID := newRequestID()

	if err := ws.WriteJSON(&types.FileUploadStartAction{
		Action:    types.FileUploadStartActionNameUploadStart,
		Size:      input.Size,
		RequestID: requestID,
		Dirname:   input.Dirname,
		Filename:  input.Filename,
		Force:     input.Force,
	}); err != nil {
		ws.Close()

		return nil, requestID, fmt.Errorf("upload start action: %w", err)
	}

	for {
		response, err := waitJSONResponse[types.FileUploadStartResponse](ctx, ws)
		if err != nil {
			ws.Close()

			return nil, requestID, fmt.Errorf("waiting start confirmation: %w", err)
		}

		if response.RequestID != requestID || response.Action != types.FileUploadStartActionNameUploadStart {
			continue
		}

		if err := response.GetError(); err != nil {
			ws.Close()

			return nil, requestID, fmt.Errorf("start confirmation: %w", err)
		}

		break
	}

	// caller should close the writer
	return &ChunkWriter{
		Conn:      ws,
		RequestID: requestID,
		written:   0,
		expected:  input.Size,
	}, requestID, nil
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
func (c *client) GetUploadTask(ctx context.Context, identifier int64) (result types.UploadTask, err error) {
	response, err := c.get(ctx, fmt.Sprintf("upload/%d", identifier), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeTaskNotFound {
			return result, ErrTaskNotFound
		}

		return result, fmt.Errorf("GET upload/%d endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("upload task from generic response: %w", err)
	}

	return result, nil
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

	responseData, err := waitWebSocketResponse[types.FileUploadChunkResponse](context.Background(), w.Conn, w.RequestID, types.FileUploadStartActionNameUploadData)
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

		errs = append(errs, finalErr)

		if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
			if !errors.Is(err, net.ErrClosed) {
				errs = append(errs, fmt.Errorf("send close message: %w", err))
			}
		}

		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close websocket: %w", err))
		}

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

		if _, err := waitWebSocketResponse[types.FileUploadFinalizeResponse](ctx, w.Conn, w.RequestID, types.FileUploadStartActionNameUploadFinalize); err != nil {
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

		if _, err := waitWebSocketResponse[types.FileUploadCancelAction](ctx, w.Conn, w.RequestID, types.FileUploadStartActionNameUploadCancel); err != nil {
			finalErr = fmt.Errorf("cancel upload confirmation: %w", err)
			return finalErr
		}
	}

	return nil
}
