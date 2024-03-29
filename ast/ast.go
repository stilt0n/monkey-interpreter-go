package ast

import (
	"bytes"
	"monkey-pl/token"
	"strings"
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

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}

func (b *BooleanLiteral) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

type StringLiteral struct {
	Value string
	Token token.Token
}

func (s *StringLiteral) expressionNode() {}

func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}

func (s *StringLiteral) String() string {
	return s.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token // prefix token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}

func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // `FUNCTION` token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode() {}

func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	return out.String()
}

type ArrayLiteral struct {
	Token    token.Token // `[` token
	Elements []Expression
}

func (a *ArrayLiteral) expressionNode() {}

func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type HashLiteral struct {
	Token token.Token // `{` token
	Pairs map[Expression]Expression
}

func (h *HashLiteral) expressionNode() {}

func (h *HashLiteral) TokenLiteral() string {
	return h.Token.Literal
}

func (h *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range h.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type IfExpression struct {
	Token       token.Token // `if` token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {}

func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

type WhileExpression struct {
	Token     token.Token // `while` token
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileExpression) expressionNode() {}

func (w *WhileExpression) TokenLiteral() string {
	return w.Token.Literal
}

func (w *WhileExpression) String() string {
	var out bytes.Buffer
	out.WriteString("while")
	out.WriteString(w.Condition.String())
	out.WriteString(" ")
	out.WriteString(w.Body.String())
	return out.String()
}

type CallExpression struct {
	Token     token.Token // `(` token
	Function  Expression  // Identifier or Function Literal
	Arguments []Expression
}

func (c *CallExpression) expressionNode() {}
func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
}
func (c *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type IndexExpression struct {
	Token token.Token // `[` token
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode() {}

func (i *IndexExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("])")
	return out.String()
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

type BlockStatement struct {
	Token      token.Token // `{` token
	Statements []Statement
}

func (b *BlockStatement) statementNode() {}

func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range b.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
