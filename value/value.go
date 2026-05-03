// Package value defines the internal representation of Lemon's values
package value

import "fmt"

type ValueType string

const (
	IntegerValType = "integer"
	BooleanValType = "boolean"
	FloatValType   = "float"
	UnitValType    = "unit"
)

type Value interface {
	Type() ValueType
	Inspect() string
}

type (
	Integer struct {
		Raw int64
	}

	Float struct {
		Raw float64
	}

	Boolean struct {
		Raw bool
	}

	Unit struct {
		// Nothing to see here.
	}
)

func (i *Integer) Type() ValueType { return IntegerValType }
func (b *Boolean) Type() ValueType { return BooleanValType }
func (f *Float) Type() ValueType   { return FloatValType }
func (u *Unit) Type() ValueType    { return UnitValType }

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Raw) }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Raw) }
func (f *Float) Inspect() string   { return fmt.Sprintf("%g", f.Raw) }
func (u *Unit) Inspect() string    { return "()" }
