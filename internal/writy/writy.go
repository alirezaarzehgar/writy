package writy

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/alirezaarzehgar/writy/cache"
	"github.com/alirezaarzehgar/writy/internal/keyval"
)

var (
	DefaultStoragePath string        = fmt.Sprintf("%s/writy.db", os.Getenv("HOME"))
	DefaultFlushCycle  time.Duration = time.Second * 5
)

type writy struct {
	logger        *slog.Logger
	storageReader *os.File
	storageWriter *os.File
	flusher       *Flusher
	cache         keyval.KeyVal
}

func New(path string, exp time.Duration) (keyval.KeyVal, error) {
	if path == "" {
		path = DefaultStoragePath
	}

	wReader, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable open file %s for writing: %w", path, err)
	}

	sReader, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable open file %s for reading: %w", path, err)
	}

	w := &writy{
		logger:        slog.Default(),
		storageReader: sReader,
		storageWriter: wReader,
		flusher:       NewFlusher(exp),
		cache:         cache.New(),
	}

	w.flusher.Run(w)

	return w, nil
}

func (w *writy) WithLogHandler(handler slog.Handler) keyval.KeyVal {
	w.logger = slog.New(handler)
	w.flusher.logger = w.logger
	w.cache.WithLogHandler(handler)
	w.logger.Debug("enable logger")
	return w
}

func (w writy) Set(key, value string) error {
	return w.cache.Set(key, value)
}

func (w writy) ForceSet(key string, value any) error {
	return w.cache.ForceSet(key, value)
}

func (w writy) Get(key string) (any, error) {
	// search in the cache
	v, err := w.cache.Get(key)
	if !cache.IsNotFound(err) {
		w.logger.Debug("key found", "key", key, "value", v, "err", err)
		return v, nil
	}

	// search in the fs
	return "", err
}

func (w writy) Del(key string) error {
	return nil
}

func (w writy) Clear() error {
	return nil
}
