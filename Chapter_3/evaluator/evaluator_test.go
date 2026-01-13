package evaluator

import (
	"Chapter_3/lexer"
	"Chapter_3/object"
	"Chapter_3/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestExclamationOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}

}
func testEval(input string) object.Object {
	lexer := lexer.NewLexer(input)
	parser := parser.NewParser(lexer)
	program := parser.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, evaluated object.Object, expected int64) bool {
	res, ok := evaluated.(*object.Integer)

	if !ok {
		t.Errorf("object is not Integer, got=%T (%+v)", evaluated, evaluated)
		return false
	}

	if res.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", res.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) bool {
	res, ok := evaluated.(*object.Boolean)

	if !ok {
		t.Errorf("object is not Boolean, got=%T (%+v)", evaluated, evaluated)
		return false
	}

	if res.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", res.Value, expected)
		return false
	}

	return true
}
