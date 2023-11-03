package main

import (
	"context"
	"fmt"
	"image"

	"github.com/dishbreak/aoc-common/lib"
	"github.com/dishbreak/aoc2019/util"
)

func main() {
	input, err := lib.GetInput("inputs/day11.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 1: %d\n", part1(input[0]))
}

func part1(input string) int {
	pgm, err := util.ParseIntcode(input)
	if err != nil {
		panic(fmt.Errorf("failed to load program: %w", err))
	}

	comp := &util.IntcodeComputer{}
	comp.Load(pgm)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputStream := make(chan int64, 1)
	output := make(chan int64)
	errStream := make(chan error)
	defer close(inputStream)

	go comp.Simulate(ctx, "hull painter", inputStream, output, errStream)
	return <-runRobot(ctx, inputStream, output)
}

func runRobot(ctx context.Context, inputStream chan<- int64, output <-chan int64) <-chan int {
	result := make(chan int)

	go func() {
		defer close(result)
		hull := make(map[image.Point]int64)
		p := image.Pt(0, 0)
		c := NewCompass()
		inputStream <- hull[p]
		outputCtr := 0
		for {
			select {
			case <-ctx.Done():
				result <- -1
				return
			case o, ok := <-output:
				if !ok {
					result <- len(hull)
					return
				}
				outputCtr++
				if outputCtr%2 == 1 {
					hull[p] = o
				} else {
					if o == 1 {
						c = c.Left
					} else {
						c = c.Right
					}
					p = p.Add(c.Direction)
					inputStream <- hull[p]
				}
			}
		}
	}()

	return result
}
