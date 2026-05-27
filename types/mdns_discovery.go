package types

import "net"

type MDNSDiscovery struct {
	Name     string            // service instance name from PTR record
	Hostname string            // target hostname from SRV record
	Port     uint16            // port from SRV record
	IPv4     []net.IP          // addresses from A records
	IPv6     []net.IP          // addresses from AAAA records
	Info     map[string]string // key=value pairs from TXT record
}
