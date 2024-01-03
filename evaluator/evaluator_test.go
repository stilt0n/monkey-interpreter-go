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
		{"-5", -5},
		{"-10", -10},
		{"--5", 5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 *2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == true", false},
		{"false == false", true},
		{"true != true", false},
		{"true != false", true},
		{"(5 > 3) == true", true},
		{"5 == true", false},
	}
	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testBooleanObject(t, evaluated, tc.expected)
	}
}

func TestBooleanOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testBooleanObject(t, evaluated, tc.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 1 }", 1},
		{"if (0) { 1 }", 1},
		{"if (1 > 2) { 100 }", nil},
		{"if (!(1 > 2)) { 100 }", 100},
		{"if (1 < 2) { 100 }", 100},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (true) { 10 } else { 20 }", 10},
	}
	for _, tc := range tests {
		evaluated := testEval(tc.input)
		output, ok := tc.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(output))
		} else {
			testNullObject(t, evaluated)
		}
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Expected object to be Boolean. Got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Expected object to have value %t. Got %t", expected, result.Value)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Expected returned object to be NULL. Got %T (%+v)", obj, obj)
		return false
	}
	return true
}
