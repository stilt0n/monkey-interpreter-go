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
	
	# NOTE: The lexer does not check for valid syntax, just valid tokens
	# these tokens are valid even if the syntax is gibberish
	
	!-/*5;
	5 < 10 > 5;
	
	if (5 < 10) {
		return true;
	} else {
		return false;
	}

	10 == 10;
	10 != 9;

	"foobar";
	"foo bar";
	[1, 2];
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
		// !-/*;
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		// 5 < 10 > 5;
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		// if (5 < 10)
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		// { return true; }
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		// else { return false; }
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		// 10 == 10;
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		// 10 != 9;
		{token.INT, "10"},
		{token.NEQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		// "foobar";
		{token.STRING, "foobar"},
		{token.SEMICOLON, ";"},
		// "foo bar";
		{token.STRING, "foo bar"},
		{token.SEMICOLON, ";"},
		// [1, 2];
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		// EOF
		{token.EOF, ""},
	}

	lex := New(input)

	for i, testTok := range tests {
		tok := lex.NextToken()
		if tok.Type != testTok.expectedType {
			t.Logf("{ id: %d, type: %q, literal: %q }", i, tok.Type, tok.Literal)
			t.Fatalf("tests[%d] - tokentype wrong. Expected=%q, got=%q", i, testTok.expectedType, tok.Type)
		}

		if tok.Literal != testTok.expectedLiteral {
			t.Logf("{ id: %d, type: %q, literal: %q }", i, tok.Type, tok.Literal)
			t.Fatalf("tests[%d] - literal wrong. Expected=%q, got=%q", i, testTok.expectedLiteral, tok.Literal)
		}
	}
}
