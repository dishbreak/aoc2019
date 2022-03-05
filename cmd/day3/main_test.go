package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPart1(t *testing.T) {
	type testCase struct {
		input    []string
		expected int
	}

	testCases := []testCase{
		{
			input: []string{
				"R8,U5,L5,D3",
				"U7,R6,D4,L4",
				"",
			},
			expected: 6,
		},
		{
			input: []string{
				"R75,D30,R83,U83,L12,D49,R71,U7,L72",
				"U62,R66,U55,R34,D71,R55,D58,R83",
			},
			expected: 159,
		},
		{
			input: []string{
				"R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51",
				"U98,R91,D20,R16,D67,R40,U7,R15,U6,R7",
			},
			expected: 135,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			actual := part1(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPart2(t *testing.T) {
	type testCase struct {
		input    []string
		expected int
	}

	testCases := []testCase{
		{
			input: []string{
				"R8,U5,L5,D3",
				"U7,R6,D4,L4",
				"",
			},
			expected: 30,
		},
		{
			input: []string{
				"R75,D30,R83,U83,L12,D49,R71,U7,L72",
				"U62,R66,U55,R34,D71,R55,D58,R83",
			},
			expected: 610,
		},
		{
			input: []string{
				"R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51",
				"U98,R91,D20,R16,D67,R40,U7,R15,U6,R7",
			},
			expected: 410,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			actual := part2(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
