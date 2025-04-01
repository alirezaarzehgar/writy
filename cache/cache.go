package cache

import (
	"log/slog"
)

type StorageType map[string]any

type Cache struct {
	logger  *slog.Logger
	storage StorageType
}

func New() *Cache {
	return &Cache{
		storage: make(map[string]any),
		logger:  slog.Default(),
	}
}

func (c *Cache) SetLogHandler(handler slog.Handler) *Cache {
	c.logger = slog.New(handler)
	return c
}

func (c *Cache) Set(key, value string) error {
	if _, ok := c.storage[key]; ok {
		c.logger.Debug("duplicated key")
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
		c.logger.Debug("key not found")
		return "", notfoundError{key}
	}
	return value, nil
}

func (c *Cache) Del(key string) error {
	_, ok := c.storage[key]
	if !ok {
		c.logger.Debug("key not found")
		return notfoundError{key}
	}

	delete(c.storage, key)
	return nil
}

func (c *Cache) Clear() error {
	c.logger.Debug("clear storage", "storage", c.storage)
	c.storage = make(map[string]any)
	return nil
}

func (c *Cache) List() (StorageType, error) {
	return c.storage, nil
}
