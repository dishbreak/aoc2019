package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPart1(t *testing.T) {
	type testCase struct {
		program []int
		output  int
	}

	testCases := []testCase{
		{
			program: []int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0},
			output:  43210,
		},
		{
			program: []int{3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23, 101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0},
			output:  54321,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			o := part1(tc.program)
			assert.Equal(t, tc.output, o)
		})
	}
}

func TestExecuteConfig(t *testing.T) {

	type testCase struct {
		program  []int
		settings []int
		result   int
	}

	testCases := []testCase{
		{
			program:  []int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0},
			settings: []int{4, 3, 2, 1, 0},
			result:   43210,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			o, err := executeConfig(context.TODO(), tc.program, tc.settings)
			assert.Nil(t, err)
			assert.Equal(t, tc.result, o)
		})
	}
}

func TestSimulation(t *testing.T) {
	type testCase struct {
		program  []int
		settings []int
		result   int
	}

	testCases := []testCase{
		{
			program: []int{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26,
				27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5},
			settings: []int{9, 8, 7, 6, 5},
			result:   139629729,
		},
		{
			program: []int{3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54,
				-5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4,
				53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10},
			settings: []int{9, 7, 8, 5, 6},
			result:   18216,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			r, err := runSimulation(context.TODO(), tc.settings, tc.program)
			assert.Nil(t, err)
			assert.Equal(t, tc.result, r)
		})
	}
}
