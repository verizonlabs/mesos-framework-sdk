package resources

import (
	"errors"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/resources"
	"mesos-framework-sdk/task"
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
