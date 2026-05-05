package ast

import "lemon/internal/token"

type Temp struct {
	Token token.Token
	Value Expr
}

func (t *Temp) statementNode()       {}
func (t *Temp) TokenLiteral() string { return t.Token.Literal }
func (t *Temp) Pos() int             { return t.Token.Line }
func (t *Temp) String() string       { return "<temp>" }
