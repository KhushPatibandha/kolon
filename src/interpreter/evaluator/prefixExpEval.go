package evaluator

import "github.com/KhushPatibandha/Kolon/src/object"

func evalPrefixExpression(operator string, right object.Object) (object.Object, bool, error) {
	if operator == "!" {
		return evalBangOperatorExpression(right)
	} else {
		return evalMinusOperatorExpression(right)
	}
}

func evalMinusOperatorExpression(right object.Object) (object.Object, bool, error) {
	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}, false, nil
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}, false, nil
}

func evalBangOperatorExpression(right object.Object) (object.Object, bool, error) {
	if right == TRUE {
		return FALSE, false, nil
	} else {
		return TRUE, false, nil
	}
}
