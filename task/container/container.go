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

package container

import (
	"errors"
	"github.com/verizonlabs/mesos-framework-sdk/include/mesos_v1"
	"github.com/verizonlabs/mesos-framework-sdk/resources"
	"github.com/verizonlabs/mesos-framework-sdk/task"
	"github.com/verizonlabs/mesos-framework-sdk/task/network"
	"github.com/verizonlabs/mesos-framework-sdk/task/volume"
)

func ParseContainer(c *task.ContainerJSON) (*mesos_v1.ContainerInfo, error) {
	if c == nil {
		return nil, nil
	}

	// "No explicit network info passed in, using default host networking."
	networks, _ := network.ParseNetworkJSON(c.Network)

	var vol []*mesos_v1.Volume
	var err error
	if len(c.Volumes) > 0 {
		vol, err = volume.ParseVolumeJSON(c.Volumes)
		if err != nil {
			return nil, errors.New("Error parsing volume JSON: " + err.Error())
		}
	}

	// Default to the UCR.
	container := &mesos_v1.ContainerInfo{
		Type:         mesos_v1.ContainerInfo_MESOS.Enum(),
		NetworkInfos: networks,
		Volumes:      vol,
	}

	if c.ImageName == nil {
		return container, nil
	}

	container.Mesos = resources.CreateMesosInfo(
		resources.CreateImage(mesos_v1.Image_DOCKER.Enum(), *c.ImageName),
	)

	return container, nil
}
