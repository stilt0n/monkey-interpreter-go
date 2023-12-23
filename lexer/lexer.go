package lexer

import "monkey-pl/token"

type Lexer struct {
	input    string
	position int // points to current char
	// allows us to peek further to see what is next
	readPosition int // current reading position (after current char)
	// To keep things simple our lexer only supports ascii chars
	// if we were to support unicode we'd need to use runes here
	// since we're sticking to ascii we're using bytes.
	// The challenge with using runes here is that we'd then
	// need to account for when chars are more than a byte long.
	ch byte // current char being examined
}

func newToken(t token.TokenType, c byte) token.Token {
	return token.Token{Type: t, Literal: string(c)}
}

func New(input string) *Lexer {
	lex := &Lexer{input: input}
	// This accomplishes initializing the other vars
	lex.readChar()
	return lex
}

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

func (lex *Lexer) NextToken() token.Token {
	var tok token.Token

	switch lex.ch {
	case '=':
		tok = newToken(token.ASSIGN, lex.ch)
	case ';':
		tok = newToken(token.SEMICOLON, lex.ch)
	case '(':
		tok = newToken(token.LPAREN, lex.ch)
	case ')':
		tok = newToken(token.RPAREN, lex.ch)
	case ',':
		tok = newToken(token.COMMA, lex.ch)
	case '+':
		tok = newToken(token.PLUS, lex.ch)
	case '{':
		tok = newToken(token.LBRACE, lex.ch)
	case '}':
		tok = newToken(token.RBRACE, lex.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}
	lex.readChar()
	return tok
}
