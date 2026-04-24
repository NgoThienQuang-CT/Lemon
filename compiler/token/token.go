// Package token defines constants representing the lexical tokens of Lemon
// programming language
package token

type Token struct {
	Type   TokenType
	Pos    Position
	Lexeme string
}

func (t *Token) String() string {
	return t.Pos.String() + " [" + t.Type.String() + "] " + t.Lexeme
}
