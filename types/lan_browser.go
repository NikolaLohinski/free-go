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
	Workstation      hostType = "workstation"       // Workstation
	Laptop           hostType = "laptop"            // Laptop
	Smartphone       hostType = "smartphone"        // Smartphone
	Tablet           hostType = "tablet"            // Tablet
	Printer          hostType = "printer"           // Printer
	VGConsole        hostType = "vg_console"        // Video game console
	Television       hostType = "television"        // Televison
	NAS              hostType = "nas"               // Network-attached storage
	IPCamera         hostType = "ip_camera"         // IP camera
	IPPhone          hostType = "ip_phone"          // IP phone
	FreeboxPlayer    hostType = "freebox_player"    // Freebox Player
	FreeboxHD        hostType = "freebox_hd"        // Freebox HD
	FreeboxCrystal   hostType = "freebox_crystal"   // Freebox Crystal
	FreeboxMini      hostType = "freebox_mini"      // Freebox Mini
	FreeboxDelta     hostType = "freebox_delta"     // Freebox Delta
	FreeboxOne       hostType = "freebox_one"       // Freebox One
	FreeboxWIFI      hostType = "freebox_wifi"      // Freebox WIFI
	FreboxPop        hostType = "freebox_pop"       // Frebox Pop
	NetworkingDevice hostType = "networking_device" // Networking device
	MultimediaDevice hostType = "multimedia_device" // Multimedia device
	Car              hostType = "car"               // Connected car
	Watch            hostType = "watch"             // Smartwatch
	Light            hostType = "light"             // Light
	Outlet           hostType = "outlet"            // Connected outlet
	Appliances       hostType = "appliances"        // Household appliances
	Thermostat       hostType = "thermostat"        // Thermostat
	Shutter          hostType = "shutter"           // Electric shutter
	Other            hostType = "other"             // Other
)

type LanInterfaceHost struct {
	Active            bool                    `json:"active"`              // If true the host sends traffic to the Freebox
	Persistent        bool                    `json:"persistent"`          // If true the host is always shown even if it has not been active since the Freebox startup
	Reachable         bool                    `json:"reachable"`           // If true the host can receive traffic from the Freebox
	PrimaryNameManual bool                    `json:"primary_name_manual"` // If true the primary name has been set manually
	VendorName        string                  `json:"vendor_name"`         // Host vendor name (from the mac address)
	Type              hostType                `json:"host_type"`           // When possible, the Freebox will try to guess the host_type, but you can manually override this to the correct value
	ID                string                  `json:"id"`                  // Host id (unique on this interface)
	LastTimeReachable Timestamp               `json:"last_time_reachable"` // Last time the host was reached
	FirstActivity     Timestamp               `json:"first_activity"`      // First time the host sent traffic, or 0 (Unix Epoch) if it wasn’t seen before this field was added.
	LastActivity      Timestamp               `json:"last_activity"`       // Last time the host sent traffic
	PrimaryName       string                  `json:"primary_name"`        // Host primary name (chosen from the list of available names, or manually set by user)
	L2Ident           L2Ident                 `json:"l2ident"`             // Layer 2 network id and its type
	Names             []HostName              `json:"names"`               // List of available names, and their source
	L3Connectivities  []LanHostL3Connectivity `json:"l3connectivities"`    // List of available layer 3 network connections
	NetworkControl    *LanHostNetworkControl  `json:"network_control"`     // If device is associated with a profile, contains profile summary.

	// Undocumented fields

	AccessPoint *LanHostAccessPoint `json:"access_point"`
	DefaultName string              `json:"default_name"`
	Interface   string              `json:"interface"`
	Model       string              `json:"model"`
}

type LanHostAccessPoint struct {
	RXBytes          int64                              `json:"rx_bytes"`
	TXBytes          int64                              `json:"tx_bytes"`
	RXRate           int64                              `json:"rx_rate"`
	TXRate           int64                              `json:"tx_rate"`
	UID              string                             `json:"uid"`
	ConnectivityType LanHostAccessPointConnectivityType `json:"connectivity_type"`
	Type             string                             `json:"type"`
	MacAddress       string                             `json:"mac"`

	WifiInformation     *LanHostAccessPointWifiInformation     `json:"wifi_information"`
	EthernetInformation *LanHostAccessPointEthernetInformation `json:"ethernet_information"`
}

type LanHostAccessPointWifiInformation struct {
	Band            string `json:"band"`
	SessionDuration int64  `json:"sess_duration"`
	PhyRXRate       int64  `json:"phy_rx_rate"`
	SSID            string `json:"ssid"`
	Standard        string `json:"standard"`
	BSSID           string `json:"bssid"`
	PhyTXRate       int64  `json:"phy_tx_rate"`
	Signal          int64  `json:"signal"`
}

type LanHostAccessPointEthernetInformation struct {
	Duplex       string `json:"duplex"`
	Speed        int64  `json:"speed"`
	MaxPortSpeed int64  `json:"max_port_speed"`
	Link         string `json:"link"`
}

type LanHostAccessPointConnectivityType = string

const (
	WifiLanHostAccessPointConnectivityType     = "wifi"
	EthernetLanHostAccessPointConnectivityType = "ethernet"
)

type HostName struct {
	Name   string `json:"name"`   // Host name
	Source string `json:"source"` // source of the name
}

type l2IdentType = string

const (
	DHCP    l2IdentType = "dhcp"     // DHCP
	NetBios l2IdentType = "netbios"  // Netbios
	MDNS    l2IdentType = "mdns"     // mDNS hostname
	MDNSSRV l2IdentType = "mdns_srv" // mDNS service
	UPNP    l2IdentType = "upnp"     // UPnP
	WSD     l2IdentType = "wsd"      // WS-Discovery

	// Undocumented fields.

	MacAddress l2IdentType = "mac_address" // MAC
)

type L2Ident struct {
	ID   string      `json:"id"`   // Layer 2 id
	Type l2IdentType `json:"type"` // Type of layer 2 address
}

type lanHostL3ConnectivityAF = string

const (
	IPV4 lanHostL3ConnectivityAF = "ipv4" // IPv4 address family
	IPV6 lanHostL3ConnectivityAF = "ipv6" // IPv6 address family
)

type LanHostL3Connectivity struct {
	Address           string                  `json:"addr"`                // Layer 3 address
	Active            bool                    `json:"active"`              // is the connection active
	Reachable         bool                    `json:"reachable"`           // is the connection reachable
	LastActivity      Timestamp               `json:"last_activity"`       // last activity timestamp
	LastTimeReachable Timestamp               `json:"last_time_reachable"` // last reachable timestamp
	Type              lanHostL3ConnectivityAF `json:"af"`                  // address family
	Model             string                  `json:"model"`               // device model if known
}

type LanHostNetworkControl struct {
	ProfileID   int    `json:"profile_id"`   // Id of profile this device is associated with.
	Name        string `json:"name"`         // Name of profile this device is associated with.
	CurrentMode string `json:"current_mode"` // Mode described in Network Control Object
}
