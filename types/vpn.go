package types

// VPNServerID identifies one of the VPN server configurations exposed by the
// Freebox. Confirmed against the official [VPN Server API] documentation.
//
// [VPN Server API]: https://dev.freebox.fr/sdk/os/vpn/
type VPNServerID string

const (
	VPNServerIDPPTP          VPNServerID = "pptp"
	VPNServerIDOpenVPNRouted VPNServerID = "openvpn_routed"
	VPNServerIDOpenVPNBridge VPNServerID = "openvpn_bridge"
)

// VPNServerType is the underlying protocol of a VPN server configuration.
type VPNServerType string

const (
	VPNServerTypePPTP    VPNServerType = "pptp"
	VPNServerTypeOpenVPN VPNServerType = "openvpn"
	VPNServerTypeIPSec   VPNServerType = "ipsec"
)

// OpenVPNCipher is the cipher used by an OpenVPN server.
type OpenVPNCipher string

const (
	OpenVPNCipherBlowfish OpenVPNCipher = "blowfish"
	OpenVPNCipherAES128   OpenVPNCipher = "aes128"
	OpenVPNCipherAES256   OpenVPNCipher = "aes256"
)

// OpenVPNConfig holds the settings nested under "conf_openvpn". Only present
// when the server type is openvpn.
type OpenVPNConfig struct {
	Cipher          OpenVPNCipher `json:"cipher,omitempty"` // Cipher used by the OpenVPN server
	DisableFragment bool          `json:"disable_fragment"` // Disable the fragment configuration option
	UseTCP          bool          `json:"use_tcp"`          // Use TCP instead of UDP
}

// VPNServerConfig is the configuration of one VPN server.
// Endpoint: GET/PUT /vpn/{id}/config/ where id is one of pptp, openvpn_routed, openvpn_bridge.
type VPNServerConfig struct {
	ID          VPNServerID    `json:"id,omitempty"`           // VPN server id (read-only)
	Type        VPNServerType  `json:"type,omitempty"`         // VPN server type (read-only)
	Enabled     bool           `json:"enabled"`                // Whether the VPN server is enabled
	EnableIPv4  bool           `json:"enable_ipv4"`            // Enable IPv4 (not relevant for openvpn_bridge and pptp)
	EnableIPv6  bool           `json:"enable_ipv6"`            // Enable IPv6 (not relevant for openvpn_bridge and pptp)
	Port        int64          `json:"port"`                   // Server port (only editable when type is openvpn)
	MinPort     int64          `json:"min_port,omitempty"`     // Read-only lower bound tied to the connection's ipv4_port_range
	MaxPort     int64          `json:"max_port,omitempty"`     // Read-only upper bound tied to the connection's ipv4_port_range
	ConfOpenVPN *OpenVPNConfig `json:"conf_openvpn,omitempty"` // OpenVPN-specific settings, only available when type is openvpn
	IPStart     string         `json:"ip_start,omitempty"`     // Read-only IPv4 pool range start for clients
	IPEnd       string         `json:"ip_end,omitempty"`       // Read-only IPv4 pool range end for clients
	IP6Start    string         `json:"ip6_start,omitempty"`    // Read-only IPv6 pool range start for clients
	IP6End      string         `json:"ip6_end,omitempty"`      // Read-only IPv6 pool range end for clients
}

// VPNUserPayload is the create/update payload for a VPN user.
// Endpoint: POST /vpn/user/, PUT /vpn/user/{login}
type VPNUserPayload struct {
	Login         string `json:"login"`                    // Username (immutable after creation)
	Password      string `json:"password,omitempty"`       // Password (8-32 chars), write-only
	IPReservation string `json:"ip_reservation,omitempty"` // Optional reserved IP for this VPN user
}

// VPNUser is a VPN user account as returned by the API.
type VPNUser struct {
	Login         string `json:"login"`                    // Username
	PasswordSet   bool   `json:"password_set,omitempty"`   // Whether a password has been set (read-only)
	IPReservation string `json:"ip_reservation,omitempty"` // Reserved IP for this VPN user, if any
}
