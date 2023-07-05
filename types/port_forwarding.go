package types

type ipProtocol = string

const (
	TCP ipProtocol = "tcp"
	UDP ipProtocol = "udp"
)

type PortForwardingRulePayload struct {
	Enabled      *bool      `json:"enabled,omitempty"`
	IPProtocol   ipProtocol `json:"ip_proto,omitempty"`
	WanPortStart int64      `json:"wan_port_start,omitempty"`
	WanPortEnd   int64      `json:"wan_port_end,omitempty"`
	LanIP        string     `json:"lan_ip,omitempty"`
	LanPort      int64      `json:"lan_port,omitempty"`
	SourceIP     string     `json:"src_ip,omitempty"`
	Comment      string     `json:"comment,omitempty"`
}

type PortForwardingRule struct {
	PortForwardingRulePayload
	ID       int64            `json:"id"`
	Valid    bool             `json:"valid"`
	Hostname string           `json:"hostname"`
	Host     LanInterfaceHost `json:"host"`
}
