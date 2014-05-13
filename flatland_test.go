package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var readTests = []struct {
	Name   string
	Input  string
	Output []TestRecord
}{
	{
		Name:   "Simple",
		Input:  "111111112",
		Output: []TestRecord{{"11111111", "2"}},
	},
	{
		Name:   "Invalid",
		Input:  "11111111",
		Output: []TestRecord{{"11111111", "2"}},
	},
}

type TestRecord struct {
	One string `flat:"1,8"`
	Two string `flat:"9,9"`
}

func TestColumnsParsedSize(t *testing.T) {
	r := Reader{
		Object: TestRecord{},
	}
	fields, err := r.parseFieldTags()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(fields), 2)
}

func TestColumnName(t *testing.T) {
	r := Reader{
		Object: TestRecord{},
	}
	fields, err := r.parseFieldTags()
	field := fields[0]
	assert.Equal(t, err, nil)
	assert.Equal(t, field.Name, "One")
}

func TestSingleColumnLengthValueColumn(t *testing.T) {
	r := Reader{
		Object: TestRecord{},
	}
	fields, err := r.parseFieldTags()
	field := fields[1]
	assert.Equal(t, err, nil)
	assert.Equal(t, field.From, 9)
	assert.Equal(t, field.To, 9)
}

func TestMultiColumnLengthValueColumn(t *testing.T) {
	r := Reader{
		Object: TestRecord{},
	}
	fields, err := r.parseFieldTags()
	field := fields[0]
	assert.Equal(t, err, nil)
	assert.Equal(t, field.From, 1)
	assert.Equal(t, field.To, 8)
}
