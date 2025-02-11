package evaluator

import (
	"errors"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/object"
)

func evalAssignmentExpression(node *ast.AssignmentExpression, injectObj bool, val object.Object, env *object.Environment) (object.Object, bool, error) {
	switch node.Operator {
	case "=":
		return evalAssignOp(node, injectObj, val, env)
	case "+=":
		return evalSymbolAssignOp("+", node, env)
	case "-=":
		return evalSymbolAssignOp("-", node, env)
	case "*=":
		return evalSymbolAssignOp("*", node, env)
	case "/=":
		return evalSymbolAssignOp("/", node, env)
	case "%=":
		return evalSymbolAssignOp("%", node, env)
	default:
		return NULL, true, errors.New("only `=`, `+=`, `-=`, `*=`, `/=`, `%=` `assignment` operators supported, got: " + node.Operator)
	}
}

func assignOpHelper(node *ast.AssignmentExpression, injectObj bool, rightVal object.Object, env *object.Environment) (object.Object, object.Object, bool, error) {
	leftSideVariable, hasErr, err := getIdentifierVariable(node.Left, env)
	if err != nil {
		return NULL, NULL, hasErr, err
	}

	if injectObj {
		return leftSideVariable.Value, rightVal, false, nil
	}

	rightSideObj, hasErr, err := Eval(node.Right, env, inTesting)
	if err != nil {
		return NULL, NULL, hasErr, err
	}

	if rightSideObj.Type() == object.RETURN_VALUE_OBJ {
		rightSideObj = rightSideObj.(*object.ReturnValue).Value[0]
	}

	return leftSideVariable.Value, rightSideObj, false, nil
}

func evalSymbolAssignOp(operator string, node *ast.AssignmentExpression, env *object.Environment) (object.Object, bool, error) {
	leftSideObj, rightSideObj, hasErr, err := assignOpHelper(node, false, nil, env)
	if err != nil {
		return NULL, hasErr, err
	}

	resObj, hasErr, err := evalInfixExpression(operator, leftSideObj, rightSideObj)
	if err != nil {
		return NULL, hasErr, err
	}
	env.Update(node.Left.Value, resObj, object.VAR)
	return resObj, false, nil
}

func evalAssignOp(node *ast.AssignmentExpression, injectObj bool, rightVal object.Object, env *object.Environment) (object.Object, bool, error) {
	_, rightSideObj, hasErr, err := assignOpHelper(node, injectObj, rightVal, env)
	if err != nil {
		return NULL, hasErr, err
	}
	env.Update(node.Left.Value, rightSideObj, object.VAR)
	return rightSideObj, false, nil
}
