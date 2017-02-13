package events

import (
	"bufio"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
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
			log.Println("Failed to read length prefix from RecordIO: " + err.Error())
			continue
		}

		lengthInt, err := strconv.Atoi(strings.TrimRight(lengthStr, "\n"))
		if err != nil {
			log.Println("Failed to convert the message length: " + err.Error())
			continue
		}

		buffer := make([]byte, lengthInt)
		n, err := io.ReadFull(reader, buffer)
		if n != lengthInt {
			log.Println("Failed to read the RecordIO message: " + err.Error())
			continue
		}

		err = proto.Unmarshal(buffer, &event)
		if err != nil {
			log.Println("Failed to unmarshal event: " + err.Error())
			continue
		}

		events <- &event
	}
}
