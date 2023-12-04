/*
 * Day 04 of AoC 2023
 *
 * Idea: For each card, sort list of winning numbers and ones you have to
 * allow for an efficient (O(n)) algorithm per card
 * One could easily solve part01 and part02 in one run...
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package main

import (
	"aoc23/tools"
	"flag"
	"fmt"
	"log"
	"slices"
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

var testinput = `Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11`

func part01() {
	cnt := 0
	total := 0
	lines := getInput()
	for _, line := range lines {
		log.Printf("Processing line %v with len %v\n", cnt, len(line))
		parts := strings.Split(line, ":")
		blocks := strings.Split(parts[1], "|")
		win := tools.ReadInts(blocks[0])
		have := tools.ReadInts(blocks[1])
		slices.Sort(win)
		slices.Sort(have)
		j := 0
		val := 0
		for i := 0; i < len(win); i++ {
			for {
				if j >= len(have) {
					break
				}
				if have[j] < win[i] {
					j++
				} else if have[j] == win[i] {
					if val == 0 {
						val = 1
					} else {
						val *= 2
					}
					j++
					break
				} else {
					break
				}
			}
		}
		cnt++
		total += val
	}
	fmt.Printf("Result part 01: %v\n", total)
}

func part02() {
	cnt := 0
	total := 0
	lines := getInput()
	stack := make(map[int]int)

	for i := 0; i < len(lines); i++ {
		stack[i] = 1
	}

	for {
		line := lines[cnt]
		log.Printf("Processing line %v with len %v\n", cnt, len(line))
		parts := strings.Split(line, ":")
		blocks := strings.Split(parts[1], "|")
		win := tools.ReadInts(blocks[0])
		have := tools.ReadInts(blocks[1])
		slices.Sort(win)
		slices.Sort(have)
		j := 0
		val := 0
		for i := 0; i < len(win); i++ {
			for {
				if j >= len(have) {
					break
				}
				if have[j] < win[i] {
					j++
				} else if have[j] == win[i] {
					val++
					j++
					break
				} else {
					break
				}
			}
		}

		for i := 1; i <= val; i++ {
			if cnt+i >= len(lines) {
				break
			}
			stack[cnt+i] += stack[cnt]
		}

		cnt++
		if cnt == len(lines) {
			break
		}
	}

	for i := 0; i < len(stack); i++ {
		total += stack[i]
	}
	fmt.Printf("Result part 02: %v\n", total)
}

func getInput() []string {
	if TESTMODE {
		return tools.ReadInputString(testinput)
	} else {
		return tools.ReadInputFile(*inputfile)
	}

}
