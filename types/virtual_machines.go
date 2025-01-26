package types

import (
	"encoding/json"
	"fmt"
)

type VirtualMachinesInfo struct {
	USBUsed     bool     `json:"usb_used"`
	SATAUsed    bool     `json:"sata_used"`
	SATAPorts   []string `json:"sata_ports"`
	UsedMemory  int64    `json:"used_memory"`
	USBPorts    []string `json:"usb_ports"`
	UsedCPUs    int64    `json:"used_cpus"`
	TotalMemory int64    `json:"total_memory"`
	TotalCPUs   int64    `json:"total_cpus"`
}

type VirtualMachineDistribution struct {
	Hash string `json:"hash"`
	OS   string `json:"os"`
	URL  string `json:"url"`
	Name string `json:"name"`
}

const (
	EventSourceVM eventSource = "vm"

	EventStateChanged eventName = "state_changed"
)

type diskType = string

const (
	RawDisk   diskType = "raw"   // Raw disk data.
	QCow2Disk diskType = "qcow2" // Qcow2 image type. Usually qcow version 3. Note: not all features are supported. In particular, reference to other images is disabled.
)

type os = string

const (
	UnknownOS    os = "unknown"
	FedoraOS     os = "fedora"
	DebianOS     os = "debian"
	UbuntuOS     os = "ubuntu"
	FreebsdOS    os = "freebsd"
	OpensuseOS   os = "opensuse"
	CentosOS     os = "centos"
	JeedomOS     os = "jeedom"
	HomebridgeOS os = "homebridge"
)

type machineStatus = string

const (
	StoppedStatus  machineStatus = "stopped"
	RunningStatus  machineStatus = "running"
	StartingStatus machineStatus = "starting"
	StoppingStatus machineStatus = "stopping"
)

type VirtualMachinePayload struct {
	Name              string       `json:"name,omitempty"`
	DiskPath          Base64Path   `json:"disk_path,omitempty"` // Base64 encoded
	DiskType          diskType     `json:"disk_type,omitempty"`
	CDPath            Base64Path   `json:"cd_path,omitempty"` // Base64 encoded
	Memory            int64        `json:"memory,omitempty"`
	OS                os           `json:"os,omitempty"`
	VCPUs             int64        `json:"vcpus,omitempty"`
	EnableScreen      bool         `json:"enable_screen,omitempty"`
	BindUSBPorts      BindUSBPorts `json:"bind_usb_ports,omitempty"` // Empty string returned if no binds defined
	EnableCloudInit   bool         `json:"enable_cloudinit,omitempty"`
	CloudInitUserData string       `json:"cloudinit_userdata,omitempty"`
	CloudHostName     string       `json:"cloudinit_hostname,omitempty"`
}

type VirtualMachine struct {
	VirtualMachinePayload
	ID     int64         `json:"id"`
	Mac    string        `json:"mac"`
	Status machineStatus `json:"status"`
}

type BindUSBPorts []string

func (b *BindUSBPorts) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal bind_usb_ports: %w", err)
	}

	if value, ok := raw.(string); ok {
		if value != "" {
			return fmt.Errorf("received unexpected content for bind_usb_ports: non empty string: '%s'", value)
		}

		*b = []string{}

		return nil
	}

	if list, ok := raw.([]interface{}); ok {
		value := make([]string, len(list))

		for index, element := range list {
			if cast, ok := element.(string); ok {
				value[index] = cast

				continue
			}

			return fmt.Errorf("received unexpected list for bind_usb_ports: element '%v' is not a string", element)
		}

		*b = value

		return nil
	}

	return fmt.Errorf("received unknown type for bind_usb_ports: neither string nor list of strings: '%s'", string(data))
}
