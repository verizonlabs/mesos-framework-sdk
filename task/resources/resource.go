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

package resources

import (
	"errors"
	"github.com/verizonlabs/mesos-framework-sdk/include/mesos_v1"
	"github.com/verizonlabs/mesos-framework-sdk/resources"
	"github.com/verizonlabs/mesos-framework-sdk/task"
)

func ParseResources(res *task.ResourceJSON) ([]*mesos_v1.Resource, error) {

	// We require at least some cpu and some mem.
	if res.Cpu <= 0.00 || res.Mem <= 0.00 {
		return nil, errors.New("CPU and memory must be greater than 0.0. " +
			"Please make sure you set cpu and mem properly.")
	}

	cpu := resources.CreateResource("cpus", res.Role, res.Cpu)
	mem := resources.CreateResource("mem", res.Role, res.Mem)
	disk, err := resources.CreateDisk(res.Disk, res.Role)
	if err != nil {
		return nil, err
	}

	return []*mesos_v1.Resource{cpu, mem, disk}, nil
}
