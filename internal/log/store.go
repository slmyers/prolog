package log

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenwidth = 8
)

type store struct {
	*os.File
	mu sync.Mutex
	// TODO: study io Writer
	buf *bufio.Writer
	size uint64
}

func newStore(f *os.File) (*store, error) {
	file, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(file.Size())
	return &store{
		File: f,
		size: size,
		buf: bufio.NewWriter(f),
	}, nil
}

// Append appends the big-endian binary representation of p to the store buffer
// avoid writing to the store file frequently. This will help with the case of
// many small writes as we make less syscalls. NOTE: I'm having trouble understanding how
// this is supposed to reduce the number of writes to file as it writes the buffer every time
// TODO: read about buffered IO
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size
	// writes the length in binary into the buffer
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0,0, err
	}
	// writes the data
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	w += lenwidth
	s.size += uint64(w)
	return uint64(w), pos, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	size := make([]byte, lenwidth)
	// read our record size
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	b := make([]byte, enc.Uint64(size))
	// read record using known record size
	if _, err := s.File.ReadAt(b, int64(pos+lenwidth)); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.buf.Flush()
	if err != nil {
		return err
	}
	if err = s.File.Close(); err != nil {
		return fmt.Errorf("store: error closing file (%w)", err)
	}
	return nil
}