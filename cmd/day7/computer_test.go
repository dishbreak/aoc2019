package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputer(t *testing.T) {
	type testCase struct {
		program        []int
		input          []int
		expectedOutput int
	}

	testCases := []testCase{
		{
			program:        []int{3, 0, 4, 0, 99},
			input:          []int{1337},
			expectedOutput: 1337,
		},
		{
			program:        []int{1101, 100, -1, 0, 4, 0, 99},
			expectedOutput: 99,
		},
		{
			program:        []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			expectedOutput: 1,
			input:          []int{8},
		},
		{
			program:        []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			input:          []int{7},
			expectedOutput: 0,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			c := &IntcodeComputer{}
			c.Load(tc.program)
			o, err := c.Execute(context.TODO(), tc.input)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedOutput, o)
		})
	}

}
