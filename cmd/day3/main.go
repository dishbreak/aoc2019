package main

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"github.com/dishbreak/aoc-common/lib"
)

func main() {
	input, err := lib.GetInput("inputs/day3.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}

func abs(o int) int {
	if o > 0 {
		return o
	}
	return -1 * o
}

func dist(w image.Point) int {
	return abs(w.X) + abs(w.Y)
}

var directions = map[byte]image.Point{
	'R': {1, 0},
	'L': {-1, 0},
	'U': {0, 1},
	'D': {0, -1},
}

func part1(input []string) int {
	match := 100000000000
	board := make(map[image.Point]bool)

	cursor := image.Point{}

	for _, inst := range strings.Split(input[0], ",") {
		vector := directions[inst[0]]
		mag, _ := strconv.Atoi(inst[1:])
		for i := 0; i < mag; i++ {
			board[cursor] = true
			cursor = cursor.Add(vector)
		}
	}

	cursor.X, cursor.Y = 0, 0

	for _, inst := range strings.Split(input[1], ",") {
		vector := directions[inst[0]]
		mag, _ := strconv.Atoi(inst[1:])
		for i := 0; i < mag; i++ {
			if d := dist(cursor); board[cursor] && d != 0 {
				if d < match {
					match = d
				}
			}
			cursor = cursor.Add(vector)
		}
	}

	return match
}

func part2(input []string) int {
	match := 100000000000
	board := make(map[image.Point]int)

	cursor := image.Point{}

	step := 0
	for _, inst := range strings.Split(input[0], ",") {
		vector := directions[inst[0]]
		mag, _ := strconv.Atoi(inst[1:])
		for i := 0; i < mag; i++ {
			board[cursor] = step
			cursor = cursor.Add(vector)
			step++
		}
	}

	cursor.X, cursor.Y = 0, 0
	step = 0
	for _, inst := range strings.Split(input[1], ",") {
		vector := directions[inst[0]]
		mag, _ := strconv.Atoi(inst[1:])
		for i := 0; i < mag; i++ {
			if d := dist(cursor); board[cursor] != 0 && d != 0 {
				t := board[cursor] + step
				if t < match {
					match = t
				}
			}
			cursor = cursor.Add(vector)
			step++
		}
	}

	return match
}
