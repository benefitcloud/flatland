package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var readTests = []struct {
	Name   string
	Input  string
	Output []TestRow
}{
	{
		Name:   "Simple",
		Input:  "112",
		Output: []TestRow{{"11", "2"}},
	},
	{
		Name:   "Invalid",
		Input:  "11",
		Output: []TestRow{{"11", "2"}},
	},
}

type TestRow struct {
	One string `flat:"1..2"`
	Two string `flat:"3..3"`
}

func TestColumnsParsedSize(t *testing.T) {
	r := Reader{
		Object: TestRow{},
	}
	fields, err := r.parseFieldTags()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(fields), 2)
}

func TestColumnName(t *testing.T) {
	r := Reader{
		Object: TestRow{},
	}
	fields, err := r.parseFieldTags()
	field := fields[0]
	assert.Equal(t, err, nil)
	assert.Equal(t, field.Name, "One")
}

func TestSingleColumnLengthValueColumn(t *testing.T) {
	r := Reader{
		Object: TestRow{},
	}
	fields, err := r.parseFieldTags()
	field := fields[1]
	assert.Equal(t, err, nil)
	assert.Equal(t, field.From, 3)
	assert.Equal(t, field.To, 3)
}

func TestMultiColumnLengthValueColumn(t *testing.T) {
	r := Reader{
		Object: TestRow{},
	}
	fields, err := r.parseFieldTags()
	field := fields[0]
	assert.Equal(t, err, nil)
	assert.Equal(t, field.From, 1)
	assert.Equal(t, field.To, 2)
}

func TestRead(t *testing.T) {
	for _, tt := range readTests {
		r := NewReader(strings.NewReader(tt.Input), TestRow{})
		out, err := r.ReadAll()

		if err != nil {
			assert.Fail(t, tt.Name, err)
		}

		assert.Equal(t, out, tt.Output)
	}
}
