package parser

import (
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