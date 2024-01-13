package evaluator

import (
	"monkey-pl/object"
	"strings"
)

func toUpperCase(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. Expected 1. Got %d", len(args))
	}
	str, ok := args[0].(*object.String)
	if !ok {
		return newError("`toUpperCase` encountered unexpected type. `toUpperCase` can only be called on a string. received toUpperCase(%s)", args[0].Type())
	}
	return &object.String{Value: strings.ToUpper(str.Value)}
}

func toLowerCase(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. Expected 1. Got %d", len(args))
	}
	str, ok := args[0].(*object.String)
	if !ok {
		return newError("`toLowerCase` encountered unexpected type. `toLowerCase` can only be called on a string. received toLowerCase(%s)", args[0].Type())
	}
	return &object.String{Value: strings.ToLower(str.Value)}
}

func split(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. Expected 2. Got %d", len(args))
	}
	str, ok1 := args[0].(*object.String)
	separator, ok2 := args[1].(*object.String)
	if !ok1 || !ok2 {
		return newError("`split` expects arguments of type split(string, string). receieved split(%s, %s)", args[0].Type(), args[1].Type())
	}
	strArr := strings.Split(str.Value, separator.Value)
	objArr := []object.Object{}
	for _, s := range strArr {
		objArr = append(objArr, &object.String{Value: s})
	}
	return &object.Array{Elements: objArr}
}
