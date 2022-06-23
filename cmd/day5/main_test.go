package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteV3(t *testing.T) {
	type testCase struct {
		program        []int
		input          int
		expectedOutput []int
		expectedError  bool
	}

	testCases := []testCase{
		{
			program:        []int{3, 0, 4, 0, 99},
			input:          1337,
			expectedOutput: []int{1337},
		},
		{
			program:        []int{1101, 100, -1, 0, 4, 0, 99},
			expectedOutput: []int{99},
		},
		{
			program:        []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			expectedOutput: []int{1},
			input:          8,
		},
		{
			program:        []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			input:          7,
			expectedOutput: []int{0},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			result, err := executeV3(tc.input, tc.program)
			assert.Equal(t, tc.expectedOutput, result)
			if tc.expectedError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
