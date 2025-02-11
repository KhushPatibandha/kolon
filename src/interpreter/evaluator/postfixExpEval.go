package evaluator

import "github.com/KhushPatibandha/Kolon/src/object"

func evalPostfixExpression(operator string, left object.Object) (object.Object, bool, error) {
	if left.Type() == object.INTEGER_OBJ {
		return evalIntegerPostfixExpression(operator, left)
	} else {
		return evalFloatPostfixExpression(operator, left)
	}
}

func evalFloatPostfixExpression(operator string, left object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Float).Value
	if operator == "++" {
		return &object.Float{Value: leftVal + 1}, false, nil
	} else {
		return &object.Float{Value: leftVal - 1}, false, nil
	}
}

func evalIntegerPostfixExpression(operator string, left object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Integer).Value
	if operator == "++" {
		return &object.Integer{Value: leftVal + 1}, false, nil
	} else {
		return &object.Integer{Value: leftVal - 1}, false, nil
	}
}
