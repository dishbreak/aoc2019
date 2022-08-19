package main

import (
	"fmt"
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscover(t *testing.T) {
	type testCase struct {
		origin   image.Point
		field    []string
		expected int
	}

	testCases := []testCase{
		{
			origin: image.Point{},
			field: []string{
				"#.........",
				"...#......",
				"...#..#...",
				".####....#",
				"..#.#.#...",
				".....#....",
				"..###.#.##",
				".......#..",
				"....#...#.",
				"...#..#..#",
			},
			expected: 7,
		},
		{
			origin: image.Pt(1, 2),
			field: []string{
				"#.#...#.#.",
				".###....#.",
				".#....#...",
				"##.#.#.#.#",
				"....#.#.#.",
				".##..###.#",
				"..#...##..",
				"..##....##",
				"......#...",
				".####.###.",
			},
			expected: 35,
		},
		{
			origin: image.Pt(5, 8),
			field: []string{
				"......#.#.",
				"#..#.#....",
				"..#######.",
				".#.#.###..",
				".#..#.....",
				"..#....#.#",
				"#..#....#.",
				".##.#..###",
				"##...#..#.",
				".#....####",
			},
			expected: 33,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			space := parseSpace(tc.field)
			assert.Equal(t, tc.expected, discover(space, tc.origin))
		})
	}
}
