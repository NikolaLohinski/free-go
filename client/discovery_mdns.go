package client

import (
	"context"
	"fmt"
	"net"
	"strconv"
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

	mcastAddr := &net.UDPAddr{IP: net.ParseIP(mdnsMulticastGroup), Port: mdnsPort}

	// recvConn joins the multicast group to receive mDNS responses.
	recvConn, err := net.ListenMulticastUDP("udp4", nil, mcastAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to join mDNS multicast group: %w", err)
	}
	defer recvConn.Close()

	// sendConn is a regular unicast socket so the query carries a valid source address.
	sendConn, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		return nil, fmt.Errorf("failed to open send socket: %w", err)
	}
	defer sendConn.Close()

	deadline := time.Now().Add(timeout)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}

	if err = recvConn.SetDeadline(deadline); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			recvConn.SetDeadline(time.Now()) //nolint:errcheck
		case <-done:
		}
	}()

	if _, err = sendConn.WriteTo(packed, mcastAddr); err != nil {
		return nil, fmt.Errorf("failed to send mDNS query: %w", err)
	}

	entries := map[string]*types.MDNSDiscovery{}
	buf := make([]byte, 65536)

	for {
		n, _, err := recvConn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}

			return nil, fmt.Errorf("failed to read mDNS response: %w", err)
		}

		response := new(dns.Msg)
		if err := response.Unpack(buf[:n]); err != nil {
			continue
		}

		parseMDNSResponse(response, entries)
	}

	if err := ctx.Err(); err != nil {
		return collectMDNSEntries(entries), err
	}

	return collectMDNSEntries(entries), nil
}

func parseMDNSResponse(msg *dns.Msg, entries map[string]*types.MDNSDiscovery) {
	allRRs := make([]dns.RR, 0, len(msg.Answer)+len(msg.Ns)+len(msg.Extra))
	allRRs = append(allRRs, msg.Answer...)
	allRRs = append(allRRs, msg.Ns...)
	allRRs = append(allRRs, msg.Extra...)

	for _, rr := range allRRs {
		switch r := rr.(type) {
		case *dns.PTR:
			if r.Hdr.Name == freeboxMDNSService {
				if _, ok := entries[r.Ptr]; !ok {
					entries[r.Ptr] = &types.MDNSDiscovery{
						Name: strings.TrimSuffix(r.Ptr, "."+freeboxMDNSService),
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
						parseMDNSTXTField(e, kv[0], kv[1])
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

func parseMDNSTXTField(e *types.MDNSDiscovery, key, value string) {
	switch key {
	case "uid":
		e.UID = value
	case "device_type":
		e.DeviceType = value
	case "api_version":
		e.APIVersion = value
	case "api_domain":
		e.APIDomain = value
	case "api_base_url":
		e.APIBaseURL = value
	case "https_port":
		if p, err := strconv.Atoi(value); err == nil {
			e.HTTPSPort = p
		}
	case "https_available":
		e.HTTPSAvailable = value == "1"
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
