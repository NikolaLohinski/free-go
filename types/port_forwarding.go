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
	Comment      string     `json:"comment"`                  // comment

	// The following fields are read-only/computed on reads (and shadowed by
	// PortForwardingRule's own fields below for that purpose), but the Freebox
	// API rejects fw/redir/ writes that omit them: creating a rule expects
	// zero-valued placeholders, while updating one requires the rule's current
	// host binding to be echoed back verbatim. Host is an interface{} because
	// it must marshal as an empty string on create but as the full host object
	// on update.
	ID       int64  `json:"id"`
	Hostname string `json:"hostname"`
	Host     any    `json:"host"`
	Valid    bool   `json:"valid"`
}

type PortForwardingRule struct {
	PortForwardingRulePayload
	ID       int64             `json:"id"`       // forwarding id
	Host     *LanInterfaceHost `json:"host"`     // forwarding target host information
	Hostname string            `json:"hostname"` // forwarding target host name

	// Undocumented fields

	Valid bool `json:"valid"`
}
