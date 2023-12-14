/*
 * Day 14 of AoC 2023
 *
 * Idea: Thi was great fun in part 2. Part 1 relatively easy, just
 * proper array handling and shifting, took me only about 20 mins.
 * Part 2 was harder, as I quickly realized that there must be cycles
 * (and thus, modulo) involved, but I did not know about any proper way to
 * identify such cycles. A bit of Internet search revealed Floyd's algorithm
 * see, e.g., https://en.wikipedia.org/wiki/Cycle_detection
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package main

import (
	"aoc23/tools"
	"bytes"
	"flag"
	"fmt"
	"log"
	"strings"
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

var testinput = `O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....
`

type Board struct {
	values [][]byte
	rows   int
	cols   int
	cycled int
}

func (b Board) String() string {
	var buf strings.Builder
	for i, l := range b.values {
		for _, c := range l {
			buf.WriteByte(c)
		}
		if i < b.rows-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func (b Board) valuation() int {
	value := 0
	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
			c := b.values[i][j]
			if c == 'O' {
				value += (b.rows - i)
			}
		}
	}
	return value
}

func (b *Board) north() {
	for j := 0; j < b.cols; j++ {
		last_free := -1
		for i := 0; i < b.rows; i++ {
			c := b.values[i][j]
			if c == '.' && last_free == -1 {
				// fmt.Printf("free: %v\n", i)
				last_free = i
			} else if c == 'O' && last_free != -1 {
				// fmt.Printf("dump: %v\n", i)
				b.values[last_free][j] = 'O'
				b.values[i][j] = '.'
				last_free++
			} else if c == '#' {
				// fmt.Printf("rock: %v\n", i)
				last_free = -1
			}
		}
	}
}

func (b *Board) south() {
	for j := 0; j < b.cols; j++ {
		last_free := -1
		for i := b.rows - 1; i >= 0; i-- {
			c := b.values[i][j]
			if c == '.' && last_free == -1 {
				// fmt.Printf("free: %v\n", i)
				last_free = i
			} else if c == 'O' && last_free != -1 {
				// fmt.Printf("dump: %v\n", i)
				b.values[last_free][j] = 'O'
				b.values[i][j] = '.'
				last_free--
			} else if c == '#' {
				// fmt.Printf("rock: %v\n", i)
				last_free = -1
			}
		}
	}
}

func (b *Board) west() {
	for i := 0; i < b.rows; i++ {
		last_free := -1
		for j := 0; j < b.cols; j++ {
			c := b.values[i][j]
			if c == '.' && last_free == -1 {
				// fmt.Printf("free: %v\n", i)
				last_free = j
			} else if c == 'O' && last_free != -1 {
				// fmt.Printf("dump: %v\n", i)
				b.values[i][last_free] = 'O'
				b.values[i][j] = '.'
				last_free++
			} else if c == '#' {
				// fmt.Printf("rock: %v\n", i)
				last_free = -1
			}
		}
	}
}

func (b *Board) east() {
	for i := 0; i < b.rows; i++ {
		last_free := -1
		for j := b.cols - 1; j >= 0; j-- {
			c := b.values[i][j]
			if c == '.' && last_free == -1 {
				// fmt.Printf("free: %v\n", i)
				last_free = j
			} else if c == 'O' && last_free != -1 {
				// fmt.Printf("dump: %v\n", i)
				b.values[i][last_free] = 'O'
				b.values[i][j] = '.'
				last_free--
			} else if c == '#' {
				// fmt.Printf("rock: %v\n", i)
				last_free = -1
			}
		}
	}
}

func (b *Board) cycle() *Board {
	b.north()
	b.west()
	b.south()
	b.east()
	b.cycled++
	return b
}

func (b Board) equals(other Board) bool {
	if b.rows != other.rows || b.cols != other.cols {
		return false
	}
	for i := 0; i < b.rows; i++ {
		if !bytes.Equal(b.values[i], other.values[i]) {
			return false
		}
	}
	return true

}

func (b Board) copy() Board {
	cb := Board{}
	cb.rows = b.rows
	cb.cols = b.cols
	cb.cycled = b.cycled
	cb.values = make([][]byte, len(b.values))
	for i := range b.values {
		cb.values[i] = append([]byte(nil), b.values[i]...)
	}
	return cb
}

func makeBoard(lines []string) Board {
	board := Board{}
	board.rows = len(lines)
	board.values = make([][]byte, len(lines))
	for i, line := range lines {
		board.cols = len(line)
		board.values[i] = make([]byte, len(line))
		for j, c := range []byte(line) {
			board.values[i][j] = c
		}
	}
	return board
}

func part01() {
	startTime := time.Now()
	total := 0
	lines := getInput()
	board := makeBoard(lines)
	board.north()
	total = board.valuation()
	elapsed := time.Since(startTime)
	fmt.Printf("Result part 01 (%v): %v\n\n", elapsed, total)
}

func part02() {

	startTime := time.Now()
	total := 0
	lines := getInput()
	board := makeBoard(lines)

	// Floyd's cycle detection algorithm, see https://en.wikipedia.org/wiki/Cycle_detection
	tortoise := board.copy()
	hare := board.copy()

	numCycles := 1000000000

	tortoise.cycle()
	hare.cycle().cycle()
	idx := 1
	for !tortoise.equals(hare) {
		tortoise.cycle()
		hare.cycle().cycle()
		idx++
		// fmt.Printf("idx: %v, (%v, %v)\n", idx, tortoise.valuation(), hare.valuation())
	}
	fmt.Printf("Found cycle after %v steps (%v, %v) - %v=%v\n", idx, tortoise.cycled, hare.cycled, tortoise.valuation(), hare.valuation())

	mu := 0
	tortoise = board.copy()
	for !tortoise.equals(hare) {
		tortoise.cycle()
		hare.cycle()
		mu++
	}

	lam := 1
	hare.cycle()
	for !tortoise.equals(hare) {
		hare.cycle()
		lam++
	}

	// ok, after mu steps we run into a cycle of length lam
	// numCycles = mu + n*lam - m
	// --> numCyles - mu = n*lam - m
	// m = (numCycles - mu)%lam
	// --> iterate mu + m times to get to the right result

	m := (numCycles - mu) % lam

	fmt.Printf("mu = %v, lam = %v, m = %v => %v cycles needed\n", mu, lam, m, mu+m)

	for i := 0; i < mu+m; i++ {
		board.cycle()
	}
	total = board.valuation()

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
