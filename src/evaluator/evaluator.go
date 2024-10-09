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
	case *ast.InfixExpression:
		left, err := Eval(node.Left)
		if err {
			return nil, true
		}
		right, err := Eval(node.Right)
		if err {
			return nil, true
		}
		return evalInfixExpression(node.Operator, left, right)
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

func evalInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool) {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ || left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		leftVal := 0.0
		rightVal := 0.0
		if left.Type() == object.INTEGER_OBJ {
			leftVal = float64(left.(*object.Integer).Value)
			rightVal = right.(*object.Float).Value
		} else {
			leftVal = left.(*object.Float).Value
			rightVal = float64(right.(*object.Integer).Value)
		}
		return evalFloatInfixExpression(operator, &object.Float{Value: leftVal}, &object.Float{Value: rightVal})
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.CHAR_OBJ && right.Type() == object.CHAR_OBJ:
		return evalCharInfixExpression(operator, left, right)
	default:
		return NULL, true
	}
}

func evalCharInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool) {
	leftVal := left.(*object.Char).Value
	rightVal := right.(*object.Char).Value
	leftVal = leftVal[1 : len(leftVal)-1]
	rightVal = rightVal[1 : len(rightVal)-1]

	switch operator {
	case "+":
		return &object.String{Value: "\"" + leftVal + rightVal + "\""}, false
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}, false
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}, false
	default:
		return NULL, true
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool) {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	leftVal = leftVal[1 : len(leftVal)-1]
	rightVal = rightVal[1 : len(rightVal)-1]

	switch operator {
	case "+":
		return &object.String{Value: "\"" + leftVal + rightVal + "\""}, false
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}, false
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}, false
	default:
		return NULL, true
	}
}

func evalFloatInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool) {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}, false
	case "-":
		return &object.Float{Value: leftVal - rightVal}, false
	case "/":
		return &object.Float{Value: leftVal / rightVal}, false
	case "*":
		return &object.Float{Value: leftVal * rightVal}, false
	case "%":
		return &object.Float{Value: float64(int(leftVal) % int(rightVal))}, false
	case ">":
		return &object.Boolean{Value: leftVal > rightVal}, false
	case "<":
		return &object.Boolean{Value: leftVal < rightVal}, false
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}, false
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}, false
	case "<=":
		return &object.Boolean{Value: leftVal <= rightVal}, false
	case ">=":
		return &object.Boolean{Value: leftVal >= rightVal}, false
	default:
		return NULL, true
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool) {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}, false
	case "-":
		return &object.Integer{Value: leftVal - rightVal}, false
	case "/":
		return &object.Integer{Value: leftVal / rightVal}, false
	case "*":
		return &object.Integer{Value: leftVal * rightVal}, false
	case "%":
		return &object.Integer{Value: leftVal % rightVal}, false
	case "&":
		return &object.Integer{Value: leftVal & rightVal}, false
	case "|":
		return &object.Integer{Value: leftVal | rightVal}, false
	case ">":
		return &object.Boolean{Value: leftVal > rightVal}, false
	case "<":
		return &object.Boolean{Value: leftVal < rightVal}, false
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}, false
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}, false
	case "<=":
		return &object.Boolean{Value: leftVal <= rightVal}, false
	case ">=":
		return &object.Boolean{Value: leftVal >= rightVal}, false
	default:
		return NULL, true
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool) {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}, false
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}, false
	case "&&":
		return &object.Boolean{Value: leftVal && rightVal}, false
	case "||":
		return &object.Boolean{Value: leftVal || rightVal}, false
	default:
		return NULL, true
	}
}
