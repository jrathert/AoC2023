/*
 * Matrix implementation
 * Allows to store a byte field
 *
 * Helper type "Position" to navigate more easy
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package tools

import (
	"bytes"
	"strings"
)

type Position [2]int

func (pos Position) IsBelowOf(other Position) bool {
	return pos[1] == other[1]+1
}

func (pos Position) IsAboveOf(other Position) bool {
	return pos[1] == other[1]-1
}

func (pos Position) IsLeftOf(other Position) bool {
	return pos[0] == other[0]-1
}

func (pos Position) IsRightOf(other Position) bool {
	return pos[0] == other[0]+1
}

type Matrix struct {
	fields [][]byte
	rows   int
	cols   int
}

func NewMatrix(rows, cols int) Matrix {
	m := Matrix{}
	m.rows = rows
	m.cols = cols
	m.fields = make([][]byte, m.rows)
	for i := range m.fields {
		m.fields[i] = make([]byte, m.cols)
	}
	return m
}

func (m Matrix) Copy() Matrix {
	ret := Matrix{}
	ret.fields = make([][]byte, m.rows)
	for i := range ret.fields {
		ret.fields[i] = make([]byte, m.cols)
		copy(ret.fields[i], m.fields[i])
	}
	return ret
}

func (m Matrix) String() string {
	var b strings.Builder
	b.Grow(m.rows * m.cols)
	for i := 0; i < m.rows; i++ {
		b.Write(m.fields[i])
		b.WriteString("\n")
	}
	return b.String()
}

func (m Matrix) NonZeroString() string {
	var b strings.Builder
	b.Grow(m.rows * m.cols)
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Cols(); j++ {
			v, _ := m.Value(j, i)
			if v == 0 {
				b.WriteString(".")
			} else {
				b.WriteString("#")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m *Matrix) Reset() {
	m.fields = nil
	m.rows = 0
	m.cols = 0
}

func (m *Matrix) AddLine(s string) {
	if m.fields == nil {
		m.fields = make([][]byte, 0)
	}
	m.fields = append(m.fields, []byte(s))
	m.rows++
	m.cols = len(s)
}

func (m Matrix) Cols() int {
	return m.cols
}

func (m Matrix) Rows() int {
	return m.rows
}

func (m Matrix) SumValues() int {
	val := 0
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			val += int(m.fields[i][j])
		}
	}
	return val
}

func (m Matrix) CountNonZero() int {
	val := 0
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			if m.fields[i][j] != 0 {
				val++
			}
		}
	}
	return val
}

func (m *Matrix) FindField(c byte) (Position, bool) {
	for y := 0; y < m.rows; y++ {
		x := bytes.IndexByte(m.fields[y], c)
		if x >= 0 {
			pos := Position{x, y}
			return pos, true
		}
	}
	return Position{-1, -1}, false
}

func (m Matrix) Value(x, y int) (byte, bool) {
	if x < 0 || x >= m.cols || y < 0 || y >= m.rows {
		return '?', false
	}
	return m.fields[y][x], true
}

func (m Matrix) ValueAtPos(pos Position) (byte, bool) {
	return m.Value(pos[0], pos[1])
}

func (m *Matrix) SetValue(x, y int, c byte) bool {
	if x < 0 || x >= m.cols || y < 0 || y >= m.rows {
		return false
	}

	m.fields[y][x] = c
	return true
}

func (m *Matrix) SetValueAtPos(pos Position, c byte) bool {
	return m.SetValue(pos[0], pos[1], c)
}

func (m Matrix) LeftOf(pos Position) (Position, bool) {
	if pos[0] > 0 {
		return Position{pos[0] - 1, pos[1]}, true
	} else {
		return pos, false
	}
}

func (m Matrix) RightOf(pos Position) (Position, bool) {
	if pos[0] < m.cols-1 {
		return Position{pos[0] + 1, pos[1]}, true
	} else {
		return pos, false
	}
}

func (m Matrix) AboveOf(pos Position) (Position, bool) {
	if pos[1] > 0 {
		return Position{pos[0], pos[1] - 1}, true
	} else {
		return pos, false
	}
}

func (m Matrix) BelowOf(pos Position) (Position, bool) {
	if pos[1] < m.rows-1 {
		return Position{pos[0], pos[1] + 1}, true
	} else {
		return pos, false
	}
}
