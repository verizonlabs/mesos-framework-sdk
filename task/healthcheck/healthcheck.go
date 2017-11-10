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

package healthcheck

import (
	"errors"
	"fmt"
	"github.com/verizonlabs/mesos-framework-sdk/include/mesos_v1"
	"github.com/verizonlabs/mesos-framework-sdk/task"
	"github.com/verizonlabs/mesos-framework-sdk/utils"
	"strings"
)

var (
	NoHealthCheckType      error = errors.New("No Health check type specified.")
	InvalidHealthCheckType error = errors.New("Invalid health check type, accepted values are tcp, http, command.")
	InvalidPortRange       error = errors.New(fmt.Sprintf("Invalid port range given, %v - %v accepted", MIN_PORT, MAX_PORT))
	UnsupportedScheme      error = errors.New("Unsupported scheme, supported schemes are http, https")
	NoHTTPPath             error = errors.New("No http path given, must give at a minimum a path to hit for http.")
	NoTCPHealthCheck       error = errors.New("No TCP health check was defined")
	NoHTTPHealthCheck      error = errors.New("No HTTP health check was defined")
	NoCommandHealthCheck   error = errors.New("No error health check was defined")
)

const (
	MINIMUM_INTERVAL_SECONDS     float64 = 15.0
	MINIMUM_TIMEOUT_SECONDS      float64 = 1.0
	MINIMUM_GRACE_PERIOD_SECONDS float64 = 1.0
	MINIMUM_CONSECUTIVE_FAILURES uint32  = 1.0
	MIN_PORT                     int     = 0
	MAX_PORT                     int     = 65535
)

func ParseHealthCheck(json *task.HealthCheckJSON, c *mesos_v1.CommandInfo) (*mesos_v1.HealthCheck, error) {
	if json == nil {
		return nil, nil
	}

	if json.Type == nil {
		return nil, NoHealthCheckType
	}

	hc := &mesos_v1.HealthCheck{}
	switch strings.ToLower(*json.Type) {
	case "tcp":
		if json.Tcp == nil {
			return nil, NoTCPHealthCheck
		}
		hc.Type = mesos_v1.HealthCheck_TCP.Enum()

		tcp, err := parseTcpHealthCheck(json.Tcp)
		if err != nil {
			return nil, err
		}

		hc.Tcp = tcp
	case "http":
		if json.Http == nil {
			return nil, NoHTTPHealthCheck
		}

		hc.Type = mesos_v1.HealthCheck_HTTP.Enum()
		http, err := parseHTTPHealthCheck(json.Http)
		if err != nil {
			return nil, err
		}

		hc.Http = http
	case "command":
		hc.Type = mesos_v1.HealthCheck_COMMAND.Enum()
		hc.Command = c
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
