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

func TestBooleanLiteralExpression(t *testing.T) {
	lex := lexer.New("true;")
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

	testLiteralExpression(t, statement.Expression, true)
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"-false;", "-", false},
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

		if !testLiteralExpression(t, expr.Right, tt.value) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"true == false;", true, "==", false},
		{"true + 5;", true, "+", 5},
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

		testInfixExpression(t, statement.Expression, tt.leftValue, tt.operator, tt.rightValue)
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
		{
			"4 > 2 == true",
			"((4 > 2) == true)",
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

func testBooleanLiteral(t *testing.T, expr ast.Expression, v bool) bool {
	boolean, ok := expr.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("Expected expression to be type *ast.BooleanLiteral. Got %T", expr)
		return false
	}

	if boolean.Value != v {
		t.Errorf("Expected tested boolean's value to be %t. Got %t", v, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", v) {
		t.Errorf("Expected boolean's token literal to be '%t'. Got '%s'", v, boolean.TokenLiteral())
	}
	return true
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("Expected expression to be an Identifier. Got %T", expr)
		return false
	}

	if ident.Value != value {
		t.Errorf("Expected identifier value to be %s. Got %s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("Expected identifier TokenLiteral to be %s. Got %s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}
	t.Errorf("testLiteralExpression received an unhandled expression type. Received: %T", expr)
	return false
}

func testInfixExpression(t *testing.T, expr ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Expected expression to be an Infix Expression. Got %T(%s)", expr, expr)
		return false
	}

	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}

	if opExpr.Operator != operator {
		t.Errorf("Expected operator to be '%s'. Got %q", operator, opExpr.Operator)
		return false
	}

	if !testLiteralExpression(t, opExpr.Right, right) {
		return false
	}
	return true
}
