package fs

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"ipfs-share/utils"
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