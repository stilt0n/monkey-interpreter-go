package evaluator

import (
	"fmt"
	"monkey-pl/ast"
	"monkey-pl/object"
)

// for boolean literals there's no reason to recreate them
// when we can just point to the same object
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return objectFromBool(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	}
	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	// Note: these pointer comparisons work now because ints get evaluated above
	// but if we add more datatypes (e.g. strings) we may need to change this
	case operator == "==":
		return objectFromBool(left == right)
	case operator == "!=":
		return objectFromBool(left != right)
	// This check takes place after equality check because equality checking
	// two objects of different types is legal (but always false).
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	lval := left.(*object.Integer).Value
	rval := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: lval + rval}
	case "-":
		return &object.Integer{Value: lval - rval}
	case "*":
		return &object.Integer{Value: lval * rval}
	case "/":
		if rval == 0 {
			return newError("illegal operation: divide by zero")
		}
		return &object.Integer{Value: lval / rval}
	case "<":
		return objectFromBool(lval < rval)
	case ">":
		return objectFromBool(lval > rval)
	case "==":
		return objectFromBool(lval == rval)
	case "!=":
		return objectFromBool(lval != rval)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(expr *ast.IfExpression) object.Object {
	condition := Eval(expr.Condition)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(expr.Consequence)
	} else if expr.Alternative != nil {
		return Eval(expr.Alternative)
	} else {
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		// Note: this means that in Monkey `0` is NOT falsey
		return FALSE
	}
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func objectFromBool(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func isTruthy(obj object.Object) bool {
	if obj == NULL || obj == FALSE {
		return false
	}
	return true
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
