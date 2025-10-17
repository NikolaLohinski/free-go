package types

type ipProtocol = string

const (
	TCP ipProtocol = "tcp" // TCP
	UDP ipProtocol = "udp" // UDP
)

type PortForwardingRulePayload struct {
	Enabled      *bool      `json:"enabled,omitempty"` // is forwarding enabled
	IPProtocol   ipProtocol `json:"ip_proto,omitempty"`
	WanPortStart int64      `json:"wan_port_start,omitempty"` // forwarding range start
	WanPortEnd   int64      `json:"wan_port_end,omitempty"`   // forwarding range end
	LanIP        string     `json:"lan_ip,omitempty"`         // forwarding target on LAN
	LanPort      int64      `json:"lan_port,omitempty"`       // forwarding target start port on LAN, (last port is lan_port + wan_port_end - wan_port_start)
	SourceIP     string     `json:"src_ip,omitempty"`         // if src_ip == 0.0.0.0 this rule will apply to any src ip otherwise it will only apply to the specified ip address
	Comment      string     `json:"comment,omitempty"`        // comment
}

type PortForwardingRule struct {
	PortForwardingRulePayload
	ID       int64             `json:"id"`       // forwarding id
	Host     *LanInterfaceHost `json:"host"`     // forwarding target host information
	Hostname string            `json:"hostname"` // forwarding target host name

	// Undocumented fields

	Valid bool `json:"valid"`
}
