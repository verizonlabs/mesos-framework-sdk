package events

import (
	"bufio"
	"github.com/golang/protobuf/proto"
	"io"
	sched "mesos-framework-sdk/include/scheduler"
	"strconv"
	"strings"
)

func Loop(data io.ReadCloser, events chan<- *sched.Event) {
	var event sched.Event
	reader := bufio.NewReader(data)

	for {
		lengthStr, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		lengthInt, err := strconv.Atoi(strings.TrimRight(lengthStr, "\n"))
		if err != nil {
			continue
		}

		buffer := make([]byte, lengthInt)
		n, err := io.ReadFull(reader, buffer)
		if n != lengthInt {
			continue
		}

		err = proto.Unmarshal(buffer, &event)
		if err != nil {
			continue
		}

		events <- &event
	}
}
