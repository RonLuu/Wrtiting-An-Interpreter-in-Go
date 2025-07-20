package ast

import "Chapter_2/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	StatementNode()
}

type Expression interface {
	Node
	ExpressionNode()
}

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

type LetStatement struct {
	Token      token.Token // The LET token
	Variable   *Variable
	Expression Expression
}

func (ls *LetStatement) StatementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Variable struct {
	Token token.Token // The VARIABLE token
	Value string
}

func (variable *Variable) ExpressionNode()      {}
func (variable *Variable) TokenLiteral() string { return variable.Token.Literal }
