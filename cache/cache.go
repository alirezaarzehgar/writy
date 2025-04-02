package cache

import "log/slog"

type StorageType map[string]any

type Cache struct {
	storage StorageType
}

func New() *Cache {
	return &Cache{storage: make(map[string]any)}
}

func (c *Cache) Set(key, value string) error {
	if _, ok := c.storage[key]; ok {
		slog.Debug("duplicated key")
		return duplicatedError{key}
	}
	return c.ForceSet(key, value)
}

func (c *Cache) ForceSet(key string, value any) error {
	c.storage[key] = value
	return nil
}

func (c *Cache) Get(key string) (any, error) {
	value, ok := c.storage[key]
	if !ok {
		slog.Debug("key not found")
		return "", notfoundError{key}
	}
	return value, nil
}

func (c *Cache) Del(key string) error {
	_, ok := c.storage[key]
	if !ok {
		slog.Debug("key not found")
		return notfoundError{key}
	}

	delete(c.storage, key)
	return nil
}

func (c *Cache) Clear() error {
	slog.Debug("clear storage", "storage", c.storage)
	c.storage = make(map[string]any)
	return nil
}

func (c *Cache) List() (StorageType, error) {
	return c.storage, nil
}
