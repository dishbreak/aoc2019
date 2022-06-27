package main

import (
	"context"
	"errors"
	"fmt"
)

type IntcodeComputer struct {
	program []int
}

func (i *IntcodeComputer) Load(program []int) {
	i.program = make([]int, len(program))
	copy(i.program, program)
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

func (i *IntcodeComputer) Execute(ctx context.Context, input []int) (int, error) {
	inputStream := make(chan int, len(input))
	output := make(chan int, 1)
	errStream := make(chan error)

	defer close(inputStream)
	defer close(output)

	for _, val := range input {
		inputStream <- val
	}

	go i.Simulate(ctx, inputStream, output, errStream)

	result := -1
	for {
		select {
		case <-ctx.Done():
			return -1, ctx.Err()
		case o, ok := <-output:
			if !ok {
				return -1, errors.New("simulation hung up")
			}
			result = o
		case err := <-errStream:
			return result, err
		}
	}
}

func (i *IntcodeComputer) Simulate(ctx context.Context, input <-chan int, output chan<- int, errStream chan<- error) {
	lastOpcode := -1
	for pc := 0; pc < len(i.program); {
		select {
		case <-ctx.Done():
			errStream <- ctx.Err()
			return
		default:
		}
		s := i.program[pc]
		opcode := s % 100

		switch opcode {
		case 1:
			inputs := loadInputs(pc, s, i.program)
			i.program[i.program[pc+3]] = inputs[0] + inputs[1]
			pc = pc + 4
		case 2:
			inputs := loadInputs(pc, s, i.program)
			i.program[i.program[pc+3]] = inputs[0] * inputs[1]
			pc = pc + 4
		case 3:
			i.program[i.program[pc+1]] = <-input
			pc = pc + 2
		case 4:
			o := i.program[pc+1]
			if k := s / 100; k > 0 {
				output <- o
			} else {
				output <- i.program[o]
			}
			pc = pc + 2
		case 5:
			inputs := loadInputs(pc, s, i.program)
			pc = pc + 3
			if inputs[0] != 0 {
				pc = inputs[1]
			}
		case 6:
			inputs := loadInputs(pc, s, i.program)
			pc = pc + 3
			if inputs[0] == 0 {
				pc = inputs[1]
			}
		case 7:
			inputs := loadInputs(pc, s, i.program)
			output := i.program[pc+3]
			i.program[output] = 0
			if inputs[0] < inputs[1] {
				i.program[output] = 1
			}
			pc = pc + 4
		case 8:
			inputs := loadInputs(pc, s, i.program)
			output := i.program[pc+3]
			i.program[output] = 0
			if inputs[0] == inputs[1] {
				i.program[output] = 1
			}
			pc = pc + 4
		case 99:
			if lastOpcode != 4 {
				errStream <- fmt.Errorf("unexpected halt")
			}
			close(errStream)
			return
		default:
			errStream <- fmt.Errorf("unexpected instruction: %d", opcode)
		}
		lastOpcode = opcode
	}
	errStream <- fmt.Errorf("unexpected EOF")
}
