package proxycore

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

var privateRanges []net.IPNet

var lookupIP = net.LookupIP

var metadataIPs = []net.IP{
	net.ParseIP("169.254.169.254"),
	net.ParseIP("169.254.170.2"),
	net.ParseIP("100.100.100.200"),
	net.ParseIP("fd00:ec2::254"),
}

var blockedHostnames = map[string]bool{
	"localhost":                 true,
	"metadata":                 true,
	"metadata.google.internal": true,
	"metadata.internal":        true,
}

func init() {
	cidrs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"0.0.0.0/8",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
		"100.64.0.0/10", // carrier-grade NAT
		"198.18.0.0/15", // benchmarking
		"224.0.0.0/4",   // multicast
		"240.0.0.0/4",   // reserved
		"64:ff9b::/96",  // NAT64, embeds IPv4
		"2002::/16",     // 6to4, embeds IPv4
	}
	for _, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err == nil {
			privateRanges = append(privateRanges, *ipnet)
		}
	}
}

func ValidateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("scheme %q not allowed, use http or https", u.Scheme)
	}

	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("empty hostname")
	}

	lower := strings.TrimSuffix(strings.ToLower(host), ".")
	if blockedHostnames[lower] {
		return fmt.Errorf("hostname %q is blocked", host)
	}
	if strings.HasSuffix(lower, ".local") || strings.HasSuffix(lower, ".localhost") || strings.HasSuffix(lower, ".internal") {
		return fmt.Errorf("hostname %q is blocked (private suffix)", host)
	}

	if ip := net.ParseIP(lower); ip != nil {
		return validateIP(ip)
	}

	ips, err := lookupIP(lower)
	if err != nil {
		return fmt.Errorf("cannot resolve %q, refusing: %w", host, err)
	}
	for _, resolved := range ips {
		if err := validateIP(resolved); err != nil {
			return fmt.Errorf("hostname %q resolves to blocked IP: %w", host, err)
		}
	}
	return nil
}

func validateIP(ip net.IP) error {
	for _, meta := range metadataIPs {
		if ip.Equal(meta) {
			return fmt.Errorf("IP %s is a cloud metadata endpoint", ip)
		}
	}
	for _, cidr := range privateRanges {
		if cidr.Contains(ip) {
			return fmt.Errorf("IP %s is in private range %s", ip, cidr.String())
		}
	}
	if ip.IsUnspecified() || ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return fmt.Errorf("IP %s is not routable", ip)
	}
	return nil
}
