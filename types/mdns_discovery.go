package types

import "net"

type MDNSDiscovery struct {
	Name           string   // service instance name from PTR record
	Hostname       string   // target hostname from SRV record
	Port           uint16   // port from SRV record
	IPv4           []net.IP // addresses from A records
	IPv6           []net.IP // addresses from AAAA records
	UID            string   // unique device identifier
	DeviceType     string   // e.g. "FreeboxServer1,2"
	APIVersion     string   // e.g. "4.0"
	APIDomain      string   // domain for remote API access, e.g. "xyz.fbxos.fr"
	APIBaseURL     string   // API root path, e.g. "/api/"
	HTTPSPort      int      // port for remote HTTPS access; 0 if not present in TXT record
	HTTPSAvailable bool     // whether HTTPS is available for remote access
}
