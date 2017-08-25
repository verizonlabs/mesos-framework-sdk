// Copyright 2017 Verizon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package volume

import (
	"errors"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/task"
	"strings"

	"github.com/golang/protobuf/proto"
)

func ParseVolumeJSON(volumes []task.VolumesJSON) ([]*mesos_v1.Volume, error) {
	mesosVolumes := []*mesos_v1.Volume{}
	for _, volume := range volumes {
		v := mesos_v1.Volume{}
		if volume.Mode == nil {
			v.Mode = mesos_v1.Volume_RW.Enum()
		} else {
			if strings.ToLower(*volume.Mode) == "RO" {
				v.Mode = mesos_v1.Volume_RO.Enum()
			} else if strings.ToLower(*volume.Mode) == "RW" {
				v.Mode = mesos_v1.Volume_RW.Enum()
			}
		}
		// Logical XOR to tell if both are set or not.
		if (volume.ContainerPath == nil) != (volume.HostPath == nil) {
			// Fail parsing and pass back error.
			return nil, errors.New("Both container and host path must be set.")
		}
		if volume.ContainerPath != nil {
			v.ContainerPath = volume.ContainerPath
		}
		if volume.HostPath != nil {
			v.HostPath = volume.HostPath
		}

		if (volume.Source != nil) && (volume.Source.Type != nil) {
			src := mesos_v1.Volume_Source{}
			if strings.ToLower(*volume.Source.Type) == "docker" {
				src.Type = mesos_v1.Volume_Source_DOCKER_VOLUME.Enum()
				src.DockerVolume = ParseDockerVolumeJSON(&volume.Source.DockerVolume)
			} else {
				src.Type = mesos_v1.Volume_Source_SANDBOX_PATH.Enum()
				sandbox := mesos_v1.Volume_Source_SandboxPath{}
				sandbox.Type = mesos_v1.Volume_Source_SandboxPath_SELF.Enum()
				sandbox.Path = proto.String(".")

				v.Source = &src
				v.Source.SandboxPath = &sandbox
			}
		} else {
			src := mesos_v1.Volume_Source{}
			src.Type = mesos_v1.Volume_Source_SANDBOX_PATH.Enum()
			sandbox := mesos_v1.Volume_Source_SandboxPath{}
			sandbox.Type = mesos_v1.Volume_Source_SandboxPath_SELF.Enum()
			sandbox.Path = proto.String(".")

			v.Source = &src
			v.Source.SandboxPath = &sandbox
		}

		mesosVolumes = append(mesosVolumes, &v)
	}
	return mesosVolumes, nil
}

func ParseDockerVolumeJSON(dockerVolume *task.DockerVolumeJSON) *mesos_v1.Volume_Source_DockerVolume {
	source := mesos_v1.Volume_Source_DockerVolume{}
	// Do we only want to support certain drivers?
	if dockerVolume.Driver != nil {
		source.Driver = dockerVolume.Driver
	}
	if len(dockerVolume.DriverOptions) > 0 {
		params := []*mesos_v1.Parameter{}
		for _, options := range dockerVolume.DriverOptions {
			for k, v := range options {
				p := mesos_v1.Parameter{}
				p.Key, p.Value = proto.String(k), proto.String(v)
				params = append(params, &p)
			}
		}
		source.DriverOptions.Parameter = params
	}
	if dockerVolume.Name != nil {
		source.Name = dockerVolume.Name
	}

	return &source
}
