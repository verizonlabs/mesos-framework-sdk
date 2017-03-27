package etcd

import (
	"errors"
	"mesos-framework-sdk/persistence"
)

// The Storage interface takes in an "Engine" whicsh contains the client
// to the backend and a description of the driver.
type EtcdEngine struct {
	engine Etcd
	driver string
}

// Return a new ETCD engine.
func NewEtcdEngine(engine *Etcd) *EtcdEngine {
	return &EtcdEngine{
		engine: engine,
		driver: persistence.ETCD_ENGINE_NAME,
	}
}

// Variadic k/v create
func (e *EtcdEngine) Create(key string, args ...string) error {
	// Single k,v pair
	if len(args) == 1 {
		e.engine.Create(key, args[0])
		// Multiple k,v pair
	} else if len(args) > 2 {
		if len(args)%2 == 0 {
			e.engine.Create(key, args[0]) // First set of args.
			for i := 1; i < len(args); i += 2 {
				e.engine.Create(args[i], args[i+1]) // Next two args is k,v
			}
		} else {
			// Each key needs a value, so disregard odd numbered variadic arguments.
			return errors.New("Multiple creates must be an even number of arguments.")
		}
	} else if len(args) == 0 {
		return errors.New("No value given for key.")
	}
	return nil
}

// Variadic k/v read
func (e *EtcdEngine) Read(r ...string) (results []string, _ error) {
	if len(r) == 1 {
		results = append(results, e.engine.Read(r[0]))
	} else if len(r) >= 2 {
		for _, v := range r {
			results = append(results, e.engine.Read(v))
		}
	} else if len(r) == 0 {
		return nil, errors.New("No read parameters passed in.")
	}
	return results, nil
}

// Variadic k,v update.
func (e *EtcdEngine) Update(key string, args ...string) error {
	// Single k,v pair
	if len(args) == 1 {
		e.engine.Update(key, args[0])
		// Multiple k,v pair
	} else if len(args) > 2 {
		if len(args)%2 == 0 {
			e.engine.Update(key, args[0]) // First set of args.
			for i := 1; i < len(args); i += 2 {
				e.engine.Update(args[i], args[i+1]) // Next two args is k,v
			}
		} else {
			// Each key needs a value, so disregard odd numbered variadic arguments.
			return errors.New("Multiple updates must be an even number of arguments.")
		}
	} else if len(args) == 0 {
		return errors.New("No value given for key.")
	}
	return nil
}

func (e *EtcdEngine) Delete(key string, args ...string) error {
	if len(args) == 0 {
		e.engine.Delete(key)
	} else if len(args) > 0 {
		e.engine.Delete(key)     // delete first key
		for _, k := range args { // then remaining keys
			e.engine.Delete(k)
		}
	}
	return nil
}

func (e *EtcdEngine) Driver() string {
	return e.driver
}
