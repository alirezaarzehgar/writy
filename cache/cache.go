package cache

import (
	"log/slog"
	"os"
)

type Cache struct {
	logger  *slog.Logger
	storage map[string]any
}

func New() Cache {
	return Cache{
		storage: make(map[string]any),
		logger:  slog.Default(),
	}
}

func (c *Cache) SetLogOptions(opts slog.HandlerOptions) {
	c.logger = slog.New(slog.NewTextHandler(os.Stderr, &opts))
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
		return &notfoundError{key}
	}

	delete(c.storage, key)
	return nil
}

func (c *Cache) Clear() {
	c.logger.Debug("clear storage")
	c.storage = make(map[string]any)
}
