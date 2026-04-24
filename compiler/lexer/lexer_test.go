package lexer

import (
	"lemon/compiler/token"
	"os"
	"testing"
)

func TestLexToken(t *testing.T) {
	tests := []struct {
		expectdLine   int
		expectdColumn int
		expectdType   token.TokenType
		expectdLexeme string
	}{
		{1, 1, token.LET, "let"},
		{1, 5, token.IDENT, "x"},
		{1, 7, token.ASSIGN, "="},
		{1, 9, token.INT, "10"},
		{1, 12, token.PLUS, "+"},
		{1, 14, token.FLOAT, "10.9"},
		{1, 18, token.SEMICOLON, ";"},
		{2, 1, token.LET, "let"},
		{2, 5, token.IDENT, "y"},
		{2, 7, token.ASSIGN, "="},
		{2, 9, token.STRING, "hello"},
		{2, 16, token.SEMICOLON, ";"},
		{3, 1, token.LPAREN, "("},
		{3, 2, token.LPAREN, "("},
		{3, 3, token.RPAREN, ")"},
		{3, 4, token.RBRACE, "}"},
		{3, 5, token.RPAREN, ")"},
		{3, 6, token.LBRACE, "{"},
		{3, 7, token.RBRACE, "}"},
		{3, 8, token.RPAREN, ")"},
		{3, 9, token.RBRACE, "}"},
		{3, 10, token.RPAREN, ")"},
		{3, 11, token.LPAREN, "("},
		{3, 12, token.LPAREN, "("},
		{8, 1, token.LET, "let"},
		{8, 5, token.IDENT, "c"},
		{8, 7, token.ASSIGN, "="},
		{8, 9, token.TRUE, "true"},
		{9, 1, token.ILLEGAL, "#"},
		{9, 2, token.ILLEGAL, "$"},
		{9, 3, token.ILLEGAL, "$"},
		{10, 1, token.INT, "5"},
		{10, 3, token.LSS, "<"},
		{10, 5, token.INT, "5"},
		{10, 7, token.GTR, ">"},
		{10, 9, token.INT, "5"},
		{10, 11, token.NEQ, "!="},
		{10, 14, token.INT, "5"},
		{10, 16, token.EQL, "=="},
		{10, 19, token.INT, "5"},
		{12, 1, token.LOOP, "loop"},
		{12, 6, token.BREAK, "break"},
		{12, 12, token.CONTINUE, "continue"},
		{12, 20, token.SEMICOLON, ";"},
		{13, 1, token.IF, "if"},
		{13, 4, token.TRUE, "true"},
		{13, 9, token.OR, "||"},
		{13, 12, token.FALSE, "false"},
		{13, 18, token.LBRACE, "{"},
		{14, 5, token.INT, "10"},
		{14, 8, token.GEQ, ">="},
		{14, 11, token.INT, "9"},
		{14, 13, token.LEQ, "<="},
		{14, 16, token.INT, "9"},
		{15, 1, token.RBRACE, "}"},
		{15, 3, token.ELSE, "else"},
		{15, 8, token.LBRACE, "{"},
		{16, 5, token.RETURN, "return"},
		{16, 12, token.IDENT, "b"},
		{16, 13, token.LPAREN, "("},
		{16, 14, token.IDENT, "x"},
		{16, 15, token.COMMA, ","},
		{16, 17, token.IDENT, "y"},
		{16, 18, token.COMMA, ","},
		{16, 20, token.FN, "fn"},
		{16, 22, token.LPAREN, "("},
		{16, 23, token.IDENT, "x"},
		{16, 24, token.COMMA, ","},
		{16, 26, token.IDENT, "y"},
		{16, 27, token.RPAREN, ")"},
		{16, 29, token.LBRACE, "{"},
		{16, 31, token.IDENT, "x"},
		{16, 33, token.PLUS, "+"},
		{16, 35, token.IDENT, "y"},
		{16, 37, token.RBRACE, "}"},
		{16, 39, token.RPAREN, ")"},
		{16, 40, token.SEMICOLON, ";"},
		{17, 1, token.RBRACE, "}"},
		{18, 1, token.EOF, ""},
	}

	filename := "../../test/random.lemon"
	source, _ := os.ReadFile(filename)
	l := New(filename, string(source))

	for i, test := range tests {
		tok := l.LexToken()

		if tok.Pos.Line != test.expectdLine {
			t.Fatalf("Test[%d] - token's line is wrong. Expected %d, got %d.",
				i, tok.Pos.Line, test.expectdLine)
		}

		if tok.Pos.Column != test.expectdColumn {
			t.Fatalf("Test[%d] - token's column is wrong. Expected %d, got %d.",
				i, tok.Pos.Column, test.expectdColumn)
		}

		if tok.Type != test.expectdType {
			t.Fatalf("Test[%d] - token's type is wrong. Expected %s, got %s.",
				i, tok.Type.String(), test.expectdType.String())
		}

		if tok.Lexeme != test.expectdLexeme {
			t.Fatalf("Test[%d] - token's lexeme is wrong. Expected %s, got %s.",
				i, tok.Lexeme, test.expectdLexeme)
		}
	}
}
