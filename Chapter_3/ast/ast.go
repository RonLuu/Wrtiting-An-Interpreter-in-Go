package ast

import (
	"Chapter_3/token"
	"bytes"
	"strings"
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

type IfExpression struct {
	Token               token.Token // if token
	ConditionExpression Expression
	Consequence         *BlockStatement
	Alternative         *BlockStatement
}

func (ifExpression *IfExpression) TokenLiteral() string { return ifExpression.Token.Literal }
func (ifExpression *IfExpression) expressionNode()      {}
func (ifExpression *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ifExpression.ConditionExpression.String())
	out.WriteString(" ")
	out.WriteString(ifExpression.Consequence.String())

	if ifExpression.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ifExpression.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (functionLiteral *FunctionLiteral) TokenLiteral() string { return functionLiteral.Token.Literal }
func (functionLiteral *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (callExpression *CallExpression) TokenLiteral() string { return callExpression.Token.Literal }
func (callExpression *CallExpression) expressionNode()      {}
func (callExpression *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range callExpression.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(callExpression.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type Statement interface {
	Node
	statementNode()
}

type BlockStatement struct {
	Token      token.Token // {
	Statements []Statement
}

func (blockStatement *BlockStatement) TokenLiteral() string { return blockStatement.Token.Literal }
func (blockStatement *BlockStatement) statementNode()       {}
func (blockStatement *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range blockStatement.Statements {
		out.WriteString(s.String())
	}
	return out.String()
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
