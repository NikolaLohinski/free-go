package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/nikolalohinski/free-go/types"
)

const (
	actionNotification = "notification"
	actionRegister     = "register"
)

type registerAction struct {
	RequestID string   `json:"request_id,omitempty"`
	Action    string   `json:"action"`
	Events    []string `json:"events"`
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
	header := http.Header{}
	if err := c.withSession(ctx)(&http.Request{
		Header: header,
	}); err != nil {
		return nil, fmt.Errorf("failed to get a session: %w", err)
	}

	url := *c.base
	url.Scheme = "ws"

	if strings.ToLower(c.base.Scheme) == "https" {
		url.Scheme = "wss"
	}

	url.Path = url.Path + "/ws/event"

	ws, dialResponse, err := websocket.DefaultDialer.Dial(url.String(), header)
	if err != nil {
		return nil, fmt.Errorf("dialing websocket returned a status %s: %w", dialResponse.Status, err)
	}

	registerActionPayload := registerAction{
		Action: actionRegister,
		Events: make([]string, len(events)),
	}
	for i, event := range events {
		registerActionPayload.Events[i] = fmt.Sprintf("%s_%s", event.Source, event.Name)
	}

	if err := ws.WriteJSON(registerActionPayload); err != nil {
		return nil, fmt.Errorf("failed to register action: %w", err)
	}

	var response registerResponse
	if err := ws.ReadJSON(&response); err != nil {
		return nil, fmt.Errorf("failed to read register response from websocket: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("registering to websocket notifications failed with error %s: %s", response.ErrorCode, response.Message)
	}

	channel := make(chan types.Event)

	go func() {
		var err error
		defer func() {
			if err != nil {
				channel <- types.Event{
					Error: fmt.Errorf("encountered error while handling the event notification: %w", err),
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
			select {
			case <-ctx.Done():
				err = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					err = fmt.Errorf("failed to gracefully close websocket connection: %w", err)
				}

				return
			default:
				var eventPayload types.EventNotification
				if err = ws.ReadJSON(&eventPayload); err != nil {
					if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
						err = fmt.Errorf("failed to read message from websocket: %w", err)
					}

					return
				}

				if !eventPayload.Success || eventPayload.Action != actionNotification {
					err = fmt.Errorf("received unexpected event payload with success=%t and action=%s", eventPayload.Success, eventPayload.Action)

					return
				}
				channel <- types.Event{
					Notification: eventPayload,
				}
			}
		}
	}()

	return channel, nil
}
