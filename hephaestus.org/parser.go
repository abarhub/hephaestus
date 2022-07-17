package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Function struct {
	name        string
	instruction []Instruction
}

type Instruction struct {
	variable string
	valeur   *Expression
}

type ExprCode int

const (
	EXPR_CODE_INT ExprCode = iota
	EXPR_CODE_VAR
	EXPR_CODE_ADD
	EXPR_CODE_SUB
)

type Expression struct {
	code      ExprCode
	valeurInt int
	variable  string
	left      *Expression
	right     *Expression
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

	s := "void main () { x=5;y=18;z=x;t=x+8;}"
	funct, err := NewParser(strings.NewReader(s)).Parse2()

	if err != nil {
		fmt.Printf("error : %v\n", err)
	} else {
		fmt.Printf("ok %v\n", funct)
		interpreter := NewInterpreter(funct)
		err = interpreter.interpreter()
		if err != nil {
			fmt.Printf("error : %v\n", err)
		}
	}
}

func (p *Parser) parseExpr() (*Expression, error) {
	var expr Expression
	tok, lit := p.scanIgnoreWhitespace()
	if tok == NUMBER {
		intVar, err := strconv.Atoi(lit)
		if err != nil {
			return nil, fmt.Errorf("found %q, expected number", lit)
		} else {
			expr = Expression{code: EXPR_CODE_INT, valeurInt: intVar}
		}
	} else if tok == IDENT {
		expr = Expression{code: EXPR_CODE_VAR, variable: lit}
	} else {
		return nil, fmt.Errorf("found %q, expected number", lit)
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok == ADD {
		expr2, err := p.parseExpr()
		if err != nil {
			return nil, fmt.Errorf("expected expression for add: %s", err)
		} else {
			var expr3 *Expression
			expr3 = new(Expression)
			expr3.code = EXPR_CODE_ADD
			expr3.left = &expr
			expr3.right = expr2
			//expr3 := Expression{code: EXPR_CODE_ADD, left: &expr, right: expr2}
			//expr = expr3
			return expr3, nil
		}
	} else if tok == SUB {
		expr2, err := p.parseExpr()
		if err != nil {
			return nil, fmt.Errorf("expected expression for sub: %s", err)
		} else {
			expr = Expression{code: EXPR_CODE_SUB, left: &expr, right: expr2}
			//expr = expr3
			//if(expr3.right.code==)
			//return &expr3, nil
		}
	} else {
		p.unscan()
	}
	return &expr, nil
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

		expr, err := p.parseExpr()
		if err != nil {
			return nil, fmt.Errorf("invalid expression: %s", err)
		} else {
			instr.valeur = expr
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
	tmp := p.s.Scan()
	tok, lit = tmp.tok, tmp.lit

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
