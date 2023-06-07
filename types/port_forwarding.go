package types

type PortForwardingRule struct {
	Enabled      bool               `json:"enabled"`
	Comment      string             `json:"comment"`
	ID           int64              `json:"id"`
	Valid        bool               `json:"valid"`
	SourceIP     string             `json:"src_ip"`
	Hostname     string             `json:"hostname"`
	LanPort      int64              `json:"lan_port"`
	LanIP        string             `json:"lan_ip"`
	ProtoIP      string             `json:"ip_proto"`
	WanPortEnd   int64              `json:"wan_port_end"`
	WanPortStart int64              `json:"wan_port_start"`
	Host         PortForwardingHost `json:"host"`
}

type PortForwardingHost struct {
	Active            bool                           `json:"active"`
	Persistent        bool                           `json:"persistent"`
	Reachable         bool                           `json:"reachable"`
	PrimaryNameManual bool                           `json:"primary_name_manual"`
	VendorName        string                         `json:"vendor_name"`
	Type              string                         `json:"host_type"`
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

type PortForwardingL2Ident struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type PortForwardingL3Connectivity struct {
	Address           string    `json:"addr"`
	Active            bool      `json:"active"`
	Reachable         bool      `json:"reachable"`
	LastActivity      Timestamp `json:"last_activity"`
	LastTimeReachable Timestamp `json:"last_time_reachable"`
	Type              string    `json:"af"`
}
