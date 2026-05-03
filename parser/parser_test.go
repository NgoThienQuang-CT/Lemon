package parser

import (
	"fmt"
	"strconv"
	"testing"

	"lemon/ast"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedMutable    bool
		expectedValue      any
	}{
		{"let x = 9;", "x", false, 9},
		{"let mut y = 6.9;", "y", true, 6.9},
		{"let z = true;", "z", false, true},
		{"let name = lemon;", "name", false, "lemon"},
	}

	for _, tt := range tests {
		prog, errors := ParseProgram(tt.input)
		checkParserErrors(t, errors)

		if prog == nil {
			t.Fatalf("ParseProgram() return nil")
		}

		if len(prog.Statements) != 1 {
			t.Fatalf("prog.Statements does not contain 1 statements, got %d instead",
				len(prog.Statements))
		}

		if prog.LastValue != nil {
			t.Fatalf("prog.LastValue is not nil, got %T (%+v) instead",
				prog.LastValue, prog.LastValue)
		}

		stmt := prog.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedMutable) {
			return
		}

		val := stmt.(*ast.LetStmt).Value
		if !testLiteral(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "funny_variable;"

	prog, errors := ParseProgram(input)
	checkParserErrors(t, errors)

	if len(prog.Statements) != 1 {
		t.Fatalf("prog.Statements does not contain 1 statement, got %d instead",
			len(prog.Statements))
	}

	if prog.LastValue != nil {
		t.Fatalf("prog.LastValue is not nil, got %T (%+v) instead",
			prog.LastValue, prog.LastValue)
	}

	stmt, ok := prog.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("prog.Statements[0] is not *ast.ExprStmt, got %T instead",
			prog.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Ident)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Ident, got %T instead",
			stmt.Expression)
	}

	if ident.Value != "funny_variable" {
		t.Errorf("ident.Value is not '%s', got '%s' instead",
			"funny_variable", ident.Value)
	}

	if ident.TokenLiteral() != "funny_variable" {
		t.Errorf("ident.TokenLiteral() is not '%s', got '%s' instead",
			"funny_variable", ident.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"

	prog, errors := ParseProgram(input)
	checkParserErrors(t, errors)

	if len(prog.Statements) != 1 {
		t.Fatalf("prog.Statements does not contain 1 statement, got %d instead",
			len(prog.Statements))
	}

	if prog.LastValue != nil {
		t.Fatalf("prog.LastValue is not nil, got %T (%+v) instead",
			prog.LastValue, prog.LastValue)
	}

	stmt, ok := prog.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("prog.Statements[0] is not *ast.ExprStmt, got %T instead",
			prog.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntLit)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IntLit, got %T instead",
			stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value is not %d, got %d instead",
			5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() is not '%s', got '%s' instead",
			"5", literal.TokenLiteral())
	}
}

func TestFloatLiteral(t *testing.T) {
	input := "6.9;"

	prog, errors := ParseProgram(input)
	checkParserErrors(t, errors)

	if len(prog.Statements) != 1 {
		t.Fatalf("prog.Statements does not contain 1 statement, got %d instead",
			len(prog.Statements))
	}

	if prog.LastValue != nil {
		t.Fatalf("prog.LastValue is not nil, got %T (%+v) instead",
			prog.LastValue, prog.LastValue)
	}

	stmt, ok := prog.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("prog.Statements[0] is not *ast.ExprStmt, got %T instead",
			prog.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.FloatLit)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FloatLit, got %T instead",
			stmt.Expression)
	}

	if literal.Value != 6.9 {
		t.Errorf("literal.Value is not %f, got %f instead",
			6.9, literal.Value)
	}

	if literal.TokenLiteral() != "6.9" {
		t.Errorf("literal.TokenLiteral() is not '%s', got '%s' instead",
			"6.9", literal.TokenLiteral())
	}
}

func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		prog, errors := ParseProgram(tt.input)
		checkParserErrors(t, errors)

		if len(prog.Statements) != 1 {
			t.Fatalf("prog.Statements does not contain 1 statement, got %d instead",
				len(prog.Statements))
		}

		if prog.LastValue != nil {
			t.Fatalf("prog.LastValue is not nil, got %T (%+v) instead",
				prog.LastValue, prog.LastValue)
		}

		stmt, ok := prog.Statements[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("prog.Statements[0] is not *ast.ExprStmt, got %T instead",
				prog.Statements[0])
		}

		literal, ok := stmt.Expression.(*ast.BoolLit)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.BoolLit, got %T instead",
				stmt.Expression)
		}

		if literal.Value != tt.expectedBoolean {
			t.Errorf("literal.Value is not %t, got %t instead",
				tt.expectedBoolean, literal.Value)
		}

		if literal.TokenLiteral() != fmt.Sprintf("%t", tt.expectedBoolean) {
			t.Errorf("literal.TokenLiteral() is not '%t', got '%s' instead",
				tt.expectedBoolean, literal.TokenLiteral())
		}
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range tests {
		prog, errors := ParseProgram(tt.input)
		checkParserErrors(t, errors)

		if len(prog.Statements) != 1 {
			t.Fatalf("prog.Statements does not contain 1 statement, got %d instead",
				len(prog.Statements))
		}

		if prog.LastValue != nil {
			t.Fatalf("prog.LastValue is not nil, got %T (%+v) instead",
				prog.LastValue, prog.LastValue)
		}

		stmt, ok := prog.Statements[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("prog.Statements[0] is not *ast.ExprStmt, got %T instead",
				prog.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.PrefixExpr)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.PrefixExpr, got %T instead",
				stmt.Expression)
		}

		if expr.Op != tt.operator {
			t.Fatalf("expr.Op is not '%s', got '%s' instead",
				tt.operator, expr.Op)
		}

		if !testLiteral(t, expr.Right, tt.value) {
			return
		}
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		leftVal  any
		operator string
		rightVal any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
	}

	for _, tt := range tests {
		prog, errors := ParseProgram(tt.input)
		checkParserErrors(t, errors)

		if len(prog.Statements) != 1 {
			t.Fatalf("prog.Statements does not contain 1 statement, got %d instead",
				len(prog.Statements))
		}

		if prog.LastValue != nil {
			t.Fatalf("prog.LastValue is not nil, got %T (%+v) instead",
				prog.LastValue, prog.LastValue)
		}

		stmt, ok := prog.Statements[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("prog.Statements[0] is not *ast.ExprStmt, got %T instead",
				prog.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftVal, tt.operator, tt.rightVal) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a = b = 6.7 <= 7.9 || 6.9 >= 9.7",
			"(= a (= b (|| (<= 6.7 7.9) (>= 6.9 9.7))))",
		},
		{
			"-a * b",
			"(* (- a) b)",
		},
		{
			"!-a",
			"(! (- a))",
		},
		{
			"a + b + c",
			"(+ (+ a b) c)",
		},
		{
			"a + b - c",
			"(- (+ a b) c)",
		},
		{
			"a * b * c",
			"(* (* a b) c)",
		},
		{
			"a * b / c",
			"(/ (* a b) c)",
		},
		{
			"a + b / c",
			"(+ a (/ b c))",
		},
		{
			"a + b * c + d / e - f",
			"(- (+ (+ a (* b c)) (/ d e)) f)",
		},
		{
			"3 + 4; -5 * 5",
			"(+ 3 4) (* (- 5) 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"(== (> 5 4) (< 3 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"(!= (< 5 4) (> 3 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"(== (+ 3 (* 4 5)) (+ (* 3 1) (* 4 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"(== (> 3 5) false)",
		},
		{
			"3 < 5 == true",
			"(== (< 3 5) true)",
		},
		{
			"1 + (2 + 3) + 4",
			"(+ (+ 1 (+ 2 3)) 4)",
		},
		{
			"(5 + 5) * 2",
			"(* (+ 5 5) 2)",
		},
		{
			"2 / (5 + 5)",
			"(/ 2 (+ 5 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(* (* (+ 5 5) 2) (+ 5 5))",
		},
		{
			"-(5 + 5)",
			"(- (+ 5 5))",
		},
		{
			"!(true == true)",
			"(! (== true true))",
		},
		{
			"a + add(b * c) + d",
			"(+ (+ a (add (* b c))) d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"(add a b 1 (* 2 3) (+ 4 5) (add 6 (* 7 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"(add (+ (+ (+ a b) (/ (* c d) f)) g))",
		},
	}

	for _, tt := range tests {
		prog, errors := ParseProgram(tt.input)
		checkParserErrors(t, errors)

		output := prog.String()
		if output != tt.expected {
			t.Errorf("expected %q, got %q instead",
				tt.expected, output)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if x > y { x }"

	prog, errors := ParseProgram(input)
	checkParserErrors(t, errors)

	if len(prog.Statements) != 0 {
		t.Fatalf("prog.Statements does not contain 0 statement, got %d instead",
			len(prog.Statements))
	}

	if prog.LastValue == nil {
		t.Fatalf("prog.LastValue is not ast.Expr, got nil instead")
	}

	expr, ok := prog.LastValue.(*ast.IfExpr)
	if !ok {
		t.Fatalf("prog.LastValue is not *ast.IfExpr, got %T instead",
			prog.LastValue)
	}

	if !testInfixExpression(t, expr.Cond, "x", ">", "y") {
		return
	}

	if len(expr.Body.Statements) != 0 {
		t.Fatalf("expr.Body.Statements does not contain 0 statement, got %d instead",
			len(expr.Body.Statements))
	}

	if expr.Body.LastValue == nil {
		t.Fatalf("expr.Body.LastValue is not ast.Expr, got nil instead")
	}

	if !testIdentifier(t, expr.Body.LastValue, "x") {
		return
	}

	if expr.Branch != nil {
		t.Fatalf("expr.Branch is not nil, got %T (%+v) instead",
			expr.Branch, expr.Branch)
	}
}

func TestIfExpressionWithBranches(t *testing.T) {
	input := "if x > y { let z = x + y; z } else if x == y { 0 } else { y }"

	prog, errors := ParseProgram(input)
	checkParserErrors(t, errors)

	if len(prog.Statements) != 0 {
		t.Fatalf("prog.Statements does not contain 0 statement, got %d instead",
			len(prog.Statements))
	}

	if prog.LastValue == nil {
		t.Fatalf("prog.LastValue is not ast.Expr, got nil instead")
	}

	expr, ok := prog.LastValue.(*ast.IfExpr)
	if !ok {
		t.Fatalf("prog.LastValue is not *ast.IfExpr, got %T instead",
			prog.LastValue)
	}

	if !testInfixExpression(t, expr.Cond, "x", ">", "y") {
		return
	}

	if len(expr.Body.Statements) != 1 {
		t.Fatalf("expr.Body.Statements does not contain 1 statement, got %d instead",
			len(expr.Body.Statements))
	}

	if expr.Body.LastValue == nil {
		t.Fatalf("expr.Body.LastValue is not ast.Expr, got nil instead")
	}

	stmt, ok := expr.Body.Statements[0].(*ast.LetStmt)
	if !ok {
		t.Fatalf("expr.Body.Statements[0] is not *ast.LetStmt, got %T instead",
			expr.Body.Statements[0])
	}

	if !testLetStatement(t, stmt, "z", false) {
		return
	}

	//letval := stmt.Value
	//if !testInfixExpression(t, letval, "x", "+", "y") {
	//	return
	//}

	if !testIdentifier(t, expr.Body.LastValue, "z") {
		return
	}

	branch, ok := expr.Branch.(*ast.IfExpr)
	if !ok {
		t.Fatalf("expr.Branch is not *ast.IfExpr, got %T instead",
			expr.Branch)
	}

	if !testInfixExpression(t, branch.Cond, "x", "==", "y") {
		return
	}

	if len(branch.Body.Statements) != 0 {
		t.Fatalf("branch.Body.Statements does not contain 0 statement, got %d instead",
			len(branch.Body.Statements))
	}

	if expr.Body.LastValue == nil {
		t.Fatalf("branch.Body.LastValue is not ast.Expr, got nil instead")
	}

	if !testIntegerLiteral(t, branch.Body.LastValue, 0) {
		return
	}

	last, ok := branch.Branch.(*ast.BlockExpr)
	if !ok {
		t.Fatalf("branch.branch is not *ast.BlockExpr, got %T instead",
			branch.Branch)
	}

	if len(last.Statements) != 0 {
		t.Fatalf("last.Statements does not contain 0 statement, got %d instead",
			len(last.Statements))
	}

	if last.LastValue == nil {
		t.Fatalf("last.LastValue is not ast.Expr, got nil instead")
	}

	if !testIdentifier(t, last.LastValue, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := "fn(x, y) { x + y; }"

	prog, errors := ParseProgram(input)
	checkParserErrors(t, errors)

	if len(prog.Statements) != 0 {
		t.Fatalf("prog.Statements does not contain 0 statement, got %d instead",
			len(prog.Statements))
	}

	if prog.LastValue == nil {
		t.Fatalf("prog.LastValue is not ast.Expr, got nil instead")
	}

	fnlit, ok := prog.LastValue.(*ast.FnLit)
	if !ok {
		t.Fatalf("prog.LastValue is not *ast.FnLit, got %T instead",
			prog.LastValue)
	}

	if len(fnlit.Params) != 2 {
		t.Fatalf("fnlit.Params does not contain 2 parameters, got %d instead",
			len(fnlit.Params))
	}

	if !testLiteral(t, fnlit.Params[0], "x") {
		return
	}

	if !testLiteral(t, fnlit.Params[1], "y") {
		return
	}

	if len(fnlit.Body.Statements) != 1 {
		t.Fatalf("fnlit.Body.Statements does not contain 1 statement, got %d instead",
			len(fnlit.Body.Statements))
	}

	if fnlit.Body.LastValue != nil {
		t.Fatalf("prog.LastValue is nil, got %T instead",
			fnlit.Body.LastValue)
	}

	bodyStmt, ok := fnlit.Body.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("fnlit.Body.Statements[0] is not *ast.ExprStmt, got %T instead",
			fnlit.Body.Statements[0])
	}

	if !testInfixExpression(t, bodyStmt.Expression, "x", "+", "y") {
		return
	}
}

func TestFunctionParameters(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		prog, errors := ParseProgram(tt.input)
		checkParserErrors(t, errors)

		if len(prog.Statements) != 1 {
			t.Fatalf("prog.Statements does not contain 1 statements, got %d instead",
				len(prog.Statements))
		}

		if prog.LastValue != nil {
			t.Fatalf("prog.LastValue is not nil, got %T (%+v) instead",
				prog.LastValue, prog.LastValue)
		}

		stmt, ok := prog.Statements[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("prog.Statements[0] is not *ast.ExprStmt, got %T instead",
				prog.Statements[0])
		}

		fn, ok := stmt.Expression.(*ast.FnLit)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.FnLit, got %T instead",
				stmt.Expression)
		}

		if len(fn.Params) != len(tt.expectedParams) {
			t.Errorf("fn.Params does not contain %d parameters, got %d instead",
				len(tt.expectedParams), len(fn.Params))
		}

		for i, ident := range tt.expectedParams {
			testIdentifier(t, fn.Params[i], ident)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	prog, errors := ParseProgram(input)
	checkParserErrors(t, errors)

	if len(prog.Statements) != 1 {
		t.Fatalf("prog.Statements does not contain 1 statements, got %d instead",
			len(prog.Statements))
	}

	if prog.LastValue != nil {
		t.Fatalf("prog.LastValue is not nil, got %T (%+v) instead",
			prog.LastValue, prog.LastValue)
	}

	stmt, ok := prog.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("prog.Statements[0] is not *ast.ExprStmt, got %T instead",
			prog.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.CallExpr)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.CallExpr, got %T instead",
			stmt.Expression)
	}

	if !testIdentifier(t, expr.Func, "add") {
		return
	}

	if !testLiteral(t, expr.Args[0], 1) {
		return
	}

	if !testInfixExpression(t, expr.Args[1], 2, "*", 3) {
		return
	}

	if !testInfixExpression(t, expr.Args[2], 4, "+", 5) {
		return
	}
}

func TestReturnExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{"return 9;", 9},
		{"return 5.6;", 5.6},
		{"return true;", true},
		{"return foo;", "foo"},
	}

	for _, tt := range tests {
		prog, errors := ParseProgram(tt.input)
		checkParserErrors(t, errors)

		if len(prog.Statements) != 1 {
			t.Fatalf("prog.Statements does not contain 1 statements, got %d instead",
				len(prog.Statements))
		}

		if prog.LastValue != nil {
			t.Fatalf("prog.LastValue is not nil, got %T (%+v) instead",
				prog.LastValue, prog.LastValue)
		}

		stmt, ok := prog.Statements[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("prog.Statements[0] is not *ast.ExprStmt, got %T instead",
				prog.Statements[0])
		}

		ret, ok := stmt.Expression.(*ast.ReturnExpr)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.ReturnExpr, got %T instead",
				stmt.Expression)
		}

		if ret.TokenLiteral() != "return" {
			t.Fatalf("ret.TokenLiteral() is not 'return', got '%s' instead",
				ret.TokenLiteral())
		}

		if !testLiteral(t, ret.Value, tt.expectedValue) {
			return
		}
	}
}

func TestWhileExpression(t *testing.T) {
	input := "while x > 0 { continue; break 8; }"

	prog, errors := ParseProgram(input)
	checkParserErrors(t, errors)

	if len(prog.Statements) != 0 {
		t.Fatalf("prog.Statements does not contain 0 statement, got %d instead",
			len(prog.Statements))
	}

	if prog.LastValue == nil {
		t.Fatalf("prog.LastValue is not ast.Expr, got nil instead")
	}

	expr, ok := prog.LastValue.(*ast.WhileExpr)
	if !ok {
		t.Fatalf("prog.LastValue is not *ast.IfExpr, got %T instead",
			prog.LastValue)
	}

	if !testInfixExpression(t, expr.Cond, "x", ">", 0) {
		return
	}

	if len(expr.Body.Statements) != 2 {
		t.Fatalf("expr.Body.Statements does not contain 2 statement, got %d instead",
			len(expr.Body.Statements))
	}

	if expr.Body.LastValue != nil {
		t.Fatalf("expr.Body.LastValue is not nil, got %T (%+v) instead",
			expr.Body.LastValue, expr.Body.LastValue)
	}

	stmt1, ok := expr.Body.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("expr.Body.Statements[0] is not *ast.ExprStmt, got %T instead",
			expr.Body.Statements[0])
	}

	conti, ok := stmt1.Expression.(*ast.ContinueExpr)
	if !ok {
		t.Fatalf("stmt1.Expression is not *ast.ContinueExpr, got %T instead",
			stmt1.Expression)
	}

	if conti.TokenLiteral() != "continue" {
		t.Fatalf("conti.TokenLiteral() is not 'continue', got %T instead",
			conti.TokenLiteral())
	}

	stmt2, ok := expr.Body.Statements[1].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("expr.Body.Statements[1] is not *ast.ExprStmt, got %T instead",
			expr.Body.Statements[1])
	}

	breakExpr, ok := stmt2.Expression.(*ast.BreakExpr)
	if !ok {
		t.Fatalf("stmt2.Expression is not *ast.BreakExpr, got %T instead",
			stmt2.Expression)
	}

	if breakExpr.TokenLiteral() != "break" {
		t.Fatalf("breakExpr.TokenLiteral() is not 'break', got %T instead",
			breakExpr.TokenLiteral())
	}

	if !testLiteral(t, breakExpr.Value, 8) {
		return
	}
}

func checkParserErrors(t *testing.T, errors []string) {
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Error(msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, stmt ast.Stmt, name string, mutable bool) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral is not 'let', got '%s' instead",
			stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStmt)
	if !ok {
		t.Errorf("stmt is not *ast.Stmt, got %T instead", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value is not '%s', got '%s' instead",
			name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() is not '%s', got '%s' instead",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	if letStmt.Mutable != mutable {
		t.Errorf("letStmt.Mutable is not %t, got %t instead",
			mutable, letStmt.Mutable)
		return false
	}

	return true
}

func testLiteral(t *testing.T, lit ast.Expr, expected any) bool {
	switch expected := expected.(type) {
	case int:
		return testIntegerLiteral(t, lit, int64(expected))
	case int64:
		return testIntegerLiteral(t, lit, expected)
	case float32:
		return testFloatLiteral(t, lit, float64(expected))
	case float64:
		return testFloatLiteral(t, lit, expected)
	case string:
		return testIdentifier(t, lit, expected)
	case bool:
		return testBooleanLiteral(t, lit, expected)
	}

	t.Errorf("type of lit not handled, got %T instead", lit)
	return false
}

func testIntegerLiteral(t *testing.T, intlit ast.Expr, value int64) bool {
	integer, ok := intlit.(*ast.IntLit)
	if !ok {
		t.Errorf("intlit is not *ast.IntLit, got %T instead",
			intlit)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value is not %d, got %d instead",
			value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() is not '%d', got '%s' instead",
			value, integer.TokenLiteral())
		return false
	}

	return true
}

func testFloatLiteral(t *testing.T, floatlit ast.Expr, value float64) bool {
	float, ok := floatlit.(*ast.FloatLit)
	if !ok {
		t.Errorf("floatlit is not *ast.FloatLit, got %T instead",
			floatlit)
		return false
	}

	if float.Value != value {
		t.Errorf("float.Value is not %g, got %g instead",
			value, float.Value)
		return false
	}

	litVal, err := strconv.ParseFloat(float.TokenLiteral(), 64)
	if err != nil {
		t.Errorf("float.TokenLiteral() is not valid, got %s instead",
			float.TokenLiteral())
		return false
	}
	if litVal != value {
		t.Errorf("float.TokenLiteral()'s value is not %g, got %g instead",
			value, litVal)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, ident ast.Expr, value string) bool {
	id, ok := ident.(*ast.Ident)
	if !ok {
		t.Errorf("ident is not *ast.Ident, got %T instead",
			ident)
		return false
	}

	if id.Value != value {
		t.Errorf("id.Value is not '%s', got '%s' instead",
			value, id.Value)
		return false
	}

	if id.TokenLiteral() != value {
		t.Errorf("id.TokenLiteral() is not '%s', got '%s' instead",
			value, id.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, boollit ast.Expr, value bool) bool {
	boolean, ok := boollit.(*ast.BoolLit)
	if !ok {
		t.Errorf("boollit is not *ast.BoolLit, got %T instead",
			boollit)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value is not %t, got %t instead",
			value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral() is not '%t', got '%s' instead",
			value, boolean.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, expr ast.Expr, left any, operator string, right any) bool {
	infix, ok := expr.(*ast.InfixExpr)
	if !ok {
		t.Errorf("expr is not *ast.InfixExpr, got %T instead", expr)
		return false
	}

	if !testLiteral(t, infix.Left, left) {
		return false
	}

	if infix.Op != operator {
		t.Errorf("infix.Op is not '%s', got '%s' instead",
			operator, infix.Op)
		return false
	}

	if !testLiteral(t, infix.Right, right) {
		return false
	}

	return true
}
