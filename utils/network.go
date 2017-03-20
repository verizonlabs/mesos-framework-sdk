package utils

import "net"

// Gets the IPs corresponding to the specified interface.
// This can be used for when you need to determine your network state for leader election, etc.
func GetIP(iface string) ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)

	for _, i := range ifaces {
		if i.Name != iface {
			continue
		}

		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			ips = append(ips, ip.String())
		}
	}

	return ips, err
}
