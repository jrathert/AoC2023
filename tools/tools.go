/*
 * Tools
 *
 * Several helper functions
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package tools

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// sum integer values
func SumInts(vals []int) int {
	sum := 0
	for _, v := range vals {
		sum += v
	}
	return sum
}

// Reverse a string
func ReverseStr(s string) string {
	byte_str := []rune(s)
	for i, j := 0, len(byte_str)-1; i < j; i, j = i+1, j-1 {
		byte_str[i], byte_str[j] = byte_str[j], byte_str[i]
	}
	return string(byte_str)
}

// Convert known integer string to int
func Str2Int(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

// Read input - return array of strings
// If parameter s is empty or equal to "input.txt"
// the file "input.txt" is read. Else read string s
func ReadInput(s string) []string {
	if len(s) == 0 || s == "input.txt" {
		s = "input.txt"
		return ReadInputFile(s)
	} else {
		return ReadInputString(s)
	}
}

// Read input from file - return array of strings
func ReadInputFile(fname string) []string {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	input := bufio.NewScanner(file)
	return scan(input)
}

// Read input from multiline string - return array of strings
func ReadInputString(s string) []string {
	input := bufio.NewScanner(strings.NewReader(s))
	return scan(input)
}

// Read all ints in a string and return their values as []int
func ReadInts(s string) []int {
	re := regexp.MustCompile(`[0-9]+`)
	elems := re.FindAllString(s, -1)
	values := make([]int, len(elems))
	for i, val := range elems {
		values[i] = Str2Int(val)
	}
	return values
}

// Read all ints in a string and return their values as []int
func ReadSignedInts(s string) []int {
	re := regexp.MustCompile(`[\+\-0-9]+`)
	elems := re.FindAllString(s, -1)
	values := make([]int, len(elems))
	for i, val := range elems {
		values[i] = Str2Int(val)
	}
	return values
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

// internally used
func scan(scanner *bufio.Scanner) []string {
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
