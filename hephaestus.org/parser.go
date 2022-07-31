package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type TypeCode int

const (
	TYPE_INT TypeCode = iota
	TYPE_VOID
	TYPE_STRING
)

type Type struct {
	code TypeCode
}

type Function struct {
	ReturnType  Type
	Name        string
	Instruction []Instruction
}

type Instruction struct {
	Variable string
	Valeur   *Expression
}

type ExprCode int

const (
	EXPR_CODE_INT ExprCode = iota
	EXPR_CODE_VAR
	EXPR_CODE_ADD
	EXPR_CODE_SUB
	EXPR_CODE_STR
)

type Expression struct {
	code         ExprCode
	valeurInt    int
	variable     string
	valeurString string
	left         *Expression
	right        *Expression
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

	s := "void main () { x=5;y=18;z=x;t=x+8;v=\"abc\";}"
	p := NewParser(strings.NewReader(s))
	funct, err := p.Parse2()

	if err != nil {
		fmt.Printf("error : %v\n", err)
	} else {
		err = p.Checker(funct)
		if err != nil {
			fmt.Printf("error : %v\n", err)
		} else {
			fmt.Printf("ok %v\n", funct)
			interpreter := NewInterpreter(funct)
			_, err = interpreter.interpreter()
			if err != nil {
				fmt.Printf("error : %v\n", err)
			}
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
	} else if tok == STRING_LITERAL {
		expr = Expression{code: EXPR_CODE_STR, valeurString: lit}
	} else {
		return nil, fmt.Errorf("found %q, expected number", lit)
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok == ADD || tok == SUB {
		expr2, err := p.parseExpr()
		if err != nil {
			return nil, fmt.Errorf("expected expression for add: %s", err)
		} else {
			var expr3 *Expression
			expr3 = new(Expression)
			if tok == ADD {
				expr3.code = EXPR_CODE_ADD
			} else if tok == SUB {
				expr3.code = EXPR_CODE_SUB
			} else {
				return nil, fmt.Errorf("expected operator + or -")
			}
			expr3.left = &expr
			expr3.right = expr2
			return expr3, nil
		}
	} else {
		p.unscan()
	}
	return &expr, nil
}

func (p *Parser) parseType() (*Type, error) {
	var res *Type
	if tok, lit := p.scanIgnoreWhitespace(); tok == VOID {
		res = new(Type)
		res.code = TYPE_VOID
		return res, nil
	} else if tok == INT {
		res = new(Type)
		res.code = TYPE_INT
		return res, nil
	} else if tok == STRING {
		res = new(Type)
		res.code = TYPE_STRING
		return res, nil
	} else {
		return nil, fmt.Errorf("found %q, expected type", lit)
	}
}

func (p *Parser) parseInstr(funct *Function) (*Instruction, error) {

	for {

		instr := &Instruction{}

		if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
			return nil, fmt.Errorf("found %q, expected identifier", lit)
		} else {
			instr.Variable = lit
		}

		if tok, lit := p.scanIgnoreWhitespace(); tok != EQUALS {
			return nil, fmt.Errorf("found %q, expected =", lit)
		}

		expr, err := p.parseExpr()
		if err != nil {
			return nil, fmt.Errorf("invalid expression: %s", err)
		} else {
			instr.Valeur = expr
		}

		if tok, lit := p.scanIgnoreWhitespace(); tok != SEMICOLON {
			return nil, fmt.Errorf("found %q, expected ';'", lit)
		}

		funct.Instruction = append(funct.Instruction, *instr)

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

	typeReturn, err := p.parseType()
	if err != nil {
		return nil, err
	}
	funct.ReturnType = *typeReturn

	if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
		return nil, fmt.Errorf("found %q, expected main", lit)
	} else {
		funct.Name = lit
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

	_, err = p.parseInstr(funct)
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
