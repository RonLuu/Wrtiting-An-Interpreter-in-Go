package ast

import "Chapter_2/token"

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

type Program struct {
	Statements []Statement
}

type Variable struct {
	Token token.Token // The VARIABLE token
	Value string
}

type LetStatement struct {
	Token      token.Token // The LET token
	Variable   *Variable
	Expression Expression
}

type ReturnStatement struct {
	Token       token.Token // The RETURN token
	ReturnValue Expression
}

func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (variable *Variable) expressionNode()      {}
func (variable *Variable) TokenLiteral() string { return variable.Token.Literal }
