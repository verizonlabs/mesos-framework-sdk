package volume

import (
	"github.com/pkg/errors"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/task"
	"strings"
)

func ParseVolumeJSON(volumes []task.VolumesJSON) ([]*mesos_v1.Volume, error) {
	mesosVolumes := []*mesos_v1.Volume{}
	for _, volume := range volumes {
		v := mesos_v1.Volume{}
		if strings.ToLower(volume.Mode) == "RO" {
			v.Mode = mesos_v1.Volume_RO.Enum()
		} else if strings.ToLower(volume.Mode) == "RW" {
			v.Mode = mesos_v1.Volume_RW.Enum()
		} else {
			v.Mode = mesos_v1.Volume_RW.Enum()
		}
		// Logical XOR to tell if both are set or not.
		if volume.ContainerPath == nil != volume.HostPath == nil {
			// Fail parsing and pass back error.
			return nil, errors.New("Both container and host path must be set.")
		}
		if volume.ContainerPath != nil {
			v.ContainerPath = volume.ContainerPath
		}
		if volume.HostPath != nil {
			v.HostPath = volume.HostPath
		}

		if volume.Source != nil {
			if volume.Source.Type != nil && strings.ToLower(volume.Source.Type) == "docker" {
				v.Source.Type = mesos_v1.Volume_Source_DOCKER_VOLUME.Enum()
				v.Source.DockerVolume = ParseDockerVolumeJSON(volume.Source.DockerVolume)
			} else {
				if v.Source.SandboxPath.Type == nil {
					v.Source.Type = mesos_v1.Volume_Source_SandboxPath_SELF.Enum()
				} else if strings.ToLower(v.Source.SandboxPath.Type) == "parent" {
					v.Source.Type = mesos_v1.Volume_Source_SandboxPath_PARENT.Enum()
				} else {
					// Default to self.
					v.Source.Type = mesos_v1.Volume_Source_SandboxPath_SELF.Enum()
				}
			}
		}
		mesosVolumes = append(mesosVolumes, &v)
	}
	return mesosVolumes, nil
}

func ParseDockerVolumeJSON(dockerVolume *task.DockerVolumeJSON) (source *mesos_v1.Volume_Source_DockerVolume) {
	// Do we only want to support certain drivers?
	if dockerVolume.Driver != nil {
		source.Driver = dockerVolume.Driver
	}
	if len(dockerVolume.DriverOptions) > 0 {
		params := mesos_v1.Parameters{}
		for k, v := range dockerVolume.DriverOptions {
			p := mesos_v1.Parameter{}
			p.Key, p.Value = k, v
			params = append(params, &p)
		}
		source.DriverOptions.Parameter = params
	}
	if dockerVolume.Name != nil {
		source.Name = dockerVolume.Name
	}

	return source
}
