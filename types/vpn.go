package types

// OpenVPNServerConfig is the OpenVPN server configuration.
// Endpoint: GET/PUT /vpn/openvpn/
// Note: field names are from the Freebox OS API — verify against a live Freebox with
// DevTools if any mismatch is observed.
type OpenVPNServerConfig struct {
	Enabled       bool   `json:"enabled"`                  // Whether the OpenVPN server is enabled
	ServerPort    int64  `json:"server_port"`              // UDP port the server listens on (default 1194)
	ServerIP      string `json:"server_ip"`                // VPN subnet IP address (e.g. "10.8.0.0")
	ServerMask    string `json:"server_mask"`              // VPN subnet mask (e.g. "255.255.255.0")
	PushDefaultGW bool   `json:"push_default_gw"`          // Whether to push the default gateway to clients
	PushDHCP      bool   `json:"push_dhcp"`                // Whether to push DHCP settings to clients
	CA            string `json:"ca,omitempty"`             // CA certificate in PEM format (read-only, set by Freebox)
	Cert          string `json:"cert,omitempty"`           // Server certificate in PEM format (read-only, set by Freebox)
}

// VPNUserPayload is the create/update payload for a VPN user.
// Endpoint: POST /vpn/user/, PUT /vpn/user/{login}
type VPNUserPayload struct {
	Login       string `json:"login"`                    // Username (immutable after creation)
	Password    string `json:"password"`                 // Password
	Description string `json:"description,omitempty"`   // Optional description
}

// VPNUser is a VPN user account as returned by the API.
type VPNUser struct {
	VPNUserPayload
}
