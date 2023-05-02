package Utilities

import (
	cache "github.com/aleksiumish/in-memory-cache"
)

var c *cache.Cache

func SetKey(key string, data interface{}) bool {

	if c == nil {
		c = cache.NewCache()
	}
	if GetKeyValue(key) != nil {
		c.Delete(key)
	}
	c.Set(key, data)
	return true
}

func GetKeyValue(Key string) interface{} {
	if c == nil {
		return nil
	}
	return c.Get(Key)
}
