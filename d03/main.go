/*
 * Day 03 of AoC 2023
 *
 * Idea: Build a "rolling window" using three lines rolling over input lines
 * (add a leading and trailing empty line) and a "calculator" processing
 * these three lines
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
	"strings"
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

var testinput = `467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`

func part01() {
	total := processWith(addLineNumbers)
	fmt.Printf("Part 01 total: %v\n", total)
}

func part02() {
	total := processWith(getGearValues)
	fmt.Printf("Part 02 total: %v\n", total)
}

// a calculator is a function that gets three lines of input, analyzes the middle
// line and returns some calculated value for the line, e.g.,
//
// part 1: sum up all numbers on a line that are next to markers
// part 2: multiply all pairs adjacent to a "*" on a line
type calculator func([3]string) int

// general input processor, calling a specific calculator
func processWith(fn calculator) int {
	var input []string
	if TESTMODE {
		input = tools.ReadInputString(testinput)
	} else {
		input = tools.ReadInputFile(*inputfile)
	}

	buf := [3]string{}
	cnt := 0
	var noop string
	total := 0
	for _, line := range input {

		// initialize buf
		if cnt == 0 {
			noop = strings.Repeat(".", len(line))
			for i := 0; i < 3; i++ {
				buf[i] = noop
			}
		}

		buf[0] = buf[1]
		buf[1] = buf[2]
		buf[2] = line

		if cnt > 0 {
			num := fn(buf)
			log.Printf("%d: Added value: %v\n", cnt, num)
			total += num
		}

		cnt++
	}
	buf[0] = buf[1]
	buf[1] = buf[2]
	buf[2] = noop

	num := fn(buf)
	log.Printf("%d: Added value: %v\n", cnt, num)
	total += num

	return total
}

// calculator for part 01
func addLineNumbers(line [3]string) int {
	re := regexp.MustCompile(`[0-9]+`)
	positions := re.FindAllStringIndex(line[1], -1)
	if positions == nil {
		return 0
	}
	total := 0
	re = regexp.MustCompile(`[^0-9.]`)
	for _, p := range positions {
		min := p[0]
		max := p[1]
		val := tools.Str2Int(line[1][min:max])
		log.Printf("  Found %d at pos %v ", val, p)
		if min > 0 {
			if re.MatchString(line[1][min-1 : min]) {
				total += val
				log.Println("- adding")
				continue
			}
			min--
		}
		if max < len(line[1])-1 {
			if re.MatchString(line[1][max : max+1]) {
				total += val
				log.Println("- adding")
				continue
			}
			max++
		}
		before := re.MatchString(line[0][min:max])
		if before {
			total += val
			log.Println("- adding")
			continue
		}
		after := re.MatchString(line[2][min:max])
		if after {
			total += val
			log.Println("- adding")
			continue
		}
		log.Println("- ignoring")
	}
	return total
}

// calculator for part 02
func getGearValues(line [3]string) int {
	re := regexp.MustCompile(`\*`)
	positions := re.FindAllStringIndex(line[1], -1)
	if positions == nil {
		return 0
	}
	total := 0
	re = regexp.MustCompile(`[0-9]+`)

	// go through all occurences of "*"
	for _, p := range positions {
		log.Printf("Examining position %v\n", p[0])
		adjacents := []int{}

		// first examine line itself
		numbers := re.FindAllStringIndex(line[1], -1)
		for _, n := range numbers {
			if n[1] == p[0] || n[0] == p[1] {
				val := tools.Str2Int(line[1][n[0]:n[1]])
				adjacents = append(adjacents, val)
				log.Printf("  Adding number %v (at [%v, %v], adjacent to %v)\n", val, n[0], n[1], p[0])
			}
		}
		for _, i := range []int{0, 2} {
			numbers := re.FindAllStringIndex(line[i], -1)
			for _, n := range numbers {
				if p[0] >= n[0]-1 && p[0] <= n[1] {
					val := tools.Str2Int(line[i][n[0]:n[1]])
					adjacents = append(adjacents, val)
					log.Printf("  Adding number %v (at [%v, %v], adjacent to %v)\n", val, n[0], n[1], p[0])
				}
			}
		}

		if len(adjacents) == 2 {
			total += adjacents[0] * adjacents[1]
		}
	}
	return total
}
