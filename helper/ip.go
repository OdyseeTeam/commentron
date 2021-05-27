package helper

import (
	"net"
	"net/http"
	"strings"
)

// GetIPAddressForRequest returns the real IP address of the request
func GetIPAddressForRequest(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
				// bad address, go to next
				continue
			}
			return ip
		}
	}

	ipParts := strings.Split(r.RemoteAddr, ":")
	ip := strings.Join(ipParts[:len(ipParts)-1], ":")

	if ip == "[::1]" {
		return "127.0.0.1"
	}
	return ip
}
