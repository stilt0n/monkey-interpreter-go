package ast

import (
	"bytes"
	"monkey-pl/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Statement wrapper for expressions
type ExpressionStatement struct {
	Token token.Token // first token of expression
	// Self-explanatory but this holds the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
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

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
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

func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Value int64
	Token token.Token
}

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
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

func (let *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(let.TokenLiteral() + " ")
	out.WriteString(let.Name.String())
	out.WriteString(" = ")
	if let.Value != nil {
		out.WriteString(let.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	ReturnValue Expression
	Token       token.Token // RETURN token
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral())
	if rs.ReturnValue != nil {
		out.WriteString(" " + rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}
