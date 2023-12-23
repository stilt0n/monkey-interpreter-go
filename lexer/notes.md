# Lexer notes

The purpose of a lexer is to turn source code into something that won't be a pain in the ass to work with.

There will be two sequential changes we make to accomplish this:

Plain Text Source Code -> Tokens -> Abstract Syntax Tree

A lexer handles turning source code into tokens.

- Lexer is short for "lexical analysis"
- Lexers can also be called tokenizers

## Example of lexing

```
Source code:

let x = 5 + 5;
```

Turns into

```
[
  LET,
  IDENTIFIER("X"),
  EQUAL_SIGN,
  INTEGER(5),
  PLUS_SIGN,
  INTEGER(5),
  SEMICOLON
]
```

Note that not ever lexer is going to make the same choices for what is a Token. The book points out that some lexers would not bother converting 5 to an int until parsing, for example.

We ignore whitespace in our tokens because it's not part of the language. If we wanted to do something similar to Python, we would suddenly need to handle whitespace.

Lexers can also attatch metadata to a token like line numbers. This is common in production ready interpreters because if you want to be able to give a stack trace (or at least tell where an error occured) you need to keep track of this information somehow.

- Identifiers
- Ints
- Keywords
- Special Characters
- Invalid syntax
- End of file

We're going to pass strings rather than files to the lexer. In prod this would obviously need to read a file instead. This could be a pretty simply improvement to add when the book is finished.
