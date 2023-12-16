/*
 * Stack implementation
 * Nothing to say
 *
 * Helper type "Position" to navigate more easy
 *
 * MIT License, Copyright (c) 2023 Jonas Rathert
 */
package tools

type node[T any] struct {
	value T
	prev  *node[T]
}

type Stack[T any] struct {
	top  *node[T]
	size int
}

func (s *Stack[T]) Push(elem T) {
	n := node[T]{elem, s.top}
	s.top = &n
	s.size++
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.size == 0 {
		var zero T
		return zero, false
	}
	n := s.top
	s.top = n.prev
	s.size--
	return n.value, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if s.size == 0 {
		var zero T
		return zero, false
	}
	n := s.top
	return n.value, true
}

func (s Stack[T]) Size() int {
	return s.size
}
