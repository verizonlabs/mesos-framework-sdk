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

	var cpu = resources.CreateResource("cpus", res.Role, res.Cpu)
	var mem = resources.CreateResource("mem", res.Role, res.Mem)
	var disk, err = resources.CreateDisk(res.Disk, res.Role)
	if err != nil {
		return nil, err
	}
	r = append(r, cpu, mem, disk)
	return r, nil
}
