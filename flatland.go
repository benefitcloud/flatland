package main

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrEmptyFieldTagCount = errors.New("this struct has no 'flat' field tags")
)

type field struct {
	Name string
	From int
	To   int
}

// A Reader reads records from a fixed width encoded file using a given struct
// as a template
type Reader struct {
	*bufio.Scanner
	Object interface{}
	Fields []field
}

// NewReader returns a new Reader that reads from r
func NewReader(r io.Reader, obj interface{}) *Reader {
	return &Reader{bufio.NewScanner(r), obj, []field{}}
}

func (r *Reader) parseFieldTags() (fs []field, err error) {
	typ := reflect.TypeOf(r.Object)
	for i := 0; i < typ.NumField(); i++ {
		if f := r.parseFieldTag(typ.Field(i)); f.Name != "" {
			fs = append(fs, f)
		}
	}

	if len(fs) < 1 {
		err = ErrEmptyFieldTagCount
	}

	return
}

func (r *Reader) parseFieldTag(sf reflect.StructField) (f field) {
	if fieldVal := sf.Tag.Get("flat"); fieldVal != "" {
		coords := strings.Split(fieldVal, ",")
		f.Name = sf.Name
		f.From, _ = strconv.Atoi(coords[0])

		if len(coords) < 2 {
			f.To, _ = strconv.Atoi(coords[0])
		} else {
			f.To, _ = strconv.Atoi(coords[1])
		}
	}

	return
}
