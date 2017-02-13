package events

import (
	"bufio"
	"github.com/golang/protobuf/proto"
	"io"
	sched "mesos-framework-sdk/include/scheduler"
	"strconv"
	"strings"
)

func Loop(data io.ReadCloser, events chan<- *sched.Event) error {
	var event sched.Event
	reader := bufio.NewReader(data)

	for {
		lengthStr, err := reader.ReadString('\n')
		if err != nil {
			events <- nil
		}

		lengthInt, err := strconv.Atoi(strings.TrimRight(lengthStr, "\n"))
		if err != nil {
			events <- nil
		}

		buffer := make([]byte, lengthInt)
		n, err := io.ReadFull(reader, buffer)
		if n != lengthInt {
			events <- nil
		}

		err = proto.Unmarshal(buffer, &event)
		if err != nil {
			events <- nil
		}

		events <- &event
	}
}
