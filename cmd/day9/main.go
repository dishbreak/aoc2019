package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dishbreak/aoc-common/lib"
)

func main() {
	input, err := lib.GetInput("inputs/day9.txt")
	if err != nil {
		panic(err)
	}

	parts := strings.Split(input[0], ",")
	program := make([]int64, len(parts))

	for i, val := range parts {
		parsed, _ := strconv.Atoi(val)
		program[i] = int64(parsed)
	}

	fmt.Printf("Part 1: %d\n", part1(program))
	fmt.Printf("Part 2: %d\n", part2(program))
}

func boost(program []int64, input int64) int64 {
	i := &IntcodeComputer{}
	i.Load(program)
	result, err := i.Execute(context.Background(), []int64{input})
	if err != nil {
		panic(err)
	}
	return result
}

const (
	testCode  int64 = 1
	boostCode int64 = 2
)

func part1(input []int64) int64 {
	return boost(input, testCode)
}

func part2(input []int64) int64 {
	return boost(input, boostCode)
}
