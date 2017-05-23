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
