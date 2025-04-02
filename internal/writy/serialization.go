package writy

import (
	"bufio"
	"fmt"
	"io"
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
