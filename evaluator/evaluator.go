// Package evaluator implements the tree-walking execution engine for Lemon.
package evaluator

import (
	"fmt"

	"lemon/ast"
	"lemon/value"
)

var (
	UNIT  = &value.Unit{}
	TRUE  = &value.Boolean{Inner: true}
	FALSE = &value.Boolean{Inner: false}
)

type evaluator struct {
	scope *value.Scope
	line  int
}

func Eval(node ast.Node, scope *value.Scope) value.Value {
	e := evaluator{scope: scope}
	return e.eval(node)
}

func (e *evaluator) eval(node ast.Node) value.Value {
	switch node := node.(type) {
	case *ast.Prog:
		return e.evalProg(node)
	case *ast.LetStmt:
		e.line = node.Token.Line
		return e.evalLetStmt(node)
	case *ast.ExprStmt:
		e.line = node.Token.Line
		return e.evalExprStmt(node)
	case *ast.Ident:
		e.line = node.Token.Line
		return e.evalIdent(node)
	case *ast.IntLit:
		e.line = node.Token.Line
		return e.evalInteger(node)
	case *ast.FloatLit:
		e.line = node.Token.Line
		return e.evalFloat(node)
	case *ast.BoolLit:
		e.line = node.Token.Line
		return e.evalBoolean(node)
	case *ast.PrefixExpr:
		e.line = node.Token.Line
		return e.evalPrefixExpr(node)
	case *ast.InfixExpr:
		e.line = node.Token.Line
		return e.evalInfixExpr(node)
	case *ast.BlockExpr:
		e.line = node.Token.Line
		return e.evalBlockExpr(node)
	case *ast.IfExpr:
		e.line = node.Token.Line
		return e.evalIfExpr(node)
	case *ast.ReturnExpr:
		e.line = node.Token.Line
		return e.evalReturnExpr(node)
	case *ast.FnLit:
		e.line = node.Token.Line
		return e.evalFunction(node)
	case *ast.CallExpr:
		e.line = node.Token.Line
		return e.evalCallExpr(node)
	case *ast.WhileExpr:
		e.line = node.Token.Line
		return e.evalWhileExpr(node)
	case *ast.BreakExpr:
		e.line = node.Token.Line
		return e.evalBreakExpr(node)
	case *ast.ContinueExpr:
		e.line = node.Token.Line
		return e.evalContinueExpr(node)
	default:
		return nil
	}
}

func boolToBoolValue(input bool) *value.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isNumber(val value.Value) bool {
	switch val.(type) {
	case *value.Integer, *value.Float:
		return true
	default:
		return false
	}
}

func numberValueToFloat(val value.Value) float64 {
	switch v := val.(type) {
	case *value.Integer:
		return float64(v.Inner)
	case *value.Float:
		return v.Inner
	default:
		return 0.0
	}
}

func isTruthy(val value.Value) bool {
	switch val {
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func (e *evaluator) newError(format string, a ...any) value.Value {
	msg := fmt.Sprintf("[line %d] Runtime Error: ", e.line)
	msg += fmt.Sprintf(format, a...)
	return &value.Error{Message: msg}
}

func isError(val value.Value) bool {
	if val != nil {
		return val.Type() == value.ErrorValType
	}
	return false
}

func isReturn(val value.Value) bool {
	if val != nil {
		return val.Type() == value.ReturnValType
	}
	return false
}

func (e *evaluator) evalProg(node *ast.Prog) value.Value {
	var result value.Value = UNIT

	for _, stmt := range node.Statements {
		result = e.eval(stmt)

		switch result := result.(type) {
		case *value.Return:
			return result.Inner
		case *value.Error:
			return result
		}
	}

	if node.LastValue != nil {
		result = e.eval(node.LastValue)
		switch result := result.(type) {
		case *value.Return:
			return result.Inner
		case *value.Error:
			return result
		}
	}

	return result
}

func (e *evaluator) evalLetStmt(node *ast.LetStmt) value.Value {
	val := e.eval(node.Value)
	if isError(val) || val.Type() == value.ReturnValType {
		return val
	}

	_, ok := e.scope.Define(node.Name.Value, node.Mutable, val)

	if !ok {
		return e.newError(
			"cannot re-define mutable variable '%s'",
			node.Name.Value,
		)
	}

	return UNIT
}

func (e *evaluator) evalExprStmt(node *ast.ExprStmt) value.Value {
	result := e.eval(node.Expression)
	if isError(result) || isReturn(result) ||
		result.Type() == value.BreakValType || result.Type() == value.ContinueValType {
		return result
	}

	return UNIT
}

func (e *evaluator) evalIdent(node *ast.Ident) value.Value {
	val, ok := e.scope.Get(node.Value)
	if !ok {
		return e.newError(
			"identifier '%s' is not defined in this scope",
			node.Value,
		)
	}

	return val
}

func (e *evaluator) evalBoolean(node *ast.BoolLit) value.Value {
	return boolToBoolValue(node.Value)
}

func (e *evaluator) evalInteger(node *ast.IntLit) value.Value {
	return &value.Integer{Inner: node.Value}
}

func (e *evaluator) evalFloat(node *ast.FloatLit) value.Value {
	return &value.Float{Inner: node.Value}
}

func (e *evaluator) evalPrefixExpr(node *ast.PrefixExpr) value.Value {
	right := e.eval(node.Right)
	if isError(right) || right.Type() == value.ReturnValType {
		return right
	}

	switch node.Op {
	case "!":
		return e.evalLogicalNegation(right)
	case "-":
		return e.evalNumericalNegation(right)
	default:
		return e.newError(
			"operator '%s' cannot be applied to type %s",
			node.Op, right.Type(),
		)
	}
}

func (e *evaluator) evalLogicalNegation(right value.Value) value.Value {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return e.newError(
			"operator '!' cannot be applied to type %s",
			right.Type(),
		)
	}
}

func (e *evaluator) evalNumericalNegation(right value.Value) value.Value {
	if right.Type() == value.IntegerValType {
		val := right.(*value.Integer).Inner
		return &value.Integer{Inner: -val}
	}

	if right.Type() == value.FloatValType {
		val := right.(*value.Float).Inner
		return &value.Float{Inner: -val}
	}

	return e.newError(
		"operator '-' cannot be applied to type %s",
		right.Type(),
	)
}

func (e *evaluator) evalInfixExpr(node *ast.InfixExpr) value.Value {
	right := e.eval(node.Right)
	if isError(right) || right.Type() == value.ReturnValType {
		return right
	}

	if node.Op == "=" {
		return e.evalAssignment(node.Left, right)
	}

	left := e.eval(node.Left)
	if isError(left) || left.Type() == value.ReturnValType {
		return left
	}

	switch {
	case left.Type() == value.BooleanValType && right.Type() == value.BooleanValType:
		return e.evalBooleanInfixExpr(node.Op, left, right)
	case left.Type() == value.IntegerValType && right.Type() == value.IntegerValType:
		return e.evalIntegerInfixExpr(node.Op, left, right)
	case isNumber(left) && isNumber(right):
		return e.evalNumberInfixExpr(node.Op, left, right)
	case node.Op == "==":
		return boolToBoolValue(left == right)
	case node.Op == "!=":
		return boolToBoolValue(left != right)
	default:
		return e.newError(
			"operator '%s' cannot be applied to types %s and %s",
			node.Op, left.Type(), right.Type(),
		)
	}
}

func (e *evaluator) evalAssignment(node ast.Node, val value.Value) value.Value {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return e.newError(
			"invalid assignment target (not an identifier)",
		)
	}

	updated, err := e.scope.Set(ident.Value, val)
	switch err {
	case value.UNDEFINED:
		return e.newError(
			"identifier '%s' is not defined in this scope",
			ident.Value,
		)
	case value.IMMUTABLE:
		return e.newError(
			"cannot assign to immutable variable '%s'",
			ident.Value,
		)
	}

	return updated
}

func (e *evaluator) evalBooleanInfixExpr(operator string, left value.Value, right value.Value) value.Value {
	l := left.(*value.Boolean).Inner
	r := right.(*value.Boolean).Inner

	switch operator {
	case "||":
		return boolToBoolValue(l || r)
	case "&&":
		return boolToBoolValue(l && r)
	default:
		return e.newError(
			"operator '%s' cannot be applied to types %s and %s",
			operator, left.Type(), right.Type(),
		)
	}
}

func (e *evaluator) evalIntegerInfixExpr(operator string, left value.Value, right value.Value) value.Value {
	l := left.(*value.Integer).Inner
	r := right.(*value.Integer).Inner

	switch operator {
	case "+":
		return &value.Integer{Inner: l + r}
	case "-":
		return &value.Integer{Inner: l - r}
	case "*":
		return &value.Integer{Inner: l * r}
	case "/":
		if l == 0 {
			return e.newError("division by zero")
		}
		return &value.Integer{Inner: l / r}
	case "%":
		if l == 0 {
			return e.newError("division by zero")
		}
		return &value.Integer{Inner: l % r}
	case ">":
		return boolToBoolValue(l > r)
	case ">=":
		return boolToBoolValue(l >= r)
	case "<":
		return boolToBoolValue(l < r)
	case "<=":
		return boolToBoolValue(l <= r)
	case "!=":
		return boolToBoolValue(l != r)
	case "==":
		return boolToBoolValue(l == r)
	default:
		return e.newError(
			"operator '%s' cannot be applied to types %s and %s",
			operator, left.Type(), right.Type(),
		)
	}
}

func (e *evaluator) evalNumberInfixExpr(operator string, left value.Value, right value.Value) value.Value {
	l := numberValueToFloat(left)
	r := numberValueToFloat(right)

	switch operator {
	case "+":
		return &value.Float{Inner: l + r}
	case "-":
		return &value.Float{Inner: l - r}
	case "*":
		return &value.Float{Inner: l * r}
	case "/":
		if l == 0.0 {
			return e.newError("division by zero")
		}
		return &value.Float{Inner: l / r}
	case ">":
		return boolToBoolValue(l > r)
	case ">=":
		return boolToBoolValue(l >= r)
	case "<":
		return boolToBoolValue(l < r)
	case "<=":
		return boolToBoolValue(l <= r)
	case "!=":
		return boolToBoolValue(l != r)
	case "==":
		return boolToBoolValue(l == r)
	default:
		return e.newError(
			"operator '%s' cannot be applied to types %s and %s",
			operator, left.Type(), right.Type(),
		)
	}
}

func (e *evaluator) evalBlockExpr(node *ast.BlockExpr) value.Value {
	blockScope := value.NewEnclosedScope(e.scope)
	oldScope := e.scope
	e.scope = blockScope

	defer func() { e.scope = oldScope }()

	var result value.Value = UNIT

	for _, stmt := range node.Statements {
		result = e.eval(stmt)

		if result != nil {
			rtype := result.Type()
			if rtype == value.ErrorValType || rtype == value.ReturnValType ||
				result.Type() == value.BreakValType || result.Type() == value.ContinueValType {
				return result
			}
		}
	}

	if node.LastValue != nil {
		result = e.eval(node.LastValue)

		rtype := result.Type()
		if rtype == value.ErrorValType || rtype == value.ReturnValType ||
			result.Type() == value.BreakValType || result.Type() == value.ContinueValType {
			return result
		}
	}

	return result
}

func (e *evaluator) evalIfExpr(node *ast.IfExpr) value.Value {
	condition := e.eval(node.Cond)
	if isError(condition) || condition.Type() == value.ReturnValType {
		return condition
	}

	if isTruthy(condition) {
		return e.eval(node.Body)
	}

	if node.Branch != nil {
		return e.eval(node.Branch)
	}

	return UNIT
}

func (e *evaluator) evalReturnExpr(node *ast.ReturnExpr) value.Value {
	return &value.Return{Inner: e.eval(node.Value)}
}

func (e *evaluator) evalFunction(node *ast.FnLit) value.Value {
	return &value.Function{
		Params: node.Params,
		Body:   node.Body,
		Scope:  e.scope,
	}
}

func (e *evaluator) evalCallExpr(node *ast.CallExpr) value.Value {
	function := e.eval(node.Func)
	if isError(function) || isReturn(function) {
		return function
	}

	fnName := "<anonymous>"
	if ident, ok := node.Func.(*ast.Ident); ok {
		fnName = ident.Value
	}

	args := e.evalExprList(node.Args)
	if len(args) == 1 && (isError(args[0]) || isReturn(args[0])) {
		return args[0]
	}

	return e.applyFunction(fnName, function, args)
}

func (e *evaluator) evalExprList(list []ast.Expr) []value.Value {
	var result []value.Value

	for _, expr := range list {
		evaluated := e.eval(expr)

		if isError(evaluated) || isReturn(evaluated) {
			return []value.Value{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func (e *evaluator) applyFunction(name string, fn value.Value, args []value.Value) value.Value {
	function, ok := fn.(*value.Function)
	if !ok {
		return e.newError(
			"%s is not a function",
			name,
		)
	}

	if len(function.Params) > len(args) {
		return e.newError(
			"not enough arguments in call to %s, expected %d, got %d instead",
			name, len(function.Params), len(args),
		)
	}

	if len(function.Params) < len(args) {
		return e.newError(
			"too many arguments in call to %s, expected %d, got %d instead",
			name, len(function.Params), len(args),
		)
	}

	extendedScope := e.extendFunctionScope(function, args)
	oldScope := e.scope
	e.scope = extendedScope

	defer func() { e.scope = oldScope }()

	result := e.eval(function.Body)

	return e.unwrapReturn(result)
}

func (e *evaluator) extendFunctionScope(fn *value.Function, args []value.Value) *value.Scope {
	scope := value.NewEnclosedScope(fn.Scope)

	for i, param := range fn.Params {
		scope.Define(param.Value, false, args[i])
	}

	return scope
}

func (e *evaluator) unwrapReturn(val value.Value) value.Value {
	if retVal, ok := val.(*value.Return); ok {
		return retVal.Inner
	}
	return val
}

func (e *evaluator) evalWhileExpr(node *ast.WhileExpr) value.Value {
	var lastval value.Value = UNIT

	for {
		condition := e.eval(node.Cond)
		if isError(condition) || isReturn(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result := e.eval(node.Body)

		if isError(result) || isReturn(result) {
			return result
		}

		if result.Type() == value.BreakValType {
			result := result.(*value.Break)
			if result.Inner != nil {
				lastval = result.Inner
			}

			break
		}

		if result.Type() == value.ContinueValType {
			continue
		}

		lastval = result
	}

	return lastval
}

func (e *evaluator) evalBreakExpr(node *ast.BreakExpr) value.Value {
	return &value.Break{Inner: e.eval(node.Value)}
}

func (e *evaluator) evalContinueExpr(node *ast.ContinueExpr) value.Value {
	return &value.Continue{}
}
