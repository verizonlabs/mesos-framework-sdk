package recordio

import (
	"bufio"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
	sched "mesos-framework-sdk/include/scheduler"
	"strconv"
	"strings"
)

// Decode continually reads and constructs events from the Mesos stream.
func Decode(data io.ReadCloser, events chan<- *sched.Event) error {
	reader := bufio.NewReader(data)

	for {
		var event sched.Event
		lengthStr, err := reader.ReadString('\n')
		if err != nil {
			return errors.New("Failed to read RecordIO message length: " + err.Error())
		}

		lengthInt, err := strconv.Atoi(strings.TrimRight(lengthStr, "\n"))
		if err != nil {
			return errors.New("RecordIO message length is not a number: " + err.Error())
		}

		buffer := make([]byte, lengthInt)
		n, err := io.ReadFull(reader, buffer)
		if n != lengthInt {
			return errors.New("Amount of bytes read does not match the RecordIO message length")
		}

		err = proto.Unmarshal(buffer, &event)
		if err != nil {
			return errors.New("Failed to decode event: " + err.Error())
		}

		events <- &event
	}
}
