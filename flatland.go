package main

import (
	"bufio"
	"bytes"
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
	Object interface{}
	Fields []field
	column int
	line   int
	field  bytes.Buffer
	r      *bufio.Reader
}

// NewReader returns a new Reader that reads from r
func NewReader(r io.Reader, obj interface{}) *Reader {
	return &Reader{
		Object: obj,
		r:      bufio.NewReader(r),
	}
}

// Read reads on record from r
func (r *Reader) Read() (record interface{}, err error) {
	for {
		record, err = r.parseRecord()
		if record != nil {
			break
		}
	}

	return
}

// ReadAll reads all the remaining records from r
func (r *Reader) ReadAll() (records []interface{}, err error) {
	for {
		record, err := r.Read()

		if err == io.EOF {
			return records, nil
		}

		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}
}

func (r *Reader) readRune() (rune, error) {
	r1, _, err := r.r.ReadRune()
	// Handle \r\n here.  We make the simplifying assumption that
	// anytime \r is followed by \n that it can be folded to \n.
	// We will not detect files which contain both \r\n and bare \n.
	if r1 == '\r' {
		r1, _, err = r.r.ReadRune()
		if err == nil {
			if r1 != '\n' {
				r.r.UnreadRune()
				r1 = '\r'
			}
		}
	}
	r.column++
	return r1, err
}

func (r *Reader) unreadRune() {
	r.r.UnreadRune()
	r.column--
}

func (r *Reader) parseRecord() (record interface{}, err error) {
	r.line++
	r.column = -1

	if r.Fields == nil {
		r.Fields, err = r.parseFieldTags()
		if err != nil {
			return nil, err
		}
	}

	_, _, err = r.r.ReadRune()

	if err != nil {
		return nil, err
	}

	r.r.UnreadRune()

	for {
		haveField, delim, err := r.parseField()
		if haveField {
			fields = append(fields, r.field.String())
		}
		if delim == '\n' || err == io.EOF {
			return fields, err
		} else if err != nil {
			return nil, err
		}
	}
}

func (r *Reader) parseField() (haveField bool, delim rune, err error) {
	r.field.Reset()

	r1, err := r.readRune()
	if err == io.EOF && r.column != 0 {
		return true, 0, err
	}
	if err != nil {
		return false, 0, err
	}

	return true, r1, nil
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
		coords := strings.Split(fieldVal, "..")
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
