package utils

import (
	"errors"
	"net"
)

const (
	IPV4FIRSTBIT = 10
	IPv4Bits     = 32
)

// Gathers the internal network as defined.
// This will not work if there are multiple
// 10.0.0.0/24's on a host.
// NOTE (tim): Talk to mike c about how the new
// /25 networks will work, as well as overlay
// networks.
func GetInternalNetworkInterface(subnet int) (net.IP, error) {
	interfaces, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, interFace := range interfaces {
		ip, net, err := net.ParseCIDR(interFace.String())
		if err != nil {
			return nil, err
		}
		ones, bits := net.Mask.Size()
		// If it's v4
		if bits <= IPv4Bits {
			// Is this a /24 network?
			if ones == subnet {
				// IP is padded to the left for ipv6.
				// First bit for ipv4 starts at 12th index.
				// NOTE (tim): This magic 10 number will
				// have to change going forward to support
				// new networking designs
				if ip[12] == byte(10) {
					return ip, nil
				}
			}
		}
	}
	return nil, errors.New("No IPv4 addresses found in 10.x.x.x/8 network")
}
