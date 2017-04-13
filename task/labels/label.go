package labels

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/include/mesos"
)

func ParseLabels(lbs []map[string]string) (*mesos_v1.Labels, error) {
	labels := make([]*mesos_v1.Label, 0)
	if lbs != nil {
		for _, labelList := range lbs {
			for k, v := range labelList {
				if k != "" || v != "" {
					label := &mesos_v1.Label{
						Key:   proto.String(k),
						Value: proto.String(v),
					}
					labels = append(labels, label)
				} else {
					return nil, errors.New("Empty key or value passed in.")
				}
			}
		}
	}
	return &mesos_v1.Labels{Labels: labels}, nil
}
