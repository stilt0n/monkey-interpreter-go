package evaluator

import "monkey-pl/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. Expected 1. Got %d.", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", arg.Type())
			}
		},
	},
	"first": {Fn: first},
	"rest":  {Fn: rest},
	"last":  {Fn: last},
	"push":  {Fn: push},
}

func first(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. Expected 1. Got %d", len(args))
	}
	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("`first` encountered an unexpected type. `first` can only be called on an array. recieved first(%T)", args[0])
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
		return newError("`rest` encountered an unexpected type. `rest` can only be called on an array. recieved rest(%T)", args[0])
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
		return newError("`last` encountered an unexpected type. `last` can only be called on an array. recieved last(%T)", args[0])
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
		return newError("first argument to `push` should be an array. recieved rest(%T)", args[0])
	}
	item := args[1]
	elements := arr.Elements
	return &object.Array{Elements: append(elements, item)}
}
