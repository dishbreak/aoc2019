package main

import (
	"context"
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
	output := -1
	lastOpcode := -1
	for pc := 0; pc < len(i.program); {
		select {
		case <-ctx.Done():
			return output, ctx.Err()
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
			i.program[i.program[pc+1]] = input[0]
			input = input[1:]
			pc = pc + 2
		case 4:
			o := i.program[pc+1]
			if k := s / 100; k > 0 {
				output = o
			} else {
				output = i.program[o]
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
			if lastOpcode == 4 {
				return output, nil
			}
			return output, fmt.Errorf("unexpected halt")
		default:
			return output, fmt.Errorf("unexpected instruction: %d", opcode)
		}
		lastOpcode = opcode
	}
	return output, fmt.Errorf("unexpected EOF")
}
