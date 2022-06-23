package main

import (
	"fmt"
	"strconv"

	"github.com/dishbreak/aoc-common/lib"
)

func main() {
	input, err := lib.GetInput("inputs/day1.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}

func part1(input []string) int {
	acc := 0
	for _, mStr := range input {
		// assume clean input.
		mass, _ := strconv.Atoi(mStr)
		acc += fuelCalcV1(mass)
	}
	return acc
}

func fuelCalcV1(mass int) int {
	div := mass / 3
	if div <= 2 {
		return 0
	}
	return div - 2
}

func part2(input []string) int {
	acc := 0
	for _, mStr := range input {
		// assume clean input.
		mass, _ := strconv.Atoi(mStr)
		acc += fuelCalcV2(mass)
	}
	return acc
}

func fuelCalcV2(mass int) int {
	acc := 0
	for fuel := fuelCalcV1(mass); fuel > 0; fuel = fuelCalcV1(fuel) {
		acc += fuel
	}
	return acc
}
