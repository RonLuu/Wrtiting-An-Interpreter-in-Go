package parser

import (
	"fmt"
	"testing"

	"Chapter_2/ast"
	"Chapter_2/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	checkParserError(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returns nil")
	}
	program.TokenLiteral()
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, it contains %d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		currentStatement := program.Statements[i]
		if !testLetStatement(t, currentStatement, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, identifier string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStatement, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStatement.Name.Value != identifier {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", identifier, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != identifier {
		t.Errorf("s.Name not '%s'. got=%s", identifier, letStatement.Name)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return add(15);
`
	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()
	checkParserError(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, it contains %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStatement, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt is not *ast.ReturnStatement. got=%t", stmt)
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	checkParserError(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("Program doesn't have enough statements. got=%d", len(program.Statements))
	}

	currentStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	currentIdentifier, ok := currentStatement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("currentStatement.Expression is not ast.Identifier. got=%T",
			currentStatement.Expression)
	}

	if currentIdentifier.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", currentIdentifier.Value)
	}

	if currentIdentifier.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			currentIdentifier.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	checkParserError(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("Program doesn't have enough statements. got=%d", len(program.Statements))
	}

	currentStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	currentIntegerLiteral, ok := currentStatement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("currentStatement.Expression is not ast.IntegerLiteral. got=%T",
			currentStatement.Expression)
	}

	if currentIntegerLiteral.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, currentIntegerLiteral.Value)
	}

	if currentIntegerLiteral.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			currentIntegerLiteral.TokenLiteral())
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()
		checkParserError(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.RightExpression, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()
		checkParserError(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an ExpressionStatement. got=%T\n", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("statement is not an InfixExpression. got=%T\n", statement.Expression)
		}

		if !testLiteralExpression(t, expression.LeftExpression, tt.leftValue) {
			return
		}

		if expression.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.RightExpression, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}
	for _, tt := range tests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()
		checkParserError(t, parser)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	checkParserError(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("Program doesn't have enough statements. got=%d", len(program.Statements))
	}

	currentStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	currentIdentifier, ok := currentStatement.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("currentStatement.Expression is not ast.Boolean. got=%T",
			currentStatement.Expression)
	}

	if currentIdentifier.Value != true {
		t.Errorf("ident.Value not %s. got=%s", "true", currentIdentifier.String())
	}

	if currentIdentifier.TokenLiteral() != "true" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			currentIdentifier.TokenLiteral())
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x< y) {x}`
	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	checkParserError(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifExpression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", statement.Expression)
	}

	if !testInfixExpression(t, ifExpression.ConditionExpression, "x", "<", "y") {
		return
	}

	if len(ifExpression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(ifExpression.Consequence.Statements))
	}

	consequence, ok := ifExpression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", ifExpression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if ifExpression.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", ifExpression.Alternative)
	}

}

func TestIfElseExpression(t *testing.T) {
	input := `if (x< y) {x} else {y}`
	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	checkParserError(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifExpression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", statement.Expression)
	}

	if !testInfixExpression(t, ifExpression.ConditionExpression, "x", "<", "y") {
		return
	}

	if len(ifExpression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(ifExpression.Consequence.Statements))
	}

	consequence, ok := ifExpression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", ifExpression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := ifExpression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", ifExpression.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}

}

func TestFunctioniteralParsing(t *testing.T) {
	input := `fn(x,y){x+y;}`

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	checkParserError(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	functionLiteral, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not ast.FunctionLiteral. got=%T", statement.Expression)
	}
	if len(functionLiteral.Parameters) != 2 {
		t.Fatalf("functionLiteral parameters wrong. want 2, got=%d\n", len(functionLiteral.Parameters))
	}

	testLiteralExpression(t, functionLiteral.Parameters[0], "x")
	testLiteralExpression(t, functionLiteral.Parameters[1], "y")

	if len(functionLiteral.Body.Statements) != 1 {
		t.Fatalf("functionLiteral.Body.Statements has not 1 statements. got=%d\n", len(functionLiteral.Body.Statements))
	}

	bodyStmt, ok := functionLiteral.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", functionLiteral.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		lexer := lexer.NewLexer(tt.input)
		parser := NewParser(lexer)
		program := parser.ParseProgram()
		checkParserError(t, parser)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func testIntegerLiteral(t *testing.T, expression ast.Expression, value int64) bool {
	integerLiteral, ok := expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expression is not IntegerLiteral expression. got=%t", expression)
		return false
	}

	if integerLiteral.Value != value {
		t.Errorf("integerLiteral value is not %d, got=%d", value, integerLiteral.Value)
		return false
	}

	if integerLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integerLiteral.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("Expression is not an identifier expression, got=%T", expression)
		return false
	}

	if identifier.Value != value {
		t.Errorf("Identifier value is not %s, got=%s", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("Identifier Token Literal is not %s, got=%s", value, identifier.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	boolean, ok := expression.(*ast.Boolean)
	if !ok {
		t.Errorf("Expression is not a boolean expression, got=%T", expression)
	}

	if boolean.Value != value {
		t.Errorf("Boolean expression not %t, got=%t", value, boolean.Value)
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, boolean.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(v))
	case int64:
		return testIntegerLiteral(t, expression, v)
	case string:
		return testIdentifier(t, expression, v)
	case bool:
		return testBooleanLiteral(t, expression, v)
	}
	t.Errorf("type of exp not handled. got=%T", expression)
	return false
}

func testInfixExpression(t *testing.T, expression ast.Expression, left interface{}, operator string, right interface{}) bool {
	infixExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", expression, expression)
		return false
	}
	if !testLiteralExpression(t, infixExpression.LeftExpression, left) {
		return false
	}
	if infixExpression.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, infixExpression.Operator)
		return false
	}
	if !testLiteralExpression(t, infixExpression.RightExpression, right) {
		return false
	}
	return true
}

func checkParserError(t *testing.T, parser *Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
