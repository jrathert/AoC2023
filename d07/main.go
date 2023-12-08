/*
 * Day 07 of AoC 2023
 *
 * Idea: CamelCards -
 *  Build a map of (string) -> number of cards
 *    Add to it while reading
 *  Reduce (take out entries < 2)
 *  Calc rank:
 *   - len 1 w/5 (1), w/4 (2), w/3 (4), w/2 (6)
 *   - len 2 w/2+3 (3), w/2+2 (5)
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package main

import (
	"aoc23/tools"
	"flag"
	"fmt"
	"log"
	"sort"
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

var testinput = `32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483`

type Hand struct {
	raw    string
	values map[string]int
	rank   int
	bet    int
}

type Game []Hand

var cardranks1 = map[byte]int{
	'A': 13,
	'K': 12,
	'Q': 11,
	'J': 10,
	'T': 9,
	'9': 8,
	'8': 7,
	'7': 6,
	'6': 5,
	'5': 4,
	'4': 3,
	'3': 2,
	'2': 1,
}

var cardranks2 = map[byte]int{
	'A': 13,
	'K': 12,
	'Q': 11,
	'T': 10,
	'9': 9,
	'8': 8,
	'7': 7,
	'6': 6,
	'5': 5,
	'4': 4,
	'3': 3,
	'2': 2,
	'J': 1,
}

func rank(hand *Hand, part int) int {
	switch part {
	case 1:
		return rank1(hand)
	case 2:
		return rank2(hand)
	default:
		return -1
	}
}

func rank1(hand *Hand) int {
	if hand.rank > 0 {
		return hand.rank
	}
	maxcnt := 0
	for _, v := range hand.values {
		if v > maxcnt {
			maxcnt = v
		}
	}
	switch maxcnt {
	case 5:
		hand.rank = 1
	case 4:
		hand.rank = 2
	case 3:
		hand.rank = 5 - len(hand.values) // works as hand.values is reduced to elems wiwith cnt > 1
	case 2:
		hand.rank = 7 - len(hand.values) // works as hand.values is reduced to elems wiwith cnt > 1
	case 0, 1:
		hand.rank = 7

	}
	return hand.rank
}

func rank2(hand *Hand) int {
	if hand.rank > 0 {
		return hand.rank
	}
	maxcnt := 0
	cntJ := hand.values["J"]
	delete(hand.values, "J")
	vals := []int{0, 0, 0, 0, 0, 0}
	for k, v := range hand.values {
		if k != "J" {
			vals[v]++
			if v > maxcnt {
				maxcnt = v
			}
		}
	}
	switch maxcnt {
	case 5:
		hand.rank = 1
	case 4:
		hand.rank = 2 - cntJ // can be 4 or 4 + J
	case 3:
		if cntJ > 0 {
			hand.rank = 3 - cntJ // either 1 or 2 J -> rank 2 or 1
		} else {
			hand.rank = 5 - len(hand.values) // either 2 or 1 -> rank 3 (FH) or 4 (triple)
		}
	case 2:
		if cntJ > 1 {
			// 3 J -> 1, 2 J -> 2
			hand.rank = 4 - cntJ
		} else if cntJ == 1 {
			// 1 J -> either 3 (FH) or 4 (triple)
			if len(hand.values) == 2 {
				// must both be 2, as others were eliminated
				// with J -> FH
				hand.rank = 3
			} else {
				hand.rank = 4
			}
		} else {
			hand.rank = 7 - len(hand.values)
		}
	case 1:
		fmt.Println("ERR Should not happen as these entries were eliminated")
	case 0: // only single values, all eliminated
		switch cntJ {
		case 4, 5:
			hand.rank = 1 // full
		case 3:
			hand.rank = 2 // 4
		case 2:
			hand.rank = 4 // triple
		case 1:
			hand.rank = 6
		case 0:
			hand.rank = 7
		}
	}
	return hand.rank
}

// should one be before two?
func cmpHands(one Hand, two Hand, part int) bool {
	if rank(&one, part) > rank(&two, part) {
		return true
	} else if rank(&one, part) == rank(&two, part) {
		val := cmp(one.raw, two.raw, part)
		if val > 0 {
			return true
		}
	}
	return false
}

// neg if one is higher than two
func cmp(one string, two string, part int) int {
	cr := getCardranks(part)
	for i := 0; i < len(one); i++ {
		v := cr[one[i]]
		w := cr[two[i]]
		if v != w {
			return w - v
		}
	}
	return 0
}

func getCardranks(part int) map[byte]int {
	switch part {
	case 1:
		return cardranks1
	case 2:
		return cardranks2
	default:
		return nil
	}
}

func part01() {
	// cnt := 0
	total := 0
	lines := getInput()
	allhands := []Hand{}

	for _, line := range lines {
		// log.Printf("Processing line %v with len %v\n", cnt, len(line))
		hand := Hand{}
		parts := strings.Split(line, " ")
		hand.bet = tools.Str2Int(parts[1])
		hand.values = map[string]int{}
		hand.raw = parts[0]
		for _, s := range parts[0] {
			hand.values[string(s)]++
		}
		for k := range hand.values {
			if hand.values[k] < 2 {
				delete(hand.values, k)
			}
		}
		rank(&hand, 1)
		allhands = append(allhands, hand)
		// fmt.Printf("%v: Hand: %v - rank: %v - bet: %v\n", cnt, hand.raw, rank(&hand), hand.bet)

	}
	sort.Slice(allhands, func(i, j int) bool { return cmpHands(allhands[i], allhands[j], 1) })

	for i := range allhands {
		h := allhands[i]
		// fmt.Printf("%v: Hand: %v - rank: %v - bet: %v\n", i, h.raw, rank(&h), h.bet)
		total += (i + 1) * h.bet
	}
	fmt.Printf("Result part 01: %v\n", total)
}

func part02() {
	total := 0
	lines := getInput()
	allhands := []Hand{}

	for _, line := range lines {
		// log.Printf("Processing line %v with len %v\n", cnt, len(line))
		hand := Hand{}
		parts := strings.Split(line, " ")
		hand.bet = tools.Str2Int(parts[1])
		hand.values = map[string]int{}
		hand.raw = parts[0]
		for _, s := range parts[0] {
			hand.values[string(s)]++
		}
		for k := range hand.values {
			if k != "J" && hand.values[k] < 2 {
				delete(hand.values, k)
			}
		}
		rank2(&hand)
		allhands = append(allhands, hand)
		// fmt.Printf("%v: Hand: %v - rank: %v - bet: %v\n", cnt, hand.raw, rank(&hand), hand.bet)

	}
	sort.Slice(allhands, func(i, j int) bool { return cmpHands(allhands[i], allhands[j], 2) })

	for i := range allhands {
		h := allhands[i]
		fmt.Printf("%v: Hand: %v - rank: %v - bet: %v\n", i, h.raw, rank2(&h), h.bet)
		total += (i + 1) * h.bet
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
