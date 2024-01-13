package evaluator

import (
	"monkey-pl/object"
	"strings"
)

func first(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. Expected 1. Got %d", len(args))
	}
	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("`first` encountered an unexpected type. `first` can only be called on an array. received first(%s)", args[0].Type())
	}

	elements := arr.Elements
	if len(elements) == 0 {
		return NULL
	}
	return elements[0]
}

func rest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. Expected 1. Got %d.", len(args))
	}
	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("`rest` encountered an unexpected type. `rest` can only be called on an array. received rest(%s)", args[0].Type())
	}
	elements := arr.Elements
	if len(elements) < 2 {
		return NULL
	}
	return &object.Array{Elements: elements[1:]}
}

func last(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. Expected 1. Got %d.", len(args))
	}
	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("`last` encountered an unexpected type. `last` can only be called on an array. received last(%s)", args[0].Type())
	}
	elements := arr.Elements
	if len(elements) == 0 {
		return NULL
	}
	return elements[len(elements)-1]
}

func push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. Expected 2. Got %d.", len(args))
	}
	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("first argument to `push` should be an array. received rest(%s)", args[0].Type())
	}
	item := args[1]
	elements := arr.Elements
	return &object.Array{Elements: append(elements, item)}
}

func join(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. Expected 2. Got %d", len(args))
	}
	arr, ok1 := args[0].(*object.Array)
	separator, ok2 := args[1].(*object.String)
	if !ok1 || !ok2 {
		return newError("`join` expected arguments of type join(array, string). received join(%s, %s)", args[0].Type(), args[1].Type())
	}
	elements := arr.Elements
	stringArr := []string{}
	for _, el := range elements {
		stringElement, ok := el.(*object.String)
		if !ok {
			return newError("can only join an array that is all strings")
		}
		stringArr = append(stringArr, stringElement.Value)
	}
	return &object.String{Value: strings.Join(stringArr, separator.Value)}
}
