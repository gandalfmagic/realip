package realip

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func Get(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-Ip")
	netIP := net.ParseIP(extractIP(ip))
	if netIP != nil {
		return ip, nil
	}

	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")
	for _, ip = range splitIps {
		realIP := ip

		netIP = net.ParseIP(extractIP(ip))
		if netIP != nil {
			return realIP, nil
		}
	}

	netIP = net.ParseIP(extractIP(r.RemoteAddr))
	if netIP != nil {
		return r.RemoteAddr, nil
	}

	return "", fmt.Errorf("no valid source ip found in the request headers")
}

// Return the ip address from an "ip:port" formatted string
// If the string is already a single ip value, returns that
func extractIP(ip string) string {
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '.':
			if last(ip, ':') < 0 {
				return ip
			}

			ip, _, _ = net.SplitHostPort(ip)
			return ip
		case ':':
			// Normalize IPv6 if needed
			if ip[0] != '[' {
				return ip
			}

			ip, _, _ = net.SplitHostPort(ip)
			return ip
		}
	}

	return ip
}

// Index of rightmost occurrence of b in s.
func last(s string, b byte) int {
	i := len(s)
	for i--; i >= 0; i-- {
		if s[i] == b {
			break
		}
	}
	return i
}
