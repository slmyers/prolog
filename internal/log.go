package server

import (
	"fmt"
	"sync"
)

type Log struct {
	mu sync.Mutex
	records []Record
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Append(record Record) (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	record.offset = uint64(len(l.records))
	l.records = append(l.records, record)
	return record.offset, nil
}

func (l *Log) Read(offset uint64) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if offset >= uint64(len(l.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return l.records[offset], nil
}

//Record why do we need offset?
type Record struct {
	Value  []byte `json:"value,omitempty"`
	offset uint64 `json:"offset,omitempty"`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")