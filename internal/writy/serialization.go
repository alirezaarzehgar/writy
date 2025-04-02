package writy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"sync"
)

var lk sync.RWMutex

type storageDecoder struct {
	f *os.File
}

func newStorageDecoder(f *os.File) *storageDecoder {
	return &storageDecoder{f: f}
}

func (s storageDecoder) Decode(off int64) (any, error) {
	lk.RLock()
	defer lk.RUnlock()

	s.f.Seek(off, io.SeekStart)

	reader := bufio.NewReader(s.f)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// Remove newline for Unquote
	line = line[:len(line)-1]
	line, err = strconv.Unquote(line)
	if err != nil {
		return nil, err
	}

	return line, err
}

type storageEncoder struct {
	f *os.File
}

func newStorageEncoder(f *os.File) *storageEncoder {
	return &storageEncoder{f: f}
}

func (s storageEncoder) Encode(key string, value any) int64 {
	lk.Lock()
	defer lk.Unlock()

	line := strconv.Quote(fmt.Sprint(value)) + "\n"
	s.f.WriteString(line)
	offset, _ := s.f.Seek(0, io.SeekCurrent)
	return offset
}

type indexEncoder struct {
	f *os.File
}

func newIndexEncoder(f *os.File) *indexEncoder {
	return &indexEncoder{f: f}
}

func (s indexEncoder) Encode(key string, offset int64) error {
	lk.Lock()
	defer lk.Unlock()

	index := []any{key, offset, 0}
	return json.NewEncoder(s.f).Encode(index)
}

type index struct {
	Key       string
	Offset    int64
	IsDeleted bool
}

const (
	// index data: ["key", offset, is_deleted]
	IndexKey       = 0
	IndexOffset    = 1
	IndexIsDeleted = 2
)

type indexDecoder struct {
	f      *os.File
	scnr   *bufio.Scanner
	logger *slog.Logger
}

func newIndexDecoder(f *os.File, l *slog.Logger) *indexDecoder {
	f.Seek(0, io.SeekStart)
	scnr := bufio.NewScanner(bufio.NewReader(f))
	return &indexDecoder{f: f, scnr: scnr, logger: l}
}

func (s indexDecoder) Decode() index {
	lk.RLock()
	defer lk.RUnlock()

	var indexLine []any
	err := json.Unmarshal(s.scnr.Bytes(), &indexLine)
	if err != nil {
		s.logger.Warn("failed to unmarshal line", "error", err)
	}

	fkey := fmt.Sprint(indexLine[IndexKey])
	foff, err := strconv.ParseInt(fmt.Sprint(indexLine[IndexOffset]), 0, 64)
	if err != nil {
		s.logger.Warn("failed to parse int", "error", err)
	}
	isDel, err := strconv.ParseBool(fmt.Sprint(indexLine[IndexIsDeleted]))
	if err != nil {
		s.logger.Warn("failed to parse bool", "error", err)
	}

	return index{Key: fkey, Offset: foff, IsDeleted: isDel}
}

func (s indexDecoder) Scan() bool {
	return s.scnr.Scan()
}
