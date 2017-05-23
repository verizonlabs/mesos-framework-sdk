package command

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/task"
)

func ParseCommandInfo(cmd *task.CommandJSON) (*mesos_v1.CommandInfo, error) {
	if cmd == nil {
		return nil, errors.New("Empty commandInfo.")
	}
	mesosCmd := &mesos_v1.CommandInfo{}
	uriList := []*mesos_v1.CommandInfo_URI{}
	if cmd.Cmd != nil {
		// NOTE (tim): Should we split on white space and use first arg as "value" and the remainder as args?
		// A command value isn't required since the commandInfo can be used just to fetch URI's
		mesosCmd.Value = cmd.Cmd
	}

	if cmd.Environment != nil {
		vars := []*mesos_v1.Environment_Variable{}
		for _, env := range cmd.Environment.Variables {
			for k, v := range env {
				vars = append(vars, &mesos_v1.Environment_Variable{
					Name:  proto.String(k),
					Value: proto.String(v),
				})
			}
		}
		mesosCmd.Environment = &mesos_v1.Environment{
			Variables: vars,
		}
	}

	if len(cmd.Uris) > 0 {
		// create all the URI'
		for _, uri := range cmd.Uris {
			uriList = append(uriList, &mesos_v1.CommandInfo_URI{
				Value:      uri.Uri,
				Executable: uri.Execute,
				Extract:    uri.Extract,
			})
		}
		mesosCmd.Uris = uriList
	}

	if len(mesosCmd.Uris) == 0 && cmd.Cmd == nil {
		return nil, errors.New("CommandInfo is empty even though a command JSON param was passed in.")
	}
	return mesosCmd, nil
}
