package types

type LanInfo struct {
	Name      string `json:"name"`
	HostCount int    `json:"host_count"`
}

const (
	EventSourceLANHost eventSource = "lan_host"

	EventHostL3AddrReachable   eventName = "l3addr_reachable"
	EventHostL3AddrUnreachable eventName = "l3addr_unreachable"
)

type hostType = string

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

type LanInterfaceHost struct {
	Active            bool                   `json:"active"`
	Persistent        bool                   `json:"persistent"`
	Reachable         bool                   `json:"reachable"`
	PrimaryNameManual bool                   `json:"primary_name_manual"`
	VendorName        string                 `json:"vendor_name"`
	Type              hostType               `json:"host_type"`
	Interface         string                 `json:"interface"`
	ID                string                 `json:"id"`
	LastTimeReachable Timestamp              `json:"last_time_reachable"`
	FirstActivity     Timestamp              `json:"first_activity"`
	LastActivity      Timestamp              `json:"last_activity"`
	PrimaryName       string                 `json:"primary_name"`
	DefaultName       string                 `json:"default_name"`
	L2Ident           L2Ident                `json:"l2ident"`
	Names             []HostName             `json:"names"`
	L3Connectivities  []L3Connectivity       `json:"l3connectivities"`
	NetworkControl    *LanHostNetworkControl `json:"network_control"`
}

type HostName struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

type l2IdentType = string

const (
	DHCP       l2IdentType = "dhcp"
	NetBios    l2IdentType = "netbios"
	MDNS       l2IdentType = "mdns"
	MDNSSRV    l2IdentType = "mdns_srv"
	UPNP       l2IdentType = "upnp"
	WSD        l2IdentType = "wsd"
	MacAddress l2IdentType = "mac_address"
)

type L2Ident struct {
	ID   string      `json:"id"`
	Type l2IdentType `json:"type"`
}

type af = string

const (
	IPV4 af = "ipv4"
	IPV6 af = "ipv6"
)

type L3Connectivity struct {
	Address           string    `json:"addr"`
	Active            bool      `json:"active"`
	Reachable         bool      `json:"reachable"`
	LastActivity      Timestamp `json:"last_activity"`
	LastTimeReachable Timestamp `json:"last_time_reachable"`
	Type              af        `json:"af"`
}

type LanHostNetworkControl struct {
	ProfileID   int    `json:"profile_id"`   // Id of profile this device is associated with.
	Name        string `json:"name"`         // Name of profile this device is associated with.
	CurrentMode string `json:"current_mode"` // Mode described in Network Control Object
}
