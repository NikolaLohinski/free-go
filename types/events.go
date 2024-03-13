package types

import "encoding/json"

type EventDescription struct {
	Source string
	Name   string
}

type Event struct {
	Notification EventNotification
	Error        error
}

type EventNotification struct {
	Action  string          `json:"action"`
	Success bool            `json:"success"`
	Source  string          `json:"source"`
	Event   string          `json:"event"`
	Result  json.RawMessage `json:"result"`
}
