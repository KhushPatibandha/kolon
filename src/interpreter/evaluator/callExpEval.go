package evaluator

import (
	"errors"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/object"
)

func evalCallArgs(args []ast.Expression, env *object.Environment) ([]object.Object, bool, error) {
	var res []object.Object
	for _, e := range args {
		evaluated, hasErr, err := Eval(e, env, inTesting)
		if err != nil {
			return nil, hasErr, err
		}
		if evaluated.Type() == object.RETURN_VALUE_OBJ {
			evaluated = evaluated.(*object.ReturnValue).Value[0]
		}
		res = append(res, evaluated)
	}
	return res, false, nil
}

func applyFunction(fn object.Object, args []object.Object) (object.Object, bool, error) {
	function, ok := fn.(*object.Function)
	if !ok {
		builtin, ok := fn.(*object.Builtin)
		if ok {
			return builtin.Fn(args...)
		} else {
			return NULL, true, errors.New("not a function: " + string(fn.Type()))
		}
	}

	for i, param := range function.Parameters {
		function.Env.Set(param.ParameterName.Value, args[i], object.VAR)
	}
	evaluated, hasErr, err := Eval(function.Body, function.Env, inTesting)
	if err != nil {
		return NULL, hasErr, err
	}

	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		mustReturn = false
		if returnValue.Value == nil && function.ReturnType == nil {
			return nil, false, nil
		}
		return returnValue, false, nil
	}

	return evaluated, false, nil
}
