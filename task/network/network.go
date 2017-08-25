// Copyright 2017 Verizon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package network

import (
	"errors"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/task"
	"strings"

	"github.com/golang/protobuf/proto"
)

// Parse NetworkJSON into a list of Networkwork Infos.
func ParseNetworkJSON(networks []task.NetworkJSON) ([]*mesos_v1.NetworkInfo, error) {
	if len(networks) == 0 {
		return []*mesos_v1.NetworkInfo{}, errors.New("Empty list of networks passed in.")
	}

	networkInfos := []*mesos_v1.NetworkInfo{}
	// Iterate over each network
	for _, network := range networks {
		n := &mesos_v1.NetworkInfo{}
		if network.Name != nil {
			n.Name = network.Name
		}
		if len(network.Groups) > 0 {
			n.Groups = network.Groups
		}
		if len(network.IpAddresses) > 0 {
			n.IpAddresses = ParseNetworkJSONIpAddresses(network.IpAddresses)
		}
		if len(network.Labels) > 0 {
			n.Labels = ParseNetworkJSONLabels(network.Labels)
		}
		if len(network.PortMapping) > 0 {
			// Gather all ports
			n.PortMappings = ParseNetworkJSONPortMapping(network.PortMapping)
		}
		networkInfos = append(networkInfos, n)
	}
	return networkInfos, nil
}

// Parses Ip addresses out of the network json struct
func ParseNetworkJSONIpAddresses(ipaddrs []task.IpAddressJSON) (ips []*mesos_v1.NetworkInfo_IPAddress) {
	for _, ipaddr := range ipaddrs {
		ip := &mesos_v1.NetworkInfo_IPAddress{}
		ip.IpAddress = ipaddr.IP
		if strings.ToLower(*ipaddr.Protocol) == "ipv4" {
			ip.Protocol = mesos_v1.NetworkInfo_IPv4.Enum()
		} else if strings.ToLower(*ipaddr.Protocol) == "ipv6" {
			ip.Protocol = mesos_v1.NetworkInfo_IPv6.Enum()
		} else {
			ip.Protocol = nil
		}
		ips = append(ips, ip)
	}
	return ips
}

// Parse all labels in the network JSON.
func ParseNetworkJSONLabels(labels []map[string]string) *mesos_v1.Labels {
	labelList := []*mesos_v1.Label{}
	for _, label := range labels {
		l := &mesos_v1.Label{}
		for k, v := range label {
			l.Key, l.Value = proto.String(k), proto.String(v)
		}
		labelList = append(labelList, l)
	}
	return &mesos_v1.Labels{Labels: labelList}
}

func ParseNetworkJSONPortMapping(portMap []*task.PortMapping) (portMapList []*mesos_v1.NetworkInfo_PortMapping) {
	for _, portMap := range portMap {
		pm := &mesos_v1.NetworkInfo_PortMapping{}
		portMap.ContainerPort, portMap.HostPort, portMap.Protocol = pm.ContainerPort, pm.HostPort, pm.Protocol
		portMapList = append(portMapList, pm)
	}
	return portMapList
}
