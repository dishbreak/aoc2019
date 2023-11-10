package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPart1(t *testing.T) {
	type testCase struct {
		input          string
		steps          int
		expectedEnergy int
	}

	testCases := []testCase{
		{
			input: `<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>`,
			steps:          10,
			expectedEnergy: 179,
		},
		{
			input: `<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>,`,
			steps:          100,
			expectedEnergy: 1940,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			r := strings.NewReader(tc.input)
			assert.Equal(t, tc.expectedEnergy, simulate(r, tc.steps))
		})
	}
}

func TestPart2(t *testing.T) {
	type testCase struct {
		input string
		steps int
	}

	testCases := []testCase{
		{
			input: `<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>`,
			steps: 2772,
		},
		{
			input: `<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>,`,
			steps: 4686774924,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			r := strings.NewReader(tc.input)
			assert.Equal(t, tc.steps, part2(r))
		})
	}
}
