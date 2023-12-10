/*
 * Day 10 of AoC 2023
 *
 * Idea: Build a matrix/maze. Identify start field, its type and one of the two
 * next fields. Then follow the loop, keep track of visitied fields and measure
 * distance. Last find inner fields by counting the edges.
 *
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package main

import (
	"aoc23/tools"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

var TESTMODE = true
var nt = flag.Bool("nt", false, "exec no test mode")
var inputfile = flag.String("f", "input.txt", "name of input file")

func main() {
	flag.Parse()
	if *nt {
		TESTMODE = false
	}
	log.SetPrefix("  ")
	log.SetFlags(0)

	part01()
	part02()
}

var testinput = `7-F7-
.FJ|7
SJLL7
|F--J
LJ.LJ`

func getNumTest(num int) string {
	if num == 1 {
		return `...........
		.S-------7.
		.|F-----7|.
		.||.....||.
		.||.....||.
		.|L-7.F-J|.
		.|..|.|..|.
		.L--J.L--J.
		...........`
	} else if num == 2 {
		return `..........
		.S------7.
		.|F----7|.
		.||OOOO||.
		.||OOOO||.
		.|L-7F-J|.
		.|II||II|.
		.L--JL--J.
		..........`
	} else if num == 3 {
		return `.F----7F7F7F7F-7....
		.|F--7||||||||FJ....
		.||.FJ||||||||L7....
		FJL7L7LJLJ||LJ.L-7..
		L--J.L7...LJS7F-7L7.
		....F-J..F7FJ|L7L7L7
		....L7.F7||L7|.L7L7|
		.....|FJLJ|FJ|F7|.LJ
		....FJL-7.||.||||...
		....L---J.LJ.LJLJ...`
	} else if num == 4 {
		return `FF7FSF7F7F7F7F7F---7
		L|LJ||||||||||||F--J
		FL-7LJLJ||||||LJL-77
		F--JF--7||LJLJ7F7FJ-
		L---JF-JLJ.||-FJLJJ7
		|F|F-JF---7F7-L7L|7|
		|FFJF7L7F-JF7|JL---7
		7-L-JL7||F7|L7F-7F7|
		L.L7LFJ|||||FJL7||LJ
		L7JLJL-JLJLJL--JLJ.L`
	} else {
		return testinput
	}
}

func part01() {
	startTime := time.Now()

	lines := getInput()
	var maze tools.Matrix
	for _, line := range lines {
		maze.AddLine(strings.TrimSpace(line))
	}
	// fmt.Printf("%v\n", maze)

	start, ok := maze.FindField('S')
	if !ok {
		fmt.Println("Can not find startpos!")
		return
	}
	fmt.Printf("start pos: %v\n", start)

	distance := -1
	total := 0
	for i := 0; i < 4; i++ {
		var pos tools.Position
		var ok bool
		if i == 0 {
			pos, ok = maze.LeftOf(start)
		} else if i == 1 {
			pos, ok = maze.RightOf(start)
		} else if i == 2 {
			pos, ok = maze.AboveOf(start)
		} else if i == 3 {
			pos, ok = maze.BelowOf(start)
		}
		distance = 1

		fmt.Printf("new try (%v): %v, ok: %v\n", i, pos, ok)
		if !ok {
			continue
		}

		var last = start
		for {
			newpos, ok := move(maze, last, pos)
			// fmt.Printf("  newpos: %v, ok: %v\n", newpos, ok)

			if ok {
				last = pos
				pos = newpos
				distance++
				if pos == start {
					distance /= 2
					break
				}
			} else {
				distance = -1
				break
			}
		}
		fmt.Printf("Round %v: greatest distance = %v\n", i, distance)
		if distance > total {
			total = distance
		}

	}
	//total := max(lengths[0], lengths[1], lengths[2], lengths[3])

	elapsed := time.Since(startTime)
	fmt.Printf("Result part 01 (%v): max distance %v\n\n", elapsed, total)
}

func part02() {
	startTime := time.Now()
	// build maze
	lines := getInput(getNumTest(1))
	var maze tools.Matrix
	for _, line := range lines {
		maze.AddLine(strings.TrimSpace(line))
	}
	// fmt.Printf("%v\n", maze)

	// find start position within maze
	start, ok := maze.FindField('S')
	if !ok {
		fmt.Println("Can not find startpos!")
		return
	}
	fmt.Printf("start pos: %v\n", start)

	// create a copy to track coverage by the loop later
	var copyMaze tools.Matrix
	line := strings.Repeat(".", maze.Cols())
	for i := 0; i < maze.Rows(); i++ {
		copyMaze.AddLine(line)
	}
	copyMaze.SetValueAtPos(start, '+')
	// fmt.Printf("%v\n", copytools.Maze)

	// replace the start character with the appropriate one and return
	// one of the two potential next positions
	pos := replaceStartChar(&maze, start)
	fmt.Printf("next pos: %v\n", pos)

	// assume we have gone from "last" (=start) to "pos", i.e. distance is already 1
	var distance = 1
	var last = start

	// now continue stepping to the next position depending on
	// char at the postion pos until you reach start again
	for {
		newpos, ok := move(maze, last, pos)
		if ok {
			copyMaze.SetValueAtPos(pos, '+')
			last = pos
			pos = newpos
			distance++
			if pos == start {
				distance /= 2
				break
			}
		}
	}

	// fmt.Printf("%v\n", copytools.Maze)

	// now identify inner points. Coming from the oiter border, inner points
	// can be identified by an uneven number of "crossings" - a crossing
	// is either a '|' or a combination of 'F*J' or 'L*7'
	// of course, we only need to consider fields we did not visit on the loop -
	// info on this we find in the copytools.Maze, that we also use to write
	// back info on outer ('O') or inner ('I') fields
	inner := 0
	lastCross := byte(' ')
	for y := 0; y < maze.Rows(); y++ {
		counter := 0
		for x := 0; x < maze.Cols(); x++ {
			cc, _ := copyMaze.Value(x, y)
			if cc == '.' {
				if counter%2 == 0 {
					copyMaze.SetValue(x, y, 'O')
				} else {
					copyMaze.SetValue(x, y, 'I')
					inner++
				}
			} else {
				c, _ := maze.Value(x, y)
				if c == '|' {
					counter++
				} else if c == 'F' || c == 'L' {
					lastCross = c
				} else if c == 'J' {
					if lastCross == 'F' {
						counter++
					}
					lastCross = ' '
				} else if c == '7' {
					if lastCross == 'L' {
						counter++
					}
					lastCross = ' '
				}
			}
		}
	}
	// fmt.Printf("%v\n", copytools.Maze)

	elapsed := time.Since(startTime)
	fmt.Printf("Result part 02 (%v): distance = %v; inner points = %v\n\n", elapsed, distance, inner)
}

func replaceStartChar(maze *tools.Matrix, start tools.Position) tools.Position {
	var pos, ret tools.Position
	var ok bool
	var mask = 0
	pos, ok = maze.LeftOf(start)
	if ok {
		c, _ := maze.ValueAtPos(pos)
		if c == '-' || c == 'L' || c == 'F' {
			ret = pos
			mask += 1
		}
	}
	pos, ok = maze.RightOf(start)
	if ok {
		c, _ := maze.ValueAtPos(pos)
		if c == '-' || c == 'J' || c == '7' {
			ret = pos
			mask += 2
		}
	}

	if mask < 3 {
		pos, ok = maze.AboveOf(start)
		if ok {
			c, _ := maze.ValueAtPos(pos)
			if c == '|' || c == 'F' || c == '7' {
				ret = pos
				mask += 4
			}
		}
		if mask < 5 {
			pos, ok = maze.BelowOf(start)
			if ok {
				c, _ := maze.ValueAtPos(pos)
				if c == '|' || c == 'L' || c == 'J' {
					ret = pos
					mask += 8
				}
			}
		}
	}

	var repl byte = 'S'
	if mask == 3 {
		repl = '-'
	} else if mask == 5 {
		repl = 'J'
	} else if mask == 9 {
		repl = '7'
	} else if mask == 6 {
		repl = 'L'
	} else if mask == 10 {
		repl = 'F'
	} else if mask == 12 {
		repl = '|'
	}
	fmt.Printf("Replacing start char with '%c'\n", repl)
	maze.SetValueAtPos(start, repl)
	return ret
}

func move(m tools.Matrix, last tools.Position, current tools.Position) (tools.Position, bool) {
	c, _ := m.ValueAtPos(current)
	// fmt.Printf("Examining %v -> %v (%c) %v\n", old, this, c, corners)
	if c == '|' {
		if last.IsBelowOf(current) {
			return m.AboveOf(current)
		} else if last.IsAboveOf(current) {
			return m.BelowOf(current)
		} else {
			fmt.Printf("Could not get here: %v -> %v (%c)\n", last, current, c)
		}
	} else if c == '-' {
		if last.IsLeftOf(current) {
			return m.RightOf(current)
		} else if last.IsRightOf(current) {
			return m.LeftOf(current)
		} else {
			fmt.Printf("Could not get here: %v -> %v (%c)\n", last, current, c)
		}
	} else if c == 'L' {
		if last.IsAboveOf(current) {
			return m.RightOf(current)
		} else if last.IsRightOf(current) {
			return m.AboveOf(current)
		} else {
			fmt.Printf("Could not get here: %v -> %v (%c)\n", last, current, c)
		}
	} else if c == 'J' {
		if last.IsAboveOf(current) {
			return m.LeftOf(current)
		} else if last.IsLeftOf(current) {
			return m.AboveOf(current)
		} else {
			fmt.Printf("Could not get here: %v -> %v (%c)\n", last, current, c)
		}
	} else if c == '7' {
		if last.IsBelowOf(current) {
			return m.LeftOf(current)
		} else if last.IsLeftOf(current) {
			return m.BelowOf(current)
		} else {
			fmt.Printf("Could not get here: %v -> %v (%c)\n", last, current, c)
		}
	} else if c == 'F' {
		if last.IsBelowOf(current) {
			return m.RightOf(current)
		} else if last.IsRightOf(current) {
			return m.BelowOf(current)
		} else {
			fmt.Printf("Could not get here: %v -> %v (%c)\n", last, current, c)
		}
	} else {
		fmt.Printf("  Why am I here: %v -> %v (%c)\n", last, current, c)
	}
	return current, false
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
