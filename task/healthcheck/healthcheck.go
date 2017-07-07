package healthcheck

import (
	"errors"
	"fmt"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/task"
	"strings"
	"mesos-framework-sdk/utils"
)

var (
	NoHealthCheckType      error = errors.New("No Health check type specified.")
	InvalidHealthCheckType error = errors.New("Invalid health check type, accepted values are tcp, http, command.")
	InvalidPortRange       error = errors.New(fmt.Sprintf("Invalid port range given, %v - %v accepted", MIN_PORT, MAX_PORT))
	UnsupportedScheme      error = errors.New("Unsupported scheme, supported schemes are http, https")
	NoHTTPPath             error = errors.New("No http path given, must give at a minimum a path to hit for http.")
)

const (
	MINIMUM_INTERVAL_SECONDS     float64 = 15.0
	MINIMUM_TIMEOUT_SECONDS      float64 = 1.0
	MINIMUM_GRACE_PERIOD_SECONDS float64 = 1.0
	MINIMUM_CONSECUTIVE_FAILURES uint32  = 1.0
	MIN_PORT                     int     = 0
	MAX_PORT                     int     = 65535
)

func ParseHealthCheck(json *task.HealthCheckJSON) (*mesos_v1.HealthCheck, error) {
	if json.Type == nil {
		return nil, NoHealthCheckType
	}

	hc := &mesos_v1.HealthCheck{}
	switch strings.ToLower(*json.Type) {
	case "tcp":
		hc.Type = mesos_v1.HealthCheck_TCP.Enum()
		tcp, err := parseTcpHealthCheck(json.Tcp)
		if err != nil {
			return nil, err
		}
		hc.Tcp = tcp
	case "http":
		hc.Type = mesos_v1.HealthCheck_HTTP.Enum()
		http, err := parseHTTPHealthCheck(json.Http)
		if err != nil {
			return nil, err
		}
		hc.Http = http
	case "command":
		hc.Type = mesos_v1.HealthCheck_COMMAND.Enum()
		cmd, err := parseCommandHealthCheck(json.Command)
		if err != nil {
			return nil, err
		}
		hc.Command = cmd
	default:
		return nil, InvalidHealthCheckType
	}

	if json.TimeoutSeconds != nil && *json.TimeoutSeconds >= MINIMUM_TIMEOUT_SECONDS {
		hc.TimeoutSeconds = json.TimeoutSeconds
	}
	if json.IntervalSeconds != nil && *json.IntervalSeconds >= MINIMUM_INTERVAL_SECONDS {
		hc.IntervalSeconds = json.IntervalSeconds
	}
	if json.GracePeriodSeconds != nil && *json.GracePeriodSeconds >= MINIMUM_GRACE_PERIOD_SECONDS {
		hc.GracePeriodSeconds = json.GracePeriodSeconds
	}
	if json.ConsecutiveFailures != nil && *json.ConsecutiveFailures >= MINIMUM_CONSECUTIVE_FAILURES {
		hc.ConsecutiveFailures = json.ConsecutiveFailures
	}

	return hc, nil
}

func parseTcpHealthCheck(json *task.TCPHealthCheck) (*mesos_v1.HealthCheck_TCPCheckInfo, error) {
	tcp := &mesos_v1.HealthCheck_TCPCheckInfo{}
	if json.Port > MIN_PORT && json.Port < MAX_PORT {
		tcp.Port = utils.ProtoUint32(uint32(json.Port))
	} else {
		return nil, InvalidPortRange
	}
	return tcp, nil
}

func parseHTTPHealthCheck(json *task.HTTPHealthCheck) (*mesos_v1.HealthCheck_HTTPCheckInfo, error) {
	http := &mesos_v1.HealthCheck_HTTPCheckInfo{}
	if json.Scheme != nil {
		switch strings.ToLower(*json.Scheme) {
		case "http", "https":
			http.Scheme = json.Scheme
		default:
			return nil, UnsupportedScheme
		}
	} else {
		// Assume HTTPS.
		json.Scheme = utils.ProtoString("https")
	}

	if json.Path != nil {
		http.Path = json.Path
	} else {
		// We need a path to hit
		return nil, NoHTTPPath
	}

	if json.Port != nil {
		http.Port = utils.ProtoUint32(uint32(*json.Port))
	}
	// What statuses are accepted.
	if len(json.Statuses) > 0 {
		http.Statuses = json.Statuses
	}

	return http, nil
}

func parseCommandHealthCheck(json *task.CommandJSON) (*mesos_v1.CommandInfo, error) {
	cmd := &mesos_v1.CommandInfo{}

	if json.Cmd != nil {
		cmd.Value = json.Cmd
	}

	if len(json.Uris) > 0 {
		uris := []*mesos_v1.CommandInfo_URI{}
		for _, u := range json.Uris {
			uris = append(uris, &mesos_v1.CommandInfo_URI{
				Extract:    u.Extract,
				Executable: u.Execute,
				Value:      u.Uri,
			})
		}
		cmd.Uris = uris
	}

	if json.Environment != nil {
		env := &mesos_v1.Environment{
			Variables: make([]*mesos_v1.Environment_Variable, 0),
		}
		for _, kv := range json.Environment.Variables {
			for k, v := range kv {
				env.Variables = append(env.Variables, &mesos_v1.Environment_Variable{
					Name:  utils.ProtoString(k),
					Value: utils.ProtoString(v),
				})
			}
		}
		cmd.Environment = env
	}

	return cmd, nil
}
