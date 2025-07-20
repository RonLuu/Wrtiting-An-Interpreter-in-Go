package parser

import (
	"Chapter_2/ast"
	"Chapter_2/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y =;
let 838383;
`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() return nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedVariable string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedVariable) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors\n", len(errors))

	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, stmt ast.Statement, expectedVariable string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got = %q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got ")
		return false
	}

	if letStmt.Variable.TokenLiteral() != expectedVariable {
		t.Errorf("letStmt.Variable.Value not '%s'. got = %s", expectedVariable, letStmt.Variable.Value)
		return false
	}

	if letStmt.Variable.TokenLiteral() != expectedVariable {
		t.Errorf("letStmt.Variable not '%s'. got = %s", expectedVariable, letStmt.Variable)
		return false
	}

	return true
}
