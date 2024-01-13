package lexer

import (
	"errors"
	"monkey-pl/token"
)

/*
Position here points to ehte next character and readingPosition is the current one.
I am following the book's naming convention here, but think currentPosition and
nextPosition might be what I'd prefer. I'm holding off on renaming these until I
have a bigger picture view of things.

`ch` uses byte instead of rune because we are assuming that all input will be ascii chars

Ascii chars are much easier to work with because we don't need to account for cases
when a single character can be multiple bytes long.
*/
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	lex := &Lexer{input: input}
	// This accomplishes initializing the other vars
	lex.readChar()
	return lex
}

func (lex *Lexer) NextToken() token.Token {
	var tok token.Token

	lex.skipWhitespace()

	// Need to skip comments outside of the switch statement so
	// that our ability to return a correct token isn't disrupted
	// this is a loop to deal with multiple lines of comments
	for lex.ch == '#' {
		lex.skipComment()
		lex.skipWhitespace()
	}

	switch lex.ch {
	case '=':
		tok = newToken(token.ASSIGN, lex.ch)
		if lex.peekChar() == '=' {
			tok.Type = token.EQ
			tok.Literal = "=="
			lex.readChar()
		}
	case '+':
		tok = newToken(token.PLUS, lex.ch)
	case '-':
		tok = newToken(token.MINUS, lex.ch)
	case '!':
		tok = newToken(token.BANG, lex.ch)
		if lex.peekChar() == '=' {
			tok.Type = token.NEQ
			tok.Literal = "!="
			lex.readChar()
		}
	case '/':
		tok = newToken(token.SLASH, lex.ch)
	case '*':
		tok = newToken(token.ASTERISK, lex.ch)
	case '<':
		tok = newToken(token.LT, lex.ch)
	case '>':
		tok = newToken(token.GT, lex.ch)
	case ';':
		tok = newToken(token.SEMICOLON, lex.ch)
	case '(':
		tok = newToken(token.LPAREN, lex.ch)
	case ')':
		tok = newToken(token.RPAREN, lex.ch)
	case ',':
		tok = newToken(token.COMMA, lex.ch)
	case '{':
		tok = newToken(token.LBRACE, lex.ch)
	case '}':
		tok = newToken(token.RBRACE, lex.ch)
	case '[':
		tok = newToken(token.LBRACKET, lex.ch)
	case ']':
		tok = newToken(token.RBRACKET, lex.ch)
	case '"':
		literal, err := lex.readString()
		if err != nil {
			tok = newToken(token.ILLEGAL, lex.ch)
		}
		tok.Type = token.STRING
		tok.Literal = literal
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isAsciiLetter(lex.ch) {
			tok.Literal = lex.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isAsciiDigit(lex.ch) {
			tok.Type = token.INT
			tok.Literal = lex.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lex.ch)
		}
	}
	lex.readChar()
	return tok
}

// Non-exported methods
func (lex *Lexer) readChar() {
	if lex.readPosition >= len(lex.input) {
		// ASCII code for "NUL"
		lex.ch = 0
	} else {
		lex.ch = lex.input[lex.readPosition]
	}
	lex.position = lex.readPosition
	lex.readPosition++
}

func (lex *Lexer) peekChar() byte {
	if lex.readPosition >= len(lex.input) {
		return 0
	} else {
		return lex.input[lex.readPosition]
	}
}

func (lex *Lexer) readIdentifier() string {
	startPosition := lex.position
	for isAsciiLetter(lex.ch) {
		lex.readChar()
	}
	return lex.input[startPosition:lex.position]
}

func (lex *Lexer) readString() (string, error) {
	startPosition := lex.position + 1
	for {
		lex.readChar()
		if lex.ch == '"' {
			break
		}
		if lex.ch == 0 {
			return "", errors.New("unexpected EOF in string")
		}
	}
	return lex.input[startPosition:lex.position], nil
}

func (lex *Lexer) readNumber() string {
	startPosition := lex.position
	for isAsciiDigit(lex.ch) {
		lex.readChar()
	}
	return lex.input[startPosition:lex.position]
}

// Skips comments. All comments are single line, but the logic
// for a multiline syntax would be almost identical
func (lex *Lexer) skipComment() {
	for lex.ch != '\n' && lex.ch != 0 {
		lex.readChar()
	}
}

func (lex *Lexer) skipWhitespace() {
	for lex.ch == ' ' || lex.ch == '\t' || lex.ch == '\n' || lex.ch == '\r' {
		lex.readChar()
	}
}

// Utility functions
func newToken(t token.TokenType, c byte) token.Token {
	return token.Token{Type: t, Literal: string(c)}
}

func isAsciiLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isAsciiDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
