package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dishbreak/aoc-common/lib"
)

func main() {
	input, err := lib.GetInput("inputs/day2.txt")
	if err != nil {
		panic(err)
	}

	parts := strings.Split(input[0], ",")
	program := make([]int, len(parts))
	for idx, val := range parts {
		parsed, _ := strconv.Atoi(val)
		program[idx] = parsed
	}

	fmt.Printf("Part 1: %d\n", part1(program))
	fmt.Printf("Part 2: %d\n", part2(program))
}

func part1(input []int) int {

	program := make([]int, len(input))
	copy(program, input)

	// to fix the 1202 error, per the puzzle instructions
	program[1] = 12
	program[2] = 2

	program = executeV1(program)
	return program[0]
}

func part2(input []int) int {
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input); j++ {
			prog := make([]int, len(input))
			copy(prog, input)
			if executeV2(prog, i, j, 19690720) {
				return 100*i + j
			}
		}
	}
	return -1
}

func executeV1(program []int) []int {
	for i := 0; i < len(program); i = i + 4 {
		switch program[i] {
		case 1:
			program[program[i+3]] = program[program[i+1]] + program[program[i+2]]
		case 2:
			program[program[i+3]] = program[program[i+1]] * program[program[i+2]]
		case 99:
			return program
		default:
			return program
		}
	}
	return program
}

func executeV2(program []int, noun, verb, target int) bool {
	program[1] = noun
	program[2] = verb
	result := executeV1(program)
	return result[0] == target
}
