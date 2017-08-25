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

package structures

import (
	"sync"
)

type DistributedMap interface {
	Set(key, value interface{}) DistributedMap
	Get(key interface{}) interface{}
	Delete(key interface{})
	Iterate() <-chan Item
	Length() int
}

type ConcurrentMap struct {
	data map[interface{}]interface{}
	sync.RWMutex
}

type Item struct {
	Key   interface{}
	Value interface{}
}

// Returns a new ConcurrentMap.
func NewConcurrentMap(size ...int) DistributedMap {
	if len(size) == 1 {
		return &ConcurrentMap{
			data: make(map[interface{}]interface{}, size[0]),
		}
	}

	return &ConcurrentMap{
		data: make(map[interface{}]interface{}),
	}
}

// Sets a value with an associated key.
func (c *ConcurrentMap) Set(key, value interface{}) DistributedMap {
	c.Lock()
	defer c.Unlock()

	c.data[key] = value

	return c
}

// Gets the value associated with the specified key.
func (c *ConcurrentMap) Get(key interface{}) interface{} {
	c.RLock()
	defer c.RUnlock()

	return c.data[key]
}

// Removes a value from the map.
func (c *ConcurrentMap) Delete(key interface{}) {
	c.Lock()
	defer c.Unlock()

	delete(c.data, key)
}

// Safely iterates over the map.
// Provides the key/values to a channel that is returned for use by the client.
func (c *ConcurrentMap) Iterate() <-chan Item {
	c.RLock()
	ch := make(chan Item, len(c.data))

	go func() {
		for key, value := range c.data {
			ch <- Item{
				Key:   key,
				Value: value,
			}
		}
		c.RUnlock()

		close(ch)
	}()

	return ch
}

// Gives the number of items in the map.
func (c *ConcurrentMap) Length() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.data)
}
