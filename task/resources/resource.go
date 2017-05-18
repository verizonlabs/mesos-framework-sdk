package resources

import (
	"errors"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/resources"
	"mesos-framework-sdk/task"
)

func ParseResources(res *task.ResourceJSON) ([]*mesos_v1.Resource, error) {
	r := make([]*mesos_v1.Resource, 0)
	// We require at least some cpu and some mem.
	if res.Cpu <= 0.00 || res.Mem <= 0.00 {
		return nil, errors.New("CPU or Memory must be greater than 0.0. " +
			"Please make sure you set cpu and mem properly.")
	}

	var cpu = resources.CreateCpu(res.Cpu, res.Role)
	var mem = resources.CreateMem(res.Mem, res.Role)
	var disk = resources.CreateDisk(res.Disk, res.Role)
	r = append(r, cpu, mem, disk)
	return r, nil
}
