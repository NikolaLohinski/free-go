package types

import "encoding/json"

type EventDescription struct {
	Source eventSource
	Name   eventName
}

type Event struct {
	Notification EventNotification
	Error        error
}

type EventNotification struct {
	Action  string          `json:"action"`
	Success bool            `json:"success"`
	Source  eventSource     `json:"source"`
	Event   eventName       `json:"event"`
	Result  json.RawMessage `json:"result"`
}

type (
	eventSource string
	eventName   string
)
