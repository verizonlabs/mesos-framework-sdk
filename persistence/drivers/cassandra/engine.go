package cassandra

import (
	"mesos-framework-sdk/persistence"
)

// Cassandra Engine
type CassandraEngine struct {
	engine *Cassandra
	driver string
}

// NOTE (tim): If we give the ability to create custom tables/ queries
// we will have to support all types of advanced queries for end users.
// It would make more sense to enforce a sensible table structure
// to hold framework data and punt the implementation to end users.
// Any additional data needed for a custom framework would have to
// use the lower level API if they need/want that level of control.
// Guide to follow for table design:
//
// http://www.datastax.com/dev/blog/basic-rules-of-cassandra-data-modeling
//
// Data we'll need to store:
// Task Information: TaskID (Primary key) and give back taskInfo, timers, filters (for launching),
//   retry attempts, ownership information.
// Framework information: FrameworkId, subscription time, Framework name.
// High-level queries we'll need to satisfy (Separate tables each):
//   What is the name of the framework we are using?
//   What is the Id of the framework we are using?
//   What is the TaskId of this TaskName?
//   What is the number of retries on this Task to launch?
//   What are the filters, if any, on this task?
//   How frequent should we retry to launch this task if it fails?
//   How long has this framework been up?
//   Who is the leader
//   ...etc
func NewCassandraEngine(engine *Cassandra) *CassandraEngine {
	return &CassandraEngine{
		engine: engine,
		driver: persistence.CASSANDRA_ENGINE_NAME,
	}
}

// Interface type
func (c *CassandraEngine) Create(table string, args ...string) error {
	// TODO (tim): We will need to consider a good
	// table design before finishing this driver
	// c.engine.Create(table, args)
	return nil
}

func (c *CassandraEngine) Read(...string) ([]string, error) {
	//c.engine.Read()
	return nil, nil
}

func (c *CassandraEngine) Update(string, ...string) error {
	//c.engine.Update()
	return nil
}

func (c *CassandraEngine) Delete(string, ...string) error {
	//c.engine.Delete()
	return nil
}

func (s *CassandraEngine) Driver() string {
	return s.driver
}
