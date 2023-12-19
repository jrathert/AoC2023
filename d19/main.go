/*
 * Day 19 of AoC 2023
 *
 * Idea: Part 1 was easy, just building the appropriate data structures
 * with maps and lists. As usual, parsing with regexp took quite some time - next
 * time I'll do it via read-by-character. ;-)
 * Part 2 was more tricky using recursion. Took me some time to build the
 * correct sum...
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

// l:= list of workflows that have only R
// while l not empty
//    replace item in other workflows with R
//    check if new item needs to be included in l

var testinput = `px{a<2006:qkq,m>2090:A,rfg}
pv{a>1716:R,A}
lnx{m>1548:A,A}
rfg{s<537:gd,x>2440:R,A}
qs{s>3448:A,lnx}
qkq{x<1416:A,crn}
crn{x>2662:A,R}
in{s<1351:px,qqz}
qqz{s>2770:qs,m<1801:hdj,R}
gd{a>3333:R,R}
hdj{m>838:A,pv}

{x=787,m=2655,a=1222,s=2876}
{x=1679,m=44,a=2067,s=496}
{x=2036,m=264,a=79,s=2244}
{x=2461,m=1339,a=466,s=291}
{x=2127,m=1623,a=2188,s=1013}`

type Rule struct {
	param  string
	op     string
	val    int
	result string
}

func (r Rule) String() string {
	return fmt.Sprintf("%v%v%v:%v", r.param, r.op, r.val, r.result)
}

type Workflow struct {
	name     string
	rules    []Rule
	fallback string
}

func (wf Workflow) String() string {
	return fmt.Sprintf("%v: %v | %v", wf.name, wf.rules, wf.fallback)
}

func (wf Workflow) isDead() bool {
	for _, r := range wf.rules {
		if r.result != "R" {
			return false
		}
	}
	return wf.fallback == "R"
}

type WorkflowMap map[string]*Workflow

type Part map[string]int

func (pt Part) value() int {
	val := 0
	for _, v := range pt {
		val += v
	}
	return val
}

func apply(in *Part, wf *Workflow) (string, bool) {
	for _, r := range wf.rules {
		if v, ok := (*in)[r.param]; ok {
			if r.op == "<" && v < r.val {
				return r.result, true
			} else if r.op == ">" && v > r.val {
				return r.result, true
			}
		}
	}
	return wf.fallback, true
}

type PartRange map[string][2]int

func (pr PartRange) copy() PartRange {
	ret := make(PartRange)
	for k, v := range pr {
		ret[k] = v
	}
	return ret
}

func (pr PartRange) value() int {
	val := 1
	for _, v := range pr {
		val *= v[1] - v[0] + 1
	}
	return val
}

func applyRange(pr PartRange, wfm *WorkflowMap, wfname string, totals int, idx int) int {

	// buf := strings.Repeat("  ", idx)
	sum := 0

	wf := (*wfm)[wfname]
	// fmt.Printf("%vApplying wf '%v' (%v) with PR %v (totals: %v)\n", buf, wfname, wf, pr, totals)

	for _, r := range wf.rules {
		rng := pr[r.param]
		if rng[0] <= r.val && r.val <= rng[1] {
			// we need to split it up
			var r1, r2 [2]int
			if r.op == "<" {
				r1 = [2]int{rng[0], r.val - 1} // affected by rule
				r2 = [2]int{r.val, rng[1]}     // not affected by rule, need to jump to next
			} else { // >
				r2 = [2]int{rng[0], r.val}     // not affected by rule, need to jump to next
				r1 = [2]int{r.val + 1, rng[1]} // affected by rule
			}

			pr[r.param] = r2

			// npr is the range affected by this rule
			npr := pr.copy()
			npr[r.param] = r1
			if r.result == "A" {
				// sum up range values
				v := npr.value()
				// fmt.Printf("%v  %v: + Adding up %v: %v\n", buf, wfname, npr, v)
				sum += v
			} else if r.result == "R" {
				// fmt.Printf("%v  %v: - Ignoring %v\n", buf, wfname, npr)
				// ignore range
			} else {
				sum += applyRange(npr, wfm, r.result, totals, idx+1)
			}

			// next rule
		}
	}

	// all rules passed, apply last PartRange
	if wf.fallback == "A" {
		v := pr.value()
		// fmt.Printf("%v  %v: * Adding up %v: %v\n", buf, wfname, pr, v)
		sum += v
	} else if wf.fallback == "R" {
		// fmt.Printf("%v  %v: - Ignoring %v\n", buf, wfname, pr)
	} else {
		v := applyRange(pr, wfm, wf.fallback, totals, idx)
		sum += v
	}
	return sum
}

func makeWorkflow(s string) *Workflow {
	wf := Workflow{}
	re1 := regexp.MustCompile(`([a-z]+)\{((?:[a-z]+[<>]\d+:[a-zA-Z]+,)+)([a-zA-Z]+)\}`)
	elems := re1.FindAllStringSubmatch(s, -1)
	wf.name = elems[0][1]
	wf.fallback = elems[0][3]
	re2 := regexp.MustCompile(`([a-z]+)([<>])(\d+):([a-zA-Z]+)`)
	elems2 := re2.FindAllStringSubmatch(elems[0][2], -1)
	wf.rules = make([]Rule, len(elems2))
	for i := range elems2 {
		// r := Rule{}
		wf.rules[i].param = elems2[i][1]
		wf.rules[i].op = elems2[i][2]
		wf.rules[i].val = tools.Str2Int(elems2[i][3])
		wf.rules[i].result = elems2[i][4]
		// wf.rules[i] = r
	}
	return &wf
}

func loadWorkflows(lines []string) WorkflowMap {
	workflows := WorkflowMap(make(map[string]*Workflow))
	for _, line := range lines {
		if len(line) == 0 {
			break
		}
		wf := makeWorkflow(line)
		workflows[wf.name] = wf
	}
	return workflows
}

func makePart(s string) *Part {
	re := regexp.MustCompile(`([a-z])=(\d+)`)
	parts := re.FindAllStringSubmatch(s, -1)
	input := Part(make(map[string]int))
	for i := range parts {
		input[parts[i][1]] = tools.Str2Int(parts[i][2])
	}
	return &input
}

func loadParts(lines []string) []*Part {
	allParts := []*Part{}
	for _, line := range lines {
		allParts = append(allParts, makePart(line))
	}
	return allParts
}

func part01() {
	startTime := time.Now()
	total := 0
	lines := getInput()
	var breakLine int
	for len(lines[breakLine]) > 0 {
		breakLine++
	}

	workflows := loadWorkflows(lines)
	parts := loadParts(lines[breakLine+1:])
	fmt.Printf("Read %v workflows, %v parts\n", len(workflows), len(parts))

	for _, in := range parts {
		val := "in"
		for val != "R" && val != "A" {
			wf, ok := workflows[val]
			if ok {
				ret, chk := apply(in, wf)
				if chk {
					val = ret
				}
			}
		}
		if val == "A" {
			inp := *in
			val := inp.value()
			// fmt.Printf("Input %v yields to %v\n", inp, val)
			total += val
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Result part 01 (%v): %v\n\n", elapsed, total)
}

func part02() {
	startTime := time.Now()
	lines := getInput()
	workflows := loadWorkflows(lines)
	fmt.Printf("Read %v workflows\n", len(workflows))

	startRange := PartRange{
		"x": [2]int{1, 4000},
		"m": [2]int{1, 4000},
		"a": [2]int{1, 4000},
		"s": [2]int{1, 4000},
	}

	// test applyRange
	// wf := workflows["px"]
	// prl := applyRange(pr, wf)
	// fmt.Printf("Workflow: %v\n", wf)
	// fmt.Printf("Result: %v\n", prl)

	wfname := "in"
	total := applyRange(startRange, &workflows, wfname, 0, 0)

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
