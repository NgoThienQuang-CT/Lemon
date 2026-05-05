// Package parser implements a parser Lemon Go source files
package parser

import (
	"fmt"
	"strconv"

	"lemon/internal/ast"
	"lemon/internal/lexer"
	"lemon/internal/token"
)

const (
	_ int = iota
	LOWEST
	ASSIGN
	OR
	AND
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:  ASSIGN,
	token.OR:      OR,
	token.AND:     AND,
	token.EQL:     EQUALS,
	token.NEQ:     EQUALS,
	token.LSS:     LESSGREATER,
	token.GTR:     LESSGREATER,
	token.LEQ:     LESSGREATER,
	token.GEQ:     LESSGREATER,
	token.PLUS:    SUM,
	token.MINUS:   SUM,
	token.STAR:    PRODUCT,
	token.SLASH:   PRODUCT,
	token.PERCENT: PRODUCT,
	token.LPAREN:  CALL,
}

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

type parser struct {
	lexer     lexer.Lexer
	errors    []string
	panicMode bool

	current token.Token
	next    token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	inloop int
}

func (p *parser) init(source string) {
	p.lexer.Init(source)
	p.errors = []string{}
	p.panicMode = false

	// Fill up the current and next token (non-llegal) of the parser
	p.advance() // the first token go to p.next
	p.advance() // the first token is transfered to p.current and the second token go to p.next

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.LCURLY, p.parseBlockExpression)
	p.registerPrefix(token.FN, p.parseFunctionLiteral)
	p.registerPrefix(token.RETURN, p.parseReturnExpression)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.EQL, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LSS, p.parseInfixExpression)
	p.registerInfix(token.GTR, p.parseInfixExpression)
	p.registerInfix(token.LEQ, p.parseInfixExpression)
	p.registerInfix(token.GEQ, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.STAR, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.PERCENT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
}

func (p *parser) errorAt(tok token.Token, msg string) {
	if p.panicMode {
		return
	}
	p.panicMode = true

	emsg := fmt.Sprintf("[line %d] Error", tok.Line)
	switch tok.Type {
	case token.EOF:
		emsg += " at end"
	case token.ILLEGAL:
		break
	default:
		emsg += (" at '" + tok.Literal + "'")
	}
	emsg += (": " + msg + "\n")

	p.errors = append(p.errors, emsg)
}

func (p *parser) errorAtCurrent(msg string) {
	p.errorAt(p.current, msg)
}

func (p *parser) errorAtNext(msg string) {
	p.errorAt(p.next, msg)
}

func (p *parser) advance() {
	p.current = p.next

	// Skip and report all illegal token until hit a legit token
	for {
		p.next = p.lexer.NextToken()

		if p.next.Type != token.ILLEGAL {
			break
		}

		p.errorAtNext(p.current.Literal)
	}
}

func (p *parser) currentIs(t token.TokenType) bool {
	return p.current.Type == t
}

func (p *parser) nextIs(t token.TokenType) bool {
	return p.next.Type == t
}

func (p *parser) consume(t token.TokenType, message string) bool {
	if !p.nextIs(t) {
		msg := fmt.Sprintf("expected %s after %s, but found %s instead.",
			t, p.current.Type, p.next.Type)
		if message != "" {
			p.errorAtNext(message)
			return false
		}
		p.errorAtNext(msg)
		return false
	}

	p.advance()

	return true
}

func (p *parser) noPrefixParseFnError() {
	msg := fmt.Sprintf("no prefix parse function for '%s' found", p.current.Type)
	p.errorAtCurrent(msg)
}

func isRightAssociativity(t token.Token) bool {
	return t.IsOneOf(token.ASSIGN)
}

func (p *parser) precedenceAtCurrent() int {
	if precedence, ok := precedences[p.current.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *parser) precedenceAtNext() int {
	if precedence, ok := precedences[p.next.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *parser) registerPrefix(ttype token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[ttype] = fn
}

func (p *parser) registerInfix(ttype token.TokenType, fn infixParseFn) {
	p.infixParseFns[ttype] = fn
}

func isBlockBased(node ast.Node) bool {
	switch node.(type) {
	case *ast.IfExpr, *ast.BlockExpr, *ast.WhileExpr:
		return true
	default:
		return false
	}
}

func ParseProgram(source string) (*ast.Prog, []string) {
	var p parser
	p.init(source)

	return p.parseProgram(), p.errors
}

func (p *parser) parseProgram() *ast.Prog {
	prog := &ast.Prog{
		Statements: []ast.Stmt{},
	}

	for !p.currentIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt == nil {
			p.advance()
			continue
		}

		if last, ok := stmt.(*ast.Temp); ok {
			if last.Value == nil {
				return prog
			}
			prog.LastValue = last.Value
			p.advance()
			return prog
		}

		prog.Statements = append(prog.Statements, stmt)
		p.advance()
	}

	return prog
}

func (p *parser) synchronize() {
	p.panicMode = false
	p.advance()

	for !p.nextIs(token.EOF) {
		if p.currentIs(token.SEMICOLON) {
			return
		}

		switch p.next.Type {
		case token.LET, token.FN, token.WHILE, token.IF, token.RETURN, token.BREAK, token.CONTINUE:
			return
		}

		p.advance()
	}
}

func (p *parser) parseStatement() ast.Stmt {
	var stmt ast.Stmt

	switch p.current.Type {
	case token.LET:
		stmt = p.parseLetStatement()
	default:
		stmt = p.parseExpressionStatement()
	}

	if p.panicMode {
		p.synchronize()
	}

	return stmt
}

func (p *parser) parseLetStatement() ast.Stmt {
	stmt := &ast.LetStmt{Token: p.current}

	stmt.Mutable = false
	if p.nextIs(token.MUT) {
		p.advance()
		stmt.Mutable = true
	}

	if !p.consume(token.IDENT, "") {
		return nil
	}

	stmt.Name = &ast.Ident{Token: p.current, Value: p.current.Literal}

	if !p.consume(token.ASSIGN, "") {
		return nil
	}

	p.advance()

	stmt.Value = p.parseExpression(LOWEST)

	if !p.consume(token.SEMICOLON, "expected ';' after declaration") {
		return nil
	}

	return stmt
}

func (p *parser) parseExpressionStatement() ast.Stmt {
	current := p.current
	stmt := &ast.ExprStmt{Token: p.current}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.nextIs(token.SEMICOLON) {
		p.advance()
		return stmt
	}

	if p.nextIs(token.RCURLY) || p.nextIs(token.EOF) {
		return &ast.Temp{Token: current, Value: stmt.Expression}
	}

	if isBlockBased(stmt.Expression) {
		return stmt
	}

	if !p.consume(token.SEMICOLON, "expected ';' after expression") {
		return nil
	}

	return stmt
}

func (p *parser) parseExpression(precedence int) ast.Expr {
	prefix := p.prefixParseFns[p.current.Type]
	if prefix == nil {
		p.noPrefixParseFnError()
		return nil
	}
	left := prefix()

	for !p.nextIs(token.SEMICOLON) && precedence < p.precedenceAtNext() {
		infix := p.infixParseFns[p.next.Type]
		if infix == nil {
			return left
		}

		p.advance()

		left = infix(left)
	}

	return left
}

func (p *parser) parseIdentifier() ast.Expr {
	return &ast.Ident{
		Token: p.current,
		Value: p.current.Literal,
	}
}

func (p *parser) parseIntegerLiteral() ast.Expr {
	lit := &ast.IntLit{Token: p.current}

	val, err := strconv.ParseInt(p.current.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.current.Literal)
		p.errorAtCurrent(msg)
		return nil
	}

	lit.Value = val

	return lit
}

func (p *parser) parseFloatLiteral() ast.Expr {
	lit := &ast.FloatLit{Token: p.current}

	val, err := strconv.ParseFloat(p.current.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.current.Literal)
		p.errorAtCurrent(msg)
		return nil
	}

	lit.Value = val

	return lit
}

func (p *parser) parseBooleanLiteral() ast.Expr {
	return &ast.BoolLit{
		Token: p.current,
		Value: p.currentIs(token.TRUE),
	}
}

func (p *parser) parseStringLiteral() ast.Expr {
	return &ast.StringLit{
		Token: p.current,
		Value: p.current.Literal,
	}
}

func (p *parser) parsePrefixExpression() ast.Expr {
	expr := &ast.PrefixExpr{
		Token: p.current,
		Op:    p.current.Literal,
	}

	p.advance()

	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *parser) parseInfixExpression(left ast.Expr) ast.Expr {
	expr := &ast.InfixExpr{
		Token: p.current,
		Op:    p.current.Literal,
		Left:  left,
	}

	precedence := p.precedenceAtCurrent()
	if isRightAssociativity(p.current) {
		precedence -= 1
	}
	p.advance()

	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *parser) parseGroupedExpression() ast.Expr {
	// Eat the '('
	p.advance()

	expr := p.parseExpression(LOWEST)

	if !p.consume(token.RPAREN, "") {
		return nil
	}

	return expr
}

func (p *parser) parseBlockExpression() ast.Expr {
	block := &ast.BlockExpr{
		Token:      p.current,
		Statements: []ast.Stmt{},
	}

	// Eat the '{'
	p.advance()

	for !p.currentIs(token.RCURLY) && !p.currentIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt == nil {
			p.advance()
			continue
		}

		if last, ok := stmt.(*ast.Temp); ok {
			block.LastValue = last.Value
			p.advance()
			return block
		}

		block.Statements = append(block.Statements, stmt)
		p.advance()
	}

	if !p.currentIs(token.RCURLY) {
		p.errorAtCurrent("missing closing brace after block")
		return nil
	}

	return block
}

func (p *parser) parseIfExpression() ast.Expr {
	expr := &ast.IfExpr{Token: p.current}

	// Eat the 'if'
	p.advance()
	expr.Cond = p.parseExpression(LOWEST)
	if expr.Cond == nil {
		p.errorAtCurrent("missing condition in if expression")
		return nil
	}

	if !p.consume(token.LCURLY, "") {
		return nil
	}

	body, ok := p.parseBlockExpression().(*ast.BlockExpr)
	if !ok {
		return nil
	}
	expr.Body = body

	if p.nextIs(token.ELSE) {
		p.advance()

		if p.nextIs(token.IF) {
			p.advance()
			expr.Branch = p.parseIfExpression()
			return expr
		}

		if !p.consume(token.LCURLY, "") {
			return nil
		}

		expr.Branch = p.parseBlockExpression()
	}
	return expr
}

func (p *parser) parseFunctionParameters() []*ast.Ident {
	idents := []*ast.Ident{}

	// fn() - no param
	if p.nextIs(token.RPAREN) {
		p.advance()
		return idents
	}

	p.advance()

	// parse the first param
	ident := &ast.Ident{Token: p.current, Value: p.current.Literal}
	idents = append(idents, ident)

	for p.nextIs(token.COMMA) {
		p.advance() // go to comma
		p.advance() // go to ident

		ident := &ast.Ident{Token: p.current, Value: p.current.Literal}
		idents = append(idents, ident)
	}

	if !p.consume(token.RPAREN, "") {
		return nil
	}

	return idents
}

func (p *parser) parseFunctionLiteral() ast.Expr {
	lit := &ast.FnLit{Token: p.current}

	if !p.consume(token.LPAREN, "") {
		return nil
	}

	lit.Params = p.parseFunctionParameters()

	if !p.consume(token.LCURLY, "") {
		return nil
	}

	body, ok := p.parseBlockExpression().(*ast.BlockExpr)
	if !ok {
		return nil
	}
	lit.Body = body

	return lit
}

func (p *parser) parseCallArguments() []ast.Expr {
	args := []ast.Expr{}

	// no arg
	if p.nextIs(token.RPAREN) {
		p.advance()
		return args
	}

	// first arg
	p.advance()
	args = append(args, p.parseExpression(LOWEST))

	for p.nextIs(token.COMMA) {
		p.advance() // go to comma
		p.advance() // go to next arg

		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.consume(token.RPAREN, "") {
		return nil
	}

	return args
}

func (p *parser) parseCallExpression(function ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Token: p.current,
		Func:  function,
		Args:  p.parseCallArguments(),
	}
}

func (p *parser) parseReturnExpression() ast.Expr {
	expr := &ast.ReturnExpr{Token: p.current}

	p.advance()

	expr.Value = p.parseExpression(LOWEST)

	return expr
}

func (p *parser) parseWhileExpression() ast.Expr {
	p.inloop++

	expr := &ast.WhileExpr{Token: p.current}

	p.advance()
	expr.Cond = p.parseExpression(LOWEST)
	if expr.Cond == nil {
		p.errorAtCurrent("missing condition in while expression")
		return nil
	}

	if !p.consume(token.LCURLY, "") {
		return nil
	}

	body, ok := p.parseBlockExpression().(*ast.BlockExpr)
	if !ok {
		return nil
	}
	expr.Body = body

	p.inloop--
	return expr
}

func (p *parser) parseBreakExpression() ast.Expr {
	if p.inloop == 0 {
		p.errorAtCurrent("cannot use 'break' outside of a loop")
		return nil
	}

	expr := &ast.BreakExpr{Token: p.current}

	p.advance()

	expr.Value = p.parseExpression(LOWEST)

	return expr
}

func (p *parser) parseContinueExpression() ast.Expr {
	if p.inloop == 0 {
		p.errorAtCurrent("cannot use 'continue' outside of a loop")
		return nil
	}

	return &ast.ContinueExpr{Token: p.current}
}
