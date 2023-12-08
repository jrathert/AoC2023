/*
 * Day 05 of AoC 2023
 *
 * Idea: That was a tough one! Part 01 super easy, just process the different seeds through
 * the maps (make sure to not process any seed twice, there is no overlap)
 * For part 02 process full slices and make sure to respect any overlaps appropriately.
 * Took me some time to derive a proper alog for part 02. As usual, overlapping slices confuse me.
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package main

import (
	"aoc23/tools"
	"flag"
	"fmt"
	"log"
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

	// testfilter()
	// testsplit()
	part01()
	part02()
}

var testinput = `seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4`

type triple [3]int
type tuple [2]int

func part01() {

	lines := getInput()
	lines = append(lines, "")

	seeds := tools.ReadInts(strings.Split(lines[0], ":")[1])
	maps := readmap(lines[1:])

	minval := -1
	for _, s := range seeds {
		val := s
		for _, m := range maps {
			for _, c := range m {
				n := mapval(c, val)
				if n != val {
					val = n
					break
				}
			}
		}
		if minval == -1 || val < minval {
			minval = val
		}
		log.Printf("Seed %v - result %v - minval %v\n", s, val, minval)
	}

	fmt.Printf("Result part 01: %v\n", minval)
}

func part02() {
	lines := getInput()
	lines = append(lines, "")

	vals := tools.ReadInts(strings.Split(lines[0], ":")[1])
	seeds := []tuple{}
	for i := 0; i < len(vals); i += 2 {
		seeds = append(seeds, tuple{vals[i], vals[i+1]})
	}
	maps := readmap(lines[1:])

	minval := -1
	for _, s := range seeds {
		ranges := make([]tuple, 1)
		ranges[0] = s

		// For each seedmap go through all lines ("filters")
		// if some range is mapped, keep it aside ("processed"), use
		// unmapped ranges ("newranges") as input for next filter
		// process until last line is reached - then reset "processed" and start
		// with new seedmap
		cnt := 0
		for _, seedmap := range maps {
			// start a new map
			// log.Printf("Start examining map %v...\n", cnt)
			processed := make([]tuple, 0)

			// iterate over all lines of map
			for _, filter := range seedmap {

				newranges := make([]tuple, 0)
				for _, rng := range ranges {
					mapped, nomapped := applyfilter(filter, rng)
					// log.Printf("Filter %v by %v: %v | %v\n", rng, filter, mapped, nomapped)
					if mapped != nil {
						processed = append(processed, tuple(mapped))
					}
					newranges = append(newranges, nomapped...)
				}
				ranges = newranges
			}
			// take all processesd as well as remaining non-processed as start point for next seedmap
			ranges = append(processed, ranges...)
			cnt++
		}

		// we end up with a number of ranges and need to identify the minimum start value
		val := -1
		for i := 0; i < len(ranges); i++ {
			v := ranges[i][0]
			if val == -1 || val > v {
				val = v
			}
		}
		// if that is less than the overall minimum value, keep it
		if minval == -1 || minval > val {
			minval = val
		}
		log.Printf("Seed %v - val %v - minval %v\n", s, val, minval)
	}

	fmt.Printf("Result part 02: %v\n", minval)
}

// applies the filter [newstart, oldstart, length] to the range [start, length]
// E.g.:   applyfilter([0, 25, 5][20, 20]) will create [0, 5] and [[20,5],[30,10]]
// If there is no match, it returns an empty set
//
//	[ [20, 5], [0, 5], [30, 10] ]
func applyfilter(filter triple, rng tuple) ([]int, []tuple) {

	diff := filter[1] - filter[0]
	if rng[0]+rng[1] <= filter[1] || rng[0] >= filter[1]+filter[2] {
		// no match, fully left or fully right
		return nil, []tuple{rng}
	} else if rng[0] >= filter[1] && rng[0]+rng[1] <= filter[1]+filter[2] {
		// full coverage, range is completely within filter and will just be shifted
		return []int{rng[0] - diff, rng[1]}, nil
	} else {
		llen := filter[1] - rng[0]
		rlen := rng[0] + rng[1] - (filter[1] + filter[2])
		start := max(rng[0], filter[1])
		end := min(rng[0]+rng[1], filter[1]+filter[2])
		len := end - start
		mapped := []int{start - diff, len}
		var nomap []tuple = make([]tuple, 0)
		// if both are > 0, the filter is completly within the range, we have three parts
		// if only one of them is > 0, we have the range
		// both < 0 can not happen - that would mean we have full coverage which is covered above
		if llen > 0 {
			// range starts before filter
			nomap = append(nomap, tuple{rng[0], llen})
		}
		if rlen > 0 {
			// range goes beyond filter
			nomap = append(nomap, [2]int{rng[0] + rng[1] - rlen, rlen})
		}
		return mapped, nomap
	}
}

func mapval(mapper [3]int, val int) int {
	if val >= mapper[1] && val < mapper[1]+mapper[2] {
		v := mapper[0] + (val - mapper[1])
		return v
	} else {
		return val
	}
}

func readmap(lines []string) [][]triple {
	maps := [][]triple{}
	mode := 0 // 0 waiting, 1 awake
	logline := ""
	var currmap []triple
	for _, line := range lines {
		if mode == 0 && strings.Contains(line, "map:") {
			logline = line
			currmap = []triple{}
			mode = 1
		} else if mode == 1 && len(line) != 0 {
			currmap = append(currmap, triple(tools.ReadInts(line)))
		} else if mode == 1 && len(line) == 0 {
			maps = append(maps, currmap)
			mode = 0
			log.Printf("Reading %v %v entries", logline, len(currmap))
		}
	}
	return maps
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
