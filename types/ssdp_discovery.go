package types

type SSDPDiscovery struct {
	Location string // URL to UPnP device description (LOCATION header)
	USN      string // unique service name (USN header)
	Server   string // server string (SERVER header)
}
