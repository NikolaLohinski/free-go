package types

type DHCPStaticLeaseInfo struct {
	ID       string           `json:"id"`       // DHCP static lease object id
	Mac      string           `json:"mac"`      // Host mac address`
	Comment  string           `json:"comment"`  // an optional comment
	Hostname string           `json:"hostname"` // hostname matching the mac address
	IP       string           `json:"ip"`       // IPv4 to assign to the host
	Host     LanInterfaceHost `json:"host"`     // LAN host information from LAN browser (refer to LanHost documentation)
}

type DHCPStaticLeasePayload struct {
	Mac      string `json:"mac,omitempty"`      // Host mac address`
	Comment  string `json:"comment,omitempty"`  // an optional comment
	Hostname string `json:"hostname,omitempty"` // hostname matching the mac address
	IP       string `json:"ip,omitempty"`       // IPv4 to assign to the host
}
