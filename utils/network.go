package utils

import "net"

// Gets the IPs corresponding to the specified interface.
// This can be used for when you need to determine your network state for leader election, etc.
func GetIPs(iface string) (map[string]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ips := make(map[string]string)

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

			if ip.To4() != nil {
				ips["tcp4"] = ip.String()
			} else {
				ips["tcp6"] = ip.String()
			}
		}
	}

	return ips, err
}
