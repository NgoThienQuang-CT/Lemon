// Package lexer implements a lexer for Lemon programming language
// The lexer will take a string of source and tokenized it
package lexer

import "lemon/compiler/token"

type Lexer struct {
	source   string
	filename string

	line   int
	column int

	start   int
	current int
}

func New(filename string, source string) *Lexer {
	return &Lexer{
		source:   source,
		filename: filename,

		line:   1,
		column: 1,

		start:   0,
		current: 0,
	}
}

func (l *Lexer) LexToken() token.Token {
	l.skipWhitespaces()
	l.start = l.current

	if l.isAtEnd() {
		return l.makeToken(token.EOF)
	}

	c := l.advance()
	if isLetter(c) {
		return l.scanIdentifier()
	}

	if isDigit(c) {
		return l.scanNumber()
	}

	switch c {
	// one character symbol
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
	case ',':
		return l.makeToken(token.COMMA)
	case '.':
		return l.makeToken(token.PERIOD)
	case ';':
		return l.makeToken(token.SEMICOLON)
	case '(':
		return l.makeToken(token.LPAREN)
	case ')':
		return l.makeToken(token.RPAREN)
	case '{':
		return l.makeToken(token.LBRACE)
	case '}':
		return l.makeToken(token.RBRACE)
	// two characters symbol
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
	case '&':
		if l.match('&') {
			return l.makeToken(token.AND)
		}
		return l.makeToken(token.ILLEGAL)
	case '|':
		if l.match('|') {
			return l.makeToken(token.OR)
		}
		return l.makeToken(token.ILLEGAL)
	// string
	case '"':
		return l.scanString()
	}

	return l.makeToken(token.ILLEGAL)
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current+1]
}

func (l *Lexer) advance() byte {
	if l.isAtEnd() {
		return 0
	}

	if l.peek() == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	l.current++

	return l.source[l.current-1]
}

func (l *Lexer) match(expected byte) bool {
	if l.isAtEnd() || l.peek() != expected {
		return false
	}
	l.advance()
	return true
}

func (l *Lexer) skipWhitespaces() {
	for {
		c := l.peek()
		switch c {
		case ' ', '\r', '\t', '\n':
			l.advance()

		// skip comment
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

func (l *Lexer) makeToken(tokenType token.TokenType) token.Token {
	return token.Token{
		Type: tokenType,
		Pos: token.Position{
			Filename: l.filename,
			Offset:   l.start,
			Line:     l.line,
			Column:   l.column - l.current + l.start,
		},
		Lexeme: l.source[l.start:l.current],
	}
}

func (l *Lexer) scanIdentifier() token.Token {
	for isLetter(l.peek()) {
		l.advance()
	}
	return l.makeToken(token.KeywordsLookup(l.source[l.start:l.current]))
}

func (l *Lexer) scanNumber() token.Token {
	for isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() != '.' && !isDigit(l.peekNext()) {
		return l.makeToken(token.INT)
	}

	// eat the period char
	l.advance()
	for isDigit(l.peek()) {
		l.advance()
	}

	return l.makeToken(token.FLOAT)
}

func (l *Lexer) scanString() token.Token {
	for l.peek() != '"' && !l.isAtEnd() {
		l.advance()
	}

	// eat the close quote
	l.advance()

	tok := l.makeToken(token.STRING)
	tok.Lexeme = l.source[l.start+1 : l.current-1]
	return tok
}
