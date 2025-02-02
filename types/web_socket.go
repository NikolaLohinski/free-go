package types

import "fmt"

type WebSocketAction string
type WebSocketRequestID int64

// https://dev.freebox.fr/sdk/os/#WebSocketResponse
type WebSocketResponse[R interface{}] struct {
	RequestID WebSocketRequestID `json:"request_id,omitempty"` // If you set a request_id in your WebSocketRequest, the same request_id will be returned in the associated response
	Action    WebSocketAction    `json:"action,omitempty"`     // The action of the response. (It may be omitted if the request does not expect any action)
	Success   bool               `json:"success"`              // Indicates if the request was successful
	Result    R                  `json:"result,omitempty"`     // The result of the request. (It may be omitted if the request does not expect any result)
	ErrorCode string             `json:"error_code"`           // In case of request error, this error_code provides information about the error. The possible error_code values are documented for each API.
	Message   string             `json:"msg"`                  // In cas of error, provides a French error message relative to the error
}

func (r *WebSocketResponse[R]) GetError() error {
	if r.Success {
		return nil
	}
	return &WebSocketResponseError{
		ErrorCode: r.ErrorCode,
		Message:   r.Message,
	}
}

type WebSocketResponseError struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"msg"`
}

func (e *WebSocketResponseError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Message)
}
