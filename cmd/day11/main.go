package main

import (
	"context"
	"fmt"
	"image"
	"math"
	"strings"

	"github.com/dishbreak/aoc-common/lib"
	"github.com/dishbreak/aoc2019/util"
)

func main() {
	input, err := lib.GetInput("inputs/day11.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 1: %d\n", part1(input[0]))
	fmt.Printf("Part 2:\n%s\n", part2(input[0]))
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
	return len(<-runRobot(ctx, inputStream, output, 0))
}

func runRobot(ctx context.Context, inputStream chan<- int64, output <-chan int64, startingValue int64) <-chan map[image.Point]int64 {
	result := make(chan map[image.Point]int64)

	go func() {
		defer close(result)
		hull := make(map[image.Point]int64)
		p := image.Pt(0, 0)
		hull[p] = startingValue
		c := NewCompass()
		inputStream <- hull[p]
		outputCtr := 0
		for {
			select {
			case <-ctx.Done():
				result <- nil
				return
			case o, ok := <-output:
				if !ok {
					result <- hull
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

func part2(input string) string {
	pgm, err := util.ParseIntcode(input)
	if err != nil {
		panic(fmt.Errorf("failed to parse program: %w", err))
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
	hull := <-runRobot(ctx, inputStream, output, 1)

	minPt := image.Pt(math.MaxInt, math.MaxInt)
	maxPt := image.Pt(math.MinInt, math.MinInt)

	for p := range hull {
		if p.X < minPt.X {
			minPt.X = p.X
		} else if p.X > maxPt.X {
			maxPt.X = p.X
		}
		if p.Y < minPt.Y {
			minPt.Y = p.Y
		} else if p.Y > maxPt.Y {
			maxPt.Y = p.Y
		}
	}

	var sb strings.Builder

	for y := maxPt.Y; y >= minPt.Y; y-- {
		for x := maxPt.X; x >= minPt.X; x-- {
			if hull[image.Pt(x, y)] == 0 {
				sb.WriteRune(' ')
				continue
			}
			sb.WriteRune('#')
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}
