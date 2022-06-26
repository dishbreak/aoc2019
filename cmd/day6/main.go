package main

import (
	"fmt"

	"github.com/dishbreak/aoc-common/lib"
)

func main() {
	input, err := lib.GetInput("inputs/day6.txt")
	if err != nil {
		panic(err)
	}

	input = input[:len(input)-1]
	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}

func part1(input []string) int {
	graph := BuildGraph(input)
	q := make([]*Planet, 0)

	com := graph["COM"]

	q = append(q, com.Satellites...)

	acc := 0
	for len(q) != 0 {
		p := q[len(q)-1]
		q = q[:len(q)-1]

		p.Orbits = p.Parent.Orbits + 1
		acc += p.Orbits

		q = append(q, p.Satellites...)
	}
	return acc
}

func part2(input []string) int {
	graph := BuildGraph(input)

	san := graph["SAN"]

	for p, steps := san.Parent, 0; p != nil; p, steps = p.Parent, steps+1 {
		p.Visited = true
		p.Steps = steps
	}

	you := graph["YOU"]
	for p, steps := you.Parent, 0; p != nil; p, steps = p.Parent, steps+1 {
		if p.Visited {
			return p.Steps + steps
		}
	}
	return -1
}
