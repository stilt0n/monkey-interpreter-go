# Notes on what I'm learning

## Interpreters vs Compilers

Basics are review but worht rehashing:

- Interpreter are able to evaluate source code without needing to create an intermediate result
- Compilers turn source code into an intermediate result which can later be executed

There are a wide variety of interpreters.

- Simple interpreters that don't bother parsing
- Interpreters that parse, build AST and evaluate (tree-walking interpreter)
- Interpreters that use JIT compilation to compile input into machine code (browser JavaScript compilers do this with a few different compilers)

This interpreter will be a tree walking interpreter

## On Monkey

PL designed for the book

Taken from book:

- C-like syntax
- Variable bindings
- ints and bools
- arithmetic expressions
- built-in functions
- first-class and higher-order functions
- closures
- string data structure
- array data structure
- hash data structure

### Variable binding

```
let age = 1;

let name = "Monkey";

let result = 10 * (20/2);
```

### Data types and structures

Suports ints, bools and strings.
Arrays:

```
let arr = [1,2,3,4,5];
```

Hash maps:

```
let map = { "key": "value", "num": 80 };
```

Access:

```
arr[0]; // => 1

map["key"]; // => "value"
```

Can also bind functions to names with `let`

```
let add = fn(a,b) { return a + b };
```

We'll also support implicit returns

```
let add = fn(a,b) { a + b }
```

Function calls:

```
add(1,2);
```

recursion:

```
// Obviously assumes x is not negative
let fib = fn(x) {
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      1
    } else {
      fib(x - 1) + fib(x - 2)
    }
  }
};
```

Higher order functions (functions that take other funcs as args)

```
let twice = fn(f, x) {
  return f(f(x));
};

let addTwo = fn(x) {
  return x + 2;
}

twice(addTwo, 2); // => 6
```

## High level parts

- Lexer
- Parser
- Abstract Syntax Tree
- internal object system
- evaluator

Will build these in order.
