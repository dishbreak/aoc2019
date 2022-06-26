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
