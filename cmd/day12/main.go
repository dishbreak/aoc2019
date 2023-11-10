package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

func main() {
	f, err := os.Open("inputs/day12.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Printf("Part 1: %d\n", part1(f))
	f.Seek(0, 0)
	fmt.Printf("Part 2: %d\n", part2(f))
}

type frame struct {
	pos, vel [][3]int
}

func add(one, other [3]int) [3]int {
	var result [3]int
	for i := range other {
		result[i] = one[i] + other[i]
	}

	return result
}

func abs(i int) int {
	if i < 0 {
		return -1 * i
	}
	return i
}

func sum(in [3]int) int {
	return abs(in[0]) + abs(in[1]) + abs(in[2])
}

func (f *frame) interact() *frame {
	nf := &frame{
		pos: make([][3]int, len(f.pos)),
		vel: make([][3]int, len(f.vel)),
	}
	copy(nf.pos, f.pos)
	copy(nf.vel, f.vel)

	for j := 0; j < len(f.pos); j++ {
		for i := 0; i < len(f.pos); i++ {
			if i == j {
				continue
			}
			for dim := 0; dim < 3; dim++ {
				switch {
				case f.pos[j][dim] < f.pos[i][dim]:
					nf.vel[j][dim]++
				case f.pos[j][dim] > f.pos[i][dim]:
					nf.vel[j][dim]--
				default: // values are equal, take no action.
				}
			}
		}
	}

	for i := range nf.pos {
		nf.pos[i] = add(nf.pos[i], nf.vel[i])
	}

	return nf
}

func (f *frame) totalEnergy() (acc int) {
	for i := range f.pos {
		acc += sum(f.pos[i]) * sum(f.vel[i])
	}
	return
}

func (f *frame) String() string {
	var sb strings.Builder
	for i := range f.pos {
		sb.WriteString(
			fmt.Sprintf(
				"pos=<x=%3d, y=%3d, z=%3d>, vel=<x=%3d, y=%3d, z=%3d>\n",
				f.pos[i][0], f.pos[i][1], f.pos[i][2],
				f.vel[i][0], f.vel[i][1], f.vel[i][2],
			),
		)
	}

	return sb.String()
}

func part1(r io.Reader) int {
	return simulate(r, 1000)
}

func simulate(r io.Reader, iterations int) int {
	f := &frame{
		pos: load(r),
	}
	f.vel = make([][3]int, len(f.pos))

	// fmt.Println(f)

	for i := 0; i < iterations; i++ {
		f = f.interact()
		// fmt.Println(f)
	}

	return f.totalEnergy()
}

func load(r io.Reader) [][3]int {
	lines := make([][3]int, 0)

	s := bufio.NewScanner(r)
	for s.Scan() {
		lines = append(lines, FromString(s.Text()))
	}

	return lines
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func part2(r io.Reader) int {
	coords := load(r)

	reports := make(chan int)

	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func(dim int) {
			defer wg.Done()
			vel := make([]int, len(coords))
			pos := make([]int, len(coords))

			for j := range coords {
				pos[j] = coords[j][dim]
			}

			seen := make(map[string]bool)
			ct := 0
			for {
				hash := fmt.Sprint(pos, vel)
				if seen[hash] {
					break
				}
				seen[hash] = true
				for i := 0; i < len(vel); i++ {
					for j := 0; j < len(vel); j++ {
						if i == j {
							continue
						}
						switch {
						case pos[i] < pos[j]:
							vel[i]++
						case pos[i] > pos[j]:
							vel[i]--
						default:
						}
					}
				}

				for i := range pos {
					pos[i] += vel[i]
				}

				ct++
			}
			reports <- ct
		}(i)
	}

	go func() {
		wg.Wait()
		close(reports)
	}()

	cyclePts := make([]int, 0)
	for report := range reports {
		cyclePts = append(cyclePts, report)
	}

	return LCM(cyclePts[0], cyclePts[1], cyclePts[2])
}
