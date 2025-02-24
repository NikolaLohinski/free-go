package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/nikolalohinski/free-go/types"
)

func (c *client) webSocket(ctx context.Context, endpoint string) (*websocket.Conn, error) {
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

	url.Path = url.Path + endpoint

	ws, dialResponse, err := websocket.DefaultDialer.DialContext(ctx, url.String(), header)
	if err != nil {
		return nil, fmt.Errorf("dialing websocket returned a status %s: %w", dialResponse.Status, err)
	}

	return ws, nil
}

func waitWebSocketResponse[R interface{}](ctx context.Context, ws *websocket.Conn, requestID types.WebSocketRequestID, action types.WebSocketAction) (*R, error) {
	for {
		message, err := waitJSONResponse[types.WebSocketResponse[R]](ctx, ws)
		if err != nil {
			return nil, err
		}

		if (message.RequestID == requestID) && message.Action == action {
			if err := message.GetError(); err != nil {
				return nil, fmt.Errorf("response: %w", err)
			}

			return &message.Result, nil
		}
	}
}

func waitJSONResponse[T interface{}](ctx context.Context, ws *websocket.Conn) (*T, error) {
	content, err := waitResponse(ctx, ws, websocket.TextMessage)
	if err != nil {
		return nil, err
	}

	result := new(T)

	if err := json.Unmarshal(content, result); err != nil {
		var wsResponse types.WebSocketResponse[interface{}]
		if json.Unmarshal(content, &wsResponse) == nil {
			if err := wsResponse.GetError(); err != nil {
				return nil, fmt.Errorf("unexpected error from server: %w", wsResponse.GetError())
			}
		}

		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}

func waitResponse(ctx context.Context, ws *websocket.Conn, expectedType int) ([]byte, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("cancelled: %w", ctx.Err())
		default:
			messageType, data, err := ws.ReadMessage()
			if err != nil {
				if expectedType == websocket.CloseMessage && websocket.IsUnexpectedCloseError(err) {
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
