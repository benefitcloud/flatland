package flatland

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var template = [][]int{
	{1, 8},
	{9, 9},
}

var readTests = []struct {
	Name   string
	Input  string
	Output [][]string
}{
	{
		Name:  "Simple",
		Input: "111111112",
		Output: [][]string{
			{"11111111", "2"},
		},
	},
	{
		Name:  "Invalid",
		Input: "11111111",
		Output: [][]string{
			{"11111111", "2"},
		},
	},
}

func TestScanLine(t *testing.T) {
	tt := readTests[0]
	r := NewReader(strings.NewReader(tt.Input), template)
	obj, _ := r.ScanLine()

	assert.Equal(t, obj, tt.Output[0])
}

func TestScanAll(t *testing.T) {
	for _, tt := range readTests {
		r := NewReader(strings.NewReader(tt.Input), template)
		objs, err := r.ScanAll()

		if err != nil {
			assert.Equal(t, err, ErrInvalidRecordLength)
		} else {
			assert.Equal(t, objs, tt.Output)
		}
	}
}
