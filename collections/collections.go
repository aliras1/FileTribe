package collections

import (
	"sync"
	"github.com/golang/glog"
)


type ICollectionItem interface {
	Id() IIdentifier
}

type ConcurrentCollection struct {
	lock sync.RWMutex
	list []ICollectionItem
}

func NewConcurrentCollection() *ConcurrentCollection {
	return &ConcurrentCollection{
		list: []ICollectionItem{},
	}
}

func (c *ConcurrentCollection) Append(item ICollectionItem) {
	c.lock.Lock()
	defer c.lock.Unlock()

	//if (len(c.list) > 0) && (reflect.TypeOf(item) != reflect.TypeOf(c.list[0])) {
	//	glog.Error("could not append to concurrent collection: invalid type")
	//	return
	//}

	for _, elem := range c.list {
		if elem.Id().Equal(item.Id()) {
			glog.Error("could not append to concurrent collection: item already exists")
			return
		}
	}

	c.list = append(c.list, item)
}

func (c *ConcurrentCollection) Iterator() <- chan interface{} {
	channel := make(chan interface{})

	go func() {
		c.lock.RLock()
		defer c.lock.RUnlock()
		defer close(channel)

		for _, item := range c.list {
			channel <- item
		}
	}()

	return channel
}

func (c *ConcurrentCollection) DeleteItem(item interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for i, elem := range c.list {
		if elem == item {
			c.list = append(c.list[:i], c.list[i+1: ]...)
			return
		}
	}
}

func (c *ConcurrentCollection) DeleteWithId(id IIdentifier) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for i, elem := range c.list {
		if elem.Id().Equal(id) {
			c.list = append(c.list[:i], c.list[i+1: ]...)
			return
		}
	}
}

func (c *ConcurrentCollection) Get(id IIdentifier) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, elem := range c.list {
		if elem.Id().Equal(id) {
			return elem
		}
	}

	return nil
}

func (c *ConcurrentCollection) FirstOrDefault(id IIdentifier) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if len(c.list) == 0 {
		return nil
	}

	if id == nil {
		return c.list[0]
	}

	for _, elem := range c.list {
		if elem.Id().Equal(id) {
			return elem
		}
	}

	return c.list[0]
}

func (c *ConcurrentCollection) Count() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return len(c.list)
}