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

package ha

type Status string

// Define a list of states an HA node can be in.
const (
	Election  Status = "Election"
	Leading   Status = "Leading"
	Talking   Status = "Talking"
	Listening Status = "Listening"
)

// Interface for a single Node in an HA configuration.
type Node interface {
	Name() (string, error)   // How do we uniquely identify this node.
	Status() (Status, error) // What is the status of this node?
	Communicate()            // Our communication channel
	Election()               // Kick off an election.
	CreateLeader() error
	GetLeader() (string, error)
}
