/*
 * Day 16 of AoC 2023
 *
 * Idea: You basically just follow the beams (using a stack once it is split up)
 * until they either reach the border or run into a circle. When the stack is
 * empty, you are done
 * Only smart idea needed was how to identify the "circle" as one needed to
 * check not only the field visited, but from what direction. I used a copy of
 * the grid and a bitfield to track from where it was visited.
 * For part 2 I just used brute fors. Was not worth to spend more time thinking
 * about good caching
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package main

import (
	"aoc23/tools"
	"flag"
	"fmt"
	"log"
	"time"
)

var TESTMODE = true
var nt = flag.Bool("nt", false, "exec no test mode")
var inputfile = flag.String("f", "input.txt", "name of input file")
var part1only = flag.Bool("1", false, "run part 1 only")
var part2only = flag.Bool("2", false, "run part 2 only")

func main() {
	flag.Parse()
	if *nt {
		TESTMODE = false
	}
	log.SetPrefix("  ")
	log.SetFlags(0)

	if !*part2only {
		part01()
	}
	if !*part1only {
		part02()
	}
}

var testinput = `.|...\....
|.-.\.....
.....|-...
........|.
..........
.........\
..../.\\..
.-.-/..|..
.|....-|.\
..//.|....`

const (
	North byte = 1
	East  byte = 2
	West  byte = 4
	South byte = 8
)

type Beam struct {
	pos       tools.Position
	direction byte
}

func (b *Beam) advance(grid *tools.Matrix) bool {
	if b.direction == North {
		if b.pos[1] == 0 {
			return false
		} else {
			b.pos[1]--
			return true
		}
	} else if b.direction == East {
		if b.pos[0] == grid.Cols()-1 {
			return false
		} else {
			b.pos[0]++
			return true
		}
	} else if b.direction == South {
		if b.pos[1] == grid.Rows()-1 {
			return false
		} else {
			b.pos[1]++
			return true
		}
	} else if b.direction == West {
		if b.pos[0] == 0 {
			return false
		} else {
			b.pos[0]--
			return true
		}
	}
	return false
}

func alreadyVisited(b *Beam, chk *tools.Matrix) bool {
	c, _ := chk.ValueAtPos(b.pos)
	if c&b.direction > 0 {
		return true
	} else {
		chk.SetValueAtPos(b.pos, c+b.direction)
		return false
	}
}

func followBeam(startBeam *Beam, grid *tools.Matrix) int {
	check := tools.NewMatrix(grid.Rows(), grid.Cols())

	var allBeams tools.Stack[*Beam]
	allBeams.Push(startBeam)

	b, _ := allBeams.Pop()

	for b != nil {
		ok := processBeam(b, grid, &check, &allBeams)
		if !ok {
			b, _ = allBeams.Pop()
		}
	}
	return check.CountNonZero()
}

func run(grid *tools.Matrix, part int) int {
	if part == 1 {
		startBeam := Beam{tools.Position{-1, 0}, East}
		// startBeam := Beam{tools.Position{3, -1}, South}
		return followBeam(&startBeam, grid)
	} else if part == 2 {
		retval := 0

		for i := 0; i < grid.Rows(); i++ {
			startBeam1 := Beam{tools.Position{-1, i}, East}
			startBeam2 := Beam{tools.Position{grid.Cols(), i}, West}
			val1 := followBeam(&startBeam1, grid)
			val2 := followBeam(&startBeam2, grid)
			retval = max(val1, val2, retval)
		}

		for j := 0; j < grid.Cols(); j++ {
			startBeam1 := Beam{tools.Position{j, -1}, South}
			startBeam2 := Beam{tools.Position{j, grid.Rows()}, North}
			val1 := followBeam(&startBeam1, grid)
			val2 := followBeam(&startBeam2, grid)
			retval = max(val1, val2, retval)
		}
		return retval
	} else {
		return -1
	}
}

func processBeam(b *Beam, grid *tools.Matrix, chk *tools.Matrix, bs *tools.Stack[*Beam]) bool {

	// pos := b.pos
	ok := b.advance(grid)
	// fmt.Printf("Advancing beam from %v,%v -> %v,%v\n", pos[1], pos[0], b.pos[1], b.pos[0])
	if !ok {
		return false
	}
	if alreadyVisited(b, chk) {
		return false
	}

	c, _ := grid.ValueAtPos(b.pos)
	if c == '.' {
	} else if c == '|' {
		if b.direction == East || b.direction == West {
			b.direction = North
			nb := Beam{b.pos, South}
			bs.Push(&nb)
		}
	} else if c == '-' {
		if b.direction == North || b.direction == South {
			b.direction = West
			nb := Beam{b.pos, East}
			bs.Push(&nb)
		}
	} else if c == '\\' {
		switch b.direction {
		case North:
			b.direction = West
		case East:
			b.direction = South
		case South:
			b.direction = East
		case West:
			b.direction = North
		}
	} else if c == '/' {
		switch b.direction {
		case North:
			b.direction = East
		case East:
			b.direction = North
		case South:
			b.direction = West
		case West:
			b.direction = South
		}
	}
	return true
}

func part01() {
	startTime := time.Now()
	cnt := 0
	total := 0

	var grid tools.Matrix
	lines := getInput()
	for _, line := range lines {
		grid.AddLine(line)
		cnt++
	}
	// fmt.Printf("*** Grid ***\n%v", grid)

	total = run(&grid, 1)
	elapsed := time.Since(startTime)
	fmt.Printf("Result part 01 (%v): %v\n\n", elapsed, total)
}

func part02() {
	startTime := time.Now()
	cnt := 0
	total := 0

	var grid tools.Matrix
	lines := getInput()
	for _, line := range lines {
		grid.AddLine(line)
		cnt++
	}
	// fmt.Printf("*** Grid ***\n%v", grid)

	total = run(&grid, 2)
	elapsed := time.Since(startTime)
	fmt.Printf("Result part 02 (%v): %v\n\n", elapsed, total)
}

func getInput(inputs ...string) []string {
	if TESTMODE {
		input := testinput
		if len(inputs) > 0 {
			input = inputs[0]
		}
		return tools.ReadInputString(input)
	} else {
		return tools.ReadInputFile(*inputfile)
	}

}
