package asm

import (
	"bytes"
	"fmt"
)

func newParser() *Parser {
	return &Parser{
		scanner: Scanner{},
		asm:     &Asm{},
	}
}

func (p Parser) parse(src *bytes.Buffer) error {
	var peek []rune // unconsumed runes
	var err error
	p.scanner.src = src
	p.scanner.scan()
	p.scanner.skipComments()
	if err := p.expectSectionHeader(); err != nil {
		return err
	}
	switch p.section {
	case Text:
		if peek, err = p.parseText(); err != nil {
			return err
		}
	case Data:
	case Bss:
	}
	if p.compare(peek, "section") {
		// parse next section
	}
	return nil
}

func (p Parser) error(e string) error {
	return fmt.Errorf("%d:%d %s", p.scanner.y, p.scanner.x, e)
}

func (p *Parser) optional(r rune) bool {
	if p.scanner.ch == r {
		return true
	}
	return false
}

func (p *Parser) runes() []rune {
	var v []rune
	p.scanner.skipWhitespace()
	for ; p.scanner.ch != ' ' && p.scanner.ch != '\t' && p.scanner.ch != ',' && p.scanner.ch != '\n' && p.scanner.ch != -1; p.scanner.scan() {
		v = append(v, p.scanner.ch)
	}
	return v
}

func (p *Parser) expect(v rune) error {
	if p.scanner.ch != v {
		return p.error(fmt.Sprintf("expected %v got %v", v, p.scanner.ch))
	}
	p.scanner.scan()
	return nil
}

// section .text \n
// section .data \n
// section .bss \n
func (p *Parser) expectSectionHeader() error {
	p.scanner.skip()
	var v []rune
	v = p.runes()
	switch true {
	case p.compare(v, "section"):
		v = p.runes()
		switch true {
		case p.compare(v, ".text"):
			p.section = Text
		case p.compare(v, ".data"):
			p.section = Data
		case p.compare(v, ".bss"):
			p.section = Bss
		default:
			return p.error("expected .text .data or .bss after section")
		}
	default:
		return p.error("expected section")
	}
	return p.expect('\n')
}

func (p *Parser) compare(a []rune, s string) bool {
	e := []rune(s)

	if len(a) != len(e) {
		return false
	}
	for i, v := range a {
		if e[i] != v {
			return false
		}
	}
	return true
}

// global _start
func (p *Parser) parseText() ([]rune, error) {
	p.asm.TextSection()

	p.scanner.skipComments()
	// optional
	v := p.runes()
	if p.compare(v, "global") {
		v = p.runes()
		if p.compare(v, "_start") {
			p.asm.t.start = "_start"
		}
	} else {
		return nil, p.error(fmt.Sprintf("unhandled %s", string(v)))
	}
	p.scanner.skipWhitespace()
	if err := p.expect('\n'); err != nil {
		return nil, err
	}
	p.scanner.skipComments()
	return p.text()
}

// [._]{0,1}<string>:\n
// op <string>\n
// op <string>, <string>\n
// op <string>, <string>, <string>\n
func (p *Parser) text() ([]rune, error) {
	// keep track of line numbers, where an op counts as a line:
	var line int
	text := &p.asm.t

	for {
		p.scanner.skipComments()
		v := p.runes()
		switch true {
		case len(v) == 0:
			return nil, nil // nothing to do

		case p.compare(v, "section"):
			return v, nil

		case v[len(v)-1] == ':':
			if v[0] == '.' {
				// a local label e.g. .LBL_01:
				text.labels[string(v[0:len(v)-1])] = line
			} else if v[0] == '_' {
				if p.compare(v, "_start:") {
					// a global label, _start currently supported
					text.labels[string(v[0:len(v)-1])] = line
				} else {
					return nil, p.error(fmt.Sprintf("global label support for %s not implemented", string(v)))
				}
			} else {
				// function label e.g. fib:
				text.functions[string(v[0:len(v)-1])] = line
			}
			p.scanner.skipWhitespace()
			if err := p.expect('\n'); err != nil {
				return nil, err
			}
			continue
		default:
			// should be op
			var inst inst
			var comma bool
			var nl bool
			i := 0
			// op label
			// op $1, %r1
			// op %r1, %r2, %r3
			for ; ; v = p.runes() {
				if i >= 4 {
					return nil, p.error("too many parts")
				}
				inst[i] = string(v)
				p.scanner.skipWhitespace()

				if i > 0 {
					// expect comma or newline
					comma = p.optional(',')
					nl = p.optional('\n')
					if !(comma || nl) {
						return nil, p.error("expected comma or newline")
					}
					p.scanner.scan()
					if nl {
						line++
						break
					}
				}
				i++
			}
			p.asm.t.insts = append(p.asm.t.insts, inst)
			continue
		}
	}
}
