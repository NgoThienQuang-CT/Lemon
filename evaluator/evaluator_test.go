package evaluator

import (
	"testing"

	"lemon/parser"
	"lemon/value"
)

func TestEvalInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) value.Value {
	program, _ := parser.ParseProgram(input)
	return Eval(program)
}

func testIntegerObject(t *testing.T, val value.Value, expected int64) bool {
	result, ok := val.(*value.Integer)
	if !ok {
		t.Errorf("val is not *value.Integer, got %T (%+v) instead", val, val)
		return false
	}

	if result.Raw != expected {
		t.Errorf("result is not %d, got %d instead", expected, result.Raw)
		return false
	}

	return true
}
