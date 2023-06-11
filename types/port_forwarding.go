package types

type ipProtocol string

const (
	TCP ipProtocol = "tcp"
	UDP ipProtocol = "udp"
)

type PortForwardingRulePayload struct {
	Enabled      bool       `json:"enabled"`
	IPProtocol   ipProtocol `json:"ip_proto"`
	WanPortStart int64      `json:"wan_port_start"`
	WanPortEnd   int64      `json:"wan_port_end"`
	LanIP        string     `json:"lan_ip"`
	LanPort      int64      `json:"lan_port"`
	SourceIP     string     `json:"src_ip"`
	Comment      string     `json:"comment"`
}

type PortForwardingRule struct {
	PortForwardingRulePayload `json:",squash"`   //nolint:staticcheck
	ID                        int64              `json:"id"`
	Valid                     bool               `json:"valid"`
	Hostname                  string             `json:"hostname"`
	Host                      PortForwardingHost `json:"host"`
}

type hostType string

const (
	Workstation      hostType = "workstation"
	Laptop           hostType = "laptop"
	Smartphone       hostType = "smartphone"
	Tablet           hostType = "tablet"
	Printer          hostType = "printer"
	VGConsole        hostType = "vg_console"
	Television       hostType = "television"
	NAS              hostType = "nas"
	IPCamera         hostType = "ip_camera"
	IPPhone          hostType = "ip_phone"
	FreeboxPlayer    hostType = "freebox_player"
	FreeboxHD        hostType = "freebox_hd"
	FreeboxCrystal   hostType = "freebox_crystal"
	FreeboxMini      hostType = "freebox_mini"
	FreeboxDelta     hostType = "freebox_delta"
	FreeboxOne       hostType = "freebox_one"
	FreeboxWIFI      hostType = "freebox_wifi"
	FreboxPop        hostType = "freebox_pop"
	NetworkingDevice hostType = "networking_device"
	MultimediaDevice hostType = "multimedia_device"
	Car              hostType = "car"
	Other            hostType = "other"
)

type PortForwardingHost struct {
	Active            bool                           `json:"active"`
	Persistent        bool                           `json:"persistent"`
	Reachable         bool                           `json:"reachable"`
	PrimaryNameManual bool                           `json:"primary_name_manual"`
	VendorName        string                         `json:"vendor_name"`
	Type              hostType                       `json:"host_type"`
	Interface         string                         `json:"interface"`
	ID                string                         `json:"id"`
	LastTimeReachable Timestamp                      `json:"last_time_reachable"`
	FirstActivity     Timestamp                      `json:"first_activity"`
	LastActivity      Timestamp                      `json:"last_activity"`
	PrimaryName       string                         `json:"primary_name"`
	DefaultName       string                         `json:"default_name"`
	L2Ident           PortForwardingL2Ident          `json:"l2ident"`
	Names             []PortForwardingName           `json:"names"`
	L3Connectivities  []PortForwardingL3Connectivity `json:"l3connectivities"`
}

type PortForwardingName struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

type l2IdentType string

const (
	DHCP    l2IdentType = "dhcp"
	NetBios l2IdentType = "netbios"
	MDNS    l2IdentType = "mdns"
	MDNSSRV l2IdentType = "mdns_srv"
	UPNP    l2IdentType = "upnp"
	WSD     l2IdentType = "wsd"
)

type PortForwardingL2Ident struct {
	ID   string      `json:"id"`
	Type l2IdentType `json:"type"`
}

type af string

const (
	IPV4 af = "ipv4"
	IPV6 af = "ipv6"
)

type PortForwardingL3Connectivity struct {
	Address           string    `json:"addr"`
	Active            bool      `json:"active"`
	Reachable         bool      `json:"reachable"`
	LastActivity      Timestamp `json:"last_activity"`
	LastTimeReachable Timestamp `json:"last_time_reachable"`
	Type              af        `json:"af"`
}
