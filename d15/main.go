/*
 * Day 15 of AoC 2023
 *
 * Idea: This was easy with just building some arrays/slices.
 * Used a somewhat inefficient approach when deleting lenses from boxes
 * (basically creating a new slice by copying data over) - but it works
 * ok with the given data. Otherwise, a linked list would be better, maybe
 * will do it later...
 *
 * Learned a lot about Go's handling of structs - these are value types, actually,
 * so everything worked only after making boxes and boxes.lenses arrays
 * of pointers!
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

var testinput = `rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7`

type lens struct {
	name  string
	value int
}

type box struct {
	name   string
	lenses []*lens
}

func makeBox(name string) box {
	b := box{name, make([]*lens, 0)}
	return b
}

func (b box) String() string {
	var buf strings.Builder
	buf.WriteString(b.name)
	buf.WriteString(": ")
	if len(b.lenses) == 0 {
		buf.WriteString("(no lenses)")
	} else {
		for i, l := range b.lenses {
			buf.WriteString(fmt.Sprintf("[%v %v]", l.name, l.value))
			if i < len(b.lenses)-1 {
				buf.WriteString(" ")
			}
		}
	}
	return buf.String()
}

func (b *box) addLens(name string, val int) int {
	if b.lenses == nil {
		b.lenses = make([]*lens, 0)
	}
	for i := 0; i < len(b.lenses); i++ {
		if b.lenses[i].name == name {
			b.lenses[i].value = val
			return i
		}
	}
	le := lens{name, val}
	b.lenses = append(b.lenses, &le)
	return len(b.lenses) - 1
}

func (b *box) removeLens(name string) int {
	for i := 0; i < len(b.lenses); i++ {
		if b.lenses[i].name == name {
			b.lenses = append(b.lenses[:i], b.lenses[i+1:]...)
			return i
		}
	}
	return -1
}

func (b box) valuate(num int) int {
	total := 0
	for i := 0; i < len(b.lenses); i++ {
		val := (num + 1) * (i + 1) * b.lenses[i].value
		total += val
	}
	return total
}

func calcHash(token string) int {
	val := 0
	for _, c := range token {
		v := int(c)
		val += v
		val *= 17
		val %= 256
	}
	return val
}

func part01() {
	startTime := time.Now()

	lines := getInput()

	total := 0
	parts := strings.Split(lines[0], ",")
	for _, p := range parts {
		val := calcHash(p)
		total += val
	}
	elapsed := time.Since(startTime)
	fmt.Printf("Result part 01 (%v): %v\n\n", elapsed, total)
}

func splitToken(token string) (string, int) {
	if token[len(token)-1] == '-' {
		return token[0 : len(token)-1], -1
	} else {
		parts := strings.Split(token, "=")
		return parts[0], tools.Str2Int(parts[1])
	}
}

func part02() {
	startTime := time.Now()

	lines := getInput()

	boxes := make([]*box, 256)
	for i := 0; i < 256; i++ {
		b := makeBox(fmt.Sprintf("Box %v", i))
		boxes[i] = &b
	}
	parts := strings.Split(lines[0], ",")
	for _, p := range parts {
		name, val := splitToken(p)
		hash := calcHash(name)
		box := boxes[hash]
		if val == -1 {
			// fmt.Printf("Removing lens %v from box %v\n", name, hash)
			box.removeLens(name)
			// fmt.Printf("Result: %v - box is now: %v\n", chk, box)
		} else {
			// fmt.Printf("Adding lens %v (%v) to box %v\n", name, val, hash)
			box.addLens(name, val)
			// fmt.Printf("Result: %v - box is now: %v\n", chk, box)
		}
	}

	total := 0
	for i, b := range boxes {
		total += b.valuate(i)
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
