/*
 * Day 08 of AoC 2023
 *
 * Idea: Very basic search for part 01, for part 02 reuse part 01 and calculate least common denominator, as brute-force would
 * take much too long
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

var testinput = `LLR

AAA = (BBB, BBB)
BBB = (AAA, ZZZ)
ZZZ = (ZZZ, ZZZ)`

var testinput2 = `LR

	11A = (11B, XXX)
	11B = (XXX, 11Z)
	11Z = (11B, XXX)
	22A = (22B, XXX)
	22B = (22C, 22C)
	22C = (22Z, 22Z)
	22Z = (22B, 22B)
	XXX = (XXX, XXX)`

type DesertMap map[string][2]string

func search(start string, match string, desert DesertMap, orders string) int {
	i := 0
	pos := start
	re := regexp.MustCompile(match)
	for {
		char := string(orders[i%len(orders)])
		if char == "L" {
			pos = desert[pos][0]
		} else {
			pos = desert[pos][1]
		}
		i++
		if i%len(orders) == 0 && re.MatchString(pos) {
			break
		}
	}
	fmt.Printf("Result for start %v and match %v: %v\n", start, match, i)
	return i
}

func buildDesert(part int) (DesertMap, string) {
	var lines []string
	if part == 1 {
		lines = getInput()
	} else {
		lines = getInput(testinput2)
	}
	cnt := 0
	orders := ""
	re := regexp.MustCompile(`([0-9A-Z]{3}) = \(([0-9A-Z]{3}), ([0-9A-Z]{3})\)`)
	var desert DesertMap = make(DesertMap)
	for _, line := range lines {
		if cnt == 0 {
			orders = line
		} else if len(line) > 0 {
			parts := re.FindAllStringSubmatch(line, -1)[0]
			desert[parts[1]] = [2]string{parts[2], parts[3]}
		}
		cnt++
	}
	return desert, orders
}

func part01() {
	desert, orders := buildDesert(1)
	total := search("AAA", `.*ZZZ`, desert, orders)
	fmt.Printf("Result part 01: %v\n", total)
}

func part02() {
	desert, orders := buildDesert(2)

	positions := []string{}
	re := regexp.MustCompile(`.*A`)
	for pos := range desert {
		if re.MatchString(pos) {
			positions = append(positions, pos)
		}
	}

	vals := make([]int, len(positions))
	for i, p := range positions {
		vals[i] = search(p, `.*Z`, desert, orders)
	}
	total := tools.LCM(vals[0], vals[1], vals[2:]...)
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
