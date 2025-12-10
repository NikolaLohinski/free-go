package types

// LanMode represents the LAN operating mode
type LanMode = string

const (
	LanModeRouter LanMode = "router" // The Freebox acts as a network router
	LanModeBridge LanMode = "bridge" // The Freebox acts as a network bridge
)

// LanConfig represents the Freebox Server LAN configuration
// https://dev.freebox.fr/sdk/os/lan/#lan-config
type LanConfig struct {
	IP          string  `json:"ip"`           // Freebox Server IPv4 address
	Name        string  `json:"name"`         // Freebox Server name
	NameDNS     string  `json:"name_dns"`     // Freebox Server DNS name
	NameMDNS    string  `json:"name_mdns"`    // Freebox Server mDNS name
	NameNetBIOS string  `json:"name_netbios"` // Freebox Server netbios name
	Mode        LanMode `json:"mode"`         // LAN mode (router or bridge)
}
