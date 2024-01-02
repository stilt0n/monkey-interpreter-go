# Parsers

Parsers take input data and build a create a data structure out of it while checking for correct syntax.

Often a parser is going to read tokens from a lexer as input data in this context.

- Data structure is typically going to be heirarchical (e.g. Trees)
- Common structures are parse trees and abstract syntax trees

## Example from book:

Consider this JavaScript:

```js
var input = '{"name": "Sir", "age": "39" }';
var output = JSON.parse(input);
console.log(output); // => { name: 'Sir', age: 28 }
console.log(output.name); // => 'Sir'
console.log(output.age); // => 39
```

In this case:

- Input is text
- Output is JavaScript Object with two fields

While this is way simpler than for a programming language it is conceptually very similar. `JSON.parse` took text and transformed it into a heirachical data structure. The data structure in the case of JSON is easy to see at a glance, while in a programming language it's less apparent.

## Abstract Syntax Trees

Most interpreters and compilers use an abstract syntax tree to represent source code.

Abstract syntax trees are called "abstract" because they don't encode all of the concrete syntax (i.e. `let five = 5;`) into the tree. At this step we can usually avoid representing stuff like whitespace, semicolons, brackets, etc. since those encode information about how to construct the tree but not about what code needs to be executed.

## Parser Generators

There are parser generators which can generate parsers for you. I'm doing this to learn how they work so it's obviously not appropriate here but this seems like a useful thing to be aware of in the future.

## Parser strategies

The parsing strategy here is a top-down parsing strategy, but bottom=up strategies exist as well. The particular parser implemented here is a Pratt parser.

## Anatomy of a Let statement

```
let x = 5; // binds 5 to name x
let y = 5 + 8 // binds y to the expression "5 + 8"
```

Programs are a series of statements. Let statements take the form:
`let <identifier> = <expression>;`

Definitions:

- Expressions produce values
- Statements do not produce values

This is not on the user end but for the parser. E.g.:

```
// This is not producing a value so it a statement
return x;
```

These aren't super rigorous definitions but good enough.

So we need two types of nodes:

- Expressions
- Statements

For a useful diagram see Page 46 of the book

## Parser Pseudocode

```js
function parseProgram() {
  // Creates root of AST
  program = newProgramASTNode();

  advanceTokens();

  while (currentToken() != EOF) {
    statement = null;

    // Find appropriate handler for token
    if (currentToken() == LET_TOKEN) {
      statement = parseLetStatement();
    } else if (currentToken() == RETURN_TOKEN) {
      statement = parseReturnStatement();
    } else if (currentToken() == IF_TOKEN) {
      statement = parseIfStatement();
    }

    // append statement to root node
    if (statement != null) {
      program.Statements.push(statement);
    }

    advanceTokens();
  }

  return program;
}
```

Parsing let statement

```js
function parseLetStatement() {
  // starts at [[let]] x = val;
  advanceTokens();
  // moves to let [[x]] = val; and parses x as identifier
  identifier = parseIdentifier();

  advanceTokens();
  // moves to let x [[=]] val; this needs to be an ASSIGN token
  // or else the syntax is invalid, but note that we can discard
  // this after the check since it's not important to the AST
  if (currentToken != ASSIGN) {
    parseError("expeced '=' after 'let <id>'");
    return null;
  }

  advanceTokens();
  // moves on to let x = [[val]];
  // Val may not be a single value, so we need to parse the expression
  value = parseExpression();

  // Now we have all the info we need to create an AST node
  // for this let statement. So we create the appropriate node
  // and then return it. This gets appended to Statements
  variableStatement = newVariableStatementASTNode();
  variableStatement.identifier = identifier;
  variableStatement.value = value;
  return variableStatement;
}
```

Parse identifier

```js
function parseIdentifier() {
  identifier = newIdentifierASTNode();
  identifier.token = currentToken();
  return identifier;
}
```

Parsing an expression

```js
function parseExpression() {
  if (currentToken() == INT) {
    if (nextToken() == PLUS) {
      return parseOperatorExpression();
    } else if (nextToken() == SEMICOLON_TOKEN) {
      return parseIntegerLiteral();
    }
  } else if (currentToken() == LEFT_PAREN) {
    return parseGroupedExpression();
  }
  // ...
}
```

Parsing operator expressions

```js
function parseOperatorExpression() {
  operatorExpression = newOperatorExpression();

  operatorExpression.left = parseIntegerLiteral();
  advanceTokens();
  operatorExpression.operator = currentToken();
  advanceTokens();
  operatorExpression.right = parseExpression();

  return operatorExpression();
}
```

- Parse program is the entry point
- Constructs root node newProgramASTNode()
- Build child nodes (statements) recursively
  - Advance tokens
  - Check current token
    - Call appropriate node construction handler
    - Error if ILLEGAL
  - Return to main loop

## Parsing expressions

Tokens get processed from left to right, but expressions aren't as straightforward.

- Operator precedence:

  - 5 \* 5 + 10
  - represented as ((5 \* 5) + 10)
  - 5 \* 5 is deeper in ast and evaluated earlier than addition
  - Parser needs to know \* has higher OP than +
  - Needs to know what to do with parens:
    - 5 _ (5 + 10) != 5 _ 5 + 10

- Expressions of same type can appear in multiple positions

  - e.g. -5 - 10
  - minus is prefix operator and infix operator

- Parens
  - (5 + 5) grouped expression
  - add(5, 5) call expression

Everything that isn't a let or return statement is an expression in Monkey.

Prefix operators:

```js
-5;
!true;
!false;
```

Infix operators (binary operators):

```js
5 + 5;
5 - 5;
5 / 5;
5 * 5;
```

comparison operators:

```js
foo == bar;
foo != bar;
foo < bar;
foo > bar;
```

parens:

```js
5 * (5 + 5)((5 + 5) * 5) * 5;
```

call expressions:

```js
add(2, 3);
add(add(2, 3), add(5, 10));
max(5, add(5, 5 * 5));
```

Identifier expressions

```js
(foo * bar) / foobar;
add(foo, bar);
```

Function literals are also expressions:

```js
let add = fn(x, y) { return x + y; };
// identifier expression function literal
fn(x, y) { return x + y }(5, 5)
// IIFE
(fn(x) { return x }(5) + 10) * 10
```

Also have if expressions:

```js
let result = if(10 > 5) { true } else { false };
```

This is a lot to handle.

## Pratt Parsing

Alternative to parsers based on context-free grammars.

Instead of associated parsing functions (e.g. parseLetStatement) w/ grammar rules we can associate them with single token types.

Each token type can have two associated parsing functions.

- Infix
- Prefix

### Terminology

Prefix operator:
In front of operand (`-5`, `--5`)
Postfix operator:
After operand (`i++`)
Infix operator:
Between operands (`a + b`)
Operator precedence:
Order of operations

## Parsing prefix expressions

Literals and Identifiers are pretty straightforward so I'm skipping them
in my notes. They should be easy enough to reason about from looking at
the code / tests.

Prefix expressions are a little trickier (but not too tricky).

They have the form:

```
<Prefix Operator><Right Expression>
```

This is going to resolve to a node that has:

```go
{
  // Other struct props
  Operator string // in TS this would be '-' | '!'
  Right Expression
}
```

So we first figure out what the operator is. Then we advance tokens and parse
the expression and point to it in the prefix expression node.

## Parsing infix expressions

We have eight infix operators:

```js
x + y;
x - y;
x * y;
x / y;
x > y;
x < y;
x == y;
x != y;
```

The general form is:

```
<left expression> <infix operator> <right expression>
```

## Example expression walkthrough

```js
1 + 2 + 3;
```

We need to nest AST nodes correctly. The AST (as a string) looks like this:

```js
1 + 2 + 3;
```

The actual tree struct is like this:

```go
&ast.InfixExpression{
  Token: token.PLUS
  Left: &ast.InfixExpression{
    Token: token.PLUS
    Left: &ast.IntegerLiteral{
      Value: 1
      Token: token.INT
    }
    Operator: "+"
    Right: &ast.IntegerLiteral{
      Value: 2
      Token: token.INT
    }
  }
  Operator: "+"
  Right: &ast.IntegerLiteral{
    Value: 3
    Token: token.INT
  }
}
```

Start:

```
key: c => current p => peek
1 + 2 + 3
^ ^
c p
```

We check if thhe current token has an associated prefixParseFn. We have one for
token.INT so we part it to an int literal.

```go
// Current state of leftExpr
&ast.IntegerLiteral{
  Value: 1,
  Token: token.INT
}
```

Now we enter the for loop.

```go
for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence()
```

This is true because we're currently at LOWEST precedence and there's
no semicolon. So we look for an infix function associated with the peek token.
Which will be `parseInfixExpression`.

```go
// Inside the for loop
infix := p.infixParseFns[p.peekToken.Type]
	if infix == nil {
		return leftExpr
	}
	p.nextToken()
	leftExpr = infix(leftExpr)
```

`infix` for token.PLUS (as said above) is not nil. So we're going to advance
tokens (to make `token.PLUS` into the current token) and then call the infix
function (`parseInfixExpression`) and reassign leftExpr to it.

This will end up creating this node:

```go
// current state of leftExpr
&ast.InfixExpression{
  Token: token.PLUS
  Left: &ast.IntegerLiteral{
    Value: 1
    Token: token.INT
  }
  Operator: "+"
  Right: &ast.IntegerLiteral{
    Value: 2
    Token: token.INT
  }
}
```

In more detail. Here is the code for `parseInfixExpression`:

```go
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}
```

We're creating an expression with the info we already know. Then we need to
parse the right expression. To do this we have to advance tokens. Then we
call `parseExpression` again with the precedence of the new current operator.

Parse expression is going to see if it can do a prefix parse. Since the current
token is `token.INT` then we can. So we parse it as an integer literal. But
the precedence (`SUM`) is the same as the precedence of the peek token. So
we won't do the loop and we return the literal for 2 which is added as `Right`
since we're doing the recursive call inside of `parseInfixExpression`.

Now we return that expression and get the node above (reposted here):

```go
// current state of leftExpr
&ast.InfixExpression{
  Token: token.PLUS
  Left: &ast.IntegerLiteral{
    Value: 1
    Token: token.INT
  }
  Operator: "+"
  Right: &ast.IntegerLiteral{
    Value: 2
    Token: token.INT
  }
}
```

In the original loop precedence is still at `LOWEST` so we continue the
loop. Since the peek token has an associated infix function we will do
another infix parse. This will follow the same process as before but
using our current tree as `Left`. Which gives us the final tree:

```go
&ast.InfixExpression{
  Token: token.PLUS
  Left: &ast.InfixExpression{
    Token: token.PLUS
    Left: &ast.IntegerLiteral{
      Value: 1
      Token: token.INT
    }
    Operator: "+"
    Right: &ast.IntegerLiteral{
      Value: 2
      Token: token.INT
    }
  }
  Operator: "+"
  Right: &ast.IntegerLiteral{
    Value: 3
    Token: token.INT
  }
}
```

Now we have

```
1 + 2 + 3;
        ^^
        cp
```

(We actually didn't have a semi in our original example and the termination
would be handled by peekPrecedence defaulting to `LOWEST`)
Since the peek token is a `token.SEMICOLON` then we terminate the loop
and return the tree (stored in `leftExpr`).

This is a pretty simple example, but the general goal is to have higher
precedence operators deeper in the AST. I don't want to add a bunch more
example but the book can be a good reference. I think it's also not very
hard to reason about how precedence would work (say if we replaced
1 + 2 + 3 with 1 + 2 \* 3) for more advanced cases because the recursion
involved makes things pretty simple.

### References

[Link to Vaughan Pratt's paper on operator precedence](https://dl.acm.org/doi/pdf/10.1145/512927.512931)

See page 88 in book for walkthrough of how it works.
