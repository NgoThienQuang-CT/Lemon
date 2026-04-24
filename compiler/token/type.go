package token

type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Identifier and basic type literals
	IDENT  // x, y, foo, bar
	INT    // 1, 2, 366
	FLOAT  // 0.9, 78.25
	STRING // "Hello, world!"

	// Operaters and delimiters
	ASSIGN // =

	PLUS    // +
	MINUS   // -
	STAR    // *
	SLASH   // /
	PERCENT // %

	AND // &&
	OR  // ||
	EQL // ==
	LSS // <
	GTR // >
	NOT // !
	NEQ // !=
	LEQ // <=
	GEQ // >=

	COMMA     // ,
	PERIOD    // .
	COLON     // :
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }

	// Keywords
	LET
	FN
	RETURN
	IF
	ELSE
	LOOP
	BREAK
	CONTINUE
	TRUE
	FALSE
	NIL
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

	ASSIGN: "ASSIGN",

	PLUS:    "PLUS",
	MINUS:   "MINUS",
	STAR:    "STAR",
	SLASH:   "SLASH",
	PERCENT: "PERCENT",

	AND: "AND",
	OR:  "OR",
	EQL: "EQL",
	LSS: "LSS",
	GTR: "GTR",
	NOT: "NOT",
	NEQ: "NEQ",
	LEQ: "LEQ",
	GEQ: "GEQ",

	COMMA:     "COMMA",
	PERIOD:    "PERIOD",
	SEMICOLON: "SEMICOLON",
	COLON:     "COLON",
	LPAREN:    "LPAREN",
	RPAREN:    "RPAREN",
	LBRACE:    "LBRACE",
	RBRACE:    "RBRACE",

	LET:      "LET",
	FN:       "FN",
	RETURN:   "RETURN",
	IF:       "IF",
	ELSE:     "ELSE",
	LOOP:     "LOOP",
	BREAK:    "BREAK",
	CONTINUE: "CONTINUE",
	TRUE:     "TRUE",
	FALSE:    "FALSE",
	NIL:      "NIL",
}

var keywords = map[string]TokenType{
	"let":      LET,
	"fn":       FN,
	"return":   RETURN,
	"if":       IF,
	"else":     ELSE,
	"loop":     LOOP,
	"break":    BREAK,
	"continue": CONTINUE,
	"true":     TRUE,
	"false":    FALSE,
	"nil":      NIL,
}

func (t TokenType) String() string {
	if t < 0 || t > TokenType(len(tokens)) {
		return "UNKNOWN"
	}
	return tokens[t]
}

func KeywordsLookup(ident string) TokenType {
	if t, isKeyword := keywords[ident]; isKeyword {
		return t
	}
	return IDENT
}
