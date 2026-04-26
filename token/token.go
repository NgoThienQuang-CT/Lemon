// Package token defines constants representing the lexical token of the Lemon
// programing language
package token

import "fmt"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Identifier and basic type literals
	IDENT  // foo
	INT    // 69
	FLOAT  // 69.420
	STRING // "onion"

	// Operators and delimiters
	ASSIGN  // =
	PLUS    // +
	MINUS   // -
	STAR    // *
	SLASH   // /
	PERCENT // %

	NOT // !
	LSS // <
	GTR // >
	EQL // ==
	NEQ // !=
	LEQ // <=
	GEQ // >=
	AND // &&
	OR  // ||

	COMMA     // ,
	DOT       // .
	COLON     // :
	SEMICOLON // ;

	LPAREN   // (
	RPAREN   // )
	LCURLY   // {
	RCURLY   // }
	LBRACKET // [
	RBRACKET // ]

	// Reserved keyword
	FN
	LET
	TRUE
	FALSE
	IF
	ELSE
	RETURN
	FOR
	BREAK
	CONTINUE
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

	ASSIGN:  "ASSIGN",
	PLUS:    "PLUS",
	MINUS:   "MINUS",
	STAR:    "STAR",
	SLASH:   "SLASH",
	PERCENT: "PERCENT",

	NOT: "NOT",
	LSS: "LSS",
	GTR: "GTR",

	EQL: "EQL",
	NEQ: "NEQ",
	LEQ: "LEQ",
	GEQ: "GEQ",
	AND: "AND",
	OR:  "OR",

	COMMA:     "COMMA",
	DOT:       "DOT",
	COLON:     "COLON",
	SEMICOLON: "SEMICOLON",

	LPAREN:   "LPAREN",
	RPAREN:   "RPAREN",
	LCURLY:   "LCURLY",
	RCURLY:   "RCURLY",
	LBRACKET: "LBRACKET",
	RBRACKET: "RBRACKET",

	FN:       "FN",
	LET:      "LET",
	TRUE:     "TRUE",
	FALSE:    "FALSE",
	IF:       "IF",
	ELSE:     "ELSE",
	RETURN:   "RETURN",
	FOR:      "FOR",
	BREAK:    "BREAK",
	CONTINUE: "CONTINUE",
}

func (ttype TokenType) String() string {
	return tokens[ttype]
}

var keywords = map[string]TokenType{
	"fn":       FN,
	"let":      LET,
	"true":     TRUE,
	"fasle":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
}

func CheckKeyword(ident string) TokenType {
	if ttype, isKeyword := keywords[ident]; isKeyword {
		return ttype
	}
	return IDENT
}

type Token struct {
	Type    TokenType
	Line    int
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("{token.%s, %d, \"%s\"},", t.Type.String(), t.Line, t.Literal)
}
