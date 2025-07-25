package parser

import (
	"Chapter_2/ast"
	"Chapter_2/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	// Create a lexer for the Parser to use
	l := lexer.NewLexer(input)
	p := NewParser(l)

	// The parser reads the program
	program := p.ParseProgram()
	// Check if there's any error after the parsing
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() return nil")
	}

	// Check if the parser doesn't read correctly three statments
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got = %d", len(program.Statements))
	}

	tests := []struct {
		expectedVariable string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	// Test if the parser read three variable correctly in let inputs
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedVariable) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, expectedVariable string) bool {
	// If the statement token is not 'let'
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got = %q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got ")
		return false
	}

	// If the variable in the let statement is not the expected variable
	if letStmt.Variable.Literal != expectedVariable {
		t.Errorf("letStmt.Variable.Literal not '%s'. got = %s", expectedVariable, letStmt.Variable.Literal)
		return false
	}

	// If the variable in the let statement is not the expected variable
	if letStmt.Variable.TokenLiteral() != expectedVariable {
		t.Errorf("letStmt.Variable not '%s'. got = %s", expectedVariable, letStmt.Variable)
		return false
	}

	return true
}
func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	// Create a lexer for the Parser to use
	l := lexer.NewLexer(input)
	p := NewParser(l)

	// The parser reads the program
	program := p.ParseProgram()

	// Check if there's any error after the parsing
	checkParserErrors(t, p)

	// Check if the parser doesn't read correctly three statments
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got = %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got = %T", returnStmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got = %q", returnStmt.TokenLiteral())
		}
	}
}

func TestVariableExpression(t *testing.T) {
	input := "foobar;"

	// Create a lexer for the Parser to use
	l := lexer.NewLexer(input)
	parser := NewParser(l)

	// The parser reads the program
	program := parser.ParseProgram()

	// Check if there's any error after the parsing
	checkParserErrors(t, parser)

	// Check if the parser doesn't read correctly one statment
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough statements. got = %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not a ast.ExpressionStatement. got = %T", program.Statements[0])
	}

	variable, ok := stmt.Expression.(*ast.Variable)
	if !ok {
		t.Fatalf("stmt.Expression is not a ast.Variable. got = %T", stmt.Expression)
	}

	if variable.Literal != "foobar" {
		t.Errorf("variable.Literal not %s. got=%s", "foobar", variable.Literal)
	}

	if variable.TokenLiteral() != "foobar" {
		t.Errorf("variable.TokenLiteral not %s. got=%s", "foobar", variable.TokenLiteral())
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	// Create a lexer for the Parser to use
	l := lexer.NewLexer(input)
	parser := NewParser(l)

	// The parser reads the program
	program := parser.ParseProgram()

	// Check if there's any error after the parsing
	checkParserErrors(t, parser)

	// Check if the parser doesn't read correctly one statment
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough statements. got = %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not a ast.ExpressionStatement. got = %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not a ast.IntegerLiteral. got = %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%s", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
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
