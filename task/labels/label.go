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

package labels

import (
	"errors"
	"github.com/verizonlabs/mesos-framework-sdk/include/mesos_v1"
	"github.com/verizonlabs/mesos-framework-sdk/utils"
)

func ParseLabels(labels map[string]string) (*mesos_v1.Labels, error) {
	if labels == nil {
		return nil, nil
	}

	l := &mesos_v1.Labels{}
	for name, value := range labels {
		if name == "" || value == "" {
			return nil, errors.New("Empty key or value passed in")
		}
		l.Labels = append(l.Labels, &mesos_v1.Label{
			Key:   utils.ProtoString(name),
			Value: utils.ProtoString(value),
		})
	}

	return l, nil
}
