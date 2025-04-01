package writy

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/alirezaarzehgar/writy/cache"
	"github.com/alirezaarzehgar/writy/internal/keyval"
)

var (
	DefaultStoragePath string        = fmt.Sprintf("%s", os.Getenv("HOME"))
	DefaultFlushCycle  time.Duration = time.Second * 5
)

type writy struct {
	logger        *slog.Logger
	storageReader *os.File
	storageWriter *os.File
	indexReader   *os.File
	indexWriter   *os.File
	flusher       *Flusher
	cache         keyval.KeyVal
}

func New(path string, exp time.Duration) (keyval.KeyVal, error) {
	if path == "" {
		path = DefaultStoragePath
	}

	storagePath := filepath.Join(path, "storage.db")
	sWriter, err := os.OpenFile(storagePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable open file %s for writing: %w", path, err)
	}
	sWriter.Seek(0, io.SeekEnd)

	sReader, err := os.Open(storagePath)
	if err != nil {
		return nil, fmt.Errorf("unable open file %s for reading: %w", path, err)
	}

	indexPath := filepath.Join(path, "index.db")
	iWriter, err := os.OpenFile(indexPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable open file %s for writing: %w", path, err)
	}

	iReader, err := os.Open(indexPath)
	if err != nil {
		return nil, fmt.Errorf("unable open file %s for reading: %w", path, err)
	}

	w := &writy{
		logger:        slog.Default(),
		storageReader: sReader,
		storageWriter: sWriter,
		indexReader:   iReader,
		indexWriter:   iWriter,
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

// NOTE: In this version our goal is performance.
// Checking fs for duplication is not suitable for us.
// Initial solution is overwrite values when duplication occures.
// Overwrites should be handled while flushing.
func (w writy) Set(key, value string) error {
	return w.ForceSet(key, value)
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

func (w writy) List() (keyval.StorageType, error) {
	// Append cache to fs
	return nil, nil
}
