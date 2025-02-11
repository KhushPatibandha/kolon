package evaluator

import (
	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/object"
)

func evalArrayValue(node *ast.ArrayValue, env *object.Environment) (object.Object, bool, error) {
	var res []object.Object
	for _, e := range node.Values {
		evaluated, hasErr, err := Eval(e, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}
		res = append(res, evaluated)
	}

	arrayObj := &object.Array{Elements: res}

	if node.Type != nil {
		arrayObj.TypeOf = node.Type.Value
		return arrayObj, false, nil
	}

	return arrayObj, false, nil
}
