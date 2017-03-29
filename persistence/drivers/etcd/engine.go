package etcd

import (
	"errors"
	"mesos-framework-sdk/persistence"
)

// The Storage interface takes in an "Engine" whicsh contains the client
// to the backend and a description of the driver.
type EtcdEngine struct {
	engine *Etcd
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
		if err := e.engine.Create(key, args[0]); err != nil {
			return err
		}
		// Multiple k,v pair
	} else if len(args) > 2 {
		if len(args)%2 == 0 {

			// First set of args.
			if err := e.engine.Create(key, args[0]); err != nil {
				return err
			}

			// Next two args is k,v
			for i := 1; i < len(args)-1; i += 2 {
				if err := e.engine.Create(args[i], args[i+1]); err != nil {
					return err
				}
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
func (e *EtcdEngine) Read(r ...string) (results []string, err error) {
	if len(r) == 1 {
		val, err := e.engine.Read(r[0])
		if err != nil {
			return results, err
		}
		results = append(results, val)
	} else if len(r) >= 2 {
		for _, v := range r {
			val, err := e.engine.Read(v)
			if err != nil {
				return results, err
			}
			results = append(results, val)
		}
	} else if len(r) == 0 {
		return nil, errors.New("No read parameters passed in.")
	}
	return results, nil
}

// Variadic k,v update.
func (e *EtcdEngine) Update(key string, args ...string) error {
	// Single k,v pair.
	if len(args) == 1 {

		// Multiple k,v pair.
		if err := e.engine.Update(key, args[0]); err != nil {
			return err
		}
	} else if len(args) > 2 {
		if len(args)%2 == 0 {

			// First set of args.
			if err := e.engine.Update(key, args[0]); err != nil {
				return err
			}

			for i := 1; i < len(args)-1; i += 2 {

				// Next two args is k,v.
				if err := e.engine.Update(args[i], args[i+1]); err != nil {
					return err
				}
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
		if err := e.engine.Delete(key); err != nil {
			return err
		}
	} else if len(args) > 0 {

		// Delete first key.
		if err := e.engine.Delete(key); err != nil {
			return err
		}
		for _, k := range args { // then remaining keys
			if err := e.engine.Delete(k); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *EtcdEngine) Driver() string {
	return e.driver
}

func (e *EtcdEngine) Engine() interface{} {
	return e.engine
}
