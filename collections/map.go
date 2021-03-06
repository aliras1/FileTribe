// Copyright (c) 2019 Laszlo Sari
//
// FileTribe is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// FileTribe is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package collections

import (
	"sync"
)

// KeyValuePair is used by the Map iterator
type KeyValuePair struct {
	Key   interface{}
	Value interface{}
}

// Map is a concurrent map
type Map struct {
	lock sync.RWMutex
	data map[interface{}]interface{}
}

// NewConcurrentMap creates a new Map
func NewConcurrentMap() *Map {
	return &Map{
		data: make(map[interface{}]interface{}),
	}
}

// Put puts a key value pair into the map
func (c *Map) Put(key interface{}, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = value
}

// Reset resets the map, losing all its previous data in the process
func (c *Map) Reset() {
	c.data = make(map[interface{}]interface{})
}

// VIterator is the value iterator of the map
func (c *Map) VIterator() <-chan interface{} {
	channel := make(chan interface{})

	go func() {
		c.lock.RLock()
		defer c.lock.RUnlock()
		defer close(channel)

		for _, v := range c.data {
			channel <- v
		}
	}()

	return channel
}

// KIterator is the key iterator of the map
func (c *Map) KIterator() <-chan interface{} {
	channel := make(chan interface{})

	go func() {
		c.lock.RLock()
		defer c.lock.RUnlock()
		defer close(channel)

		for k := range c.data {
			channel <- k
		}
	}()

	return channel
}

// KVIterator is the key-value iterator of the map
func (c *Map) KVIterator() <-chan KeyValuePair {
	channel := make(chan KeyValuePair)

	go func() {
		c.lock.RLock()
		defer c.lock.RUnlock()
		defer close(channel)

		for k, v := range c.data {
			channel <- KeyValuePair{Key: k, Value: v}
		}
	}()

	return channel
}

// Delete deletes an item with the specified key
func (c *Map) Delete(key interface{}) interface{} {
	c.lock.Lock()
	defer c.lock.Unlock()

	v := c.data[key]

	delete(c.data, key)

	return v
}

// Get gets an item from the map with the specified key
func (c *Map) Get(key interface{}) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.data[key]
}

// ToList converts the map into a list
func (c *Map) ToList() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var l []interface{}
	for _, v := range c.data {
		l = append(l, v)
	}

	return l
}

// Count returns how many elements are stored in the map
func (c *Map) Count() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return len(c.data)
}
