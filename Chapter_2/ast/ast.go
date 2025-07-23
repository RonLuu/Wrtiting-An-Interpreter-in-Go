package ast

import (
	"Chapter_2/token"
	"bytes"
)

// A node can be a Statement or an Expression
// A node must be able to return its literal
type Node interface {
	TokenLiteral() string
	String() string
}

// A Statment is a type of node
type Statement interface {
	Node
	statementNode()
}

// An Expression is a type of node
type Expression interface {
	Node
	expressionNode()
}

type Variable struct {
	Token   token.Token // The VARIABLE token
	Literal string
}

func (variable *Variable) expressionNode()      {}
func (variable *Variable) TokenLiteral() string { return variable.Token.Literal }
func (variable *Variable) String() string       { return variable.Literal }

// A LetStatment is a type of Statement
type LetStatement struct {
	Token      token.Token // The LET token
	Variable   *Variable
	Expression Expression
}

// It includes statementNode and TokenLiteral
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Variable.String())
	out.WriteString(" = ")

	if ls.Expression != nil {
		out.WriteString(ls.Expression.String())
	}

	out.WriteString(";")

	return out.String()
}

// A ReturnStatement is a type of Statement
type ReturnStatement struct {
	Token       token.Token // The RETURN token
	ReturnValue Expression
}

// It includes statementNode and TokenLiteral
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

// An ExpressionStatement is a type of Statement
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

// It includes statementNode and TokenLiteral
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	var out bytes.Buffer

	if es.Expression != nil {
		out.WriteString(es.Expression.String())
	}

	return out.String()
}

// A program is an array of Statement
type Program struct {
	Statements []Statement
}

func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
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
