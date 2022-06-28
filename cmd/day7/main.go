package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/dishbreak/aoc-common/lib"
)

func main() {
	input, err := lib.GetInput("inputs/day7.txt")
	if err != nil {
		panic(err)
	}

	program := toProgram(input[0])

	fmt.Printf("Part 1: %d\n", part1(program))
	fmt.Printf("Part 2: %d\n", part2(program))
}

func toProgram(input string) []int {
	parts := strings.Split(input, ",")
	result := make([]int, len(parts))

	for i, val := range parts {
		result[i], _ = strconv.Atoi(val)
	}

	return result
}

func part1(input []int) int {
	ctx := context.Background()
	configs := generateConfigs(ctx, []int{0, 1, 2, 3, 4})
	results := executeConfigs(ctx, input, configs)

	acc := 0
	for i := range results {
		if i > acc {
			acc = i
		}
	}
	return acc
}

func generateConfigs(ctx context.Context, values []int) <-chan []int {
	valStream := make(chan []int)
	go func(values []int) {
		defer close(valStream)
		configs := Permute(values)
		for _, config := range configs {
			select {
			case <-ctx.Done():
				return
			case valStream <- config:
				continue
			}
		}
	}(values)

	return valStream
}

func executeConfigs(ctx context.Context, program []int, input <-chan []int) <-chan int {
	valStream := make(chan int)

	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case s, ok := <-input:
					if !ok {
						return
					}
					o, err := executeConfig(ctx, program, s)
					if err != nil {
						log.Printf("error executing: %s", err)
						continue
					}
					valStream <- o
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(valStream)
	}()

	return valStream
}

func executeConfig(ctx context.Context, program []int, settings []int) (int, error) {
	signal := 0

	for _, setting := range settings {
		c := &IntcodeComputer{}
		c.Load(program)
		var err error
		signal, err = c.Execute(ctx, []int{setting, signal})
		if err != nil {
			return signal, err
		}
	}
	return signal, nil
}

func part2(program []int) int {
	ctx := context.Background()
	settings := generateConfigs(ctx, []int{5, 6, 7, 8, 9})
	results := feedbackSimulator(ctx, program, settings)

	acc := 0
	for i := range results {
		if i > acc {
			acc = i
		}
	}

	return acc
}

func feedbackSimulator(ctx context.Context, program []int, settings <-chan []int) <-chan int {
	valStream := make(chan int)

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case s, ok := <-settings:
					if !ok {
						return
					}
					val, err := runSimulation(ctx, s, program)
					if err != nil {
						log.Printf("ERROR: %s\n", err)
						continue
					}
					valStream <- val
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(valStream)
	}()

	return valStream
}

func runSimulation(ctx context.Context, settings []int, program []int) (int, error) {
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()

	errStream := make(chan error)
	defer close(errStream)

	link := make([]chan int, 6)
	for i := range link[:5] {
		link[i] = make(chan int, 1)
		defer close(link[i])
	}
	link[5] = make(chan int)

	for i, setting := range settings {
		link[i] <- setting
	}

	/*
		   -+-0->[A]-1->[B]-2->[C]-3->[D]-4->[E]--.
			|                                     | [0]
			'--------------------0----------------+
												  | [5]* Only written to on halt
												  v
	*/
	pipelineStage(ctx, "A", program, link[0], link[1], errStream, nil)     //A
	pipelineStage(ctx, "B", program, link[1], link[2], errStream, nil)     //B
	pipelineStage(ctx, "C", program, link[2], link[3], errStream, nil)     //C
	pipelineStage(ctx, "D", program, link[3], link[4], errStream, nil)     //D
	pipelineStage(ctx, "E", program, link[4], link[0], errStream, link[5]) //E

	link[0] <- 0

	for {
		select {
		case <-ctx.Done():
			return -1, ctx.Err()
		case r := <-link[5]:
			log.Println("simulation complete")
			return r, nil
		case err := <-errStream:
			return -1, err
		}
	}

}

func pipelineStage(ctx context.Context, name string, program []int, input <-chan int, output chan<- int, errs chan<- error, term chan<- int) {
	go func() {
		c := &IntcodeComputer{}
		c.Load(program)
		c.Simulate(ctx, name, input, output, errs, term, nil)
	}()
}
