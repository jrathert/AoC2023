/*
 * Day 12 of AoC 2023
 *
 * Idea: This is the hardest day for me so far. Got part 1 done recursively after
 * many different approaches,and with running into a thousand problems/bugs.
 * Finally worked, but did not scale
 * Adding a cache (6 lines of code) made part 2 work - but again took me quite
 * some time.
 * Now it is superfast. But recursion is still not my hometurf...
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

type Pump struct {
	val      int // value
	length   int // value + 1, as we need a "."
	front    int // cumulated length of all before
	back     int // cumulated length of all after
	re       *regexp.Regexp
	valcache map[int]int
}

func (p Pump) String() string {
	s := fmt.Sprintf("[%v > %v > %v (%v)]", p.front, p.length, p.back, p.front+p.length+p.back)
	return s
}

func makePumps(vals []int) []Pump {

	pumps := make([]Pump, len(vals))

	for i := range vals {
		p := Pump{}
		p.val = vals[i]
		p.length = vals[i] + 1                 // incl "."
		p.front = tools.SumInts(vals[0:i]) + i // there are i elements before i
		p.back = tools.SumInts(vals[i+1:]) + (len(vals) - (i + 1))
		p.re = regexp.MustCompile(fmt.Sprintf(`[#\?]{%v}[\?\.]`, p.val))
		p.valcache = make(map[int]int)
		pumps[i] = p
	}
	return pumps
}

var testinput = `???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1`

func numMatches(pumps []Pump, idx int, line string) int {

	// buf := strings.Repeat("  ", idx)

	p := pumps[0]

	// check if cached
	v, ok := p.valcache[idx]
	if ok {
		return v
	}

	// fmt.Printf("%vExamining %v in line %v\n", buf, p.val, line)
	if p.length+p.back > len(line) {
		// fmt.Printf("%vNot enough space (%v) in remaining string '%v'\n", buf, p.length+p.back, line)
		return 0
	}

	// fmt.Printf("Trying to match %v with '%v'\n", line, p.re.String())
	m := p.re.FindStringIndex(line)
	if m == nil {
		// fmt.Printf("%vCannot find %v in remaining string '%v'\n", buf, p.val, line)
		return 0
	} else if m[0] > 0 && strings.Contains(line[:m[0]], "#") {
		// fmt.Printf("%vCannot jump over '#' with %v in string '%v'\n", buf, p.val, line)
		return 0
	}
	newstart := m[0] + p.length
	// fmt.Printf("%vMatch in pos %v of len %v (%v)\n", buf, m[0], m[1]-m[0], newstart)

	retval := 0
	remain := len(line)
	if len(pumps) > 1 {
		if newstart+p.back > remain {
			// fmt.Printf("%vRemaining line too short/empty - but %v pump(s) remaining\n", buf, len(pumps)-1)
			return 0
		}
		// fmt.Printf("%vCalling for remaining (%v) pumps in remaining string '%v'\n", buf, len(pumps)-1, line[newstart:])
		val := numMatches(pumps[1:], idx+newstart, line[newstart:])
		if val != -1 {
			retval += val
		}
	}

	// match - now check if there is room to step one char
	// (only works if first char is  '?' and last is not '.')
	if byte(line[m[0]]) == '?' && byte(line[m[1]-1]) != '.' && m[0]+1+p.back < remain {
		// yes, there is an option, follow that path instead
		// fmt.Printf("%vTrying to shift (%v) one position in remaining string '%v'\n", buf, p.val, line[m[0]:])
		cnt := numMatches(pumps, idx+m[0]+1, line[m[0]+1:])
		if cnt != -1 {
			retval += cnt
		}
	} else if !strings.Contains(line[m[0]:m[1]], "#") && newstart+p.back < remain {
		// fmt.Printf("%vTrying to shift (%v) at end of ??? position '%v'\n", buf, p.val, line[m[0]:])
		cnt := numMatches(pumps, idx+newstart, line[newstart:])
		if cnt != -1 {
			retval += cnt
		}
	} else {
		// fmt.Printf("%vNo shift possible for (%v) in '%v'\n", buf, p.val, line[m[0]:])
	}

	if len(pumps) == 1 {
		// this was the last pump
		if newstart < remain && strings.Contains(line[newstart:], "#") {
			// fmt.Printf("%vNo more pumps, but remaining string is '%v'\n", buf, line[newstart:])
			return retval
		}

		retval += 1
	}

	// cache value for next time
	p.valcache[idx] = retval

	// fmt.Printf("%vReturning %v\n", buf, retval)
	return retval
}

func part01() {

	startTime := time.Now()
	cnt := 0
	total := 0
	lines := getInput()
	for _, line := range lines {
		// fmt.Printf("%v (%v): %v\n", cnt, len(line), line)
		parts := strings.Split(line, " ")
		vals := tools.ReadInts(parts[1])
		pumps := makePumps(vals)
		// fmt.Printf("%v\n", pumps)
		options := numMatches(pumps, 0, parts[0]+".")
		// fmt.Printf("%v: '%v' - %v arrangement(s)\n", cnt, line, options)
		// fmt.Println("==============================================================")
		total += options
		cnt++
	}
	elapsed := time.Since(startTime)
	fmt.Printf("Result part 01 (%v): %v\n\n", elapsed, total)
}

func part02() {
	startTime := time.Now()
	cnt := 0
	total := 0
	lines := getInput()
	for _, line := range lines {
		// fmt.Printf("%v (%v): %v\n", cnt, len(line), line)
		parts := strings.Split(line, " ")
		strip := parts[0]
		numbers := parts[1]
		for i := 0; i < 4; i++ {
			strip += "?" + parts[0]
			numbers += "," + parts[1]
		}
		// fmt.Printf("%v [%v] %v\n", strip, len(strip), numbers)
		vals := tools.ReadInts(numbers)
		pumps := makePumps(vals)
		// for i, p := range pumps {
		// 	fmt.Printf("%v: %v\n", i, p)
		// }
		options := numMatches(pumps, 0, strip+".")
		// fmt.Printf("%v: '%v' - %v arrangement(s)\n", cnt, line, options)
		// fmt.Println("==============================================================")
		total += options
		cnt++
	}
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
