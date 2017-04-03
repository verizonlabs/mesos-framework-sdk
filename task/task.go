package task

type ApplicationJSON struct {
	Name        string              `json:"name"`
	Resources   *ResourceJSON       `json:"resources"`
	Command     *CommandJSON        `json:"command"`
	Container   *ContainerJSON      `json:"container"`
	HealthCheck *HealthCheckJSON    `json:"healthcheck"`
	Labels      []map[string]string `json:"labels"`
	Filters     []Filter            `json:"filters"`
}

type HealthCheckJSON struct {
	Endpoint *string `json:"endpoint"`
}

type Filter struct {
	Type  string   `json:"type"`
	Value []string `json:"value"`
}

type KillJson struct {
	Name *string `json:"name"`
}

type ResourceJSON struct {
	Mem  float64 `json:"mem"`
	Cpu  float64 `json:"cpu"`
	Disk float64 `json:"disk"`
	Role string  `json:"role"`
}

type CommandJSON struct {
	Cmd  *string   `json:"cmd"`
	Uris []UriJSON `json:"uris"`
}

type ContainerJSON struct {
	ContainerType *string       `json:"type"`
	ImageName     *string       `json:"image"`
	Tag           *string       `json:"tag"`
	Network       []NetworkJSON `json:"network"`
	Volumes       []VolumesJSON `json:"volume"`
}

type VolumesJSON struct {
	ContainerPath *string           `json:"container_path"`
	HostPath      *string           `json:"host_path"`
	Mode          *string           `json:"mode"`
	Source        *VolumeSourceJSON `json:"source"`
}

type VolumeSourceJSON struct {
	Type         *string          `json:"type"`
	DockerVolume DockerVolumeJSON `json:"docker_volume"`
}

type DockerVolumeJSON struct {
	Driver        *string             `json:"driver"`
	Name          *string             `'json:"name"`
	DriverOptions []map[string]string `'json:"driveropts"`
}

type NetworkJSON struct {
	IpAddresses []IpAddressJSON     `json:"ipaddress,omitempty"`
	Name        *string             `json:"name"`
	Groups      []string            `json:"group"`
	Labels      []map[string]string `json:"labels"`
	PortMapping []*PortMapping      `json:"portmapping"`
}

type PortMapping struct {
	HostPort      *uint32 `json:"hostport"`
	ContainerPort *uint32 `json:"containerport"`
	Protocol      *string `json:"protocol"`
}

type IpAddressJSON struct {
	IP       *string `json:"ip"`
	Protocol *string `json:"protocol"`
}

type UriJSON struct {
	Uri     *string `json:"uri"`
	Extract *bool   `json:"extract"`
	Execute *bool   `json:"execute"`
}
