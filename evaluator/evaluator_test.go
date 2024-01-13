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
		{`"whoa" == "whoa"`, true},
		{`"whoa" == "wow"`, false},
		{`"whoa" != "wow"`, true},
		{`"a" < "b"`, true},
		{`"a" > "b"`, false},
		{`"abc" == -"cba"`, true},
	}
	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testBooleanObject(t, evaluated, tc.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	input := `"Hello, World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("Expected object to be a string. Got %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello, World!" {
		t.Errorf("Expected string to have value 'Hello, World!'. Got %s", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + ", " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("Expected object to be a string. Got %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello, World!" {
		t.Errorf("Expected string to have value 'Hello, World!'. Got %s", str.Value)
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

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5;, 9;", 10},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}`, 10,
		},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true - false; 5;",
			"unknown operator: BOOLEAN - BOOLEAN",
		},
		{
			"if (10 > 1) { true * false; }",
			"unknown operator: BOOLEAN * BOOLEAN",
		},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return 10 / true;
				}
				return 1;
			}`, "type mismatch: INTEGER / BOOLEAN",
		},
		{
			"2 / 0",
			"illegal operation: divide by zero",
		},
		{
			"2 / (5 - 5)",
			"illegal operation: divide by zero",
		},
		{
			"foobar;",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
		{
			`let f = fn(x) { x; }; { f: "Monkey" };`,
			"unhashable object used as a hash key: FUNCTION",
		},
		{
			`let func = fn(x) { func(x + 1); }; func(1);`,
			"maximum stack depth exceeded",
		},
		{
			`while (true) { x; };`,
			"maximum iteration count exceeded",
		},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("Expected an error object to be returned. Got %T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tc.expected {
			t.Errorf("Expected error message to be %q. Got %q", tc.expected, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let x = 5 * 5; x;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tc := range tests {
		testIntegerObject(t, testEval(tc.input), tc.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("Expected a function object to be returned. Go %T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("Expected function to have 1 parameters. Got %d", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("Expected function param to be 'x'. Got %q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("Expected function body to be %q. Got %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let double = fn(x) { x * 2; }; double(15);", 30},
		{"fn(x) { x; }(1);", 1},
		{"let add = fn(x, y) { x + y }; add(1, add(2, 3));", 6},
	}
	for _, tc := range tests {
		testIntegerObject(t, testEval(tc.input), tc.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) { x + y };
};
let addTwo = newAdder(2)
addTwo(2);
`
	testIntegerObject(t, testEval(input), 4)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("abc")`, 3},
		{`len("ab c")`, 4},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. Expected 1. Got 2."},
		{`len([1, 2, 3]);`, 3},
		{`len(["alpha", "beta", "gamma"])`, 3},
		{`len([])`, 0},
	}
	for _, tc := range tests {
		evaluated := testEval(tc.input)
		switch expected := tc.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("expected Error object. Got %T (%+v)", evaluated, evaluated)
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. Expected %q. Got %q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	results, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("Expected object to be an Array. Got %T (%+v)", evaluated, evaluated)
	}

	if len(results.Elements) != 3 {
		t.Fatalf("Expected array to have 3 elements. Got %d", len(results.Elements))
	}

	testIntegerObject(t, results.Elements[0], 1)
	testIntegerObject(t, results.Elements[1], 4)
	testIntegerObject(t, results.Elements[2], 6)
}

func TestArrayINdexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1,2,3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let arr = [1, 2, 3]; arr[2];",
			3,
		},
		{
			"let arr = [1, 2, 3]; arr[0] + arr[1] + arr[2];",
			6,
		},
		{
			"let arr = [1, 2, 3]; let i = arr[0]; arr[i];",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}
	for _, tc := range tests {
		evaluated := testEval(tc.input)
		integer, ok := tc.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}
	`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. Got %T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Expected hash to have 6 pairs. Got %d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
			testIntegerObject(t, pair.Value, expectedValue)
		}
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		integer, ok := tc.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestWhileExpression(t *testing.T) {
	input := `let i = 0; while (i < 10) { let i = i + 1; }; i;`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("expected result to be of type *object.Integer. Got %T instead.", evaluated)
	}

	if !testIntegerObject(t, result, 10) {
		return
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()

	return Eval(program, object.NewEnvironment())
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
