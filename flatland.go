package main

import (
	"bufio"
	"errors"
	"io"
)

var (
	ErrInvalidRecordLength = errors.New("record is invalid length for template")
)

// A Reader reads records from a fixed width encoded file using a given struct
// as a template
type Reader struct {
	*bufio.Scanner
	Template [][]int
	empty    bool
	eor      bool
}

// NewReader returns a new Reader that reads from r
func NewReader(r io.Reader, template [][]int) *Reader {
	return &Reader{bufio.NewScanner(r), template, false, false}
}

func (r *Reader) ScanAll() ([][]string, error) {
	var objs [][]string

	for r.Scan() {
		obj, err := r.parseRecord()

		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}

	return objs, nil
}

func (r *Reader) ScanLine() ([]string, error) {
	r.Scan()
	obj, err := r.parseRecord()

	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *Reader) EmptyLine() bool {
	return r.empty && r.eor
}

func (r *Reader) EndOfRecord() bool {
	return r.eor
}

func (r *Reader) parseRecord() ([]string, error) {
	line := r.Text()
	record := []string{}

	for _, coords := range r.Template {
		if len(line) < coords[0] || len(line) < coords[1] {
			return nil, ErrInvalidRecordLength
		}
		value := line[coords[0]-1 : coords[1]]
		record = append(record, value)
	}

	return record, nil
}
