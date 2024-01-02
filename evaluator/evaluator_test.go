package evaluator

import (
	"monkey-pl/lexer"
	"monkey-pl/object"
	"monkey-pl/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Expected object to be Integer. Got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Expected object to have value %d. Got %d", expected, result.Value)
		return false
	}

	return true
}
