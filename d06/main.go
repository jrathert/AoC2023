/*
 * Day 06 of AoC 2023
 *
 * Idea: PQ Formular
 *   The function is
 *      t = time in msec
 *      d = duration allowed
 *      m = max distance to beat
 *   distance reached is (d-t)*t = dt-t^2
 *   taking max distance (+1) into account:
 *      f(t) = -t^2 + dt - m
 *   Result is the null-values and all values in between
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package main

import (
	"aoc23/tools"
	"flag"
	"fmt"
	"log"
	"math"
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

var testinput = `Time:      7  15   30
Distance:  9  40  200`

func part01() {
	total := 1
	lines := getInput()
	d := tools.ReadInts(strings.Split(lines[0], ":")[1])
	m := tools.ReadInts(strings.Split(lines[1], ":")[1])

	for i := 0; i < len(d); i++ {
		p := float64(d[i])
		q := float64(m[i] + 1)
		// formula: -t^2+d[i]t-m[i] => t^2-d[i]t+m[i]
		// pq-formula: d[i]/2 +/- sqrt (d[i]*d[i]/4 - m[i])
		v := math.Pow(p/2.0, 2.0) - q
		first := math.Ceil(p/2.0 - math.Sqrt(v))
		sec := math.Floor(p/2.0 + math.Sqrt(v))
		log.Printf("x1, x2 = %v, %v\n", first, sec)
		// log.Printf(("%v, %v\n"), first, sec)
		total *= int(sec - first + 1)
	}

	fmt.Printf("Result part 01: %v\n", total)
}

func part02() {
	total := 1
	lines := getInput()
	d := readSeparatedInt(strings.Split(lines[0], ":")[1])
	m := readSeparatedInt(strings.Split(lines[1], ":")[1])

	p := float64(d)
	q := float64(m + 1)
	// formula: -t^2+d[i]t-m[i] => t^2-d[i]t+m[i]
	// pq-formula: d[i]/2 +/- sqrt (d[i]*d[i]/4 - m[i])
	v := math.Pow(p/2.0, 2.0) - q
	first := math.Ceil(p/2.0 - math.Sqrt(v))
	sec := math.Floor(p/2.0 + math.Sqrt(v))
	log.Printf("x1, x2 = %v, %v\n", first, sec)
	// log.Printf(("%v, %v\n"), first, sec)
	total *= int(sec - first + 1)

	fmt.Printf("Result part 02: %v\n", total)
}

func readSeparatedInt(s string) int {
	s = strings.Replace(s, " ", "", -1)
	return tools.Str2Int(s)
}

func getInput() []string {
	if TESTMODE {
		return tools.ReadInputString(testinput)
	} else {
		return tools.ReadInputFile(*inputfile)
	}

}
