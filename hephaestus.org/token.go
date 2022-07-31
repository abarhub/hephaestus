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
	STRING_LITERAL

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
	EQUALS2             // ==
	LESSER              // <
	LESSER_OR_EQUALS    // <=
	GREATER             // >
	GREATER_OR_EQUALS   // >=

	// Keywords
	VOID
	INT
	STRING
	BOOLEAN
	TRUE
	FALSE
)
