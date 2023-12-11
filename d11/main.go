/*
 * Day 11 of AoC 2023
 *
 * Idea: This was easy - whenever I see a rectangular exercise, especially one as today,
 * I ask myself what would be the best way to store relevant information
 * NOT using one or more huge arrays (i.e. a 2D matrix) as most likely you
 * run into trouble in the second part.
 * Today, just store the positions of the galaxies and for each column and row
 * how much they really represent. When calculating distances, just use this
 * vector. Part 01 took a bit longer than necessary, including the usual
 * one-off-error when calculating the distance. Part 02 was then done in 1 minute.
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package main

import (
	"aoc23/tools"
	"flag"
	"fmt"
	"log"
	"regexp"
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

var testinput = `...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....`

func part01() {
	startTime := time.Now()
	total := calcDistances(2)
	elapsed := time.Since(startTime)
	fmt.Printf("Result part 01 (%v): %v\n\n", elapsed, total)
}

func part02() {
	startTime := time.Now()
	total := calcDistances(1000000)
	elapsed := time.Since(startTime)
	fmt.Printf("Result part 02 (%v): %v\n\n", elapsed, total)
}

func partSum(start, end int, diffs []int) int {
	if end < start {
		start, end = end, start
	}
	val := 0
	for i := start; i < end; i++ {
		val += diffs[i]
	}
	return val
}

func calcDistances(replace int) int {

	var galaxies []tools.Position // all galaxies
	var cols, rows int            // size of the universe (same values)
	var xDist []int               // actual distance per column
	var yDist []int               // actual distance per row

	re := regexp.MustCompile(`#`)
	lines := getInput()
	cnt := 0
	for _, line := range lines {
		if cnt == 0 {
			rows = len(line)
			cols = rows
			xDist = make([]int, cols)
			yDist = make([]int, rows)
		}

		gx := re.FindAllIndex([]byte(line), -1)
		if len(gx) > 0 {
			yDist[cnt] = 1 // there is (at least) one galaxy in this row
		}
		for _, g := range gx {
			galaxies = append(galaxies, tools.Position{g[0], cnt})
			xDist[g[0]] = 1 // there is (at least) one galaxy in this column
		}
		cnt++
	}

	// now fix the empty rows and columns by actual number
	for i := 0; i < rows; i++ {
		if xDist[i] == 0 {
			xDist[i] = replace
		}
		if yDist[i] == 0 {
			yDist[i] = replace
		}
	}

	// fmt.Printf("Galaxies: %v\n", galaxies)
	// fmt.Printf("numGPerCol: %v\n", numGPerCol)
	// fmt.Printf("numGPerRow: %v\n", numGPerRow)

	// for each pair of galaxies, calculate the x and y diff
	// by calculating sum of respective slice in xDist/yDist vector
	// and sum them up
	total := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			a := galaxies[i]
			b := galaxies[j]
			diffX := partSum(a[0], b[0], xDist)
			diffY := partSum(a[1], b[1], yDist)
			dist := diffX + diffY
			// fmt.Printf("dist (%v)-(%v)): %v (%v, %v)\n", galaxies[i], galaxies[j], dist, diffX, diffY)
			total += dist
		}
	}
	return total
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
