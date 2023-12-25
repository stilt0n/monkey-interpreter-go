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
  - 5 * 5 + 10
  - represented as ((5 * 5) + 10)
  - 5 * 5 is deeper in ast and evaluated earlier than addition
  - Parser needs to know * has higher OP than +
  - Needs to know what to do with parens:
    - 5 * (5 + 10) != 5 * 5 + 10

- Expressions of same type can appear in multiple positions
  - e.g. -5 - 10
  - minus is prefix operator and infix operator

- Parens
  - (5 + 5) grouped expression
  - add(5, 5) call expression

Everything that isn't a let or return statement is an expression in Monkey.

Prefix operators:
```js
-5
!true
!false
```
Infix operators (binary operators):
```js
5 + 5
5 - 5
5 / 5
5 * 5
```
comparison operators:
```js
foo == bar
foo != bar
foo < bar
foo > bar
```
parens:
```js
5 * (5 + 5)
((5 + 5) * 5) * 5
```
call expressions:
```js
add(2,3)
add(add(2,3),add(5,10))
max(5, add(5, (5 * 5)))
```
Identifier expressions
```js
foo * bar / foobar
add(foo, bar)
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

