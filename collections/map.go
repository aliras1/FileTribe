package collections

import (
	"sync"
)

type Compare func(interface{}, interface{}) bool

type ICollectionItem interface {
	Id() IIdentifier
}

type KeyValuePair struct {
	Key interface{}
	Value interface{}
}

type Map struct {
	lock sync.RWMutex
	data map[interface{}]interface{}
}

func NewConcurrentMap() *Map {
	return &Map{
		data: make(map[interface{}]interface{}),
	}
}

func (c *Map) Put(key interface{}, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = value
}

func (c *Map) Reset() {
	c.data = make(map[interface{}]interface{})
}

func (c *Map) VIterator() <- chan interface{} {
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

func (c *Map) KIterator() <- chan interface{} {
	channel := make(chan interface{})

	go func() {
		c.lock.RLock()
		defer c.lock.RUnlock()
		defer close(channel)

		for k, _ := range c.data {
			channel <- k
		}
	}()

	return channel
}

func (c *Map) KVIterator() <- chan KeyValuePair {
	channel := make(chan KeyValuePair)

	go func() {
		c.lock.RLock()
		defer c.lock.RUnlock()
		defer close(channel)

		for k, v := range c.data {
			channel <- KeyValuePair{Key:k, Value:v}
		}
	}()

	return channel
}


func (c *Map) Delete(key interface{}) interface{} {
	c.lock.Lock()
	defer c.lock.Unlock()

	v := c.data[key]

	delete(c.data, key)

	return v
}

func (c *Map) Get(key interface{}) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.data[key]
}

//func (c *Map) FirstOrDefault(id IIdentifier) interface{} {
//	c.lock.RLock()
//	defer c.lock.RUnlock()
//
//	if len(c.list) == 0 {
//		return nil
//	}
//
//	if id == nil {
//		return c.list[0]
//	}
//
//	for _, elem := range c.list {
//		if elem.Id().Equal(id) {
//			return elem
//		}
//	}
//
//	return c.list[0]
//}

func (c *Map) ToList() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var l []interface{}
	for _, v := range c.data {
		l = append(l, v)
	}

	return l
}


func (c *Map) Count() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return len(c.data)
}