package parser

import (
	"fmt"
	"monkey-pl/ast"
	"monkey-pl/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y", "foobar", "y"},
	}
	for _, tt := range tests {
		pars := New(lexer.New(tt.input))
		program := pars.ParseProgram()
		checkForParserErrors(t, pars)
		if len(program.Statements) != 1 {
			t.Fatalf("Expected 1 statement. Got %d", len(program.Statements))
		}
		statement := program.Statements[0]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
		val := statement.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input      string
		expectedRV interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}
	for _, tt := range tests {
		pars := New(lexer.New(tt.input))
		program := pars.ParseProgram()
		checkForParserErrors(t, pars)
		if len(program.Statements) != 1 {
			t.Fatalf("Expected 1 statement. Got %d", len(program.Statements))
		}
		statement := program.Statements[0]
		if statement.TokenLiteral() != "return" {
			t.Fatalf("Expected return statement's token literal to be 'return'. Got %s", statement.TokenLiteral())
		}
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("Expected return statement's type to be *ast.ReturnStatement. Got %T", statement)
		}

		if !testLiteralExpression(t, returnStatement.ReturnValue, tt.expectedRV) {
			return
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

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
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

	literal, ok := statement.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("Expected statement's expression to be a StringLiteral. Got %T", statement.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("Expected statement value to be 'hello world'. Got %s", literal.Value)
	}

	if literal.TokenLiteral() != "hello world" {
		t.Errorf("Expected TokenLiteral to be 'hello world'. Got %s", literal.TokenLiteral())
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
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
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

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"
	lex := lexer.New(input)
	pars := New(lex)
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected program to have 1 statement. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected program statement to be an expression statement. Got %T", program.Statements[0])
	}

	expr, ok := statement.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("Expected statement expression to be of type *ast.IfExpression. Got %T", statement.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}

	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("Expected consequence to have 1 statement. Got %d", len(expr.Consequence.Statements))
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Expected consequence statement to be of type *ast.ExpressionStatement. Got %T", expr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expr.Alternative != nil {
		t.Errorf("Expected not to get an Alternative expression. Got %+v", expr.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"
	lex := lexer.New(input)
	pars := New(lex)
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected program to have 1 statement. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected program statement to be an expression statement. Got %T", program.Statements[0])
	}

	expr, ok := statement.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("Expected statement expression to be of type *ast.IfExpression. Got %T", statement.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}

	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("Expected consequence to have 1 statement. Got %d", len(expr.Consequence.Statements))
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Expected consequence statement to be of type *ast.ExpressionStatement. Got %T", expr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expr.Alternative == nil {
		t.Fatal("Expected to get an Alternative expression. Got nil")
	}

	alternative, ok := expr.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Expected alternative statement to be of type *ast.ExpressionStatement. Got %T", expr.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestWhileExpression(t *testing.T) {
	input := "while (x < y) { x; }"
	lex := lexer.New(input)
	pars := New(lex)
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected program to have 1 statement. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected program statement to be an expression statement. Got %T", program.Statements[0])
	}

	expr, ok := statement.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("Expected statement expression to be of type *ast.WhileExpression. Got %T", statement.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}

	if len(expr.Body.Statements) != 1 {
		t.Errorf("Expected loop body to have 1 statement. Got %d", len(expr.Body.Statements))
	}

	body, ok := expr.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected body statement to be of type *ast.ExpressionStatement")
	}

	if !testIdentifier(t, body.Expression, "x") {
		return
	}
}

func TestFunctionLiteralExpression(t *testing.T) {
	input := "fn(x, y) { x + y; }"
	lex := lexer.New(input)
	pars := New(lex)
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected program to have 1 statement. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected statement to be of type *ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expected statement expression to be of type ast.FunctionLiteral. Got %T", statement.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("Expected function literal to have 2 parameters. Got %d", len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("Expected function body to have 1 statement. Got %d", len(function.Body.Statements))
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected function body statement to be of type ast.ExpressionStatement. Got %T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		pars := New(lexer.New(tt.input))
		program := pars.ParseProgram()
		checkForParserErrors(t, pars)
		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("Expected function to have %d parameters. Got %d", len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	pars := New(lexer.New(input))
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected program to have 1 statement. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Expected program statement to be an expression statement. Got %T", program.Statements[0])
	}

	expr, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Expected statement expression to be of type *ast.CallExpression. Got %T", statement.Expression)
	}

	if !testIdentifier(t, expr.Function, "add") {
		return
	}

	if len(expr.Arguments) != 3 {
		t.Fatalf("Expected function to be called with 3 arguments. Got %d", len(expr.Arguments))
	}

	testLiteralExpression(t, expr.Arguments[0], 1)
	testInfixExpression(t, expr.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expr.Arguments[2], 4, "+", 5)
}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	pars := New(lexer.New(input))
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected program statement to be an expression statement. Got %T", program.Statements[0])
	}
	array, ok := statement.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expected an ArrayLiteral. Got %T", statement.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("Expected test array to have length 3. Got %d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestIndexExpression(t *testing.T) {
	input := "arr[1 + 1];"
	pars := New(lexer.New(input))
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	indexExpr, ok := statement.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("Expected expression to be an IndexExpression. Got %T", statement.Expression)
	}

	if !testIdentifier(t, indexExpr.Left, "arr") {
		return
	}

	if !testInfixExpression(t, indexExpr.Index, 1, "+", 1) {
		return
	}
}

func TestHashLiteralStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	lex := lexer.New(input)
	pars := New(lex)
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := statement.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Expected expression to be a HashLiteral. Got %T", statement.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. Expected 3. Got %d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("Expected key to be a StringLiteral. Got %T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParseEmptyHash(t *testing.T) {
	input := "{}"
	pars := New(lexer.New(input))
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := statement.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Expected expression to be a HashLiteral. Got %T", statement.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. Expected 0. Got %d", len(hash.Pairs))
	}
}

func TestHashLiteralWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8}`
	pars := New(lexer.New(input))
	program := pars.ParseProgram()
	checkForParserErrors(t, pars)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := statement.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Expected expression to be a HashLiteral. Got %T", statement.Expression)
	}
	if len(hash.Pairs) != 2 {
		t.Errorf("hash.Pairs has wrong length. Expected 2. Got %d", len(hash.Pairs))
	}
	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("Expected k to be a StringLiteral. Got %T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
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
