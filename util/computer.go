package util

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

func ParseIntcode(input string) ([]int64, error) {
	pts := strings.Split(input, ",")
	result := make([]int64, len(pts))

	for i, pt := range pts {
		n, err := strconv.Atoi(pt)
		if err != nil {
			return nil, fmt.Errorf("invalid token %s: %w", pt, err)
		}
		result[i] = int64(n)
	}

	return result, nil
}

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
		inputs[j] = i.getValue(flags[j], int64(j+1)+idx)
	}
	return inputs
}

type ParameterMode int64

const (
	PositionalMode ParameterMode = iota
	ImmediateMode
	RelativeMode
)

func (i *IntcodeComputer) getValue(mode ParameterMode, addr int64) int64 {
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
	defer close(inputStream)

	for _, val := range input {
		inputStream <- val
	}

	go i.Simulate(ctx, "execution", inputStream, output, errStream)

	result := int64(-1)
	for {
		select {
		case <-ctx.Done():
			return -1, ctx.Err()
		case o, ok := <-output:
			if !ok {
				return result, nil
			}
			result = o
		case err := <-errStream:
			return result, err
		}
	}
}

func (i *IntcodeComputer) resolve(mode ParameterMode, ref int64) (addr int64) {
	switch mode {
	case RelativeMode:
		addr = ref + i.rel
	default: // positional mode (immedate mode unsupported for opcode 3)
		addr = ref
	}
	return
}

func (i *IntcodeComputer) Simulate(ctx context.Context, name string, input <-chan int64, output chan<- int64, errStream chan<- error) {
	defer close(output)
	for pc := int64(0); true; {
		select {
		case <-ctx.Done():
			return
		default:
		}
		s := i.program[pc]
		opcode := s % 100
		switch opcode {
		/*
			Opcode 1 adds together numbers read from two positions and stores
			the result in a third position. The three integers immediately
			after the opcode tell you these three positions - the first two
			indicate the positions from which you should read the input values,
			and the third indicates the position at which the output should be
			stored.
		*/
		case 1:
			inputs := i.loadInputs(pc)
			destAddr := i.resolve(ParameterMode(s/10000), i.program[pc+3])
			i.program[destAddr] = inputs[0] + inputs[1]
			pc = pc + 4
		/*
			Opcode 2 works exactly like opcode 1, except it multiplies the two
			inputs instead of adding them. Again, the three integers after the
			opcode indicate where the inputs and outputs are, not their values.
		*/
		case 2:
			inputs := i.loadInputs(pc)
			destAddr := i.resolve(ParameterMode(s/10000), i.program[pc+3])
			i.program[destAddr] = inputs[0] * inputs[1]
			pc = pc + 4
		/*
			Opcode 3 takes a single integer as input and saves it to the
			position given by its only parameter. For example, the instruction
			3,50 would take an input value and store it at address 50.
		*/
		case 3:
			inputValue := <-input
			i.program[i.resolve(ParameterMode(s/100), i.program[pc+1])] = inputValue
			pc = pc + 2
		/*
			Opcode 4 outputs the value of its only parameter. For example, the
			instruction 4,50 would output the value at address 50.
		*/
		case 4:
			o := i.getValue(ParameterMode(s/100), pc+1)
			output <- o
			pc = pc + 2
		/*
			Opcode 5 is jump-if-true: if the first parameter is non-zero, it
			sets the instruction pointer to the value from the second parameter.
			Otherwise, it does nothing.
		*/
		case 5:
			inputs := i.loadInputs(pc)
			pc = pc + 3
			if inputs[0] != 0 {
				pc = inputs[1]
			}
		/*
			Opcode 6 is jump-if-false: if the first parameter is zero, it sets
			the instruction pointer to the value from the second parameter.
			Otherwise, it does nothing.
		*/
		case 6:
			inputs := i.loadInputs(pc)
			pc = pc + 3
			if inputs[0] == 0 {
				pc = inputs[1]
			}
		/*
			Opcode 7 is less than: if the first parameter is less than the
			second parameter, it stores 1 in the position given by the third
			parameter. Otherwise, it stores 0.
		*/
		case 7:
			inputs := i.loadInputs(pc)
			destAddr := i.resolve(ParameterMode(s/10000), i.program[pc+3])
			i.program[destAddr] = 0
			if inputs[0] < inputs[1] {
				i.program[destAddr] = 1
			}
			pc = pc + 4
		/*
			Opcode 8 is equals: if the first parameter is equal to the second
			parameter, it stores 1 in the position given by the third parameter.
			Otherwise, it stores 0.
		*/
		case 8:
			inputs := i.loadInputs(pc)
			destAddr := i.resolve(ParameterMode(s/10000), i.program[pc+3])
			i.program[destAddr] = 0
			if inputs[0] == inputs[1] {
				i.program[destAddr] = 1
			}
			pc = pc + 4
		/*
			Opcode 9 adjusts the relative base by the value of its only
			parameter. The relative base increases (or decreases, if the value
			is negative) by the value of the parameter.
		*/
		case 9:
			mode := ParameterMode(s / 100)
			i.rel += i.getValue(mode, pc+1)
			pc = pc + 2
		/*
			Opcode 99 halts the program.
		*/
		case 99:
			return
		default:
			errStream <- fmt.Errorf("unexpected instruction: %d", opcode)
		}
	}
	errStream <- fmt.Errorf("unexpected EOF")
}
