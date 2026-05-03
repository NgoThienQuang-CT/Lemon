package lexer

import (
	"fmt"
	"lemon/token"
	"testing"
)

var testInput = `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

@

let result = add(five, ten);
!-/*5;
5 < 10 > 5 || 4 && 6;

if 5 < 10 {
	return true;
} else {
	return false;
}

let m = while let i = 0; i <= 5; i += 1 {
	if false {
		continue;
	} else if i >= 2 {
		break 0;
	}
}

// ankle hurt

10 == 10;
10 != 9;
10 % 2;

"foo";
"foo bar";
"lemon
`

func TestScanning(t *testing.T) {
	tests := []struct {
		expectedType    token.TokenType
		expectedLine    int
		expectedLiteral string
	}{
		{token.LET, 1, "let"},
		{token.IDENT, 1, "five"},
		{token.ASSIGN, 1, "="},
		{token.INT, 1, "5"},
		{token.SEMICOLON, 1, ";"},
		{token.LET, 2, "let"},
		{token.IDENT, 2, "ten"},
		{token.ASSIGN, 2, "="},
		{token.INT, 2, "10"},
		{token.SEMICOLON, 2, ";"},
		{token.LET, 4, "let"},
		{token.IDENT, 4, "add"},
		{token.ASSIGN, 4, "="},
		{token.FN, 4, "fn"},
		{token.LPAREN, 4, "("},
		{token.IDENT, 4, "x"},
		{token.COMMA, 4, ","},
		{token.IDENT, 4, "y"},
		{token.RPAREN, 4, ")"},
		{token.LCURLY, 4, "{"},
		{token.IDENT, 5, "x"},
		{token.PLUS, 5, "+"},
		{token.IDENT, 5, "y"},
		{token.SEMICOLON, 5, ";"},
		{token.RCURLY, 6, "}"},
		{token.SEMICOLON, 6, ";"},
		{token.ILLEGAL, 8, "Unexpect character '@'."},
		{token.LET, 10, "let"},
		{token.IDENT, 10, "result"},
		{token.ASSIGN, 10, "="},
		{token.IDENT, 10, "add"},
		{token.LPAREN, 10, "("},
		{token.IDENT, 10, "five"},
		{token.COMMA, 10, ","},
		{token.IDENT, 10, "ten"},
		{token.RPAREN, 10, ")"},
		{token.SEMICOLON, 10, ";"},
		{token.NOT, 11, "!"},
		{token.MINUS, 11, "-"},
		{token.SLASH, 11, "/"},
		{token.STAR, 11, "*"},
		{token.INT, 11, "5"},
		{token.SEMICOLON, 11, ";"},
		{token.INT, 12, "5"},
		{token.LSS, 12, "<"},
		{token.INT, 12, "10"},
		{token.GTR, 12, ">"},
		{token.INT, 12, "5"},
		{token.OR, 12, "||"},
		{token.INT, 12, "4"},
		{token.AND, 12, "&&"},
		{token.INT, 12, "6"},
		{token.SEMICOLON, 12, ";"},
		{token.IF, 14, "if"},
		{token.INT, 14, "5"},
		{token.LSS, 14, "<"},
		{token.INT, 14, "10"},
		{token.LCURLY, 14, "{"},
		{token.RETURN, 15, "return"},
		{token.TRUE, 15, "true"},
		{token.SEMICOLON, 15, ";"},
		{token.RCURLY, 16, "}"},
		{token.ELSE, 16, "else"},
		{token.LCURLY, 16, "{"},
		{token.RETURN, 17, "return"},
		{token.IDENT, 17, "false"},
		{token.SEMICOLON, 17, ";"},
		{token.RCURLY, 18, "}"},
		{token.LET, 20, "let"},
		{token.IDENT, 20, "m"},
		{token.ASSIGN, 20, "="},
		{token.WHILE, 20, "while"},
		{token.LET, 20, "let"},
		{token.IDENT, 20, "i"},
		{token.ASSIGN, 20, "="},
		{token.INT, 20, "0"},
		{token.SEMICOLON, 20, ";"},
		{token.IDENT, 20, "i"},
		{token.LEQ, 20, "<="},
		{token.INT, 20, "5"},
		{token.SEMICOLON, 20, ";"},
		{token.IDENT, 20, "i"},
		{token.PLUS, 20, "+"},
		{token.ASSIGN, 20, "="},
		{token.INT, 20, "1"},
		{token.LCURLY, 20, "{"},
		{token.IF, 21, "if"},
		{token.IDENT, 21, "false"},
		{token.LCURLY, 21, "{"},
		{token.CONTINUE, 22, "continue"},
		{token.SEMICOLON, 22, ";"},
		{token.RCURLY, 23, "}"},
		{token.ELSE, 23, "else"},
		{token.IF, 23, "if"},
		{token.IDENT, 23, "i"},
		{token.GEQ, 23, ">="},
		{token.INT, 23, "2"},
		{token.LCURLY, 23, "{"},
		{token.BREAK, 24, "break"},
		{token.INT, 24, "0"},
		{token.SEMICOLON, 24, ";"},
		{token.RCURLY, 25, "}"},
		{token.RCURLY, 26, "}"},
		{token.INT, 30, "10"},
		{token.EQL, 30, "=="},
		{token.INT, 30, "10"},
		{token.SEMICOLON, 30, ";"},
		{token.INT, 31, "10"},
		{token.NEQ, 31, "!="},
		{token.INT, 31, "9"},
		{token.SEMICOLON, 31, ";"},
		{token.INT, 32, "10"},
		{token.PERCENT, 32, "%"},
		{token.INT, 32, "2"},
		{token.SEMICOLON, 32, ";"},
		{token.STRING, 34, "foo"},
		{token.SEMICOLON, 34, ";"},
		{token.STRING, 35, "foo bar"},
		{token.SEMICOLON, 35, ";"},
		{token.ILLEGAL, 37, "Unterminated string."},
	}

	var l Lexer
	l.Init(testInput)

	for i, tt := range tests {
		name := fmt.Sprintf("test [%d]", i)
		t.Run(name, func(t *testing.T) {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Fatalf("token type wrong. Expected %q, got %q",
					tt.expectedType, tok.Type)
			}

			if tok.Line != tt.expectedLine {
				t.Fatalf("token line wrong. Expected %d, got %d",
					tt.expectedLine, tok.Line)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("token literal wrong. Expected %q, got %q",
					tt.expectedLiteral, tok.Literal)
			}
		})
	}
}
