package main

import (
	"fmt"
	"image"
	"sync"

	"github.com/dishbreak/aoc-common/lib"
)

func main() {
	input, err := lib.GetInput("inputs/day10.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 1: %d\n", part1(input[:len(input)-1]))
}

type asteroidField struct {
	min, max image.Point
	r        image.Rectangle
	pts      map[image.Point]bool
	total    int
}

func parseSpace(input []string) (a asteroidField) {
	a.pts = make(map[image.Point]bool)

	for y, line := range input {
		for x, char := range line {
			if char == '.' {
				continue
			}
			a.total++
			pt := image.Pt(x, y)
			a.pts[pt] = true
		}
	}
	a.max = image.Pt(len(input[0]), len(input)).Add(image.Pt(-1, -1))

	a.r = image.Rectangle{a.min, a.max}.Inset(-1)
	return
}

func (a asteroidField) contains(p image.Point) bool {
	xInBounds := p.X >= a.min.X && p.X <= a.max.X
	yInBounds := p.Y >= a.min.Y && p.Y <= a.max.Y

	return xInBounds && yInBounds
}

func part1(input []string) int {
	space := parseSpace(input)

	results := make(chan int, space.total)
	var wg sync.WaitGroup
	wg.Add(space.total)

	for pt := range space.pts {
		go func(pt image.Point) {
			defer wg.Done()
			results <- discover(space, pt)
		}(pt)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	max := 0
	for val := range results {
		if max < val {
			max = val
		}
	}

	return max
}

type direction image.Point

var (
	north     = image.Point{0, -1}
	south     = image.Point{0, 1}
	east      = image.Point{1, 0}
	west      = image.Point{-1, 0}
	northwest = north.Add(west)
	northeast = north.Add(east)
	southwest = south.Add(west)
	southeast = south.Add(east)
)

type box struct {
	nwCorner, neCorner, swCorner, seCorner image.Point
}

func InitBox(origin image.Point) *box {
	b := &box{}
	b.neCorner = origin.Add(northeast)
	b.nwCorner = origin.Add(northwest)
	b.seCorner = origin.Add(southeast)
	b.swCorner = origin.Add(southwest)
	return b
}

func (b *box) Expand() {
	b.neCorner = b.neCorner.Add(northeast)
	b.nwCorner = b.nwCorner.Add(northwest)
	b.seCorner = b.seCorner.Add(southeast)
	b.swCorner = b.swCorner.Add(southwest)
}

func (b *box) Covers(a asteroidField) bool {
	if bMin := b.nwCorner; bMin.X >= a.min.X || bMin.Y >= a.min.Y {
		return false
	}
	if bMax := b.seCorner; bMax.X <= a.max.X || bMax.Y <= a.max.Y {
		return false
	}
	return true
}

func abs(i int) int {
	if i < 0 {
		return -1 * i
	}
	return i
}

func reduce(p image.Point) image.Point {
	if p.X == 0 {
		y := 1
		if p.Y < 0 {
			y = -1
		}
		return image.Pt(0, y)
	}

	if p.Y == 0 {
		x := 1
		if p.X < 0 {
			x = -1
		}
		return image.Pt(x, 0)
	}

	a, b := abs(p.X), abs(p.Y)

	// euclidean GCD method
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}

	return p.Div(a)
}

func discover(space asteroidField, origin image.Point) (seen int) {
	blocked := make(map[image.Point]bool)

	trace := func(start, end, direction image.Point) (found int) {
		for b := start; !b.Eq(end); b = b.Add(direction) {
			if blocked[b] || !space.pts[b] {
				continue
			}
			found++

			v := reduce(b.Sub(origin))
			for n := b.Add(v); space.contains(n); n = n.Add(v) {
				if space.pts[n] {
					blocked[n] = true
				}
			}
		}
		return
	}
	for b := InitBox(origin); !b.Covers(space); b.Expand() {
		seen += trace(b.nwCorner, b.neCorner, east)
		seen += trace(b.neCorner, b.seCorner, south)
		seen += trace(b.seCorner, b.swCorner, west)
		seen += trace(b.swCorner, b.nwCorner, north)
	}
	return
}
