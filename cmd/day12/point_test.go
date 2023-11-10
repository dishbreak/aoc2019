package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePoint(t *testing.T) {
	type testCase struct {
		input  string
		result [3]int
	}

	testCases := []testCase{
		{"<x=-1, y=0, z=2>", [3]int{-1, 0, 2}},
		{"<x=2, y=-10, z=-7>", [3]int{2, -10, -7}},
		{"<x=4, y=-8, z=8>", [3]int{4, -8, 8}},
		{"<x=3, y=5, z=-1>", [3]int{3, 5, -1}},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			assert.Equal(t, tc.result, FromString(tc.input))
		})
	}
}
