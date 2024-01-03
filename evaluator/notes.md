# Evaluation

We need to evaluate a program for it to become meaningful.

The evaluation process defines how the programming language
works.

e.g.:

```js
let num = 5;
if (num) {
  return a;
} else {
  return b;
}
```

We need to decided if 5 is truthy or falsey (or if this is an
error). Another example:

```js
let foo = fn() { return 1; }
let bar = fn() { return 2; }
baz(foo(), bar());
```

Which function gets called first in baz's args?

In most C like languages this is going to be foo. But we
could choose to do it differently.

## Evaluation strategies

Evaluation is where programming language implementations
diverge the most.

While a simple mental model for an interpreter is a program
that parses and executes code, many real-world interpreters
aren't this simple. (e.g. JavaScript uses a mix of this and
a variety of more or less optimized compilers that get
involved when certain heuristics are triggered. And even in
JavaScript the interpreter will be engine specific).

The simplest approach is to just walk the tree and do the
stuff it says. This is called a tree-walking interpreter.

Even with tree-walking interpreters there can be optimizations
or conversions to intermediate representations.

Interpreters can also traverse the AST and convert it to
Bytecode. Bytecode is an intermediate representation (IR)
of an AST.

The instructions that make up Bytecode are called opcodes.

Opcodes depend some on host language (in this case Go) but
they're generally going to be similar to assembly.

Common opcode instructions:

- push
- pop

The main difference between bytecode and assembly is that
bytecode is not native and runs on a virtual machine.

Virtual machines can have performance benefits over just
evaluating the AST and they are more portable than native code.

Not all bytecode setups involve an AST. The book points out
that this makes the line between interpreting and compiling
fuzzy.

Languages can also compile bytecode right before execution
which is called JIT (just in time) compilation. (I think
JavaScript engines generally do at least some of this).

But a language can also skip this and JIT compile an AST.

Finally, you can also traverse an AST and do JIT only when
a branch is reused a certain number of times.

### Tradeoffs in evaluation strategies

Mainly we're looking at:

- Implementation complexity
- Performance
- Portability

Simple tree-walking is the most straightforward to implement
and if the language it's built on top of is highly portable
then the tree-walking compiler will be as well. But it's also
much slower than running bytecode on a VM or running native
assembly instructions.

Compilation to native code is faster than bytecode but is
not as portable because it requires supporting multiple
machine architectures. An interpreter that uses JIT compilation
will be much more complex to make.

Languages don't need to take a fixed approach:

WebKit JavaScript engine used AST walking and direct execution
initially and switched to bytecode in 2008. WebKit now has
4 stages of JIT compilation.

Ruby also originally used tree-walking and switched to bytecode
which led to large performance gains.

For mental model building, a tree-walking interpreter is a good
approach.

## Tree-Walking Interpreter

For this project we're going to interpret directly from the
Abstract Syntax Tree rather than doing any kind of
preprocessing or compilation. The interpreter is similar to
the classic Lisp interpreter from ["The Structure and
Interpretation of Computer Programs"](https://www.amazon.com/Structure-Interpretation-Computer-Programs-Engineering/dp/0262510871)

Need to implement two things:

- Tree-Walking evaluator
- Representation of Monkey values to Go

### Psuedocode

```js
function eval(astNode) {
  if (astNode is integerliteral) {
    return astNode.integerValue
  } else if (astNode is booleanLiteral) {
    return astNode.booleanValue
  } else if (astNode is infixExpression) {
    leftEvaluated = eval(astNode.Left)
    rightEvalauted = eval(astNode.Right);
    if astNode.Operator == "+" {
      return leftEvaluated + rightEvaluated
    } else if ast.Operator == "-" {
      return leftEvaluated + rightEvaluated
    }
    // ... etc.
  }
}
```

## Representing Objects

We need to figure out what to return from eval. We need to
solve this problem:

```js
let a = 5;
// ... a bunch of other code
// How do we keep track of what `a` is in a reasonable way?
a + a;
```

Multiple ways to do this:

- Use host language datatypes directly
- Use pointers
- Mix native types and pointers
- etc.

Considerations:

- What data types does the host language support?
- What data types does the interpreted language need?
- What performance requirements does the interpreted language have?
- Are you implementing garbage collection?
- What does public API of these values look like

Other languages:

Java has primitive data types for stuff like ints, longs,
floats, booleans, chars, etc. and reference types for stuff
like compound data structures.

In Ruby everything is an object and there is no primitive
data type.

## Monkey Object system

Monkeys don't care a bunch about performance so we'll just
represent everything as an object.

Object will be an interface instead of a struct to allow
more flexibility in implementing different types of objects

### Further reading on Objects

Book recommends [Wren source code](https://github.com/wren-lang/wren)

## Self evaluating expressions

These are expressions where what you type into the
repl should be the output you receive.

### Example

Given an `*ast.IntegerLiteral` `Eval` should return an
`*object.Integer` with `Value=*ast.IntegerLiteral.Value`

## Design decisions

- Using truthy and falsey values
- In Monkey:

  - FALSE and NULL are falsey
  - All other values (including 0) are truthy

- Return statements can be used top level
- After return is called, no further statements in block run
  - The "block" here can also refer to the top level program

e.g. if we feed the interpreter:

```js
5 + 6 + 7;
return true;
// interpreter ignores this
8 * 2;
```

The last statement will not run.

## Approach to return statements

Pass a `return value` through evaluator.
When return is encountered wrap value that needs
to be returned inside object.
Then decide if we need to stop eval.

Since return values can be nested multiple block statements
deep, we will need to change our approach to evaluating blocks.

Now the statements slice of a program will need to be handled
differently than the statements slice of a block.

## Handling errors

Errors are handled in a similar fashion to return statements.

Errors, like return statements, need to stop evalution of
a series of statements.

### Stack traces

This is sort of an aside. For Monkey we are not implementing
a stack trace. But this would be possible if we had attached
info about line / column number during lexing.
