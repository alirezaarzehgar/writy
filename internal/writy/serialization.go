package writy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
)

type storageDecoder struct {
	f *os.File
}

func newStorageDecoder(f *os.File) *storageDecoder {
	return &storageDecoder{f: f}
}

func (s storageDecoder) Decode(off int64) (any, error) {
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
	offset, _ := s.f.Seek(0, io.SeekCurrent)
	line := strconv.Quote(fmt.Sprint(value)) + "\n"
	s.f.WriteString(line)
	return offset
}

type indexEncoder struct {
	f *os.File
}

func newIndexEncoder(f *os.File) *indexEncoder {
	return &indexEncoder{f: f}
}

func (s indexEncoder) Encode(key string, offset int64) error {
	index := []any{key, offset, 0}
	return json.NewEncoder(s.f).Encode(index)
}

func (s indexEncoder) Delete(nextIndOffset int64) error {
	slog.Debug("indexEncoder.Delete", "next index offset", nextIndOffset)
	_, err := s.f.WriteAt([]byte("1"), nextIndOffset-3)
	return err
}

type index struct {
	Key             string
	ValueOffset     int64
	IsDeleted       bool
	NextIndexOffset int64
}

const (
	// index data: ["key", offset, is_deleted]
	IndexKey       = 0
	IndexOffset    = 1
	IndexIsDeleted = 2
)

type indexDecoder struct {
	f             *os.File
	scnr          *bufio.Scanner
	currentOffset int64
}

func newIndexDecoder(f *os.File) *indexDecoder {
	f.Seek(0, io.SeekStart)
	scnr := bufio.NewScanner(bufio.NewReader(f))
	return &indexDecoder{f: f, scnr: scnr}
}

func (s indexDecoder) Decode() index {
	var indexLine []any
	err := json.Unmarshal(s.scnr.Bytes(), &indexLine)
	if err != nil {
		slog.Warn("failed to unmarshal line", "error", err)
		return index{}
	}

	key := fmt.Sprint(indexLine[IndexKey])
	storageOffset, err := strconv.ParseInt(fmt.Sprint(indexLine[IndexOffset]), 0, 64)
	if err != nil {
		slog.Warn("failed to parse int", "error", err)
		return index{}
	}
	isDeleted, err := strconv.ParseBool(fmt.Sprint(indexLine[IndexIsDeleted]))
	if err != nil {
		slog.Warn("failed to parse bool", "error", err)
		return index{}
	}

	return index{Key: key, ValueOffset: storageOffset, IsDeleted: isDeleted, NextIndexOffset: s.currentOffset}
}

func (s *indexDecoder) Scan() bool {
	result := s.scnr.Scan()
	s.currentOffset += int64(len(s.scnr.Text()) + 1)
	return result
}
