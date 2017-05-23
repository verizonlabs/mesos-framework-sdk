package resources

// This package contains helper methods for creating mesos types.
import (
	"errors"
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/task"
	"strings"
)

func CreateTaskInfo(
	name *string,
	uuid *mesos_v1.TaskID,
	cmd *mesos_v1.CommandInfo,
	res []*mesos_v1.Resource,
	con *mesos_v1.ContainerInfo,
	labels *mesos_v1.Labels) *mesos_v1.TaskInfo {
	return &mesos_v1.TaskInfo{
		Name:      name,
		TaskId:    uuid,
		Resources: res,
		Command:   cmd,
		Container: con,
		Labels:    labels,
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

func CreateContainerInfoForDocker(
	img *string,
	network *mesos_v1.ContainerInfo_DockerInfo_Network,
	ports []*mesos_v1.ContainerInfo_DockerInfo_PortMapping,
	params []*mesos_v1.Parameter,
	volDriver *string) *mesos_v1.ContainerInfo_DockerInfo {
	return &mesos_v1.ContainerInfo_DockerInfo{
		Image:        img,
		Network:      network,
		PortMappings: ports,
		Parameters:   params,
		VolumeDriver: volDriver,
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

// Creates a disk based on given task.Disk struct.
func CreateDisk(disk task.Disk, role string) (*mesos_v1.Resource, error) {
	// Create a diskinfo resource if required.
	d := &mesos_v1.Resource_DiskInfo{}

	// Disk must have a size.
	if disk.Size <= 0.0 {
		return nil, errors.New("Disk allocation size is 0 or less than 0.  Must be a positive float value.")
	}
	// If disk source is not nil, we have a PATH or MOUNT style disk.
	if disk.Source != nil {
		// It's either PATH or MOUNT, it cannot be mixed.
		if disk.Source.Type != nil {
			// Check if we have a PATH or MOUNT type disk source.
			if strings.ToLower(*disk.Source.Type) == "path" {
				d.Source.Type = mesos_v1.Resource_DiskInfo_Source_PATH.Enum()
				if disk.Source.Path != nil {
					d.Source.Path.Root = disk.Source.Path
				} else if disk.Source.Mount != nil {
					// Specified PATH but gave mount.
					return nil, errors.New("Disk source set to Path type, but set mount field. Please set path field instead.")
				} else {
					// Specified PATH but no fields set.
					return nil, errors.New("Disk source set to Path type but field path not set.")
				}
			} else if strings.ToLower(*disk.Source.Type) == "mount" {
				d.Source.Type = mesos_v1.Resource_DiskInfo_Source_MOUNT.Enum()
				if disk.Source.Mount != nil {
					d.Source.Mount.Root = disk.Source.Mount
				} else if disk.Source.Path != nil {
					// Specified MOUNT but gave path.
					return nil, errors.New("Mount type given, but path field set. Please set mount instead.")
				} else {
					// Mount path type given, must have Mount field set.
					return nil, errors.New("Mount type given, but no mount path set.")
				}
			} else {
				return nil, errors.New("Invalid Disk source passed in, must be MOUNT or PATH if specified.")
			}
		} else {
			// User specified a source field but not the type (required).
			return nil, errors.New("Disk source set but no type given. Valid types are MOUNT or PATH.")
		}
	} // End checking DISK SOURCE fields.

	// TODO (tim): Add in external volume capabilities.
	// disk.Volume is for external volumes.

	// Create the resource to return.
	resource := &mesos_v1.Resource{
		Name: proto.String("disk"),
		Type: mesos_v1.Value_SCALAR.Enum(),
		Scalar: &mesos_v1.Value_Scalar{
			Value: proto.Float64(disk.Size),
		},
	}

	// If we set our source to something, add DISK field to resource.
	if d.Source != nil {
		resource.Disk = d
	}

	// Set a role if one was passed in.
	if role != "" {
		resource.Role = proto.String(role)
	}
	return resource, nil
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
	if *imgType == mesos_v1.Image_DOCKER {
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
