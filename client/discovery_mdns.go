package client

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/nikolalohinski/free-go/types"
)

const (
	mdnsMulticastGroup = "224.0.0.251"
	mdnsPort           = 5353
	freeboxMDNSService = "_fbx-api._tcp.local."
)

// DiscoverMDNS sends an mDNS PTR query for _fbx-api._tcp.local and collects
// responses until timeout expires or ctx is cancelled.
func DiscoverMDNS(ctx context.Context, timeout time.Duration) ([]types.MDNSDiscovery, error) {
	msg := new(dns.Msg)
	msg.SetQuestion(freeboxMDNSService, dns.TypePTR)
	msg.RecursionDesired = false

	packed, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack mDNS query: %w", err)
	}

	addr := &net.UDPAddr{IP: net.ParseIP(mdnsMulticastGroup), Port: mdnsPort}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to join mDNS multicast group: %w", err)
	}
	defer conn.Close()

	if _, err = conn.WriteTo(packed, addr); err != nil {
		return nil, fmt.Errorf("failed to send mDNS query: %w", err)
	}

	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	entries := map[string]*types.MDNSDiscovery{}
	buf := make([]byte, 65536)

	for {
		select {
		case <-ctx.Done():
			return collectMDNSEntries(entries), ctx.Err()
		default:
		}

		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}

			return collectMDNSEntries(entries), fmt.Errorf("failed to read mDNS response: %w", err)
		}

		response := new(dns.Msg)
		if err := response.Unpack(buf[:n]); err != nil {
			continue
		}

		parseMDNSResponse(response, entries)
	}

	return collectMDNSEntries(entries), nil
}

func parseMDNSResponse(msg *dns.Msg, entries map[string]*types.MDNSDiscovery) {
	allRRs := append(append(msg.Answer, msg.Ns...), msg.Extra...)

	for _, rr := range allRRs {
		switch r := rr.(type) {
		case *dns.PTR:
			if r.Hdr.Name == freeboxMDNSService {
				if _, ok := entries[r.Ptr]; !ok {
					entries[r.Ptr] = &types.MDNSDiscovery{
						Name: strings.TrimSuffix(r.Ptr, "."+freeboxMDNSService),
						Info: map[string]string{},
					}
				}
			}
		case *dns.SRV:
			if e, ok := entries[r.Hdr.Name]; ok {
				e.Hostname = strings.TrimSuffix(r.Target, ".")
				e.Port = r.Port
			}
		case *dns.TXT:
			if e, ok := entries[r.Hdr.Name]; ok {
				for _, txt := range r.Txt {
					if kv := strings.SplitN(txt, "=", 2); len(kv) == 2 {
						e.Info[kv[0]] = kv[1]
					}
				}
			}
		}
	}

	for _, rr := range allRRs {
		name := strings.TrimSuffix(rr.Header().Name, ".")
		for _, e := range entries {
			if e.Hostname != name {
				continue
			}
			switch r := rr.(type) {
			case *dns.A:
				e.IPv4 = appendUniqueIP(e.IPv4, r.A)
			case *dns.AAAA:
				e.IPv6 = appendUniqueIP(e.IPv6, r.AAAA)
			}
		}
	}
}

func collectMDNSEntries(entries map[string]*types.MDNSDiscovery) []types.MDNSDiscovery {
	result := make([]types.MDNSDiscovery, 0, len(entries))
	for _, e := range entries {
		result = append(result, *e)
	}

	return result
}

func appendUniqueIP(ips []net.IP, ip net.IP) []net.IP {
	for _, existing := range ips {
		if existing.Equal(ip) {
			return ips
		}
	}

	return append(ips, ip)
}
