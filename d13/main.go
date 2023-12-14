/*
 * Day 13 of AoC 2023
 *
 * Idea: Simple string matching. Only challenge was to understand that
 * you MUST find a flipping element.
 * Algo for part 1 and 2 is almost the same: Go through the lines until you find
 * two that are (almost) the same, and then traverse to the borders until done (or fail).
 * For part 2 make sure you have exactly one flipping involved
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

var testinput = `#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#
`

var part2 = false

func printPattern(pattern []string) {
	for _, s := range pattern {
		fmt.Println(s)
	}
}

func isMirroredAndFlipped(s []string, idx int) bool {
	ret := 1
	haveFlipped := false
	for ret <= idx && idx+ret <= len(s) {
		s1 := s[idx+ret-1]
		s2 := s[idx-ret]
		n, _ := numDiff(s1, s2)
		if n == 0 {
			// same
			ret++
		} else {
			if !haveFlipped {
				if n == 1 {
					// fmt.Printf("Changing char %v in line %v/%v: %v\n", f, idx-ret, idx+ret-1, s2)
					haveFlipped = true
					ret++
				} else {
					// n must be > 1
					return false
				}
			} else {
				return false
			}
		}
	}
	return haveFlipped
}

func isMirrored(s []string, idx int) bool {
	ret := 1
	for ret <= idx && idx+ret <= len(s) {
		s1 := s[idx+ret-1]
		s2 := s[idx-ret]
		if s1 == s2 {
			ret++
		} else {
			return false
		}
	}
	return true
}

func getMirrorIdx(pattern []string, withFlipping bool) int {

	i := 1
	for i < len(pattern) {
		if i > 0 {
			n, _ := numDiff(pattern[i], pattern[i-1])
			if n < 2 {
				// fmt.Printf("checking line %v and %v (%v)\n", i, i-1, pattern[i])
				mirrored := false
				if !withFlipping {
					mirrored = isMirrored(pattern, i)
				} else {
					mirrored = isMirroredAndFlipped(pattern, i)
				}
				if mirrored {
					return i
				}
			}
		}
		i++
	}
	return 0
}

func transpose(s []string) []string {
	ret := make([]string, len(s[0]))
	for i := range s {
		for j := range s[i] {
			ret[j] = ret[j] + string(s[i][j])
		}
	}
	return ret
}

func numDiff(s1, s2 string) (int, int) {
	if len(s1) != len(s2) {
		return -1, -1
	}
	first := -1
	num := 0
	for i := range s1 {
		if s1[i] != s2[i] {
			if first == -1 {
				first = i
			}
			num++
		}
	}
	return num, first
}

func process(withFlipping bool) int {
	cnt := 0
	total := 0
	lines := getInput()
	lines = append(lines, "")
	pattern := []string{}
	for _, line := range lines {
		// log.Printf("Processing line %v with len %v\n", cnt, len(line))
		if len(line) == 0 {
			// process
			cnt++
			val := 0
			transposed := ""
			val = getMirrorIdx(pattern, withFlipping)
			if val != 0 {
				total += 100 * val
			} else {
				// transpose pattern
				// fmt.Println("Transposing...")
				pattern = transpose(pattern)
				val = getMirrorIdx(pattern, withFlipping)
				if val != 0 {
					transposed = "(transposed)"
					total += val
				}
			}
			fmt.Printf("%v: val %v %v\n", cnt, val, transposed)
			pattern = []string{}
		} else {
			pattern = append(pattern, line)
		}
	}
	return total
}

func part01() {
	startTime := time.Now()
	total := process(false)
	// total = cnt
	elapsed := time.Since(startTime)
	fmt.Printf("Result part 01 (%v): %v\n\n", elapsed, total)
}

func part02() {
	startTime := time.Now()
	total := process(true)
	// total = cnt
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
