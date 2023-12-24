package ast

import "monkey-pl/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// This is the root node of every Monkey AST
type Program struct {
	Statements []Statement
}

// Basically this points to the first statement in the program
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Value string
	Token token.Token // IDENT token
}

// This is an empty implementation to help with type checking
// We're making identifiers into expressions because it keeps
// the language simple and because they will generally be
// used in expressions down the line
func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

/*
Need:
  - Name of variable
  - value being bound
  - reference to associated Token
*/
type LetStatement struct {
	Name  *Identifier
	Value Expression
	Token token.Token // LET token
}

// This is an empty implementation to help type checking
func (let *LetStatement) statementNode() {}

func (let *LetStatement) TokenLiteral() string {
	return let.Token.Literal
}

type ReturnStatement struct {
	ReturnValue Expression
	Token       token.Token // RETURN token
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
