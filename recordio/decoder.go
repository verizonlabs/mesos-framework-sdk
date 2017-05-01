package recordio

import (
	"bufio"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
	"mesos-framework-sdk/include/mesos_v1_executor"
	"mesos-framework-sdk/include/mesos_v1_scheduler"
	"strconv"
	"strings"
)

// Decode continually reads and constructs events from the Mesos stream.
func Decode(data io.ReadCloser, events interface{}) error {
	reader := bufio.NewReader(data)

	for {
		lengthStr, err := reader.ReadString('\n')
		if err != nil {
			return err
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

		switch events := events.(type) {
		case chan *mesos_v1_scheduler.Event:
			var event mesos_v1_scheduler.Event
			err := proto.Unmarshal(buffer, &event)
			if err != nil {
				return errors.New("Failed to decode event: " + err.Error())
			}

			events <- &event
		case chan *mesos_v1_executor.Event:
			var event mesos_v1_executor.Event
			err := proto.Unmarshal(buffer, &event)
			if err != nil {
				return errors.New("Failed to decode event: " + err.Error())
			}

			events <- &event
		}
	}
}
