/*
 * Day 09 of AoC 2023
 *
 * Idea: Do avoid any recursions etc, just take the next line and keep track
 * of result data. Backwards is a bit more complex - there are other ways
 * that do not require the "firstval" array, but I find them difficult to
 * understand later, so kept it.
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package main

import (
	"aoc23/tools"
	"flag"
	"fmt"
	"log"
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

var testinput = `0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45`

func forwardValues(values []int) int {
	result := 0
	for {
		vlen := len(values)
		result += values[vlen-1]
		nextvals := make([]int, 0, vlen-1)
		allZeros := true
		for i := 1; i < vlen; i++ {
			diff := values[i] - values[i-1]
			nextvals = append(nextvals, diff)
			if diff != 0 {
				allZeros = false
			}
		}
		if allZeros {
			// nextvals are all zeros -> we are done
			return result
		} else {
			values = nextvals
		}
	}
}

func backwardValues(values []int) int {
	firstvals := []int{}
	for {
		vlen := len(values)
		firstvals = append(firstvals, values[0])
		nextvals := make([]int, 0, vlen-1)
		allZeros := true
		for i := 1; i < vlen; i++ {
			diff := values[i] - values[i-1]
			nextvals = append(nextvals, diff)
			if diff != 0 {
				allZeros = false
			}
		}
		if allZeros {
			// nextvals are all zeros -> we are done
			// firstvals contain all the "first values"
			result := 0
			for i := len(firstvals) - 1; i >= 0; i-- {
				result = firstvals[i] - result
			}
			return result
		} else {
			values = nextvals
		}
	}
}

func part01() {
	cnt := 0
	total := 0
	lines := getInput()

	for _, line := range lines {

		values := tools.ReadSignedInts(line)
		result := forwardValues(values)
		log.Printf("%v: result is  %v\n", cnt, result)
		total += result
		cnt++
	}
	fmt.Printf("Result part 01: %v\n", total)
}

func part02() {
	cnt := 0
	total := 0
	lines := getInput()

	for _, line := range lines {

		values := tools.ReadSignedInts(line)
		result := backwardValues(values)
		log.Printf("%v: result is  %v\n", cnt, result)
		total += result
		cnt++
	}
	fmt.Printf("Result part 02: %v\n", total)
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
