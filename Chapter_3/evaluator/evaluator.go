package evaluator

import (
	"Chapter_3/ast"
	"Chapter_3/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		{
			return &object.Integer{Value: node.Value}
		}
	case *ast.Boolean:
		{
			return nativeBoolToBooleanObject(node.Value)
		}
	case *ast.PrefixExpression:
		rightObject := Eval(node.RightExpression)
		return evalPrefixExpression(node.Operator, rightObject)
	}

	return nil
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}

func nativeBoolToBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}

	return FALSE
}

func evalPrefixExpression(operator string, rightObject object.Object) object.Object {
	switch operator {
	case "!":
		return evalExclamationOperatorExpression(rightObject)
	case "-":
		return evalNegativeOperatorExpression(rightObject)
	default:
		return NULL
	}
}

func evalExclamationOperatorExpression(rightObject object.Object) object.Object {
	switch rightObject {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalNegativeOperatorExpression(rightObject object.Object) object.Object {
	if rightObject.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := rightObject.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
