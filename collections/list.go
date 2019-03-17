package collections

import "sync"

type List struct {
	data []interface{}
	lock sync.RWMutex
}

func NewConcurrentList() *List {
	return &List{
		data: []interface{}{},
	}
}

func (l *List) Add(item interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.data = append(l.data, item)
}

func (l *List) Iterator() <- chan interface{} {
	l.lock.RLock()
	defer l.lock.RUnlock()

	channel := make(chan interface{})

	go func() {
		l.lock.RLock()
		defer l.lock.RUnlock()
		defer close(channel)

		for _, item := range l.data {
			channel <- item
		}
	}()

	return channel
}

func (l *List) Get(i int) interface{} {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i >= len(l.data) {
		return nil
	}

	return l.data[i]
}

func (l *List) Count() int {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return len(l.data)
}