// Package lexer implements a lexer for Lemon source text.
// It takes a string as source which can then be tokenized
// through repeated calls to the NextToken method.
package lexer

import (
	"fmt"

	"lemon/internal/token"
)

type Lexer struct {
	src     string
	start   int // mark the beginning index of a token in the text source
	current int // a cursor that iterate the text source byte by byte
	line    int
}

func (l *Lexer) Init(source string) {
	l.src = source
	l.start = 0
	l.current = 0
	l.line = 1
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespaces()
	l.start = l.current

	if l.isAtEnd() {
		return l.makeToken(token.EOF)
	}

	char := l.advance()
	if isLetter(char) {
		return l.scanIdentifier()
	}
	if isDigit(char) {
		return l.scanNumber()
	}

	switch char {
	case '+':
		return l.makeToken(token.PLUS)
	case '-':
		return l.makeToken(token.MINUS)
	case '*':
		return l.makeToken(token.STAR)
	case '/':
		return l.makeToken(token.SLASH)
	case '%':
		return l.makeToken(token.PERCENT)
	case '<':
		if l.match('=') {
			return l.makeToken(token.LEQ)
		}
		return l.makeToken(token.LSS)
	case '>':
		if l.match('=') {
			return l.makeToken(token.GEQ)
		}
		return l.makeToken(token.GTR)
	case '=':
		if l.match('=') {
			return l.makeToken(token.EQL)
		}
		return l.makeToken(token.ASSIGN)
	case '!':
		if l.match('=') {
			return l.makeToken(token.NEQ)
		}
		return l.makeToken(token.NOT)
	case '|':
		if l.match('|') {
			return l.makeToken(token.OR)
		}
		goto illegal
	case '&':
		if l.match('&') {
			return l.makeToken(token.AND)
		}
		goto illegal
	case ',':
		return l.makeToken(token.COMMA)
	case '.':
		return l.makeToken(token.DOT)
	case ':':
		return l.makeToken(token.COLON)
	case ';':
		return l.makeToken(token.SEMICOLON)
	case '(':
		return l.makeToken(token.LPAREN)
	case ')':
		return l.makeToken(token.RPAREN)
	case '{':
		return l.makeToken(token.LCURLY)
	case '}':
		return l.makeToken(token.RCURLY)
	case '[':
		return l.makeToken(token.LBRACKET)
	case ']':
		return l.makeToken(token.RBRACKET)
	case '"':
		return l.scanString()
	default:
		break
	}
illegal:
	return l.makeError(fmt.Sprintf("unexpect character '%c'", char))
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.src)
}

func isLetter(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (l *Lexer) advance() byte {
	if l.isAtEnd() {
		return 0
	}

	l.current += 1
	return l.src[l.current-1]
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}

	return l.src[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.isAtEnd() {
		return 0
	}

	return l.src[l.current+1]
}

func (l *Lexer) match(expected byte) bool {
	if l.isAtEnd() || l.peek() != expected {
		return false
	}

	l.current += 1
	return true
}

func (l *Lexer) skipWhitespaces() {
	for {
		char := l.peek()
		switch char {
		case ' ', '\r', '\t':
			l.advance()
		case '\n':
			l.line++
			l.advance()
		case '/':
			if l.peekNext() != '/' {
				return
			}
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
		default:
			return
		}
	}
}

func (l *Lexer) makeToken(ttype token.TokenType) token.Token {
	var lit string
	if ttype == token.STRING {
		lit = l.src[l.start+1 : l.current-1]
	} else {
		lit = l.src[l.start:l.current]
	}

	return token.Token{
		Type:    ttype,
		Line:    l.line,
		Literal: lit,
	}
}

func (l *Lexer) makeError(message string) token.Token {
	return token.Token{
		Type:    token.ILLEGAL,
		Line:    l.line,
		Literal: message,
	}
}

func (l *Lexer) scanNumber() token.Token {
	for isDigit(l.peek()) {
		l.advance()
	}

	ttype := token.INT
	// Look for a fractional part
	if l.peek() == '.' && isDigit(l.peekNext()) {
		ttype = token.FLOAT
		// Eat the '.'
		l.advance()
		for isDigit(l.peek()) {
			l.advance()
		}
	}

	return l.makeToken(ttype)
}

func (l *Lexer) scanIdentifier() token.Token {
	for isLetter(l.peek()) || isDigit(l.peek()) {
		l.advance()
	}

	return l.makeToken(token.CheckKeyword(l.src[l.start:l.current]))
}

func (l *Lexer) scanString() token.Token {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}

	if l.isAtEnd() {
		return l.makeError("unterminated string")
	}

	// Eat the closing quote.
	l.advance()
	return l.makeToken(token.STRING)
}
