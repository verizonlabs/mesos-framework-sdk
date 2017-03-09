package resources

import (
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/include/mesos"
)

/*
This package contains functions to create common protobufs with ease.
*/
// Creates a taskInfo
func CreateTaskInfo(
	name *string,
	uuid *mesos_v1.TaskID,
	cmd *mesos_v1.CommandInfo,
	res []*mesos_v1.Resource,
	con *mesos_v1.ContainerInfo) *mesos_v1.TaskInfo {
	return &mesos_v1.TaskInfo{
		Name:      name,
		TaskId:    uuid,
		Command:   cmd,
		Resources: res,
		Container: con,
	}

}

func CreateDockerContainerInfo(
	c *mesos_v1.ContainerInfo_DockerInfo,
	n []*mesos_v1.NetworkInfo,
	v []*mesos_v1.Volume,
	h *string) *mesos_v1.ContainerInfo {
	return &mesos_v1.ContainerInfo{
		Type:         mesos_v1.ContainerInfo_DOCKER.Enum(),
		Hostname:     h,
		Docker:       c,
		NetworkInfos: n,
		Volumes:      v,
	}
}

func CreateMesosContainerInfo(
	c *mesos_v1.ContainerInfo_MesosInfo,
	n []*mesos_v1.NetworkInfo,
	v []*mesos_v1.Volume,
	h *string) *mesos_v1.ContainerInfo {
	return &mesos_v1.ContainerInfo{
		Type:         mesos_v1.ContainerInfo_MESOS.Enum(),
		Hostname:     h,
		Mesos:        c,
		NetworkInfos: n,
		Volumes:      v,
	}
}

func CreateContainerInfoForMesos(img *mesos_v1.Image) *mesos_v1.ContainerInfo_MesosInfo {
	return &mesos_v1.ContainerInfo_MesosInfo{
		Image: img,
	}
}

// Creates a cpu share that is not reserved.
func CreateCpu(cpuShare float64, role string) *mesos_v1.Resource {
	resource := &mesos_v1.Resource{
		Name: proto.String("cpus"),
		Type: mesos_v1.Value_SCALAR.Enum(),
		Scalar: &mesos_v1.Value_Scalar{
			Value: proto.Float64(cpuShare),
		},
	}
	if role != "" {
		resource.Role = proto.String(role)
	}

	return resource
}

// Creates a memory share that is not reserved.
func CreateMem(memShare float64, role string) *mesos_v1.Resource {
	resource := &mesos_v1.Resource{
		Name: proto.String("mem"),
		Type: mesos_v1.Value_SCALAR.Enum(),
		Scalar: &mesos_v1.Value_Scalar{
			Value: proto.Float64(memShare),
		},
	}
	if role != "" {
		resource.Role = proto.String(role)
	}

	return resource
}

func CreateDisk(vol *mesos_v1.Volume, source *mesos_v1.Resource_DiskInfo_Source) *mesos_v1.Resource_DiskInfo {
	return &mesos_v1.Resource_DiskInfo{
		Volume: vol,
		Source: source,
	}
}

func CreateVolume(hostPath, containerPath string, image *mesos_v1.Image, source *mesos_v1.Volume_Source) *mesos_v1.Volume {
	return &mesos_v1.Volume{
		Mode:          mesos_v1.Volume_RW.Enum(),
		HostPath:      proto.String(hostPath),
		ContainerPath: proto.String(containerPath),
		Image:         image,
		Source:        source,
	}
}

func CreateImage(name string, id string, imgType *mesos_v1.Image_Type) *mesos_v1.Image {
	var img *mesos_v1.Image
	if imgType == mesos_v1.Image_DOCKER.Enum() {
		img = &mesos_v1.Image{
			Type: imgType,
			Docker: &mesos_v1.Image_Docker{
				Name: proto.String(name),
			},
		}
	} else {
		img = &mesos_v1.Image{
			Type: imgType,
			Appc: &mesos_v1.Image_Appc{
				Name: proto.String(name),
				Id:   proto.String(id),
			},
		}
	}
	return img
}

func CreateVolumeSource(source *mesos_v1.Volume_Source_Type,
	dockerVol *mesos_v1.Volume_Source_DockerVolume,
	sourcePath *mesos_v1.Volume_Source_SandboxPath) *mesos_v1.Volume_Source {

	if sourcePath.GetPath() == "" {
		return &mesos_v1.Volume_Source{
			Type:         source,
			DockerVolume: dockerVol,
		}
	} else {
		return &mesos_v1.Volume_Source{
			Type:        source,
			SandboxPath: sourcePath,
		}
	}

}

func CreateCommandInfo(
	cmd *string, args []string,
	user *string,
	uris []*mesos_v1.CommandInfo_URI,
	env *mesos_v1.Environment,
	isShell *bool) *mesos_v1.CommandInfo {

	return &mesos_v1.CommandInfo{
		Value:       cmd,
		Arguments:   args,
		User:        user,
		Environment: env,
		Shell:       isShell,
	}
}

// Assumes only cmd, uris and shell set to true.
func CreateSimpleCommandInfo(cmd *string, uris []*mesos_v1.CommandInfo_URI) *mesos_v1.CommandInfo {
	return &mesos_v1.CommandInfo{
		Value: cmd,
		Uris:  uris,
		Shell: proto.Bool(true),
	}
}

func LaunchOfferOperation(taskList []*mesos_v1.TaskInfo) *mesos_v1.Offer_Operation {
	return &mesos_v1.Offer_Operation{
		Type:   mesos_v1.Offer_Operation_LAUNCH.Enum(),
		Launch: &mesos_v1.Offer_Operation_Launch{TaskInfos: taskList},
	}
}
