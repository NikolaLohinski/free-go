package types

import (
	"encoding/json"
	"fmt"
)

type EventDescription struct {
	Source eventSource
	Name   eventName
}

func (d *EventDescription) String() string {
	return string(d.Source) + "_" + string(d.Name)
}

type Event struct {
	Notification WebSocketNotification
	Error        error
}

// Deprecated: use WebSocketNotification instead.
type EventNotification = WebSocketNotification

type WebSocketNotification struct {
	Action  WebSocketAction `json:"action"`
	Success bool            `json:"success"`
	Source  eventSource     `json:"source"`
	Event   eventName       `json:"event"`
	Result  json.RawMessage `json:"result"`
}

func (n *WebSocketNotification) VmStateChange() (*VmStateChange, error) {
	if !n.Action.Is(&EventDescription{
		Source: EventSourceVM,
		Name:   EventStateChanged,
	}) {
		return nil, fmt.Errorf("unexpected action: %s", n.Action)
	}

	result := new(VmStateChange)
	if err := json.Unmarshal(n.Result, result); err != nil {
		return nil, fmt.Errorf("unmarshal VmStateChange: %w", err)
	}

	return result, nil
}

func (n *WebSocketNotification) VmDiskTask() (*VmDiskTask, error) {
	if !n.Action.Is(&EventDescription{
		Source: EventSourceVMDisk,
		Name:   EventDiskTaskDone,
	}) {
		return nil, fmt.Errorf("unexpected action: %s", n.Action)
	}

	result := new(VmDiskTask)
	if err := json.Unmarshal(n.Result, result); err != nil {
		return nil, fmt.Errorf("unmarshal VmDiskTask: %w", err)
	}

	return result, nil
}

func (n *WebSocketNotification) LanHost() (*LanHost, error) {
	if !n.Action.Is(&EventDescription{
		Source: EventSourceLANHost,
		Name:   EventHostL3AddrReachable,
	}) && !n.Action.Is(&EventDescription{
		Source: EventSourceLANHost,
		Name:   EventHostL3AddrUnreachable,
	}) {
		return nil, fmt.Errorf("unexpected action: %s", n.Action)
	}

	result := new(LanHost)
	if err := json.Unmarshal(n.Result, result); err != nil {
		return nil, fmt.Errorf("unmarshal LanHost: %w", err)
	}

	return result, nil
}

type (
	eventSource string
	eventName   string
	eventAction string
)

type (
	VmStateChange struct {
		ID     int           `json:"id"`     // VM id.
		Status machineStatus `json:"status"` // New VM.status.
	}

	VmDiskTask struct {
		ID    int                        `json:"id"`              // Task id.
		Type  virtualMachineDiskTaskType `json:"type"`            // Type of disk operation.
		Done  bool                       `json:"done"`            // Is task done
		Error string                     `json:"error,omitempty"` // Error message if task failed.
	}

	LanHost struct {
		ID                int                   `json:"id"`                  // Host id (unique on this interface).
		PrimaryName       string                `json:"primary_name"`        // Host primary name (chosen from the list of available names, or manually set by user).
		HostType          hostType              `json:"host_type"`           // When possible, the Freebox will try to guess the host_type, but you can manually override this to the correct value.
		PrimaryNameManual bool                  `json:"primary_name_manual"` // If true the primary name has been set manually.
		L2Ident           []L2Ident             `json:"l2ident"`             // Layer 2 network id and its type
		VendorName        string                `json:"vendor_name"`         // Host vendor name (from the mac address)
		Persistent        bool                  `json:"persistent"`          // If true the host is always shown even if it has not been active since the Freebox startup
		Reachable         bool                  `json:"reachable"`           // If true the host can receive traffic from the Freebox
		LastTimeReachable Timestamp             `json:"last_time_reachable"` // Last time the host was reached
		Active            bool                  `json:"active"`              // If true the host sends traffic to the Freebox
		LastActivity      Timestamp             `json:"last_activity"`       // Last time the host sent traffic
		FirstActivity     Timestamp             `json:"first_activity"`      // First time the host sent traffic, or 0 (Unix Epoch) if it wasn’t seen before this field was added.
		Names             []HostName            `json:"names"`               // List of names associated with this host
		L3Connectivities  []L3Connectivity      `json:"l3connectivities"`    // List of available layer 3 network connections
		NetworkControl    LanHostNetworkControl `json:"network_control"`     // If device is associated with a profile, contains profile summary.
	}
)
