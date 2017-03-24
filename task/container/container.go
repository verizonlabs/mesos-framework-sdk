package container

import (
	"errors"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/resources"
	"mesos-framework-sdk/task"
	"mesos-framework-sdk/task/network"
	"mesos-framework-sdk/task/volume"
)

func ParseContainer(c *task.ContainerJSON) (*mesos_v1.ContainerInfo, error) {
	var container *mesos_v1.ContainerInfo_MesosInfo
	if c.ImageName != nil {
		container = resources.CreateContainerInfoForMesos(
			resources.CreateImage(
				*c.ImageName, "", mesos_v1.Image_DOCKER.Enum(),
			),
		)
	} else {
		return nil, errors.New("Container image name was not passed in. Please pass in a container name.")
	}

	networks, err := network.ParseNetworkJSON(c.Network)
	if err != nil {
		//"No explicit network info passed in, using default host networking."
	}

	var vol []*mesos_v1.Volume
	if len(c.Volumes) > 0 {
		vol, err = volume.ParseVolumeJSON(c.Volumes)
		if err != nil {
			return nil, errors.New("Error parsing volume JSON: " + err.Error())
		}
	}

	return resources.CreateMesosContainerInfo(container, networks, vol, nil), nil
}
