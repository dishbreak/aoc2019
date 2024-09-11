package main

import (
	"context"
	"fmt"
	"image"

	"github.com/dishbreak/aoc-common/lib"
	"github.com/dishbreak/aoc2019/util"
)

func main() {
	input, err := lib.GetInput("inputs/day13.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 1: %d\n", part1(input[0]))
	fmt.Printf("Part 2: %d\n", part2(input[0]))
}

func part1(input string) int {
	pgm, err := util.ParseIntcode(input)
	if err != nil {
		panic(err)
	}

	comp := util.IntcodeComputer{}
	comp.Load(pgm)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	output := make(chan int64)
	errStream := make(chan error)

	go comp.Simulate(ctx, "brick game", nil, output, errStream)

	gameBoard := make(map[image.Point]int)
	i := 0
	tilePt := image.Pt(-1, -1)
	for o := range output {
		switch i {
		case 0: // store the X
			tilePt.X = int(o)
		case 1: // store the Y
			tilePt.Y = int(o)
		case 2: // process based on id
			gameBoard[tilePt] = int(o)
		}
		i++
		i = i % 3
	}

	acc := 0
	for _, tileId := range gameBoard {
		if tileId == 2 {
			acc++
		}
	}

	return acc
}

func part2(input string) int {
	pgm, err := util.ParseIntcode(input)
	if err != nil {
		panic(err)
	}

	// Memory address 0 represents the number of quarters that have been
	// inserted; set it to 2 to play for free.
	pgm[0] = 2

	comp := util.IntcodeComputer{}
	comp.Load(pgm)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	joystick := make(chan int64, 1)
	output := make(chan int64)
	errStream := make(chan error)

	go comp.Simulate(ctx, "brick game", joystick, output, errStream)

	i := 0
	var ballCoords image.Point
	ballInitialized := false

	var tilePt image.Point
	var tileId int
	score := -1
	for o := range output {
		switch i {
		case 0:
			tilePt.X = int(o)
		case 1:
			tilePt.Y = int(o)
		case 2:
			tileId = int(o)
		}

		i = (i + 1) % 3
		if i != 0 {
			continue
		}

		if tilePt.X == -1 && tilePt.Y == 0 {
			score = tileId
			continue
		}

		if tileId != 4 {
			continue
		}

		if !ballInitialized {
			ballCoords = tilePt
			ballInitialized = true
			continue
		}

		travel := tilePt.Sub(ballCoords)
		joystick <- int64(travel.X)

		ballCoords = tilePt
	}

	return score
}
