// Package token defines constants representing the lexical token of the Lemon
// programing language
package token

import (
	"fmt"
	"slices"
)

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

	// Reserved keywords
	FN
	LET
	MUT
	TRUE
	FALSE
	IF
	ELSE
	RETURN
	WHILE
	BREAK
	CONTINUE
)

var tokens = [...]string{
	ILLEGAL: "illegal",
	EOF:     "EOF",

	IDENT:  "identifier",
	INT:    "integer",
	FLOAT:  "float",
	STRING: "string",

	ASSIGN:  "=",
	PLUS:    "+",
	MINUS:   "-",
	STAR:    "*",
	SLASH:   "/",
	PERCENT: "%",

	NOT: "!",
	LSS: "<",
	GTR: ">",

	EQL: "==",
	NEQ: "!=",
	LEQ: "<=",
	GEQ: ">=",
	AND: "&&",
	OR:  "||",

	COMMA:     ",",
	DOT:       ".",
	COLON:     ":",
	SEMICOLON: ";",

	LPAREN:   "(",
	RPAREN:   ")",
	LCURLY:   "{",
	RCURLY:   "}",
	LBRACKET: "[",
	RBRACKET: "]",

	FN:       "fn",
	LET:      "let",
	MUT:      "mut",
	TRUE:     "true",
	FALSE:    "false",
	IF:       "if",
	ELSE:     "else",
	RETURN:   "return",
	WHILE:    "while",
	BREAK:    "break",
	CONTINUE: "continue",
}

func (ttype TokenType) String() string {
	return tokens[ttype]
}

var keywords = map[string]TokenType{
	"fn":       FN,
	"let":      LET,
	"mut":      MUT,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"while":    WHILE,
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
	return fmt.Sprintf("%d-[%s '%s']", t.Line, t.Type.String(), t.Literal)
}

func (t Token) IsOneOf(expectedTypes ...TokenType) bool {
	return slices.Contains(expectedTypes, t.Type)
}
