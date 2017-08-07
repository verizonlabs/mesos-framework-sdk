package labels

import (
	"errors"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/utils"
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
