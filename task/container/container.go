package container

import (
	"errors"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/resources"
	"mesos-framework-sdk/task"
	"mesos-framework-sdk/task/network"
	"mesos-framework-sdk/task/volume"
	"strings"
)

func ParseContainer(c *task.ContainerJSON) (*mesos_v1.ContainerInfo, error) {
	var container *mesos_v1.ContainerInfo_MesosInfo
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
		if c.ContainerType != nil {
			if strings.ToLower(c.ContainerType) == "docker" {
				container = resources.CreateContainerInfoForDocker(
					resources.CreateImage(*c.ImageName, "", mesos_v1.Image_DOCKER.Enum()),
				)
			}
			ret = resources.CreateDockerContainerInfo(container, networks, vol, nil)
		} else {
			container = resources.CreateContainerInfoForMesos(
				resources.CreateImage(
					*c.ImageName, "", mesos_v1.Image_DOCKER.Enum(),
				),
			)
			ret = resources.CreateMesosContainerInfo(container, networks, vol, nil)
		}
	} else {
		return nil, errors.New("Container image name was not passed in. Please pass in a container name.")
	}

	return ret, nil
}
