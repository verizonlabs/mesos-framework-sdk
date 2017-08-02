package task

type ApplicationJSON struct {
	Name        string              `json:"name"`
	Instances   int                 `json:"instances"`
	Resources   *ResourceJSON       `json:"resources"`
	Command     *CommandJSON        `json:"command"`
	Container   *ContainerJSON      `json:"container"`
	HealthCheck *HealthCheckJSON    `json:"healthcheck"`
	Labels      []map[string]string `json:"labels"`
	Filters     []Filter            `json:"filters"`
	Retry       *TimeRetry          `json:"retry"`
}

type TimeRetry struct {
	Time       string `json:"time"`
	Backoff    bool   `json:"exp_backoff"`
	MaxRetries int    `json:"total_retries"`
}

type HealthCheckJSON struct {
	DelaySeconds        *float64         `json:"delay,omitempty"`
	IntervalSeconds     *float64         `json:"interval,omitempty"`
	TimeoutSeconds      *float64         `json:"timeout,omitempty"`
	ConsecutiveFailures *uint32          `json:"fails,omitempty"`
	GracePeriodSeconds  *float64         `json:"graceperiod,omitempty"`
	Type                *string          `json:"type,omitempty"`
	Command             *CommandJSON     `json:"command,omitempty"`
	Http                *HTTPHealthCheck `json:"http,omitempty"`
	Tcp                 *TCPHealthCheck  `json:"tcp,omitempty"`
	Endpoint            *string          `json:"endpoint,omitempty"`
}

type HTTPHealthCheck struct {
	Scheme   *string  `json:"scheme"`
	Port     *int32   `json:"port"`
	Path     *string  `json:"path"`
	Statuses []uint32 `json:"statuses"`
}

type TCPHealthCheck struct {
	Port int
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
	Disk Disk    `json:"disk"`
	Role string  `json:"role"`
}

type Disk struct {
	Size        float64          `json:"size"`
	Persistence *DiskPersistence `json:"persistence"`
	Volume      *VolumesJSON     `json:"volume"`
	Source      *DiskSource      `json:"source"`
}

type DiskSource struct {
	Type  *string `json:"type"`
	Path  *string `json:"path"`
	Mount *string `json:"mount"`
}

type DiskPersistence struct {
	Id        *string `json:"id"`
	Principle *string `json:"principle"`
}

type CommandJSON struct {
	Cmd         *string           `json:"cmd"`
	Uris        []UriJSON         `json:"uris"`
	Environment map[string]string `json:"environment"`
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
	DriverOptions []map[string]string `'json:"driver_opts"`
}

type NetworkJSON struct {
	IpAddresses []IpAddressJSON     `json:"ipaddress,omitempty"`
	Name        *string             `json:"name"`
	Groups      []string            `json:"group"`
	Labels      []map[string]string `json:"labels"`
	PortMapping []*PortMapping      `json:"port_mapping"`
}

type PortMapping struct {
	HostPort      *uint32 `json:"host_port"`
	ContainerPort *uint32 `json:"container_port"`
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
