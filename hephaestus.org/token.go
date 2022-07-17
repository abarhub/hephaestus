package main

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // main
	NUMBER

	// Misc characters
	ASTERISK            // *
	COMMA               // ,
	OPEN_PARENTHESIS    // (
	CLOSE_PARENTHESIS   // )
	OPEN_CURLY_BRACKET  // {
	CLOSE_CURLY_BRACKET // }
	EQUALS              // =
	SEMICOLON           // ;
	ADD                 // +
	SUB                 // -

	// Keywords
	VOID
	INT
	STRING
)
