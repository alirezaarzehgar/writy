package writy

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alirezaarzehgar/writy/cache"
)

var glk sync.RWMutex

type StorageType map[string]any

var (
	DefaultStoragePath string        = fmt.Sprintf("%s", os.Getenv("HOME"))
	DefaultFlushCycle  time.Duration = time.Second * 5
)

type Writy struct {
	storageReader *os.File
	storageWriter *os.File
	indexReader   *os.File
	indexWriter   *os.File
	flusher       *Flusher
	cache         *cache.Cache
	gc            *GarbageCollector
	w8ForDaemons  *sync.WaitGroup
}

func New(path string) (*Writy, error) {
	if path == "" {
		path = DefaultStoragePath
	}

	storagePath := filepath.Join(path, "storage.db")
	sWriter, err := os.OpenFile(storagePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable open file %s for writing: %w", path, err)
	}
	sWriter.Seek(0, io.SeekEnd)

	sReader, err := os.Open(storagePath)
	if err != nil {
		return nil, fmt.Errorf("unable open file %s for reading: %w", path, err)
	}

	indexPath := filepath.Join(path, "index.db")
	iWriter, err := os.OpenFile(indexPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable open file %s for writing: %w", path, err)
	}
	iWriter.Seek(0, io.SeekEnd)

	iReader, err := os.Open(indexPath)
	if err != nil {
		return nil, fmt.Errorf("unable open file %s for reading: %w", path, err)
	}

	w := &Writy{
		storageReader: sReader,
		storageWriter: sWriter,
		indexReader:   iReader,
		indexWriter:   iWriter,
		flusher:       newFlusher(DefaultFlushCycle),
		gc:            newGarbageCollector(),
		cache:         cache.New(),
		w8ForDaemons:  &sync.WaitGroup{},
	}

	w.flusher.Run(w)
	w.gc.Run(w, w.flusher)

	return w, nil
}

// NOTE: In this version our goal is performance.
// Checking fs for duplication is not suitable for us.
// Initial solution is ignoring duplicated records when flushing.
func (w Writy) Set(key string, value any) error {
	glk.Lock()
	defer glk.Unlock()

	return w.cache.ForceSet(key, value)
}

func (w Writy) Get(key string) any {
	glk.RLock()
	defer glk.RUnlock()

	value := w.cache.Get(key)
	if value != nil {
		return value
	}

	off := searchIndexByKey(&w, key)
	if off < 0 {
		return nil
	}

	value = getValueByOffset(&w, off)
	w.cache.ForceSet(key, value)
	return value
}

func (w Writy) Del(key string) error {
	glk.Lock()
	defer glk.Unlock()

	w.cache.Del(key)

	indDec := newIndexDecoder(w.indexReader)
	indEnc := newIndexEncoder(w.indexWriter)
	for indDec.Scan() {
		ind := indDec.Decode()
		if !ind.IsDeleted && key == ind.Key {
			slog.Debug("index found", "index", ind)
			return indEnc.Delete(ind.NextIndexOffset)
		}
	}

	return nil
}

func (w Writy) Keys() (keys []string) {
	glk.RLock()
	defer glk.RUnlock()

	indDec := newIndexDecoder(w.indexReader)
	for indDec.Scan() {
		ind := indDec.Decode()
		if !ind.IsDeleted {
			keys = append(keys, ind.Key)
		}
	}
	return
}

func (w Writy) Cleanup() {
	w.flusher.flush()
	w.gc.collect()
	w.w8ForDaemons.Wait()
}
