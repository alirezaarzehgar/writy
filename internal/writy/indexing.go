package writy

import (
	"encoding/json"
	"io"
	"sync"
)

const (
	// storage data: ["key", "value"]
	STORAGE_KEY   = 0
	STORAGE_VALUE = 1

	// index data: ["key", offset, is_deleted]
	INDEX_KEY        = 0
	INDEX_VALUE      = 1
	INDEX_IS_DELETED = 2
)

var lk sync.Mutex

func writeIndex(w *writy, k string, v any) {
	storage := w.storageWriter
	ind := w.indexWriter

	byteValue, err := json.Marshal(v)
	if err != nil {
		w.logger.Debug("can not marshal value", "error", err)
		return
	}
	w.logger.Debug("write to fs", "key", k, "value", v)

	lk.Lock()
	defer lk.Unlock()

	off, err := storage.Seek(0, io.SeekCurrent)
	if err != nil {
		w.logger.Warn("unable to get current offset", "error", err)
	}

	l := []any{k, string(byteValue)}
	if err := json.NewEncoder(storage).Encode(l); err != nil {
		w.logger.Warn("unable to encode data", "error", err, "line", l)
	}
	i := []any{k, off, 0}
	if err := json.NewEncoder(ind).Encode(i); err != nil {
		w.logger.Warn("unable to encode data", "error", err, "line", i)
	}
}

func searchIndex() {

}
