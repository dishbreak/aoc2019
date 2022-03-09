package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPossiblePasscode(t *testing.T) {
	type testCase struct {
		passcode int
		isValid  bool
	}

	testCases := []testCase{
		{
			passcode: 111111,
			isValid:  true,
		},
		{
			passcode: 223450,
			isValid:  false,
		},
		{
			passcode: 123789,
			isValid:  false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			assert.Equal(t, tc.isValid, isPossiblePasscode(tc.passcode))
		})
	}
}

func TestPossiblePasscodeV2(t *testing.T) {
	type testCase struct {
		passcode int
		isValid  bool
	}

	testCases := []testCase{
		{
			passcode: 112233,
			isValid:  true,
		},
		{
			passcode: 123444,
			isValid:  false,
		},
		{
			passcode: 111122,
			isValid:  true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			assert.Equal(t, tc.isValid, isPossiblePasscodeV2(tc.passcode))
		})
	}
}
