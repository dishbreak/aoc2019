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
