package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/nikolalohinski/free-go/types"
)

const (
	actionNotification = "notification"
	actionRegister     = "register"
)

type registerAction struct {
	RequestID types.WebSocketRequestID `json:"request_id,omitempty"`
	Action    string                   `json:"action"`
	Events    []string                 `json:"events"`
}

type registerResponse struct {
	RequestID string          `json:"request_id,omitempty"`
	Action    string          `json:"action"`
	Success   bool            `json:"success"`
	Result    json.RawMessage `json:"result,omitempty"`
	ErrorCode string          `json:"error_code"`
	Message   string          `json:"msg,omitempty"`
}

// ListenEvents implements Client.
func (c *client) ListenEvents(ctx context.Context, events []types.EventDescription) (chan types.Event, error) {
	ws, err := c.webSocket(ctx, "/ws/event")
	if err != nil {
		return nil, fmt.Errorf("websocket connection: %w", err)
	}

	registerActionPayload := registerAction{
		Action: actionRegister,
		Events: make([]string, len(events)),
	}
	for i, event := range events {
		registerActionPayload.Events[i] = fmt.Sprintf("%s_%s", event.Source, event.Name)
	}

	if err := ws.WriteJSON(registerActionPayload); err != nil {
		ws.Close()

		return nil, fmt.Errorf("failed to register action: %w", err)
	}

	var response registerResponse
	if err := ws.ReadJSON(&response); err != nil {
		ws.Close()

		return nil, fmt.Errorf("failed to read register response from websocket: %w", err)
	}

	if !response.Success {
		ws.Close()

		return nil, fmt.Errorf("registering to websocket notifications failed with error %s: %s", response.ErrorCode, response.Message)
	}

	channel := make(chan types.Event, 10)

	go func() {
		var finalErr error

		defer func() {
			if finalErr != nil {
				channel <- types.Event{
					Error: fmt.Errorf("encountered error while handling the event notification: %w", err),
				}
			}

			err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				channel <- types.Event{
					Error: fmt.Errorf("failed to gracefully close websocket connection: %w", err),
				}
			}

			if err = ws.Close(); err != nil {
				channel <- types.Event{
					Error: fmt.Errorf("closing websocket returned an error: %w", err),
				}
			}

			close(channel)
		}()

		for {
			eventPayload, err := waitEventNotification(ctx, ws, actionNotification)
			if err != nil {
				finalErr = fmt.Errorf("wait json response: %w", err)

				return
			}

			if !eventPayload.Success {
				err = fmt.Errorf("received unexpected event payload with success=%t", eventPayload.Success)

				return
			}
			channel <- types.Event{
				Notification: *eventPayload,
			}
		}
	}()

	return channel, nil
}

func waitEventNotification(ctx context.Context, ws *websocket.Conn, action types.WebSocketAction) (*types.WebSocketNotification, error) {
	for {
		message, err := waitJSONResponse[types.WebSocketNotification](ctx, ws)
		if err != nil {
			return nil, err
		}

		if message.Action == action {
			return message, nil
		}
	}
}
