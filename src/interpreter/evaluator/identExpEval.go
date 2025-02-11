package evaluator

import (
	"errors"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/object"
)

func evalIdentifier(node *ast.Identifier, env *object.Environment) (object.Object, bool, error) {
	variable, hasErr, err := getIdentifierVariable(node, env)
	if err != nil {
		return NULL, hasErr, err
	}
	return variable.Value, false, nil
}

func getIdentifierVariable(node *ast.Identifier, env *object.Environment) (*object.Variable, bool, error) {
	variable, ok := env.Get(node.Value)
	if !ok {
		builtin, ok := builtins[node.Value]
		if ok {
			return &object.Variable{Type: object.FUNCTION, Value: builtin}, false, nil
		} else {
			return &object.Variable{Type: object.VAR, Value: NULL}, true, errors.New("identifier not found: " + node.Value)
		}
	}
	return variable, false, nil
}
