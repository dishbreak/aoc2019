package util

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	type testCase struct {
		input  string
		result []int64
		failed bool
	}

	testCases := []testCase{
		{
			input: "13,355,34,344",
			result: []int64{
				13, 355, 34, 344,
			},
		},
		{
			input:  "13,fish,34,344",
			result: nil,
			failed: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			r, err := ParseIntcode(tc.input)
			assert.Equal(t, tc.result, r)
			if tc.failed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestComputer(t *testing.T) {
	type testCase struct {
		program        []int64
		input          []int64
		expectedOutput int64
	}

	testCases := []testCase{
		{
			program:        []int64{3, 0, 4, 0, 99},
			input:          []int64{1337},
			expectedOutput: 1337,
		},
		{
			program:        []int64{1101, 100, -1, 0, 4, 0, 99},
			expectedOutput: 99,
		},
		{
			program:        []int64{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			expectedOutput: 1,
			input:          []int64{8},
		},
		{
			program:        []int64{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			input:          []int64{7},
			expectedOutput: 0,
		},
		{
			program: []int64{
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
			input:          []int64{34, 14},
			expectedOutput: 0,
		},
		{
			program: []int64{
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
			input:          []int64{34, 34},
			expectedOutput: 1,
		},
		{
			program:        []int64{104, 1125899906842624, 99},
			input:          []int64{},
			expectedOutput: 1125899906842624,
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

func TestSimulate(t *testing.T) {
	program := []int64{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99}
	inputChan := make(chan int64, len(program))
	defer close(inputChan)

	outputChan := make(chan int64)
	outputs := make([]int64, 0)
	doneChan := make(chan interface{})

	go func() {
		for val := range outputChan {
			outputs = append(outputs, val)
		}
		close(doneChan)
	}()

	errChan := make(chan error)
	i := &IntcodeComputer{}
	i.Load(program)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer close(errChan)
	go func() {
		i.Simulate(ctx, "TestSimulate", inputChan, outputChan, errChan)
	}()
	for {
		select {
		case <-doneChan:
			assert.Equal(t, program, outputs)
			return
		case e := <-errChan:
			t.Fatalf("failed simulate: %s", e)
		}
	}

}

func TestGetValue(t *testing.T) {
	type testCase struct {
		program  []int64
		relBase  int64
		addr     int64
		mode     ParameterMode
		expected int64
	}

	exampleProgram := []int64{
		10, // 0
		3,  // 1
		45, // 2
		17, // 3
		12, // 4
		2,  // 5
		67, // 6
	}
	testCases := []testCase{
		{
			program:  exampleProgram,
			addr:     3,
			mode:     ImmediateMode,
			expected: 17,
		},
		{
			program:  exampleProgram,
			addr:     5,
			mode:     PositionalMode,
			expected: 45,
		},
		{
			program:  exampleProgram,
			addr:     5,
			mode:     RelativeMode,
			expected: 45,
		},
		{
			program:  exampleProgram,
			addr:     5,
			relBase:  2,
			mode:     RelativeMode,
			expected: 12,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			c := &IntcodeComputer{}
			c.Load(tc.program)
			c.rel = tc.relBase

			assert.Equal(t, tc.expected, c.getValue(tc.mode, tc.addr))
		})
	}
}
