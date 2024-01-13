package evaluator

import (
	"fmt"
	"monkey-pl/object"
)

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
	"print":       {Fn: print},
	"first":       {Fn: first},
	"rest":        {Fn: rest},
	"last":        {Fn: last},
	"push":        {Fn: push},
	"join":        {Fn: join},
	"toUpperCase": {Fn: toUpperCase},
	"toLowerCase": {Fn: toLowerCase},
	"split":       {Fn: split},
}

// Called `puts` in the book.
func print(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return NULL
}
