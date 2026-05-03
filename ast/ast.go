// Package ast declares the types used to represent syntax trees for Lemon
// programming language.
package ast

import (
	"strings"

	"lemon/token"
)

type (
	Node interface {
		TokenLiteral() string
		String() string
	}

	Stmt interface {
		Node
		statementNode()
	}

	Expr interface {
		Node
		expressionNode()
	}
)

type Prog struct {
	Statements []Stmt
	LastValue  Expr
}

func (p *Prog) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	if p.LastValue != nil {
		return p.LastValue.TokenLiteral()
	}
	return ""
}

func (p *Prog) String() string {
	var out strings.Builder

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
		out.WriteString(" ")
	}

	if p.LastValue != nil {
		out.WriteString(p.LastValue.String())
	}

	return out.String()
}

type (
	LetStmt struct {
		Token   token.Token
		Name    *Ident
		Value   Expr
		Mutable bool
	}

	ExprStmt struct {
		Token      token.Token
		Expression Expr
	}
)

func (ls *LetStmt) statementNode()  {}
func (es *ExprStmt) statementNode() {}

func (ls *LetStmt) TokenLiteral() string  { return ls.Token.Literal }
func (es *ExprStmt) TokenLiteral() string { return es.Token.Literal }

func (ls *LetStmt) String() string {
	var out strings.Builder

	out.WriteString("(let")
	if ls.Mutable {
		out.WriteString("-mut")
	}
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" ")
	out.WriteString(ls.Value.String())
	out.WriteString(")")

	return out.String()
}

func (es *ExprStmt) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type (
	Ident struct {
		Token token.Token
		Value string
	}

	BoolLit struct {
		Token token.Token
		Value bool
	}

	IntLit struct {
		Token token.Token
		Value int64
	}

	FloatLit struct {
		Token token.Token
		Value float64
	}

	StringLit struct {
		Token token.Token
		Value string
	}

	PrefixExpr struct {
		Token token.Token
		Op    string
		Right Expr
	}

	InfixExpr struct {
		Token token.Token
		Left  Expr
		Op    string
		Right Expr
	}

	IfExpr struct {
		Token  token.Token
		Cond   Expr
		Body   *BlockExpr
		Branch Expr
	}

	BlockExpr struct {
		Token      token.Token
		Statements []Stmt
		LastValue  Expr
	}

	FnLit struct {
		Token  token.Token
		Params []*Ident
		Body   *BlockExpr
	}

	CallExpr struct {
		Token token.Token
		Func  Expr
		Args  []Expr
	}

	WhileExpr struct {
		Token token.Token
		Cond  Expr
		Body  *BlockExpr
	}

	ReturnExpr struct {
		Token token.Token
		Value Expr
	}

	BreakExpr struct {
		Token token.Token
		Value Expr
	}

	ContinueExpr struct {
		Token token.Token
	}
)

func (i *Ident) expressionNode()         {}
func (bl *BoolLit) expressionNode()      {}
func (il *IntLit) expressionNode()       {}
func (fl *FloatLit) expressionNode()     {}
func (sl *StringLit) expressionNode()    {}
func (pe *PrefixExpr) expressionNode()   {}
func (ie *InfixExpr) expressionNode()    {}
func (ie *IfExpr) expressionNode()       {}
func (be *BlockExpr) expressionNode()    {}
func (fl *FnLit) expressionNode()        {}
func (ce *CallExpr) expressionNode()     {}
func (we *WhileExpr) expressionNode()    {}
func (re *ReturnExpr) expressionNode()   {}
func (be *BreakExpr) expressionNode()    {}
func (ce *ContinueExpr) expressionNode() {}

func (i *Ident) TokenLiteral() string         { return i.Token.Literal }
func (bl *BoolLit) TokenLiteral() string      { return bl.Token.Literal }
func (il *IntLit) TokenLiteral() string       { return il.Token.Literal }
func (fl *FloatLit) TokenLiteral() string     { return fl.Token.Literal }
func (sl *StringLit) TokenLiteral() string    { return sl.Token.Literal }
func (pe *PrefixExpr) TokenLiteral() string   { return pe.Token.Literal }
func (ie *InfixExpr) TokenLiteral() string    { return ie.Token.Literal }
func (ie *IfExpr) TokenLiteral() string       { return ie.Token.Literal }
func (be *BlockExpr) TokenLiteral() string    { return be.Token.Literal }
func (fl *FnLit) TokenLiteral() string        { return fl.Token.Literal }
func (ce *CallExpr) TokenLiteral() string     { return ce.Token.Literal }
func (we *WhileExpr) TokenLiteral() string    { return we.Token.Literal }
func (re *ReturnExpr) TokenLiteral() string   { return re.Token.Literal }
func (be *BreakExpr) TokenLiteral() string    { return be.Token.Literal }
func (ce *ContinueExpr) TokenLiteral() string { return ce.Token.Literal }

func (i *Ident) String() string      { return i.Token.Literal }
func (bl *BoolLit) String() string   { return bl.Token.Literal }
func (il *IntLit) String() string    { return il.Token.Literal }
func (fl *FloatLit) String() string  { return fl.Token.Literal }
func (sl *StringLit) String() string { return sl.Token.Literal }
func (pe *PrefixExpr) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(pe.Op)
	out.WriteString(" ")
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

func (ie *InfixExpr) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ie.Op)
	out.WriteString(" ")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

func (ie *IfExpr) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ie.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ie.Cond.String())
	out.WriteString(" ")
	out.WriteString(ie.Body.String())
	if ie.Branch != nil {
		out.WriteString(" ")
		out.WriteString(ie.Branch.String())
	}
	out.WriteString(")")

	return out.String()
}

func (be *BlockExpr) String() string {
	var out strings.Builder

	out.WriteString("(block")
	for _, stmt := range be.Statements {
		out.WriteString(" ")
		out.WriteString(stmt.String())
	}
	if be.LastValue != nil {
		out.WriteString(" ")
		out.WriteString(be.LastValue.String())
	}
	out.WriteString(")")

	return out.String()
}

func (fl *FnLit) String() string {
	var out strings.Builder

	params := []string{}
	for _, param := range fl.Params {
		params = append(params, param.String())
	}

	out.WriteString("(")
	out.WriteString(fl.TokenLiteral())
	out.WriteString(" (")
	out.WriteString(strings.Join(params, " "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	out.WriteString(")")

	return out.String()
}

func (ce *CallExpr) String() string {
	var out strings.Builder

	args := []string{}
	for _, arg := range ce.Args {
		args = append(args, arg.String())
	}

	out.WriteString("(")
	out.WriteString(ce.Func.String())
	out.WriteString(" ")
	out.WriteString(strings.Join(args, " "))
	out.WriteString(")")

	return out.String()
}

func (we *WhileExpr) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(we.TokenLiteral())
	out.WriteString(" (")
	out.WriteString(we.Cond.String())
	out.WriteString(") ")
	out.WriteString(we.Body.String())
	out.WriteString(")")

	return out.String()
}

func (re *ReturnExpr) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(re.TokenLiteral())
	if re.Value != nil {
		out.WriteString(" ")
		out.WriteString(re.Value.String())
	}
	out.WriteString(")")

	return out.String()
}

func (be *BreakExpr) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(be.TokenLiteral())
	if be.Value != nil {
		out.WriteString(" ")
		out.WriteString(be.Value.String())
	}
	out.WriteString(")")

	return out.String()
}

func (ce *ContinueExpr) String() string {
	return "(" + ce.TokenLiteral() + ")"
}
