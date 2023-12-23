# Monkey interpreter in Go

This is an interpreter that I am writing following Thorsten Ball's [Writing an Interpreter in Go](https://interpreterbook.com/).

I've been doing infrastructure work lately and at some point it became clear that a large number of the tools I'm responsible for managing are built around programming languages concepts. In particular, ASTs turn out to be ubiquitous in my work and are used in linting, code formatting, refactoring tools like codemods, dependency analysis, and then of course compilation and transpilation.

I don't think it's necessary to have more than a superficial understanding of ASTs to do my work, and I think many who do don't, but it seems like if I want to be able to contribute to the tooling I use or create my own, it will be important for me to understand this at a much deeper level.

I chose this particular book because:

- It's the first one I hear about that addresses my use case
- I've been learning Go for fun anyway
- Go seems like a better good language choice for my intended application because of its performance

## Extra features

I may not many features beyond what was in the book, but the ones I did are listed here:

- Comments. Can use `#` character for single line comments.
