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
		{
			program: []int{
				3,  // 00 input
				11, // 01 ->inputA var
				3,  // 02 input
				12, // 03 ->inputB var
				8,  // 04 cmp
				11, // 05 ->inputA var
				12, // 06 ->inputB var
				13, // 07 ->output var
				4,  // 08 output
				13, // 09 ->output var
				99, // 10 halt
				0,  // 11 inputA var
				0,  // 12 inputB var
				0,  // 13 output var
			},
			input:          []int{34, 14},
			expectedOutput: 0,
		},
		{
			program: []int{
				3,  // 00 input
				11, // 01 ->inputA var
				3,  // 02 input
				12, // 03 ->inputB var
				8,  // 04 cmp
				11, // 05 ->inputA var
				12, // 06 ->inputB var
				13, // 07 ->output var
				4,  // 08 output
				13, // 09 ->output var
				99, // 10 halt
				0,  // 11 inputA var
				0,  // 12 inputB var
				0,  // 13 output var
			},
			input:          []int{34, 34},
			expectedOutput: 1,
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
