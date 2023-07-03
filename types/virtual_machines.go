package types

import (
	"encoding/base64"
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

type diskType string

const (
	RawDisk   diskType = "raw"
	QCow2Disk diskType = "qcow2"
)

type os string

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

type machineStatus string

const (
	StoppedStatus  machineStatus = "stopped"
	RunningStatus  machineStatus = "running"
	StartingStatus machineStatus = "starting"
	StoppingStatus machineStatus = "stopping"
)

type VirtualMachine struct {
	ID                int64         `json:"id"`
	Name              string        `json:"name"`
	Mac               string        `json:"mac"`
	DiskPath          string        `json:"disk_path"`
	DiskType          diskType      `json:"disk_type"`
	CDPath            CDPath        `json:"cd_path"` // Base64 encoded
	Memory            int64         `json:"memory"`
	OS                os            `json:"os"`
	VCPUs             int64         `json:"vcpus"`
	Status            machineStatus `json:"status"`
	EnableScreen      bool          `json:"enable_screen"`
	BindUSBPorts      BindUSBPorts  `json:"bind_usb_ports"` // Empty string returned if no binds defined
	EnableCloudInit   bool          `json:"enable_cloudinit"`
	CloudInitUserData string        `json:"cloudinit_userdata"`
	CloudHostName     string        `json:"cloudinit_hostname"`
}

type CDPath string

func (c *CDPath) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal cd_path: %w", err)
	}

	result, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return fmt.Errorf("failed to decode '%s' from base64: %w", raw, err)
	}

	*c = CDPath(string(result))

	return nil
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
