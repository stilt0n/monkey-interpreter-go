package parser

import (
	"fmt"
	"monkey-pl/ast"
	"monkey-pl/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	lex := lexer.New(input)
	pars := New(lex)

	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	if program == nil {
		t.Fatalf("expected program not to be nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 8;
return 100;
return 928532;
`
	lex := lexer.New(input)
	pars := New(lex)

	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements. Got %d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("expected statement to be *ast.ReturnStatement. Got %T", statement)
			continue
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("expected statement.TokenLiteral() to return 'return'. Got %q", returnStatement.TokenLiteral())
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, expectedName string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("expected s.TokenLiteral() to return 'let'. Got '%q'", s.TokenLiteral())
		return false
	}
	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("expected statement to be *ast.LetStatement. Got %T", s)
		return false
	}

	if letStatement.Name.Value != expectedName {
		t.Errorf("expected letStatement.Name.Value to be %s. Got %s", expectedName, letStatement.Name.Value)
	}

	return true
}

func checkForParserErrors(t *testing.T, pars *Parser) {
	errors := pars.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parsed with %d errors", len(errors))
	for _, message := range errors {
		t.Errorf("parser error: %q", message)
	}
	t.FailNow()
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lex := lexer.New(input)
	pars := New(lex)
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected program to have 1 statement. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected statement to be an ExpressionStatement. Got %T", program.Statements[0])
	}

	ident, ok := statement.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("Expected statement's expression to be an Identifier. Got %T", statement.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("Expected statement value to be 'foobar'. Got '%s'", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("Expected TokenLiteral to be 'foobar'. Got '%s'", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	lex := lexer.New(input)
	pars := New(lex)
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected program to have 1 statement. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected statement to be an ExpressionStatement. Got %T", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Expected statement's expression to be an Integer Literal. Got %T", statement.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("Expected statement value to be 5. Got %d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("Expected TokenLiteral to be '5'. Got '%s'", literal.TokenLiteral())
	}
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		lex := lexer.New(tt.input)
		pars := New(lex)
		program := pars.ParseProgram()
		checkForParserErrors(t, pars)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected program to have 1 statement. Got %d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected statement to be an ExpressionStatement. Got %T", program.Statements[0])
		}

		expr, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Expected statement to be an ast.PrefixExrpession. Got %T", statement.Expression)
		}

		if expr.Operator != tt.operator {
			t.Fatalf("Expected expression operator to be '%s'. Got '%s'", tt.operator, expr.Operator)
		}

		if !testIntegerLiteral(t, expr.Right, tt.integerValue) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		lex := lexer.New(tt.input)
		pars := New(lex)
		program := pars.ParseProgram()
		checkForParserErrors(t, pars)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected program to have 1 statement. Got %d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected statement to be an ast.ExpressionStatement. Got %T", program.Statements[0])
		}

		expr, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("Expected expression to be an ast.InfixExpression. Got %T", statement.Expression)
		}

		if !testIntegerLiteral(t, expr.Left, tt.leftValue) {
			return
		}

		if expr.Operator != tt.operator {
			t.Fatalf("Expected statement operator to be '%s'. Got '%s'", tt.operator, expr.Operator)
		}

		if !testIntegerLiteral(t, expr.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		pars := New(lex)
		program := pars.ParseProgram()
		checkForParserErrors(t, pars)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expected expression to be type *ast.IntegerLiteral. Got %T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("Expected tested integer's value to be %d. Got %d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("Expected integer's token literal to be '%d'. Got '%s'", value, integer.TokenLiteral())
		return false
	}
	return true
}
