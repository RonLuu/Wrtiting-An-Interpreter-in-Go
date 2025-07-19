package lexer

import (
	"Chapter_2/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	// Given an input for testing
	input := `let five = 5;
let ten = 10;
let add = fn(x, y) {
x + y;
};
let result = add(five, ten);
!-/*5;
5 < 10 > 5;
if (5 < 10) {
return true;
} else {
return false;
}
10 == 10;
10 != 9;
`

	// The test should check as followed
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.VARIABLE, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.VARIABLE, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.VARIABLE, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.RLBRACKET, "("},
		{token.VARIABLE, "x"},
		{token.COMMA, ","},
		{token.VARIABLE, "y"},
		{token.RRBRACKET, ")"},
		{token.PLBRACKET, "{"},

		{token.VARIABLE, "x"},
		{token.PLUS, "+"},
		{token.VARIABLE, "y"},
		{token.SEMICOLON, ";"},

		{token.PRBRACKET, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.VARIABLE, "result"},
		{token.ASSIGN, "="},
		{token.VARIABLE, "add"},
		{token.RLBRACKET, "("},
		{token.VARIABLE, "five"},
		{token.COMMA, ","},
		{token.VARIABLE, "ten"},
		{token.RRBRACKET, ")"},
		{token.SEMICOLON, ";"},

		{token.EXCLAMATION, "!"},
		{token.MINUS, "-"},
		{token.DIV, "/"},
		{token.MULT, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.RLBRACKET, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RRBRACKET, ")"},
		{token.PLBRACKET, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.PRBRACKET, "}"},
		{token.ELSE, "else"},
		{token.PLBRACKET, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.PRBRACKET, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.NEQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	// Get a lexer from the input
	l := NewLexer(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("Failed at [%d] - wrong token type, expected %q, got %q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("Failed at [%d] - wrong literal, expected %q, got %q", i, tt.expectedLiteral, tok.Literal)
		}

	}
}
