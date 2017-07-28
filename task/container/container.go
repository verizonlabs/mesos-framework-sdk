package container

import (
	"errors"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/resources"
	"mesos-framework-sdk/task"
	"mesos-framework-sdk/task/network"
	"mesos-framework-sdk/task/volume"
	"strings"
)

func ParseContainer(c *task.ContainerJSON) (*mesos_v1.ContainerInfo, error) {
	var ret *mesos_v1.ContainerInfo
	if c == nil {
		return nil, nil
	}

	networks, err := network.ParseNetworkJSON(c.Network)
	if err != nil {
		// NOTE (tim): We don't really need an error message here.
		// Debug message:
		// "No explicit network info passed in, using default host networking."
	}

	var vol []*mesos_v1.Volume
	if len(c.Volumes) > 0 {
		vol, err = volume.ParseVolumeJSON(c.Volumes)
		if err != nil {
			return nil, errors.New("Error parsing volume JSON: " + err.Error())
		}
	}

	if c.ImageName != nil {
		// Is container type explicitly set?
		var container *mesos_v1.ContainerInfo
		if c.ContainerType != nil {
			if strings.ToLower(*c.ContainerType) == "docker" {
				container.Docker = resources.CreateDockerInfo(
					resources.CreateImage(mesos_v1.Image_DOCKER.Enum(), *c.ImageName),
					mesos_v1.ContainerInfo_DockerInfo_BRIDGE.Enum(),
					nil,
					nil,
					nil, // volume driver
				)

				ret = resources.CreateContainerInfo(container, networks, vol, nil)
			} else if strings.ToLower(*c.ContainerType) == "mesos" {
				container.Mesos = resources.CreateMesosInfo(
					resources.CreateImage(mesos_v1.Image_APPC.Enum(), *c.ImageName),
				)
				ret = resources.CreateContainerInfo(container, networks, vol, nil)
			}
		} else { // Default is MESOS
			container.Mesos = resources.CreateMesosInfo(resources.CreateImage(
				mesos_v1.Image_DOCKER.Enum(), *c.ImageName))
			ret = resources.CreateContainerInfo(container, networks, vol, nil)
		}
	} else { // No image name was provided, commandinfo only.
		// Mesos-container with no image.
		ret = &mesos_v1.ContainerInfo{
			Type:         mesos_v1.ContainerInfo_MESOS.Enum(),
			Mesos:        &mesos_v1.ContainerInfo_MesosInfo{},
			NetworkInfos: networks,
			Volumes:      vol,
		}
	}

	return ret, nil
}
