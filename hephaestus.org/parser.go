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
	TYPE_BOOLEAN
)

type InstructionCode int

const (
	INSTRUCTION_AFFECTATION InstructionCode = iota
	INSTRUCTION_CALL
)

type Type struct {
	code     TypeCode
	position *Position
}

type Function struct {
	ReturnType  Type
	Name        string
	Instruction []Instruction
	position    *Position
}

type Instruction struct {
	Code         InstructionCode
	FunctionName string
	Variable     string
	Valeur       *Expression
	Parameter    []Expression
	position     *Position
}

type ExprCode int

const (
	EXPR_CODE_INT ExprCode = iota
	EXPR_CODE_VAR
	EXPR_CODE_ADD
	EXPR_CODE_SUB
	EXPR_CODE_STR
	EXPR_CODE_LT
	EXPR_CODE_LTE
	EXPR_CODE_GT
	EXPR_CODE_GTE
	EXPR_CODE_EQU
	EXPR_CODE_TRUE
	EXPR_CODE_FALSE
)

type Expression struct {
	code         ExprCode
	valeurInt    int
	variable     string
	valeurString string
	left         *Expression
	right        *Expression
	position     *Position
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

var binaryOperation = map[Token]ExprCode{ADD: EXPR_CODE_ADD,
	SUB: EXPR_CODE_SUB, EQUALS2: EXPR_CODE_EQU, LESSER: EXPR_CODE_LT, LESSER_OR_EQUALS: EXPR_CODE_LTE,
	GREATER: EXPR_CODE_GT, GREATER_OR_EQUALS: EXPR_CODE_GTE}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func parser() {

	s := "void main () { x=5;y=18;z=x;t=x+8;v=\"abc\";print(x,y,z,t,v);}"
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
	tok, lit, pos, err := p.scanIgnoreWhitespace()
	if err != nil {
		return nil, err
	} else if tok == NUMBER {
		intVar, err := strconv.Atoi(lit)
		if err != nil {
			return nil, fmt.Errorf("invalide number %q (pos=%v)", lit, pos)
		} else {
			expr = Expression{code: EXPR_CODE_INT, valeurInt: intVar, position: pos}
		}
	} else if tok == IDENT {
		expr = Expression{code: EXPR_CODE_VAR, variable: lit, position: pos}
	} else if tok == STRING_LITERAL {
		expr = Expression{code: EXPR_CODE_STR, valeurString: lit, position: pos}
	} else if tok == TRUE {
		expr = Expression{code: EXPR_CODE_TRUE, position: pos}
	} else if tok == FALSE {
		expr = Expression{code: EXPR_CODE_FALSE, position: pos}
	} else {
		return nil, fmt.Errorf("found %q, expected number or ident or string (pos=%v)", lit, pos)
	}
	tok, lit, pos, err = p.scanIgnoreWhitespace()
	if err != nil {
		return nil, err
	} else if val, ok := binaryOperation[tok]; ok {
		expr2, err := p.parseExpr()
		if err != nil {
			return nil, fmt.Errorf("expected expression for add: %s (pos=%v)", err, pos)
		} else {
			var expr3 *Expression
			expr3 = new(Expression)
			expr3.code = val
			expr3.left = &expr
			expr3.right = expr2
			expr3.position = pos
			return expr3, nil
		}
	} else {
		p.unscan()
	}
	return &expr, nil
}

func (p *Parser) parseType() (*Type, error) {
	var res *Type
	if tok, lit, pos, err := p.scanIgnoreWhitespace(); err != nil {
		return nil, err
	} else if tok == VOID {
		res = new(Type)
		res.code = TYPE_VOID
		res.position = pos
		return res, nil
	} else if tok == INT {
		res = new(Type)
		res.code = TYPE_INT
		res.position = pos
		return res, nil
	} else if tok == STRING {
		res = new(Type)
		res.code = TYPE_STRING
		res.position = pos
		return res, nil
	} else if tok == BOOLEAN {
		res = new(Type)
		res.code = TYPE_BOOLEAN
		res.position = pos
		return res, nil
	} else {
		return nil, fmt.Errorf("found %q, expected type (pos=%v)", lit, pos)
	}
}

func (p *Parser) parseInstr(funct *Function) (*Instruction, error) {

	for {

		instr := &Instruction{}
		var name = ""
		var posStart *Position

		if tok, lit, pos, err := p.scanIgnoreWhitespace(); err != nil {
			return nil, err
		} else if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected identifier (pos=%v)", lit, pos)
		} else {
			name = lit
			posStart = pos
		}

		if tok, lit, pos, err := p.scanIgnoreWhitespace(); err != nil {
			return nil, err
		} else if tok == EQUALS {
			expr, err := p.parseExpr()
			if err != nil {
				return nil, fmt.Errorf("invalid expression: %s", err)
			} else {
				instr.Code = INSTRUCTION_AFFECTATION
				instr.Valeur = expr
				instr.Variable = name
				instr.position = posStart
			}
		} else if tok == OPEN_PARENTHESIS {
			instr.Code = INSTRUCTION_CALL
			instr.FunctionName = name
			instr.position = posStart
			var param []Expression
			end := false
			for !end {
				expr, err := p.parseExpr()
				if err != nil {
					return nil, fmt.Errorf("invalid expression: %s", err)
				} else {
					param = append(param, *expr)
					if tok, _, pos, err := p.scanIgnoreWhitespace(); err != nil {
						return nil, err
					} else if tok == COMMA {
						// on continue
					} else if tok == CLOSE_PARENTHESIS {
						end = true
					} else {
						return nil, fmt.Errorf("invalid call (pos=%v)", pos)
					}
				}
			}
			instr.Parameter = param
		} else {
			return nil, fmt.Errorf("found %q, expected = (pos=%v)", lit, pos)
		}

		if tok, lit, pos, err := p.scanIgnoreWhitespace(); err != nil {
			return nil, err
		} else if tok != SEMICOLON {
			return nil, fmt.Errorf("found %q, expected ';' (pos=%v)", lit, pos)
		}

		funct.Instruction = append(funct.Instruction, *instr)

		if tok, _, _, err := p.scanIgnoreWhitespace(); err != nil {
			return nil, err
		} else if tok == CLOSE_CURLY_BRACKET {
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

	if tok, lit, pos, err := p.scanIgnoreWhitespace(); err != nil {
		return []Function{}, err
	} else if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected main (pos=%v)", lit, pos)
	} else {
		funct.Name = lit
	}

	if tok, lit, pos, err := p.scanIgnoreWhitespace(); err != nil {
		return []Function{}, err
	} else if tok != OPEN_PARENTHESIS {
		return nil, fmt.Errorf("found %q, expected ( (pos=%v)", lit, pos)
	}

	if tok, lit, pos, err := p.scanIgnoreWhitespace(); err != nil {
		return []Function{}, err
	} else if tok != CLOSE_PARENTHESIS {
		return nil, fmt.Errorf("found %q, expected )(pos=%v)", lit, pos)
	}

	if tok, lit, pos, err := p.scanIgnoreWhitespace(); err != nil {
		return []Function{}, err
	} else if tok != OPEN_CURLY_BRACKET {
		return nil, fmt.Errorf("found %q, expected { (pos=%v)", lit, pos)
	}

	_, err = p.parseInstr(funct)
	if err != nil {
		return nil, fmt.Errorf("expected instruction: %s", err)
	}

	if tok, lit, pos, err := p.scanIgnoreWhitespace(); err != nil {
		return []Function{}, err
	} else if tok != CLOSE_CURLY_BRACKET {
		return nil, fmt.Errorf("found %q, expected } (pos=%v)", lit, pos)
	}

	return []Function{*funct}, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string, pos *Position, error error) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit, nil, nil
	}

	// Otherwise read the next token from the scanner.
	tmp, err := p.s.Scan()
	if err != nil {
		error = err
	}
	tok, lit, pos = tmp.tok, tmp.lit, &tmp.position

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string, pos *Position, err error) {
	tok, lit, pos, err = p.scan()
	if tok == WS {
		tok, lit, pos, err = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
