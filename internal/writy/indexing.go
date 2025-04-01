package writy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
)

const (
	// storage data: ["key", "value"]
	STORAGE_KEY   = 0
	STORAGE_VALUE = 1

	// index data: ["key", offset, is_deleted]
	INDEX_KEY        = 0
	INDEX_OFFSET     = 1
	INDEX_IS_DELETED = 2
)

var lk sync.RWMutex

// TODO: Redesign this and find efficient way for writing may lines once.
// This implementation is not performant.
func writeIndex(w *Writy, k string, v any) {
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

func searchIndexByKey(w *Writy, k string) int64 {
	lk.RLock()
	defer lk.RUnlock()

	w.indexReader.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(bufio.NewReader(w.indexReader))

	for scanner.Scan() {
		var indLine []any
		if err := json.Unmarshal(scanner.Bytes(), &indLine); err != nil {
			w.logger.Debug("unable to decode index line", "error", err, "line", scanner.Text())
			continue
		}

		fkey := fmt.Sprint(indLine[INDEX_KEY])
		foff, _ := strconv.ParseInt(fmt.Sprint(indLine[INDEX_OFFSET]), 0, 64)
		isDel, _ := strconv.ParseBool(fmt.Sprint(indLine[INDEX_IS_DELETED]))

		if !isDel && k == fkey {
			w.logger.Debug("found offset", "fkey", fkey, "k", k, "isdel", isDel, "offset", foff)
			return foff
		}
	}

	return -1
}

func getValueByOffset(w *Writy, off int64) any {
	off, err := w.storageReader.Seek(off, io.SeekCurrent)
	if err != nil {
		w.logger.Debug("unable seek to off position", "error", err, "offset", off)
		return nil
	}

	scanner := bufio.NewScanner(bufio.NewReader(w.storageReader))
	scanner.Scan()

	var storLine []string
	if err := json.Unmarshal(scanner.Bytes(), &storLine); err != nil {
		w.logger.Debug("failed to unmarshal storage line",
			"offset", off, "error", err, "line", scanner.Text())
		return nil
	}

	w.logger.Debug("desirable line found", "line", storLine)
	return storLine[STORAGE_VALUE]
}
