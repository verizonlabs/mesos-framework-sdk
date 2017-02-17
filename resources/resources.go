package mesos_framework_sdk

import (
	"github.com/gogo/protobuf/proto"
	"mesos-framework-sdk/include/mesos"
)

/*
This package contains functions to create common protobufs with ease.

Name   *string       `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	Type   *Value_Type   `protobuf:"varint,2,req,name=type,enum=mesos.v1.Value_Type" json:"type,omitempty"`
	Scalar *Value_Scalar `protobuf:"bytes,3,opt,name=scalar" json:"scalar,omitempty"`
	Ranges *Value_Ranges `protobuf:"bytes,4,opt,name=ranges" json:"ranges,omitempty"`
	Set    *Value_Set    `protobuf:"bytes,5,opt,name=set" json:"set,omitempty"`
	// The role that this resource is reserved for. If "*", this indicates
	// that the resource is unreserved. Otherwise, the resource will only
	// be offered to frameworks that belong to this role.
	Role *string `protobuf:"bytes,6,opt,name=role,def=*" json:"role,omitempty"`
	// If this is set, this resource was dynamically reserved by an
	// operator or a framework. Otherwise, this resource is either unreserved
	// or statically reserved by an operator via the --resources flag.
	Reservation *Resource_ReservationInfo `protobuf:"bytes,8,opt,name=reservation" json:"reservation,omitempty"`
	Disk        *Resource_DiskInfo        `protobuf:"bytes,7,opt,name=disk" json:"disk,omitempty"`
	// If this is set, the resources are revocable, i.e., any tasks or
	// executors launched using these resources could get preempted or
	// throttled at any time. This could be used by frameworks to run
	// best effort tasks that do not need strict uptime or performance
	// guarantees. Note that if this is set, 'disk' or 'reservation'
	// cannot be set.
	Revocable *Resource_RevocableInfo `protobuf:"bytes,9,opt,name=revocable" json:"revocable,omitempty"`
	// If this is set, the resources are shared, i.e. multiple tasks
	// can be launched using this resource and all of them shall refer
	// to the same physical resource on the cluster. Note that only
	// persistent volumes can be shared currently.
	Shared           *Resource_SharedInfo `protobuf:"bytes,10,opt,name=shared" json:"shared,omitempty"`
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

func CreateDisk(vol *mesos_v1.Volume, source *mesos_v1.Resource_DiskInfo_Source) {
	return &mesos_v1.Resource_DiskInfo{
		Volume: vol,
		Source: source,
	}
}

func CreateVolume(hostPath, containerPath string, image *mesos_v1.Image, source mesos_v1.Volume_Source) *mesos_v1.Volume {
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
