package asm

import "bytes"

type (
	Scanner struct {
		src *bytes.Buffer
		ch  rune
		x   int
		y   int
	}
)

func (s *Scanner) scan() {
	var err error
	var i int
	s.ch, i, err = s.src.ReadRune()
	s.x++
	if s.ch == '\n' {
		s.x++
		s.y = 0
	}
	if err != nil || i == 0 {
		s.ch = -1
	}
}

func (s *Scanner) readAll() {
	for s.ch != -1 {
		s.scan()
	}
}

func (s *Scanner) skip() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' {
		s.scan()
	}
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' {
		s.scan()
	}
	if s.ch == '#' {
		s.scan()
		for s.ch != '\n' {
			s.scan()
		}
	}
}

func (s *Scanner) skipComments() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' {
		s.scan()
	}
	if s.ch == '#' {
		s.scan()
		for s.ch != '\n' {
			s.scan()
		}
		s.scan()
		s.skip()
	}
}
