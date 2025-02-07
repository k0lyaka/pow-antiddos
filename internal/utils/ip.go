package utils

import (
	"net"
	"net/http"
	"strings"
)

func ExtractIP(r *http.Request) (string, error) {
	// reverse proxy/CF
	ip := r.Header.Get("x-real-ip")
	parsedIp := net.ParseIP(ip)

	if parsedIp != nil {
		return ip, nil
	}

	ips := r.Header.Get("x-forwarded-for")
	for _, ip := range strings.Split(ips, ",") {
		parsedIp = net.ParseIP(ip)
		if parsedIp != nil {
			return ip, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	return ip, nil
}
