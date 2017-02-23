package mesos_framework_sdk

import (
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/include/mesos"
)

/*
This package contains functions to create common protobufs with ease.
*/

// Creates a cpu share that is not reserved.
func CreateCpu(cpuShare float64, role string) *mesos_v1.Resource {
	return &mesos_v1.Resource{
		Name: proto.String("cpus"),
		Type: mesos_v1.Value_SCALAR.Enum(),
		Scalar: &mesos_v1.Value_Scalar{
			Value: proto.Float64(cpuShare),
		},
		Role: proto.String(role),
	}
}

// Creates a memory share that is not reserved.
func CreateMem(memShare float64, role string) *mesos_v1.Resource {
	return &mesos_v1.Resource{
		Name: proto.String("mem"),
		Type: mesos_v1.Value_SCALAR.Enum(),
		Scalar: &mesos_v1.Value_Scalar{
			Value: proto.Float64(memShare),
		},
		Role: proto.String(role),
	}
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
