// Package value defines the internal representation of Lemon's values
package value

import (
	"fmt"
	"strings"

	"lemon/internal/ast"
)

type ValueType string

const (
	IntegerValType  = "integer"
	BooleanValType  = "boolean"
	FloatValType    = "float"
	StringValType   = "string"
	UnitValType     = "unit"
	ReturnValType   = "return"
	FunctionValType = "function"
	BreakValType    = "break"
	ContinueValType = "continue"
	ErrorValType    = "error"
)

type Value interface {
	Type() ValueType
	Inspect() string
}

type (
	Integer struct {
		Inner int64
	}

	Float struct {
		Inner float64
	}

	Boolean struct {
		Inner bool
	}

	String struct {
		Inner string
	}

	Unit struct {
		// Nothing to see here.
	}

	Return struct {
		Inner Value
	}

	Function struct {
		Params []*ast.Ident
		Body   *ast.BlockExpr
		Scope  *Scope
	}

	Break struct {
		Inner Value
	}

	Continue struct {
		// Nothing to see here.
	}

	Error struct {
		Message string
	}
)

func (i *Integer) Type() ValueType  { return IntegerValType }
func (b *Boolean) Type() ValueType  { return BooleanValType }
func (f *Float) Type() ValueType    { return FloatValType }
func (s *String) Type() ValueType   { return StringValType }
func (u *Unit) Type() ValueType     { return UnitValType }
func (r *Return) Type() ValueType   { return ReturnValType }
func (e *Error) Type() ValueType    { return ErrorValType }
func (f *Function) Type() ValueType { return FunctionValType }
func (b *Break) Type() ValueType    { return BreakValType }
func (c *Continue) Type() ValueType { return ContinueValType }

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Inner) }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Inner) }
func (f *Float) Inspect() string   { return fmt.Sprintf("%g", f.Inner) }
func (s *String) Inspect() string  { return s.Inner }
func (u *Unit) Inspect() string    { return "()" }
func (r *Return) Inspect() string  { return r.Inner.Inspect() }
func (e *Error) Inspect() string   { return e.Message }
func (f *Function) Inspect() string {
	var out strings.Builder

	params := []string{}
	for _, param := range f.Params {
		params = append(params, param.String())
	}

	out.WriteString("(fn (")
	out.WriteString(strings.Join(params, " "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	out.WriteString(")")

	return out.String()
}
func (b *Break) Inspect() string    { return "!" }
func (c *Continue) Inspect() string { return "!" }
