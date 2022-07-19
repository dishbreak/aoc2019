package main

import (
	"context"
	"errors"
	"fmt"
)

type IntcodeComputer struct {
	program map[int64]int64
	rel     int64
}

func (i *IntcodeComputer) Load(program []int64) {
	i.program = make(map[int64]int64)
	for n, val := range program {
		i.program[int64(n)] = val
	}
}

func (i *IntcodeComputer) loadInputs(idx int64) []int64 {
	s := i.program[idx]
	s = s / 100
	flags := make([]ParameterMode, 3)
	for j := range flags {
		flags[j] = ParameterMode(s % 10)
		s = s / 10
	}
	inputs := make([]int64, 2)
	for j := range inputs {
		inputs[j] = i.getInput(flags[j], int64(j+1)+idx)
	}
	return inputs
}

type ParameterMode int64

const (
	PositionalMode ParameterMode = iota
	ImmediateMode
	RelativeMode
)

func (i *IntcodeComputer) getInput(mode ParameterMode, addr int64) int64 {
	switch mode {
	case RelativeMode:
		addr = i.program[addr]
		return i.program[i.rel+addr]
	case ImmediateMode:
		return i.program[addr]
	case PositionalMode:
		addr = i.program[addr]
		return i.program[addr]
	default:
		panic(fmt.Errorf("invalid parameter mode %d", mode))
	}
}

func (i *IntcodeComputer) Execute(ctx context.Context, input []int64) (int64, error) {
	inputStream := make(chan int64, len(input))
	output := make(chan int64, 1)
	errStream := make(chan error)
	done := make(chan interface{})
	defer close(inputStream)
	defer close(output)

	for _, val := range input {
		inputStream <- val
	}

	go i.Simulate(ctx, "execution", inputStream, output, errStream, nil, done)

	result := int64(-1)
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
		case <-done:
			return result, nil
		}
	}
}

func (i *IntcodeComputer) Simulate(ctx context.Context, name string, input <-chan int64, output chan<- int64, errStream chan<- error, term chan<- int64, done chan<- interface{}) {
	lastOutput := int64(-1)
	for pc := int64(0); true; {
		select {
		case <-ctx.Done():
			return
		default:
		}
		s := i.program[pc]
		opcode := s % 100
		switch opcode {
		case 1:
			inputs := i.loadInputs(pc)
			i.program[i.program[pc+3]] = inputs[0] + inputs[1]
			pc = pc + 4
		case 2:
			inputs := i.loadInputs(pc)
			i.program[i.program[pc+3]] = inputs[0] * inputs[1]
			pc = pc + 4
		case 3:
			inputValue := <-input
			i.program[i.program[pc+1]] = inputValue
			pc = pc + 2
		case 4:
			o := i.program[pc+1]
			if k := s / 100; k <= 0 {
				o = i.program[o]
			}
			lastOutput = o
			output <- o
			pc = pc + 2
		case 5:
			inputs := i.loadInputs(pc)
			pc = pc + 3
			if inputs[0] != 0 {
				pc = inputs[1]
			}
		case 6:
			inputs := i.loadInputs(pc)
			pc = pc + 3
			if inputs[0] == 0 {
				pc = inputs[1]
			}
		case 7:
			inputs := i.loadInputs(pc)
			output := i.program[pc+3]
			i.program[output] = 0
			if inputs[0] < inputs[1] {
				i.program[output] = 1
			}
			pc = pc + 4
		case 8:
			inputs := i.loadInputs(pc)
			output := i.program[pc+3]
			i.program[output] = 0
			if inputs[0] == inputs[1] {
				i.program[output] = 1
			}
			pc = pc + 4
		case 9:
			mode := ParameterMode(s / 100)
			i.rel = i.getInput(mode, pc+1)
			pc = pc + 2
		case 99:
			if term != nil {
				term <- lastOutput
			}
			if done != nil {
				close(done)
			}
			return
		default:
			errStream <- fmt.Errorf("unexpected instruction: %d", opcode)
		}
	}
	errStream <- fmt.Errorf("unexpected EOF")
}
