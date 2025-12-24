package ast

import (
	"Chapter_2.5/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Expression interface {
	Node
	expressionNode()
}

// Identifier is an Expression
type Identifier struct {
	Token token.Token
	Value string
}

func (Identifier *Identifier) TokenLiteral() string { return Identifier.Token.Literal }
func (Identifier *Identifier) expressionNode()      {}
func (Identifier *Identifier) String() string       { return Identifier.Value }

// IntegerLiteral is an Expression
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (intLiteral *IntegerLiteral) TokenLiteral() string { return intLiteral.Token.Literal }
func (intLiteral *IntegerLiteral) expressionNode()      {}
func (intLiteral *IntegerLiteral) String() string       { return intLiteral.Token.Literal }

type PrefixExpression struct {
	Token           token.Token
	Operator        string
	RightExpression Expression
}

func (prefixExp *PrefixExpression) TokenLiteral() string { return prefixExp.Token.Literal }
func (prefixExp *PrefixExpression) expressionNode()      {}
func (prefixExp *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(prefixExp.Operator)
	out.WriteString(prefixExp.RightExpression.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token           token.Token
	Operator        string
	LeftExpression  Expression
	RightExpression Expression
}

func (infixExp *InfixExpression) TokenLiteral() string { return infixExp.Token.Literal }
func (infixExp *InfixExpression) expressionNode()      {}
func (infixExp *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(infixExp.LeftExpression.String())
	out.WriteString(" " + infixExp.Operator + " ")
	out.WriteString(infixExp.RightExpression.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (boolean *Boolean) TokenLiteral() string { return boolean.Token.Literal }
func (boolean *Boolean) expressionNode()      {}
func (boolean *Boolean) String() string       { return boolean.Token.Literal }

type Statement interface {
	Node
	statementNode()
}

// ExpressionStatement is a statement
// consists of only one expression, e.g. x+10;
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expStatement *ExpressionStatement) TokenLiteral() string { return expStatement.Token.Literal }
func (expStatement *ExpressionStatement) statementNode()       {}
func (expStatement *ExpressionStatement) String() string {
	if expStatement.Expression != nil {
		return expStatement.Expression.String()
	}
	return ""
}

// LetStatement is an Statement
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }
func (letStatement *LetStatement) statementNode()       {}
func (letStatement *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(letStatement.TokenLiteral() + " ")
	out.WriteString(letStatement.Name.String())
	out.WriteString(" = ")
	if letStatement.Value != nil {
		out.WriteString(letStatement.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (returnStatment *ReturnStatement) TokenLiteral() string { return returnStatment.Token.Literal }
func (returnStatment *ReturnStatement) statementNode()       {}
func (returnStatment *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(returnStatment.TokenLiteral() + " ")
	if returnStatment.ReturnValue != nil {
		out.WriteString(returnStatment.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (program *Program) String() string {
	var out bytes.Buffer
	for _, s := range program.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
