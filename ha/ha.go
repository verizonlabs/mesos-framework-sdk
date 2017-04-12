package ha

// The HA package will handle electing a leader or leaders out a cluster of N nodes.
// Each node will talk to the leader(s) to make sure they are up.
// If a leader dies, election occurs and a new leader is elected.

// A node may rejoin a cluster but never leave.
//

type Status string

const (
	Election Status = "Election"
	Leading Status = "Leading"
	Talking Status = "Talking"
	Listening Status = "Listening"
)


// Any node in the structure.
// All participants in the structure will satisfy the node interface
// The node will handle talking to other nodes
// via a Communications channel.
//
// A list of seed nodes will be passed in by configuration.
//
// Once connection is made to another node, they let each other know about
//   their state.  If they are still electing, and haven't contacted all other
//   nodes, they attempt to reach all other nodes.

// Once all nodes are reached, election occurs.
// Election logic may be done via a configurable value.
// Default will be to create a random hash, and
// use that to create an ordering.
// Otherwise other values can be passed in to be hashed
// to determine "ordering".

// Once nodes are connected, they search for other nodes.
type Node interface {
	Name() (string, error)       // How do we uniquely identify this node.
	Status() (Status, error)    // What is the status of this node?
	Communicate()  // Our communication channel
	Election() // Kick off an election.
}