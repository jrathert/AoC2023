/*
 * Day DD of AoC 2023
 *
 * Idea: Text
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

var testinput = `xxx
`

func part01() {
	startTime := time.Now()
	cnt := 0
	total := 0
	lines := getInput()
	for _, line := range lines {
		log.Printf("Processing line %v with len %v\n", cnt, len(line))
		cnt++
	}
	total = cnt
	elapsed := time.Since(startTime)
	fmt.Printf("Result part 01 (%v): %v\n\n", elapsed, total)
}

func part02() {
	// startTime := time.Now()
	// cnt := 0
	// total := 0
	// lines := getInput()
	// for _, line := range lines {
	// 	log.Printf("Processing line %v with len %v\n", cnt, len(line))
	// 	total++
	// }
	// elapsed := time.Since(startTime)
	// fmt.Printf("Result part 01 (%v): %v\n\n", elapsed, total)
	fmt.Println("Part 02 not implemented yet")
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
