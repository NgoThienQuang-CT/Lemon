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
		{"-5", -5},
		{"-10", -10},
		{"10 % 3", 1},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerValue(t, evaluated, tt.expected)
	}
}

func TestEvalFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"6.7", 6.7},
		{"6.9", 6.9},
		{"2 * 3.14 * 6", 37.68},
		{"4 / 2.0", 2.0},
		{"(1 + 1.5) * 2 / 10", 0.5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatValue(t, evaluated, tt.expected)
	}
}

func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 <= 2", true},
		{"1 >= 2", false},
		{"1 <= 1", true},
		{"1 >= 1", true},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanValue(t, evaluated, tt.expected)
	}
}

func TestLogicalNegation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanValue(t, evaluated, tt.expected)
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if true { 10 }", 10},
		{"if false { 10 }", nil},
		{"if 1 { 10 }", 10},
		{"if 1 < 2 { 10 }", 10},
		{"if 1 > 2 { 10 }", nil},
		{"if 1 > 2 { 10 } else { 20 }", 20},
		{"if 1 < 2 { 10 } else { 20 }", 10},
		{"if 1 > 2 { 10 } else if 1 == 2 { 20 } else { 30 }", 30},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerValue(t, evaluated, int64(integer))
		} else {
			testUnitValue(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9", 10},
		{
			`
if 10 > 1 {
  if 10 > 1 {
    return 10;
  }

  return 1;
}
`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerValue(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input       string
		expectedMsg string
	}{
		{
			"5 + true;",
			"[line 1] Runtime Error: operator '+' cannot be applied to types integer and boolean",
		},
		{
			"5 + true; 5;",
			"[line 1] Runtime Error: operator '+' cannot be applied to types integer and boolean",
		},
		{
			"-true",
			"[line 1] Runtime Error: operator '-' cannot be applied to type boolean",
		},
		{
			"true + false;",
			"[line 1] Runtime Error: operator '+' cannot be applied to types boolean and boolean",
		},
		{
			"true + false + true + false;",
			"[line 1] Runtime Error: operator '+' cannot be applied to types boolean and boolean",
		},
		{
			"5; true + false; 5",
			"[line 1] Runtime Error: operator '+' cannot be applied to types boolean and boolean",
		},
		{
			"if (10 > 1) { true + false; }",
			"[line 1] Runtime Error: operator '+' cannot be applied to types boolean and boolean",
		},
		{
			`// line 1 here
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"[line 4] Runtime Error: operator '+' cannot be applied to types boolean and boolean",
		},
		{
			"foo",
			"[line 1] Runtime Error: identifier 'foo' is not defined in this scope",
		},
		{
			"let mut a = 9; let mut a = 0; a",
			"[line 1] Runtime Error: cannot re-define mutable variable 'a'",
		},
		{
			"let a = 9; a = a + 1",
			"[line 1] Runtime Error: cannot assign to immutable variable 'a'",
		},
		{
			"let add = fn(x, y) { x + y }; add(1)",
			"[line 1] Runtime Error: not enough arguments in call to add, expected 2, got 1 instead",
		},
		{
			"let add = fn(x, y) { x + y }; add(1, 2, 3)",
			"[line 1] Runtime Error: too many arguments in call to add, expected 2, got 3 instead",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errVal, ok := evaluated.(*value.Error)
		if !ok {
			t.Errorf("no error value is returned, got %T (%+v) instead", errVal, errVal)
			continue
		}

		if errVal.Message != tt.expectedMsg {
			t.Errorf("wrong error message, expected %s, got %s instead",
				tt.expectedMsg, errVal.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a", 5},
		{"let a = 5 * 5; a", 25},
		{"let a = 5; let b = a; b", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c", 15},
		{"let mut a = 0; a = a + 2", 2},
	}

	for _, tt := range tests {
		testIntegerValue(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionValue(t *testing.T) {
	input := "fn(x) { x + 2 }"

	evaluated := testEval(input)
	fn, ok := evaluated.(*value.Function)
	if !ok {
		t.Fatalf("evaluated is not *value.Function, got %T (%+v) instead", evaluated, evaluated)
	}

	if len(fn.Params) != 1 {
		t.Fatalf("fn.Params does not contain 1 parameter, got %+v instead", fn.Params)
	}

	if fn.Params[0].String() != "x" {
		t.Fatalf("fn.Params[0] is not 'x', got '%q' instead", fn.Params[0])
	}

	expectedBody := "(block (+ x 2))"
	if fn.Body.String() != expectedBody {
		t.Fatalf("fn.Body is not %q, got %q instead", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let foo = fn(x) { x }; foo(5)", 5},
		{"let foo = fn(x) { return x; }; foo(5)", 5},
		{"let double = fn(x) { x * 2 }; double(5)", 10},
		{"let add = fn(x, y) { x + y }; add(5, 5)", 10},
		{"let add = fn(x, y) { x + y }; add(5 + 5, add(5, 5))", 20},
		{"fn(x) { x }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerValue(t, testEval(tt.input), tt.expected)
	}
}

func TestClosure(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) { 
		x + y 
	}
};

let add2 = newAdder(2);
add2(2)`

	testIntegerValue(t, testEval(input), 4)
}

func TestLoop(t *testing.T) {
	input := `
let mut sum = 0;
{
	let mut i = 0;
	while i < 10 {
		if i == 5 { i = i + 1; continue; }
		if i == 8 { break sum; }
		sum = sum + i;
		i = i + 1;
	}
}`
	testIntegerValue(t, testEval(input), 23)
}

func testEval(input string) value.Value {
	program, _ := parser.ParseProgram(input)
	scope := value.NewScope()
	return Eval(program, scope)
}

func testUnitValue(t *testing.T, val value.Value) bool {
	if val != UNIT {
		t.Errorf("val is not UNIT, got %T (%+v)", val, val)
		return false
	}
	return true
}

func testIntegerValue(t *testing.T, val value.Value, expected int64) bool {
	result, ok := val.(*value.Integer)
	if !ok {
		t.Errorf("val is not *value.Integer, got %T (%+v) instead", val, val)
		return false
	}

	if result.Inner != expected {
		t.Errorf("result is not %d, got %d instead", expected, result.Inner)
		return false
	}

	return true
}

func testFloatValue(t *testing.T, val value.Value, expected float64) bool {
	result, ok := val.(*value.Float)
	if !ok {
		t.Errorf("val is not *value.Float, got %T (%+v) instead", val, val)
		return false
	}

	if result.Inner != expected {
		t.Errorf("result is not %g, got %g instead", expected, result.Inner)
		return false
	}

	return true
}

func testBooleanValue(t *testing.T, val value.Value, expected bool) bool {
	result, ok := val.(*value.Boolean)
	if !ok {
		t.Errorf("val is not *value.Boolean, got %T (%+v) instead", val, val)
		return false
	}

	if result.Inner != expected {
		t.Errorf("result is not %t, got %t instead", expected, result.Inner)
		return false
	}

	return true
}
