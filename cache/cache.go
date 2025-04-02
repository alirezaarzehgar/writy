package cache

import "fmt"

type StorageType map[string]any

type Cache struct {
	storage StorageType
}

func New() *Cache {
	return &Cache{storage: make(map[string]any)}
}

func (c *Cache) Set(key, value string) error {
	if _, ok := c.storage[key]; ok {
		return fmt.Errorf("duplicated: %s", key)
	}
	return c.ForceSet(key, value)
}

func (c *Cache) ForceSet(key string, value any) error {
	c.storage[key] = value
	return nil
}

func (c *Cache) Get(key string) any {
	value, ok := c.storage[key]
	if !ok {
		return nil
	}
	return value
}

func (c *Cache) Del(key string) error {
	_, ok := c.storage[key]
	if !ok {
		return fmt.Errorf("not found: %s", key)
	}

	delete(c.storage, key)
	return nil
}

func (c *Cache) Clear() {
	c.storage = make(map[string]any)
}

func (c *Cache) List() StorageType {
	return c.storage
}
