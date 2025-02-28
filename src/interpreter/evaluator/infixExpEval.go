package evaluator

import (
	"errors"

	"github.com/KhushPatibandha/Kolon/src/object"
)

func evalInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
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
	case left.Type() == object.ARRAY_OBJ && right.Type() == object.ARRAY_OBJ:
		return evalArrayInfixExpression(operator, left, right)
	default:
		return NULL, true, errors.New("invalid `infix` operation with variable types on left and right, got: `" + string(left.Type()) + "` and `" + string(right.Type()) + "`")
	}
}

func evalArrayInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.Array)

	switch operator {
	case "+":
		leftVal.Elements = append(leftVal.Elements, rightVal.Elements...)
		return leftVal, false, nil
	case "==":
		for i, element := range leftVal.Elements {
			if element.Type() != rightVal.Elements[i].Type() {
				return FALSE, false, nil
			}
			isCorrect, hasErr, err := evalInfixExpression("==", element, rightVal.Elements[i])
			if err != nil {
				return NULL, hasErr, err
			}
			if isCorrect == FALSE {
				return FALSE, false, nil
			}
		}
		return TRUE, false, nil
	case "!=":
		for i, element := range leftVal.Elements {
			if element.Type() != rightVal.Elements[i].Type() {
				return TRUE, false, nil
			}
			isCorrect, hasErr, err := evalInfixExpression("!=", element, rightVal.Elements[i])
			if err != nil {
				return NULL, hasErr, err
			}
			if isCorrect == TRUE {
				return TRUE, false, nil
			}
		}
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("can only use `+`, `==`, `!=` infix operators with 2 arrays, got: " + operator)
	}
}

func evalCharInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Char).Value
	rightVal := right.(*object.Char).Value
	leftVal = leftVal[1 : len(leftVal)-1]
	rightVal = rightVal[1 : len(rightVal)-1]

	switch operator {
	case "+":
		return &object.String{Value: "\"" + leftVal + rightVal + "\""}, false, nil
	case "==":
		if leftVal == rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "!=":
		if leftVal != rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("can only use `+`, `==`, `!=` infix operators with 2 `char`, got: " + operator)
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	leftVal = leftVal[1 : len(leftVal)-1]
	rightVal = rightVal[1 : len(rightVal)-1]

	switch operator {
	case "+":
		return &object.String{Value: "\"" + leftVal + rightVal + "\""}, false, nil
	case "==":
		if leftVal == rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "!=":
		if leftVal != rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("can only use `+`, `==`, `!=` infix operators with 2 `string`, got: " + operator)
	}
}

func evalFloatInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}, false, nil
	case "-":
		return &object.Float{Value: leftVal - rightVal}, false, nil
	case "/":
		return &object.Float{Value: leftVal / rightVal}, false, nil
	case "*":
		return &object.Float{Value: leftVal * rightVal}, false, nil
	case ">":
		if leftVal > rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "<":
		if leftVal < rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "==":
		if leftVal == rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "!=":
		if leftVal != rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "<=":
		if leftVal <= rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case ">=":
		if leftVal >= rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("can only use `+`, `-`, `*`, `/`, `>`, `<`, `<=`, `>=`, `!=`, `==` infix operators with 2 `float`, got: " + operator)
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}, false, nil
	case "-":
		return &object.Integer{Value: leftVal - rightVal}, false, nil
	case "/":
		if rightVal == 0 {
			return NULL, true, errors.New("integer divide by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}, false, nil
	case "*":
		return &object.Integer{Value: leftVal * rightVal}, false, nil
	case "%":
		return &object.Integer{Value: leftVal % rightVal}, false, nil
	case "&":
		return &object.Integer{Value: leftVal & rightVal}, false, nil
	case "|":
		return &object.Integer{Value: leftVal | rightVal}, false, nil
	case ">":
		if leftVal > rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "<":
		if leftVal < rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "==":
		if leftVal == rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "!=":
		if leftVal != rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "<=":
		if leftVal <= rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case ">=":
		if leftVal >= rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("can only use `+`, `-`, `*`, `/`, `%`, `>`, `<`, `<=`, `>=`, `!=`, `==`, `|`, `&` infix operators with 2 `int`, got: " + operator)
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "==":
		if leftVal == rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "!=":
		if leftVal != rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "&&":
		if leftVal && rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "||":
		if leftVal || rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("can only use `==`, `!=`, `&&`, `||` infix operators with 2 `bool`, got: " + operator)
	}
}
