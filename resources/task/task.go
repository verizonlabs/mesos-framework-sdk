package task

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"mesos-framework-sdk/include/mesos"
)

type ApplicationJSON struct {
	Name        string              `json:"name"`
	Resources   *ResourceJSON       `json:"resources"`
	Command     *CommandJSON        `json:"command"`
	Container   *ContainerJSON      `json:"container"`
	HealthCheck *HealthCheckJSON    `json:"healthcheck"`
	Labels      []map[string]string `json:"labels"`
}

type HealthCheckJSON struct {
	Endpoint *string `json:"endpoint"`
}

type KillJson struct {
	Name *string `json:"name"`
}

type ResourceJSON struct {
	Mem  float64 `json:"mem"`
	Cpu  float64 `json:"cpu"`
	Disk float64 `json:"disk"`
}

type CommandJSON struct {
	Cmd  *string   `json:"cmd"`
	Uris []UriJSON `json:"uris"`
}

type ContainerJSON struct {
	ImageName *string       `json:"image"`
	Tag       *string       `json:"tag"`
	Network   []NetworkJSON `json:"network"`
}

type NetworkJSON struct {
	IpAddresses []IpAddressJSON `json:"ipaddress,omitempty"`
	Name        *string         `json:"name"`
}

// Parse NetworkJSON into a list of Networkwork Infos.
func ParseNetworkJSON(networks []NetworkJSON) ([]*mesos_v1.NetworkInfo, error) {
	if len(networks) == 0 {
		return []*mesos_v1.NetworkInfo{}, errors.New("Empty list of networks passed in.")
	}
	networkSlice := []*mesos_v1.NetworkInfo{}
	for _, network := range networks {
		ips := []*mesos_v1.NetworkInfo_IPAddress{}
		for range network.IpAddresses {
			i := &mesos_v1.NetworkInfo_IPAddress{
				IpAddress: proto.String(""),
			}
			ips = append(ips, i)
		}
		n := &mesos_v1.NetworkInfo{
			IpAddresses: ips,
		}
		networkSlice = append(networkSlice, n)
	}
	return networkSlice, nil
}

type IpAddressJSON struct {
	IP       *string `json:"ip"`
	Protocol *string `json:"protocol"`
}

type UriJSON struct {
	Uri     *string `json:"uri"`
	Extract *bool   `json:"extract"`
	Execute *bool   `json:"execute"`
}
