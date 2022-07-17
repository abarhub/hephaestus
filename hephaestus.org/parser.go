package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// SelectStatement represents a SQL SELECT statement.
/*type SelectStatement struct {
	Fields    []string
	TableName string
}*/

type Function struct {
	name        string
	instruction []Instruction
}

type Instruction struct {
	variable string
	valeur   int
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func parser() {

	s := "void main () { x=5;y=18;}"
	funct, err := NewParser(strings.NewReader(s)).Parse2()

	if err != nil {
		fmt.Printf("error : %v\n", err)
	} else {
		fmt.Printf("ok %v\n", funct)
		interpreter := NewInterpreter(funct)
		interpreter.interpreter()
	}
}

func (p *Parser) parseInstr(funct *Function) (*Instruction, error) {

	for {

		instr := &Instruction{}

		if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
			return nil, fmt.Errorf("found %q, expected identifier", lit)
		} else {
			instr.variable = lit
		}

		if tok, lit := p.scanIgnoreWhitespace(); tok != EQUALS {
			return nil, fmt.Errorf("found %q, expected =", lit)
		}

		if tok, lit := p.scanIgnoreWhitespace(); tok != NUMBER {
			return nil, fmt.Errorf("found %q, expected number", lit)
		} else {
			intVar, err := strconv.Atoi(lit)
			if err != nil {
				return nil, fmt.Errorf("found %q, expected number", lit)
			} else {
				instr.valeur = intVar
			}
		}

		if tok, lit := p.scanIgnoreWhitespace(); tok != SEMICOLON {
			return nil, fmt.Errorf("found %q, expected ';'", lit)
		}

		funct.instruction = append(funct.instruction, *instr)

		if tok, _ := p.scanIgnoreWhitespace(); tok == CLOSE_CURLY_BRACKET {
			p.unscan()
			break
		} else {
			p.unscan()
		}

	}

	return nil, nil
}

func (p *Parser) Parse2() ([]Function, error) {

	funct := &Function{}

	if tok, lit := p.scanIgnoreWhitespace(); tok != VOID {
		return nil, fmt.Errorf("found %q, expected void", lit)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
		return nil, fmt.Errorf("found %q, expected main", lit)
	} else {
		funct.name = lit
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != OPEN_PARENTHESIS {
		return nil, fmt.Errorf("found %q, expected (", lit)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != CLOSE_PARENTHESIS {
		return nil, fmt.Errorf("found %q, expected )", lit)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != OPEN_CURLY_BRACKET {
		return nil, fmt.Errorf("found %q, expected {", lit)
	}

	_, err := p.parseInstr(funct)
	if err != nil {
		return nil, fmt.Errorf("expected instruction: %s", err)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != CLOSE_CURLY_BRACKET {
		return nil, fmt.Errorf("found %q, expected }", lit)
	}

	return []Function{*funct}, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
