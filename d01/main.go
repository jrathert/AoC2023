/*
 * Day 01 of AoC 2023
 *
 * Idea: Part 1 was easy with just using regexp.
 * Part 2 was more difficult due to go regexp not recognizing overlapping entries.
 * I decided to reverse the string and the regexp to make sure to find the last
 * entry.
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

func part01() {
	var input []string
	if TESTMODE {
		input = tools.ReadInputString(`1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet`)
	} else {
		input = tools.ReadInputFile(*inputfile)
	}

	re := regexp.MustCompile(`[0-9]`)
	total := 0
	cnt := 0
	for _, line := range input {
		cnt += 1
		numArr := re.FindAllString(line, -1)
		val := tools.Str2Int(numArr[0] + numArr[len(numArr)-1])
		log.Printf("%d: %s -> %d\n", cnt, line, val)
		total += val
	}

	fmt.Println(total)
}

func part02() {
	var input []string
	if TESTMODE {
		input = tools.ReadInputString(`two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen`)
	} else {
		input = tools.ReadInputFile(*inputfile)
	}

	restr := "one|two|three|four|five|six|seven|eight|nine"
	re := regexp.MustCompile(restr + "|[1-9]")
	rere := regexp.MustCompile(tools.ReverseStr(restr) + "|[1-9]")
	total := 0
	cnt := 0
	for _, line := range input {
		cnt += 1
		first := re.FindString(line)
		last := tools.ReverseStr(rere.FindString(tools.ReverseStr(line)))
		val := str2num(first)*10 + str2num(last)

		log.Printf("%d: %s %s %s -> %d\n", cnt, line, first, last, val)
		total += val
	}

	fmt.Println(total)
}

func str2num(s string) int {
	switch s {
	// case "zero":
	// 	return 0
	case "one":
		return 1
	case "two":
		return 2
	case "three":
		return 3
	case "four":
		return 4
	case "five":
		return 5
	case "six":
		return 6
	case "seven":
		return 7
	case "eight":
		return 8
	case "nine":
		return 9
	default:
		return tools.Str2Int(s)
	}
}
