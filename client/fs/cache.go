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

package fs

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/aliras1/FileTribe/utils"
)

const (
	cacheFileName = "cache.dat"
	maxEntries = 100
)

type Cache struct {
	data map[string]string
	storage *Storage
}

func NewCache(storage *Storage) (*Cache, error) {
	c := &Cache{
		storage: storage,
	}

	var data map[string]string

	if utils.FileExists(c.path()) {
		encoded, err := ioutil.ReadFile(c.path())
		if err != nil {
			return nil, errors.Wrap(err, "could not read file")
		}

		json.Unmarshal(encoded, &data)
	} else {
		data = make(map[string]string)
	}

	c.data = data

	return c, nil
}

func (c *Cache) Put(key string, value string) error {
	c.data[key] = value

	if err := c.save(); err != nil {
		return errors.Wrap(err, "could not save cache")
	}

	return nil
}

func (c *Cache) Get(key string) (value string, ok bool) {
	value, ok = c.data[key]
	return
}

func (c *Cache) save() error {
	encoded, err := c.encode()
	if err != nil {
		return errors.Wrap(err, "could not encode cache")
	}

	if err := utils.WriteFile(c.path(), encoded); err != nil {
		return errors.Wrap(err, "could not write cache file")
	}

	return nil
}

func (c *Cache) encode() ([]byte, error) {
	encoded, err := json.Marshal(c.data)
	if err != nil {
		return nil, errors.Wrap(err, "could not json marshal cache data")
	}

	return encoded, nil
}

func (c *Cache) path() string {
	return c.storage.tmpPath + cacheFileName
}