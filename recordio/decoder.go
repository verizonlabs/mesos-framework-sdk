// Copyright 2017 Verizon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package recordio

import (
	"bufio"
	"errors"
	"io"
	"mesos-framework-sdk/include/mesos_v1_executor"
	"mesos-framework-sdk/include/mesos_v1_scheduler"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
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

		buffer := make([]byte, 0, lengthInt)
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
