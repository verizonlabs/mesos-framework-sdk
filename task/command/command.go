package command

import (
	"errors"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/task"
	"mesos-framework-sdk/utils"
)

func ParseCommandInfo(cmd *task.CommandJSON) (*mesos_v1.CommandInfo, error) {
	if cmd == nil {
		return nil, errors.New("Empty commandInfo.")
	}

	mesosCmd := &mesos_v1.CommandInfo{
		Value:       cmd.Cmd,
		Environment: &mesos_v1.Environment{},
	}
	uriList := []*mesos_v1.CommandInfo_URI{}

	if cmd.Environment != nil {
		for name, value := range cmd.Environment {
			mesosCmd.Environment.Variables = append(mesosCmd.Environment.Variables, &mesos_v1.Environment_Variable{
				Name:  utils.ProtoString(name),
				Value: utils.ProtoString(value),
			})
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
