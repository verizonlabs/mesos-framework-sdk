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

func Read(data io.ReadCloser) error {
	var event sched.Event
	reader := bufio.NewReader(data)

	for {
		lengthStr, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		lengthInt, err := strconv.Atoi(strings.TrimRight(lengthStr, "\n"))
		if err != nil {
			return err
		}

		buffer := make([]byte, lengthInt)
		n, err := io.ReadFull(reader, buffer)
		if n != lengthInt {
			return errors.New("Bytes read are not equal to the message length")
		}

		err = proto.Unmarshal(buffer, &event)
		if err != nil {
			return err
		}
	}
}
