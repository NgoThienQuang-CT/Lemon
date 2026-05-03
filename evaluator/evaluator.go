// Package evaluator implements the tree-walking execution engine for Lemon.
package evaluator

import (
	"lemon/ast"
	"lemon/value"
)

var UNIT = &value.Unit{}

func Eval(node ast.Node) value.Value {
	switch node := node.(type) {
	case *ast.Prog:
		return evalProgram(node)
	case *ast.ExprStmt:
		return evalExprStatement(node)
	case *ast.IntLit:
		return &value.Integer{Raw: node.Value}
	default:
		return nil
	}
}

func evalProgram(prog *ast.Prog) value.Value {
	result := evalStatements(prog.Statements)
	if prog.LastValue != nil {
		result = Eval(prog.LastValue)
	}

	return result
}

func evalStatements(stmts []ast.Stmt) value.Value {
	var result value.Value

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func evalExprStatement(exprStmt *ast.ExprStmt) value.Value {
	Eval(exprStmt.Expression)
	return UNIT
}
