package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dishbreak/aoc-common/lib"
)

func main() {
	input, err := lib.GetInput("inputs/day5.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}

func toInstructions(input []string) []int {
	parts := strings.Split(input[0], ",")
	parsed := make([]int, len(parts))

	for i, val := range parts {
		parsed[i], _ = strconv.Atoi(val)
	}
	return parsed
}

func part1(input []string) int {
	program := toInstructions(input)

	result, err := executeV3(1, program)
	if err != nil {
		panic(err)
	}

	for i, val := range result[0 : len(result)-1] {
		if val != 0 {
			panic(fmt.Errorf("unexpected output: %d (position %d)", val, i))
		}
	}
	return result[len(result)-1]
}

func part2(input []string) int {
	program := toInstructions(input)

	result, err := executeV3(5, program)
	if err != nil {
		panic(err)
	}

	return result[0]
}

func loadInputs(i int, s int, program []int) []int {
	s = s / 100
	flags := make([]bool, 3)
	for j := range flags {
		flags[j] = (s % 10) > 0
		s = s / 10
	}
	inputs := make([]int, 2)
	for j := range inputs {
		inputs[j] = program[i+j+1]
		if flags[j] {
			continue
		}
		inputs[j] = program[inputs[j]]
	}
	return inputs
}

func executeV3(input int, program []int) ([]int, error) {
	outputs := make([]int, 0)
	lastOpcode := -1
	for i := 0; i < len(program); {
		s := program[i]
		opcode := s % 100

		switch opcode {
		case 1:
			inputs := loadInputs(i, s, program)
			program[program[i+3]] = inputs[0] + inputs[1]
			i = i + 4
		case 2:
			inputs := loadInputs(i, s, program)
			program[program[i+3]] = inputs[0] * inputs[1]
			i = i + 4
		case 3:
			program[program[i+1]] = input
			i = i + 2
		case 4:
			o := program[i+1]
			if k := s / 100; k > 0 {
				outputs = append(outputs, o)
			} else {
				outputs = append(outputs, program[o])
			}
			i = i + 2
		case 5:
			inputs := loadInputs(i, s, program)
			i = i + 3
			if inputs[0] != 0 {
				i = inputs[1]
			}
		case 6:
			inputs := loadInputs(i, s, program)
			i = i + 3
			if inputs[0] == 0 {
				i = inputs[1]
			}
		case 7:
			inputs := loadInputs(i, s, program)
			output := i + 3
			program[output] = 0
			if inputs[0] < inputs[1] {
				program[output] = 1
			}
			i = i + 4
		case 8:
			inputs := loadInputs(i, s, program)
			output := i + 3
			program[output] = 0
			if inputs[0] == inputs[1] {
				program[output] = 1
			}
			i = i + 4
		case 99:
			if lastOpcode == 4 {
				return outputs, nil
			}
			return outputs, fmt.Errorf("unexpected halt")
		default:
			return outputs, fmt.Errorf("unexpected instruction: %d", opcode)
		}
		lastOpcode = opcode
	}
	return outputs, fmt.Errorf("unexpected EOF")
}
