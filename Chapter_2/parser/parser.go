package parser

import (
	"fmt"
	"strconv"

	"Chapter_2/ast"
	"Chapter_2/lexer"
	"Chapter_2/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:        EQUALS,
	token.NEQ:       EQUALS,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.DIV:       PRODUCT,
	token.MULT:      PRODUCT,
	token.RLBRACKET: CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (parser *Parser) registerPrefix(tokenType token.TokenType, function prefixParseFn) {
	parser.prefixParseFns[tokenType] = function
}

func (parser *Parser) registerInfix(tokenType token.TokenType, function infixParseFn) {
	parser.infixParseFns[tokenType] = function
}

type Parser struct {
	lexer *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	// Given a token, return an appropriate function
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	errors []string
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.VARIABLE, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.EXCLAMATION, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanExpression)
	p.registerPrefix(token.RLBRACKET, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.MULT, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.RLBRACKET, p.parseCallExpression)
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) peekErrors(tokenType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tokenType, parser.peekToken.Type)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) peekPrecedence() int {
	precedence, ok := precedences[parser.peekToken.Type]
	if ok {
		return precedence
	}
	return LOWEST
}

func (parser *Parser) curPrecedence() int {
	precedence, ok := precedences[parser.curToken.Type]
	if ok {
		return precedence
	}
	return LOWEST
}
func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.curToken.Type != token.EOF {
		currentStatement := parser.parseStatement()
		program.Statements = append(program.Statements, currentStatement)
		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.curToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	currentStatement := &ast.LetStatement{Token: parser.curToken}

	if !parser.expectPeek(token.VARIABLE) {
		return nil
	}

	currentStatement.Name = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	parser.nextToken()
	currentStatement.Value = parser.parseExpression(LOWEST)
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}
	return currentStatement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	currentStatement := &ast.ReturnStatement{Token: parser.curToken}
	parser.nextToken()
	currentStatement.ReturnValue = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return currentStatement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))
	currentStatement := &ast.ExpressionStatement{Token: parser.curToken}
	currentStatement.Expression = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return currentStatement
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))
	prefixFn := parser.prefixParseFns[parser.curToken.Type]
	if prefixFn == nil {
		parser.noPrefixParseFnError(parser.curToken.Type)
		return nil
	}
	leftExp := prefixFn()

	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFns[parser.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		parser.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	// defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
	}
	parser.nextToken()

	expression.RightExpression = parser.parseExpression(PREFIX)
	return expression
}

func (parser *Parser) parseInfixExpression(leftExpression ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:          parser.curToken,
		Operator:       parser.curToken.Literal,
		LeftExpression: leftExpression,
	}

	precedence := parser.curPrecedence()
	parser.nextToken()
	expression.RightExpression = parser.parseExpression(precedence)
	return expression
}

func (parser *Parser) parseBooleanExpression() ast.Expression {
	return &ast.Boolean{Token: parser.curToken, Value: parser.curTokenIs(token.TRUE)}
}

func (parser *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: parser.curToken}

	if !parser.expectPeek(token.RLBRACKET) {
		return nil
	}

	parser.nextToken()
	expression.ConditionExpression = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RRBRACKET) {
		return nil
	}

	if !parser.expectPeek(token.PLBRACKET) {
		return nil
	}

	expression.Consequence = parser.parseBlockStatement()

	if parser.peekTokenIs(token.ELSE) {
		parser.nextToken()
		if !parser.expectPeek(token.PLBRACKET) {
			return nil
		}
		expression.Alternative = parser.parseBlockStatement()
	}
	return expression
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: parser.curToken}
	block.Statements = []ast.Statement{}
	parser.nextToken()

	for !parser.curTokenIs(token.PRBRACKET) && !parser.curTokenIs(token.EOF) {
		currentStatement := parser.parseStatement()
		if currentStatement != nil {
			block.Statements = append(block.Statements, currentStatement)
		}
		parser.nextToken()
	}

	return block
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	functionLiteral := &ast.FunctionLiteral{Token: parser.curToken}

	if !parser.expectPeek(token.RLBRACKET) {
		return nil
	}

	functionLiteral.Parameters = parser.parseFunctionParameters()

	if !parser.expectPeek(token.PLBRACKET) {
		return nil
	}

	functionLiteral.Body = parser.parseBlockStatement()

	return functionLiteral
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if parser.peekTokenIs(token.RRBRACKET) {
		parser.nextToken()
		return identifiers
	}

	parser.nextToken()
	identifier := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
	identifiers = append(identifiers, identifier)
	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		ident := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !parser.expectPeek(token.RRBRACKET) {
		return nil
	}
	return identifiers
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: parser.curToken, Function: function}
	exp.Arguments = parser.parseCallArguments()
	return exp
}

func (parser *Parser) parseCallArguments() []ast.Expression {
	arguments := []ast.Expression{}

	if parser.peekTokenIs(token.RRBRACKET) {
		parser.nextToken()
		return arguments
	}

	parser.nextToken()
	arguments = append(arguments, parser.parseExpression(LOWEST))

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		arguments = append(arguments, parser.parseExpression(LOWEST))
	}

	if !parser.expectPeek(token.RRBRACKET) {
		return nil
	}

	return arguments
}
func (parser *Parser) noPrefixParseFnError(tokenType token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", tokenType)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) parseIdentifier() ast.Expression {
	// defer untrace(trace("parseIdentifier"))
	return &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))
	intLiteral := &ast.IntegerLiteral{Token: parser.curToken}

	value, err := strconv.ParseInt(parser.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", parser.curToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}

	intLiteral.Value = value
	return intLiteral
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RRBRACKET) {
		return nil
	}
	return exp
}

func (parser *Parser) curTokenIs(tokenType token.TokenType) bool {
	return parser.curToken.Type == tokenType
}

func (parser *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return parser.peekToken.Type == tokenType
}

func (parser *Parser) expectPeek(tokenType token.TokenType) bool {
	if parser.peekTokenIs(tokenType) {
		parser.nextToken()
		return true
	}
	parser.peekErrors(tokenType)
	return false
}
