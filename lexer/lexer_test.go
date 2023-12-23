package lexer

import (
	"monkey-pl/token"
	"testing"
)

func SimpleTest(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lex := New(input)

	for i, testToken := range tests {
		tok := lex.NextToken()
		if tok.Type != testToken.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. Expected=%q, got=%q", i, testToken.expectedType, tok.Type)
		}

		if tok.Literal != testToken.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. Expected=%q, got=%q", i, testToken.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	# this line is a comment and should be ignored
	let ten = 10;

	let add = fn(x, y) {
		x + y; # this adds x and y and returns the result implicity
	};

	let result = add(five, ten);
	# Will a comment work at the end??
	`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// let five = 5;
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		// let ten = 10;
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		// let add = fn(x,y) { x + y; };
		// let add = fn(x,y)
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		// { x + y; };
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		// let result = add(five, ten)
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		// EOF
		{token.EOF, ""},
	}

	lex := New(input)

	for i, testTok := range tests {
		tok := lex.NextToken()
		if tok.Type != testTok.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. Expected=%q, got=%q", i, testTok.expectedType, tok.Type)
		}

		if tok.Literal != testTok.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. Expected=%q, got=%q", i, testTok.expectedLiteral, tok.Literal)
		}
	}
}
