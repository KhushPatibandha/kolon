package evaluator

import (
	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) (object.Object, bool) {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerValue:
		return &object.Integer{Value: node.Value}, false
	case *ast.FloatValue:
		return &object.Float{Value: node.Value}, false
	case *ast.StringValue:
		return &object.String{Value: node.Value}, false
	case *ast.CharValue:
		return &object.Char{Value: node.Value}, false
	case *ast.BooleanValue:
		if node.Value {
			return TRUE, false
		} else {
			return FALSE, false
		}
	case *ast.PrefixExpression:
		right, err := Eval(node.Right)
		if err {
			return nil, true
		}
		return evalPrefixExpression(node.Operator, right)
	}
	return nil, false
}

func evalStatements(stmts []ast.Statement) (object.Object, bool) {
	var result object.Object
	var err bool
	for _, statement := range stmts {
		result, err = Eval(statement)
	}
	return result, err
}

func evalPrefixExpression(operator string, right object.Object) (object.Object, bool) {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return NULL, false
	}
}

func evalMinusOperatorExpression(right object.Object) (object.Object, bool) {
	if right.Type() != object.INTEGER_OBJ && right.Type() != object.FLOAT_OBJ {
		return NULL, true
	}

	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}, false
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}, false
}

func evalBangOperatorExpression(right object.Object) (object.Object, bool) {
	switch right := right.(type) {
	case *object.Boolean:
		if right == TRUE {
			return FALSE, false
		} else {
			return TRUE, false
		}
	case *object.Null:
		return TRUE, false
	case *object.Integer, *object.Float, *object.String, *object.Char:
		return NULL, true
	default:
		return FALSE, false
	}
}
