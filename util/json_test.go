package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringify(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "struct",
			input:    struct{ Name string }{"test"},
			expected: `{"Name":"test"}`,
		},
		{
			name:     "map",
			input:    map[string]int{"a": 1},
			expected: `{"a":1}`,
		},
		{
			name:     "slice",
			input:    []int{1, 2, 3},
			expected: `[1,2,3]`,
		},
		{
			name:     "string",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "nil",
			input:    nil,
			expected: "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Stringify(tt.input))
		})
	}
}
