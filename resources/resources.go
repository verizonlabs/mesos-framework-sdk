package resources

// This package contains helper methods for creating mesos types.
import (
	"errors"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/task"
	"mesos-framework-sdk/utils"
	"strings"
)

func CreateTaskInfo(
	name *string,
	uuid *mesos_v1.TaskID,
	cmd *mesos_v1.CommandInfo,
	res []*mesos_v1.Resource,
	con *mesos_v1.ContainerInfo,
	hc *mesos_v1.HealthCheck,
	labels *mesos_v1.Labels) *mesos_v1.TaskInfo {

	return &mesos_v1.TaskInfo{
		Name:        name,
		TaskId:      uuid,
		Resources:   res,
		Command:     cmd,
		Container:   con,
		HealthCheck: hc,
		Labels:      labels,
	}

}

func CreateDockerInfo(
	img *mesos_v1.Image,
	network *mesos_v1.ContainerInfo_DockerInfo_Network,
	ports []*mesos_v1.ContainerInfo_DockerInfo_PortMapping,
	params []*mesos_v1.Parameter,
	volDriver *string) *mesos_v1.ContainerInfo_DockerInfo {

	return &mesos_v1.ContainerInfo_DockerInfo{
		Image:        img.Docker.Name,
		Network:      network,
		PortMappings: ports,
		Parameters:   params,
		VolumeDriver: volDriver,
	}
}

func CreateMesosInfo(img *mesos_v1.Image) *mesos_v1.ContainerInfo_MesosInfo {
	return &mesos_v1.ContainerInfo_MesosInfo{
		Image: img,
	}
}

func CreateResource(name, role string, value float64) *mesos_v1.Resource {
	resource := &mesos_v1.Resource{
		Name: utils.ProtoString(name),
		Type: mesos_v1.Value_SCALAR.Enum(),
		Scalar: &mesos_v1.Value_Scalar{
			Value: utils.ProtoFloat64(value),
		},
	}

	if role != "" {
		resource.Role = utils.ProtoString(role)
	}

	return resource
}

// Creates a disk based on given task.Disk struct.
func CreateDisk(disk task.Disk, role string) (*mesos_v1.Resource, error) {

	// Disk must have a size.
	if disk.Size <= 0.0 {
		return nil, errors.New("Disk allocation size is 0 or less than 0.  Must be a positive float value.")
	}

	resource := CreateResource("disk", role, disk.Size)
	if disk.Source == nil {

		// This is a root disk.
		// Root disks map to the storage on the main operating system drive that the operator has presented to the agent.
		// Data is mapped into the work_dir of the agent.
		return resource, nil
	}

	// It's either PATH or MOUNT, it cannot be mixed.
	if disk.Source.Type == nil {

		// User specified a source field but not the type (required).
		return nil, errors.New("Disk source set but no type given. Valid types are MOUNT or PATH.")
	}

	sourceType := strings.ToLower(*disk.Source.Type)

	// Check if we have a PATH or MOUNT type disk source.
	if sourceType != "path" && sourceType != "mount" {
		return nil, errors.New("Invalid Disk source passed in, must be MOUNT or PATH if specified.")
	}

	d := &mesos_v1.Resource_DiskInfo{}
	if strings.ToLower(*disk.Source.Type) == "path" {
		if disk.Source.Path == nil {

			// Specified PATH but no fields set.
			return nil, errors.New("Disk source set to Path type but field path not set.")
		}

		if disk.Source.Mount != nil {

			// Specified PATH but gave mount.
			return nil, errors.New("Disk source set to Path type, but set mount field. Please set path field instead.")
		}

		d.Source.Type = mesos_v1.Resource_DiskInfo_Source_PATH.Enum()
		d.Source.Path.Root = disk.Source.Path
	} else if strings.ToLower(*disk.Source.Type) == "mount" {
		if disk.Source.Mount == nil {

			// Mount path type given, must have Mount field set.
			return nil, errors.New("Mount type given, but no mount path set.")
		}

		if disk.Source.Path != nil {

			// Specified MOUNT but gave path.
			return nil, errors.New("Mount type given, but path field set. Please set mount instead.")
		}

		d.Source.Type = mesos_v1.Resource_DiskInfo_Source_MOUNT.Enum()
		d.Source.Mount.Root = disk.Source.Mount
	}

	// TODO (tim): Add in external volume capabilities.
	// disk.Volume is for external volumes.
	resource.Disk = d
	return resource, nil
}

func CreateVolume(hostPath, containerPath string, image *mesos_v1.Image, source *mesos_v1.Volume_Source) *mesos_v1.Volume {
	return &mesos_v1.Volume{
		Mode:          mesos_v1.Volume_RW.Enum(),
		HostPath:      utils.ProtoString(hostPath),
		ContainerPath: utils.ProtoString(containerPath),
		Image:         image,
		Source:        source,
	}
}

func CreateImage(imgType *mesos_v1.Image_Type, name string) *mesos_v1.Image {
	return &mesos_v1.Image{
		Type: imgType,
		Docker: &mesos_v1.Image_Docker{
			Name: utils.ProtoString(name),
		},
		Appc: &mesos_v1.Image_Appc{
			Name: utils.ProtoString(name),
		},
	}
}

func CreateVolumeSource(source *mesos_v1.Volume_Source_Type,
	dockerVol *mesos_v1.Volume_Source_DockerVolume,
	sourcePath *mesos_v1.Volume_Source_SandboxPath) *mesos_v1.Volume_Source {

	if sourcePath.GetPath() == "" {
		return &mesos_v1.Volume_Source{
			Type:         source,
			DockerVolume: dockerVol,
		}
	}

	return &mesos_v1.Volume_Source{
		Type:        source,
		SandboxPath: sourcePath,
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
		Shell: utils.ProtoBool(true),
	}
}

func LaunchOfferOperation(taskList []*mesos_v1.TaskInfo) *mesos_v1.Offer_Operation {
	return &mesos_v1.Offer_Operation{
		Type:   mesos_v1.Offer_Operation_LAUNCH.Enum(),
		Launch: &mesos_v1.Offer_Operation_Launch{TaskInfos: taskList},
	}
}
