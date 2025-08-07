package parser

import (
	"Chapter_2/ast"
	"Chapter_2/lexer"
	"Chapter_2/token"
	"fmt"
	"strconv"
)

type (
	// Prefix parse function is a function that return an expression
	prefixParseFn func() ast.Expression
	// Infix  parse function is a function that takes a left expression return a whole expression
	infixParseFn func(ast.Expression) ast.Expression
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

var precedences = map[token.TokenType]int{
	token.EQ:    EQUALS,
	token.NEQ:   EQUALS,
	token.LT:    LESSGREATER,
	token.GT:    LESSGREATER,
	token.PLUS:  SUM,
	token.MINUS: SUM,
	token.DIV:   PRODUCT,
	token.MULT:  PRODUCT,
}

type Parser struct {
	// A parser must
	lexer         *lexer.Lexer                      // need a lexer to read the token
	curToken      token.Token                       // remember the current token
	peekToken     token.Token                       // remember the next token
	errors        []string                          // contain all types of error when reading the program
	prefixParseFn map[token.TokenType]prefixParseFn // contain a prefix-parser-function dictionary
	infixParseFn  map[token.TokenType]infixParseFn  // contain a infix- parser-function dictionary
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
	p.registerPrefix(token.VARIABLE, p.parseVariable)            // register a parse variable function
	p.registerPrefix(token.INT, p.parseIntegerLiteral)           // register a parse integer function
	p.registerPrefix(token.EXCLAMATION, p.parsePrefixExpression) // register a parse 'not' function
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)       // register a parse 'negative' function

	// Initialise a prefix-parse-function dictionary
	p.infixParseFn = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.MULT, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)

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

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:     p.curToken,
		LeftValue: left,
		Operator:  p.curToken.Literal,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.RightValue = p.parseExpression(precedence)
	return expression
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
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
	// Get the current prefix operation
	prefix := p.prefixParseFn[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// Assign the prefix
	leftExp := prefix()
	// While the parser hasn't reached the semicolon
	// and
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFn[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
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
