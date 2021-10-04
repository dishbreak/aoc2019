package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	input  int
	result int
}

func TestFuelCalcV1(t *testing.T) {

	testCases := []testCase{
		{12, 2},
		{14, 2},
		{1969, 654},
		{100756, 33583},
		{3, 0},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			assert.Equal(t, tc.result, fuelCalcV1(tc.input))
		})
	}
}

func TestFuelCalcV2(t *testing.T) {
	testCases := []testCase{
		{14, 2},
		{1969, 966},
		{100756, 50346},
	}

	for idx, tc := range testCases {
		testName := fmt.Sprintf("test case %d", idx)
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, tc.result, fuelCalcV2(tc.input))
		})
	}
}
