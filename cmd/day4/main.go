package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dishbreak/aoc-common/lib"
)

func main() {
	input, err := lib.GetInput("inputs/day4.txt")
	if err != nil {
		panic(err)
	}

	parts := strings.Split(input[0], "-")
	start, _ := strconv.Atoi(parts[0])
	end, _ := strconv.Atoi(parts[1])

	fmt.Printf("Part 1: %d\n", part1(start, end))
	fmt.Printf("Part 2: %d\n", part2(start, end))
}

func toDigits(s int) []int {
	digits := make([]int, 0)
	for s > 0 {
		digits = append(digits, s%10)
		s = s / 10
	}
	return digits
}

func isPossiblePasscode(passcode int) bool {
	digits := toDigits(passcode)

	repeat := false
	for i := 1; i < len(digits); i++ {
		if digits[i] == digits[i-1] {
			repeat = true
		}
		if digits[i] > digits[i-1] {
			return false
		}
	}
	return repeat
}

func isPossiblePasscodeV2(passcode int) bool {
	digits := toDigits(passcode)

	counter := make([]int, 10)

	counter[digits[0]]++
	for i := 1; i < len(digits); i++ {
		counter[digits[i]]++
		if digits[i] > digits[i-1] {
			return false
		}
	}

	fmt.Println(passcode, counter)

	for _, count := range counter {
		if count == 2 {
			return true
		}
	}

	return false

}

func bruteForce(start, end int, test func(int) bool) int {
	acc := 0

	for i := start; i <= end; i++ {
		if test(i) {
			acc++
		}
	}

	return acc
}

func part1(start, end int) int {
	return bruteForce(start, end, isPossiblePasscode)
}

func part2(start, end int) int {
	return bruteForce(start, end, isPossiblePasscodeV2)
}
