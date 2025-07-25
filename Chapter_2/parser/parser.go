package parser

import (
	"Chapter_2/ast"
	"Chapter_2/lexer"
	"Chapter_2/token"
	"fmt"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// PRECEDENCE FOR OPERATION
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// A parser must
// need a lexer to read the token
// remember the current token and the next token
// contain all types of error when reading the program
// contain a prefix-parser-function dictionary
// contain a infix- parser-function dictionary
type Parser struct {
	lexer         *lexer.Lexer
	curToken      token.Token
	peekToken     token.Token
	errors        []string
	prefixParseFn map[token.TokenType]prefixParseFn
	infixParseFn  map[token.TokenType]infixParseFn
}

// Debug function
func (p *Parser) PrintParser() {
	fmt.Println("Parser:")
	p.lexer.PrintLexer()
	p.curToken.PrintToken()
	p.peekToken.PrintToken()
}

// Create a new parser
func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}
	// Set the current token
	p.nextToken()
	// Set the peek token
	p.nextToken()
	// Initialise a prefix-parse-function dictionary
	p.prefixParseFn = make(map[token.TokenType]prefixParseFn)
	// Add an entry for the prefix-parse-function dictionary
	p.registerPrefix(token.VARIABLE, p.parseVariable)
	// Initialise a prefix-parse-function dictionary
	p.infixParseFn = make(map[token.TokenType]infixParseFn)

	return p
}

// A function to move on to the next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) parseVariable() ast.Expression {
	return &ast.Variable{Token: p.curToken, Literal: p.curToken.Literal}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFn[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFn[tokenType] = fn
}
func (p *Parser) ParseProgram() *ast.Program {
	// Initialise a program
	program := &ast.Program{}
	// Initialise the statements in the program
	program.Statements = []ast.Statement{}

	// While we haven't reached the EOF token
	for p.curToken.Type != token.EOF {
		// Read the current 'statement'
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	// Decide what statement this is
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// Set the current token of the parser on the current statement
	stmt := &ast.LetStatement{Token: p.curToken}

	// Check if the next token is a VARIABLE token
	if !p.expectPeek(token.VARIABLE) {
		return nil
	}

	// Set the variable for the Let Statment
	stmt.Variable = &ast.Variable{Token: p.curToken, Literal: p.curToken.Literal}

	// Check if the next token is a ASSIGN token
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	// Skip over the Return token
	p.nextToken()

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFn[p.curToken.Type]

	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	return leftExp
}
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// If the token is as expected
// Return True and move on
// Else return False and add on Error
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(expectedToken token.TokenType) {
	msg := fmt.Sprintf("Expect the next token to be %s, got %s instead", expectedToken, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
