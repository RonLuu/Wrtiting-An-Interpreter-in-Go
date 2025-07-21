package parser

import (
	"Chapter_2/ast"
	"Chapter_2/lexer"
	"Chapter_2/token"
	"fmt"
)

type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

// Debug function
func (p *Parser) PrintParser() {
	fmt.Println("Parser:")
	p.lexer.PrintLexer()
	p.curToken.PrintToken()
	p.peekToken.PrintToken()
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}
	// Set the current token
	p.nextToken()
	// Set the peek token
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	// Initialise a program
	program := &ast.Program{}
	// Initialise all the statements in the program
	program.Statements = []ast.Statement{}

	// While we haven't reached the EOF token
	for p.curToken.Type != token.EOF {
		// Get the current 'statement'
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
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	// Check if the next token is a VARIABLE token
	if !p.expectPeek(token.VARIABLE) {
		return nil
	}

	stmt.Variable = &ast.Variable{Token: p.curToken, Value: p.curToken.Literal}

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

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

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
