package main

import (
	"bufio"
	"bytes"
	"io"
)

type ScannerRes struct {
	tok Token
	lit string
}

// Scanner represents a lexical scanner.
type Scanner struct {
	r              *bufio.Reader
	tab            []ScannerRes
	positionUnread int
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func newScannerRes(tok Token, lit string) ScannerRes {
	return ScannerRes{tok: tok, lit: lit}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() ScannerRes {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	// If we see a digit then consume as a number.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return newScannerRes(EOF, "")
	case '*':
		return newScannerRes(ASTERISK, string(ch))
	case ',':
		return newScannerRes(COMMA, string(ch))
	case '(':
		return newScannerRes(OPEN_PARENTHESIS, string(ch))
	case ')':
		return newScannerRes(CLOSE_PARENTHESIS, string(ch))
	case '{':
		return newScannerRes(OPEN_CURLY_BRACKET, string(ch))
	case '}':
		return newScannerRes(CLOSE_CURLY_BRACKET, string(ch))
	case '=':
		return newScannerRes(EQUALS, string(ch))
	case ';':
		return newScannerRes(SEMICOLON, string(ch))
	case '+':
		return newScannerRes(ADD, string(ch))
	case '-':
		return newScannerRes(SUB, string(ch))
	}

	return newScannerRes(ILLEGAL, string(ch))
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() ScannerRes {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return newScannerRes(WS, buf.String())
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() ScannerRes {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch buf.String() {
	case "void":
		return newScannerRes(VOID, buf.String())
	case "int":
		return newScannerRes(INT, buf.String())
	case "string":
		return newScannerRes(STRING, buf.String())
	}

	// Otherwise return as a regular identifier.
	return newScannerRes(IDENT, buf.String())
}

func (s *Scanner) scanNumber() ScannerRes {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return newScannerRes(NUMBER, buf.String())
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return ch >= '0' && ch <= '9' }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
