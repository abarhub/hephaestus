package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Position struct {
	line   int
	column int
	pos    int
}

type ScannerRes struct {
	tok      Token
	lit      string
	position Position
}

// Scanner represents a lexical scanner.
type Scanner struct {
	r              *bufio.Reader
	tab            []ScannerRes
	positionUnread int
	position       Position
	lastposition   *Position
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) newScannerRes(tok Token, lit string) ScannerRes {
	return ScannerRes{tok: tok, lit: lit}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (ScannerRes, error) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	// If we see a digit then consume as a number.
	if isWhitespace(ch) {
		err := s.unread()
		if err != nil {
			return ScannerRes{}, err
		}
		scan, err := s.scanWhitespace()
		return scan, err
	} else if isLetter(ch) {
		err := s.unread()
		if err != nil {
			return ScannerRes{}, err
		}
		scan, err := s.scanIdent()
		return scan, err
	} else if isDigit(ch) {
		err := s.unread()
		if err != nil {
			return ScannerRes{}, err
		}
		scan, err := s.scanNumber()
		return scan, err
	} else if ch == '"' {
		err := s.unread()
		if err != nil {
			return ScannerRes{}, err
		}
		scan, err := s.scanString()
		return scan, err
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return s.newScannerRes(EOF, ""), nil
	case '*':
		return s.newScannerRes(ASTERISK, string(ch)), nil
	case ',':
		return s.newScannerRes(COMMA, string(ch)), nil
	case '(':
		return s.newScannerRes(OPEN_PARENTHESIS, string(ch)), nil
	case ')':
		return s.newScannerRes(CLOSE_PARENTHESIS, string(ch)), nil
	case '{':
		return s.newScannerRes(OPEN_CURLY_BRACKET, string(ch)), nil
	case '}':
		return s.newScannerRes(CLOSE_CURLY_BRACKET, string(ch)), nil
	case '=':
		ch := s.read()
		if ch == '=' {
			return s.newScannerRes(EQUALS2, "=="), nil
		} else {
			err := s.unread()
			return s.newScannerRes(EQUALS, "="), err
		}
	case ';':
		return s.newScannerRes(SEMICOLON, string(ch)), nil
	case '+':
		return s.newScannerRes(ADD, string(ch)), nil
	case '-':
		return s.newScannerRes(SUB, string(ch)), nil
	case '<':
		ch := s.read()
		if ch == '=' {
			return s.newScannerRes(LESSER_OR_EQUALS, "<="), nil
		} else {
			err := s.unread()
			return s.newScannerRes(LESSER, "<"), err
		}
	case '>':
		ch := s.read()
		if ch == '=' {
			return s.newScannerRes(GREATER_OR_EQUALS, ">="), nil
		} else {
			err := s.unread()
			return s.newScannerRes(GREATER, ">"), err
		}
	}

	return s.newScannerRes(ILLEGAL, string(ch)), nil
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (ScannerRes, error) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			err := s.unread()
			if err != nil {
				return ScannerRes{}, err
			}
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return s.newScannerRes(WS, buf.String()), nil
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (ScannerRes, error) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			err := s.unread()
			if err != nil {
				return ScannerRes{}, err
			}
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch buf.String() {
	case "void":
		return s.newScannerRes(VOID, buf.String()), nil
	case "int":
		return s.newScannerRes(INT, buf.String()), nil
	case "string":
		return s.newScannerRes(STRING, buf.String()), nil
	case "boolean":
		return s.newScannerRes(BOOLEAN, buf.String()), nil
	case "true":
		return s.newScannerRes(TRUE, buf.String()), nil
	case "false":
		return s.newScannerRes(FALSE, buf.String()), nil
	}

	// Otherwise return as a regular identifier.
	return s.newScannerRes(IDENT, buf.String()), nil
}

func (s *Scanner) scanNumber() (ScannerRes, error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) {
			err := s.unread()
			if err != nil {
				return ScannerRes{}, err
			}
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return s.newScannerRes(NUMBER, buf.String()), nil
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	position := Position{}
	position.line = s.position.line
	position.column = s.position.column
	position.pos = s.position.pos
	s.lastposition = &position
	s.position.pos = s.position.pos + 1
	if ch == '\n' {
		s.position.line++
		s.position.column = 0
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() error {
	_ = s.r.UnreadRune()
	if s.lastposition == nil {
		return fmt.Errorf("no character before")
	}
	s.position = *s.lastposition
	s.lastposition = nil
	return nil
}

func (s *Scanner) scanString() (ScannerRes, error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '"' {
			_, _ = buf.WriteRune(ch)
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return s.newScannerRes(STRING_LITERAL, buf.String()[1:buf.Len()-1]), nil
}

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return ch >= '0' && ch <= '9' }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
