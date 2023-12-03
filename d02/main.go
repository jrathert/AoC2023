/*
 * Day 02 of AoC 2023
 *
 * Idea: Build a list of "Games" containing the relevant information and selected
 * sets. Then just iterate over games and calculate what is necessary.
 * Major problem was to process input with limited regexp from go.
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

var testinput = `Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
`

func part01() {

	games := buildGamesFromInput()

	maxCubes := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	sumids := 0
	for _, g := range games {
		log.Printf("Examining game %v with %v selects\n", g.number, g.numSelects)
		sumids += g.number
	loop:
		for _, s := range g.selects {
			for _, c := range [3]string{"red", "green", "blue"} {
				if s[c] > maxCubes[c] {
					log.Printf("  Oops: %v is > %v (for color %v) - breaking %v\n", s[c], maxCubes[c], c, g.number)
					sumids -= g.number
					break loop
				}
			}
		}
	}
	fmt.Printf("Total sum of valid game IDs: %v\n", sumids)
}

func part02() {

	games := buildGamesFromInput()

	sumpowers := 0
	for _, g := range games {
		log.Printf("Examining game %v with %v selects\n", g.number, g.numSelects)
		minCubes := map[string]int{
			"red":   0,
			"green": 0,
			"blue":  0,
		}

		for _, s := range g.selects {
			for _, c := range [3]string{"red", "green", "blue"} {
				if s[c] > minCubes[c] {
					log.Printf("  Oops: setting minCubes to %v (for color %v) - game %v\n", s[c], c, g.number)
					minCubes[c] = s[c]
				}
			}
		}
		sumpowers += minCubes["red"] * minCubes["green"] * minCubes["blue"]

	}
	fmt.Printf("Sum of game powers: %v\n", sumpowers)
}

type Game struct {
	number     int
	numSelects int
	selects    []map[string]int
}

func buildGamesFromInput() []Game {

	var input []string
	if TESTMODE {
		input = tools.ReadInputString(testinput)
	} else {
		input = tools.ReadInputFile(*inputfile)
	}

	games := []Game{}

	cnt := 0
	reg := regexp.MustCompile(`[0-9]+`)
	res := regexp.MustCompile(`([0-9]+) (red|green|blue)`)
	for _, line := range input {
		cnt += 1
		var game Game
		elems := strings.Split(line, ":")
		noGame := reg.FindString(elems[0])
		selects := strings.Split(elems[1], ";")

		game.number = tools.Str2Int(noGame)
		game.numSelects = len(selects)
		game.selects = make([]map[string]int, game.numSelects)
		// fmt.Printf("%d: Game %v with %v selects\n", cnt, noGame, len(selects))

		for i := 0; i < game.numSelects; i++ {
			game.selects[i] = make(map[string]int)
			r := res.FindAllStringSubmatch(selects[i], -1)
			// fmt.Printf("%d: %q\n", cnt, r)
			for j := 0; j < len(r); j++ {
				// fmt.Printf("Setting %v to %v\n", r[j][2], r[j][1])
				game.selects[i][r[j][2]] = tools.Str2Int(r[j][1])
			}
			// fmt.Printf("%d: game selects: %v\n", cnt, game.selects[i])
		}
		games = append(games, game)
	}
	return games
}
