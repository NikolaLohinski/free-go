package client

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/nikolalohinski/free-go/types"
)

const (
	ssdpMulticastGroup  = "239.255.255.250"
	ssdpPort            = 1900
	freeboxSSDPTarget   = "urn:schemas-freebox-fr:device:Freebox:1"
	ssdpMSearchTemplate = "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"MAN: \"ssdp:discover\"\r\n" +
		"MX: 3\r\n" +
		"ST: urn:schemas-freebox-fr:device:Freebox:1\r\n" +
		"\r\n"
)

// DiscoverSSDP sends a UPnP/SSDP M-SEARCH for Freebox devices and collects
// responses until timeout expires or ctx is cancelled.
func DiscoverSSDP(ctx context.Context, timeout time.Duration) ([]types.SSDPDiscovery, error) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		return nil, fmt.Errorf("failed to open UDP socket: %w", err)
	}
	defer conn.Close()

	ssdpAddr := &net.UDPAddr{IP: net.ParseIP(ssdpMulticastGroup), Port: ssdpPort}
	if _, err = conn.WriteTo([]byte(ssdpMSearchTemplate), ssdpAddr); err != nil {
		return nil, fmt.Errorf("failed to send M-SEARCH: %w", err)
	}

	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	seen := map[string]struct{}{}
	var results []types.SSDPDiscovery
	buf := make([]byte, 65536)

	for {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}

			return results, fmt.Errorf("failed to read SSDP response: %w", err)
		}

		entry, err := parseSSDPResponse(buf[:n])
		if err != nil || entry == nil {
			continue
		}

		if _, duplicate := seen[entry.USN]; !duplicate {
			seen[entry.USN] = struct{}{}
			results = append(results, *entry)
		}
	}

	return results, nil
}

func parseSSDPResponse(data []byte) (*types.SSDPDiscovery, error) {
	lines := strings.Split(string(data), "\r\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "HTTP/1.1 200") {
		return nil, nil
	}

	headers := map[string]string{}
	for _, line := range lines[1:] {
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		headers[strings.ToUpper(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
	}

	st, ok := headers["ST"]
	if !ok || st != freeboxSSDPTarget {
		return nil, nil
	}

	location, ok := headers["LOCATION"]
	if !ok || location == "" {
		return nil, fmt.Errorf("missing LOCATION header in SSDP response")
	}

	return &types.SSDPDiscovery{
		Location: location,
		USN:      headers["USN"],
		Server:   headers["SERVER"],
	}, nil
}
