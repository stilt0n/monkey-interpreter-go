package parser

import (
	"fmt"
	"monkey-pl/ast"
	"monkey-pl/lexer"
	"monkey-pl/token"
)

type Parser struct {
	lex          *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lex: l, errors: []string{}}
	// Sets currentToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	// construct the ast root node
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// advance and handle tokens
	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		// unimplemented
		return nil
	}
}

// See commented pseudo-code in `notes.md`
func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Implement expression parsing

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

// Note: I'm not sure how I feel about this abstraction but the
// book uses it a lot. I may remove it later, but it's easier to
// follow along with it in place for now.
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// Check if peek token matches expectations and advance
// tokens only when it does. Useful for tokens like ASSIGN
// which have meaning in constructing the AST but are not
// required IN the AST. Adds a peek error when peek doesn't match
func (p *Parser) expectPeek(t token.TokenType) bool {
	// Note: the book abstracts this check into a function called
	// `peekTokenIs`. If it's used a lot, it may be worth abstracting
	// but to me this seems like an unnecessary abstraction.
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, received %s", t, p.peekToken.Type)
	p.errors = append(p.errors, message)
}
