package container

import (
	"errors"
	"github.com/golang/protobuf/proto"
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
		if c.ContainerType != nil {
			var dockerContainer *mesos_v1.ContainerInfo_DockerInfo
			if strings.ToLower(*c.ContainerType) == "docker" {
				dockerContainer = resources.CreateContainerInfoForDocker(
					c.ImageName,
					mesos_v1.ContainerInfo_DockerInfo_BRIDGE.Enum(),
					[]*mesos_v1.ContainerInfo_DockerInfo_PortMapping{},
					[]*mesos_v1.Parameter{},
					proto.String(""), // volume driver
				)
			}
			ret = resources.CreateDockerContainerInfo(dockerContainer, networks, vol, nil)
		} else {
			var container *mesos_v1.ContainerInfo_MesosInfo
			container = resources.CreateContainerInfoForMesos(
				resources.CreateImage(
					*c.ImageName, "", mesos_v1.Image_DOCKER.Enum(),
				),
			)
			ret = resources.CreateMesosContainerInfo(container, networks, vol, nil)
		}
	} else {
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
