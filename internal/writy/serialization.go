package writy

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

type serializaer struct {
	fileReader *os.File
}

func newSerilizer(f *os.File) *serializaer {
	return &serializaer{fileReader: f}
}

func (s serializaer) Read(off int64) (any, error) {
	lk.RLock()
	defer lk.RUnlock()

	s.fileReader.Seek(off, io.SeekStart)

	reader := bufio.NewReader(s.fileReader)
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
