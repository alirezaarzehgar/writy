package cache

import (
	"log/slog"

	"github.com/alirezaarzehgar/writy/internal/keyval"
)

type cache struct {
	logger  *slog.Logger
	storage map[string]any
}

func New() keyval.KeyVal {
	return &cache{
		storage: make(map[string]any),
		logger:  slog.Default(),
	}
}

func (c *cache) WithLogHandler(handler slog.Handler) keyval.KeyVal {
	c.logger = slog.New(handler)
	return c
}

func (c *cache) Set(key, value string) error {
	if _, ok := c.storage[key]; ok {
		c.logger.Debug("duplicated key")
		return duplicatedError{key}
	}
	return c.ForceSet(key, value)
}

func (c *cache) ForceSet(key string, value any) error {
	c.storage[key] = value
	return nil
}

func (c *cache) Get(key string) (any, error) {
	value, ok := c.storage[key]
	if !ok {
		c.logger.Debug("key not found")
		return "", notfoundError{key}
	}
	return value, nil
}

func (c *cache) Del(key string) error {
	_, ok := c.storage[key]
	if !ok {
		c.logger.Debug("key not found")
		return &notfoundError{key}
	}

	delete(c.storage, key)
	return nil
}

func (c *cache) Clear() error {
	c.logger.Debug("clear storage")
	c.storage = make(map[string]any)
	return nil
}
