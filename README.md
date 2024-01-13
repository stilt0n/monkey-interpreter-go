# Monkey interpreter in Go

[You can try this interpreter out here](https://stilt0n.github.io/monkey-client/)

---

This is an interpreter that I am writing following Thorsten Ball's [Writing an Interpreter in Go](https://interpreterbook.com/).

I've been doing infrastructure work lately and at some point it became clear that a large number of the tools I'm responsible for managing are built around programming languages concepts. In particular, ASTs turn out to be ubiquitous in my work and are used in linting, code formatting, refactoring tools like codemods, dependency analysis, and then of course compilation and transpilation.

I don't think it's necessary to have more than a superficial understanding of ASTs to do my work, and I think many who do don't, but it seems like if I want to be able to contribute to the tooling I use or create my own, it will be important for me to understand this at a much deeper level.

I chose this particular book because:

- It's the first one I heard about that addresses my use case
- I've been learning Go for fun anyway
- Go seems like a good language choice for my intended application because of its performance

## On `notes.md` files

I've tried to write notes about what I'm learning through implementing the different parts of the interpreter. I am also writing notes as comments directly into the source code. The notes convey my understanding of things when I wrote them. They may sometimes be incorrect, and I may not always return to correct them.

#### About code blocks

When representing the abstract syntax of monkey in code blocks
I use go for syntax highlighting. For the concrete syntax, I
am using JavaScript because for syntax highlighting purposes
the syntaxes of the two languages are very close.

## Extra features

I may not implement many features beyond what was in the book, but the ones I did are listed here:

- Comments. Can use `#` character for single line comments.
- Divide by zero error handling. Returns divide by zero error when this is attempted.
- Error for function calls with incorrect number of arguments.
- Check for bottomless recursion and give error when stack depth is too deep
- Add error for unterminated strings (error is a little cryptic though)
- String comparison (==, !=, <, >)
- Negative operator in front of a string reverses it (e.g `-"abc" == "cba"`)
- split, join, toUpperCase, and toLowerCase functions

## Other stuff

I've added a server package so that this repl can be run on a
server. That way I can hopefully build a basic web text editor.

## Organization

There are two main packages:

- `cli`
- `server`

`cli` runs the repl in the command line. `server` runs an http server that can be sent code to evaluate.
