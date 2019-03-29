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