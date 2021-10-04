package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteV1(t *testing.T) {
	type testCase struct {
		input  []int
		result []int
	}

	testCases := []testCase{
		{
			input:  []int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50},
			result: []int{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50},
		},
		{
			input:  []int{1, 0, 0, 0, 99},
			result: []int{2, 0, 0, 0, 99},
		},
		{
			input:  []int{2, 3, 0, 3, 99},
			result: []int{2, 3, 0, 6, 99},
		},
		{
			input:  []int{2, 4, 4, 5, 99, 0},
			result: []int{2, 4, 4, 5, 99, 9801},
		},
		{
			input:  []int{1, 1, 1, 4, 99, 5, 6, 0, 99},
			result: []int{30, 1, 1, 4, 2, 5, 6, 0, 99},
		},
	}

	for idx, tc := range testCases {
		name := fmt.Sprintf("test case %d", idx)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.result, executeV1(tc.input))
		})
	}
}
