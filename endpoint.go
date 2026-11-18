package main

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"path/filepath"
	"strconv"
	"strings"
)

// monitorAddress selects one explicit transport; it never changes the target
// implicitly when an endpoint is unavailable.
func monitorAddress(socketPath, host string, socketExplicit bool) (string, string, error) {
	if host == "" {
		if !filepath.IsAbs(socketPath) {
			return "", "", errors.New("-socket must be an absolute path")
		}
		return "unixgram", socketPath, nil
	}
	if socketExplicit {
		return "", "", errors.New("-host and -socket cannot be used together")
	}
	name, port, err := net.SplitHostPort(host)
	if err != nil {
		name, port = host, "323"
		if strings.HasPrefix(name, "[") && strings.HasSuffix(name, "]") {
			name = name[1 : len(name)-1]
			if _, err := netip.ParseAddr(name); err != nil {
				return "", "", fmt.Errorf("invalid -host address %q", host)
			}
		} else if _, err := netip.ParseAddr(name); err != nil && strings.ContainsAny(name, ":[]") {
			return "", "", fmt.Errorf("invalid -host address %q (use host:port or [IPv6]:port)", host)
		}
	}
	if name == "" || strings.ContainsAny(name, "/\\ \t\r\n") {
		return "", "", fmt.Errorf("invalid -host address %q", host)
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 {
		return "", "", errors.New("UDP command port must be between 1 and 65535")
	}
	return "udp", net.JoinHostPort(name, strconv.Itoa(n)), nil
}
