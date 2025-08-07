package ast

// Data structure
import (
	"Chapter_2/token"
	"bytes"
)

// A Node can be a Statement or an Expression
type Node interface {
	TokenLiteral() string // A node must be able to return its token literal
	String() string       // A node must have a string representation
}

// An Expression is a type of Node
type Expression interface {
	Node
	expressionNode()
}

// A variable is a type of Expression
type Variable struct {
	Token   token.Token // The VARIABLE token
	Literal string
}

func (variable *Variable) TokenLiteral() string { return variable.Token.Literal }
func (variable *Variable) String() string       { return variable.Literal }
func (variable *Variable) expressionNode()      {}

// An Integer is a type of Expression
type IntegerLiteral struct {
	Token token.Token // The Integer token
	Value int64
}

func (integerLiteral *IntegerLiteral) TokenLiteral() string { return integerLiteral.Token.Literal }
func (integerLiteral *IntegerLiteral) String() string       { return integerLiteral.Token.Literal }
func (integerLiteral *IntegerLiteral) expressionNode()      {}

// A PrefixExpression is a type of Expression
type PrefixExpression struct {
	Token    token.Token // The prefix token: EXCLAMATION, MINUS
	Operator string      // The prefix operator: "!", "-"
	Right    Expression  // The right Expression: "!", "-"
}

func (prefixExpression *PrefixExpression) TokenLiteral() string {
	return prefixExpression.Token.Literal
}
func (prefixExpression *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(prefixExpression.Operator)
	out.WriteString(prefixExpression.Right.String())
	out.WriteString(")")
	return out.String()
}
func (prefixExpression *PrefixExpression) expressionNode() {}

// An InfixExpression is a type of Expression
type InfixExpression struct {
	Token      token.Token // The infix token: PLUS, MINUS, MULTIPLICATION
	LeftValue  Expression  // The left Expression
	Operator   string      // The infix operator: "+", "-", "*", ...
	RightValue Expression  // The right Expression
}

func (infixExpression *InfixExpression) TokenLiteral() string {
	return infixExpression.Token.Literal
}
func (infixExpression *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(infixExpression.LeftValue.String())
	out.WriteString(" " + infixExpression.Operator + " ")
	out.WriteString(infixExpression.RightValue.String())
	out.WriteString(")")
	return out.String()
}
func (InfixExpression *InfixExpression) expressionNode() {}

// A Statment is a type of Node
type Statement interface {
	Node
	statementNode()
}

// A LetStatment is a type of Statement
type LetStatement struct {
	Token      token.Token // The LET token
	Variable   *Variable
	Expression Expression
}

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
func (ls *LetStatement) statementNode() {}

// A ReturnStatement is a type of Statement
type ReturnStatement struct {
	Token       token.Token // The RETURN token
	ReturnValue Expression
}

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
func (rs *ReturnStatement) statementNode() {}

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
