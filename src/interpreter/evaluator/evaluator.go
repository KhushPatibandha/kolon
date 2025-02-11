package evaluator

import (
	"errors"
	"fmt"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/object"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

var (
	NULL       = &object.Null{}
	CONTINUE   = &object.Continue{}
	BREAK      = &object.Break{}
	TRUE       = &object.Boolean{Value: true}
	FALSE      = &object.Boolean{Value: false}
	inForLoop  = false
	mustReturn = false
	inTesting  = false
)

func Eval(node ast.Node, env *object.Environment, inTest bool) (object.Object, bool, error) {
	inTesting = inTest
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env)

	case *ast.IntegerValue:
		return &object.Integer{Value: node.Value}, false, nil
	case *ast.FloatValue:
		return &object.Float{Value: node.Value}, false, nil
	case *ast.StringValue:
		return &object.String{Value: node.Value}, false, nil
	case *ast.CharValue:
		return &object.Char{Value: node.Value}, false, nil
	case *ast.BooleanValue:
		if node.Value {
			return TRUE, false, nil
		} else {
			return FALSE, false, nil
		}
	case *ast.ArrayValue:
		return evalArrayValue(node, env)
	case *ast.HashMap:
		return evalHashMap(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.PrefixExpression:
		right, hasErr, err := Eval(node.Right, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}
		if right.Type() == object.RETURN_VALUE_OBJ {
			right = right.(*object.ReturnValue).Value[0]
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left, hasErr, err := Eval(node.Left, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}
		right, hasErr, err := Eval(node.Right, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}
		if left.Type() == object.RETURN_VALUE_OBJ {
			left = left.(*object.ReturnValue).Value[0]
		}
		if right.Type() == object.RETURN_VALUE_OBJ {
			right = right.(*object.ReturnValue).Value[0]
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PostfixExpression:
		left, hasErr, err := Eval(node.Left, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}
		if left.Type() == object.RETURN_VALUE_OBJ {
			left = left.(*object.ReturnValue).Value[0]
		}
		resObj, hasErr, err := evalPostfixExpression(node.Operator, left)
		if err != nil {
			return NULL, hasErr, err
		}
		if node.IsStmt {
			if id, ok := node.Left.(*ast.Identifier); ok {
				idVariable, hasErr, err := getIdentifierVariable(id, env)
				if err != nil {
					return NULL, hasErr, err
				}
				var isVar bool
				if idVariable.Type == object.VAR {
					isVar = true
				}
				if isVar {
					env.Update(id.Value, resObj, object.VAR)
				} else {
					return NULL, true, errors.New("trying to use postfix operation on const variable.")
				}
			}
		}
		return resObj, false, nil
	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, false, nil, env)
	case *ast.CallExpression:
		function, hasErr, err := Eval(node.Name, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}
		args, hasErr, err := evalCallArgs(node.Args, env)
		if err != nil {
			return NULL, hasErr, err
		}
		return applyFunction(function, args)
	case *ast.IndexExpression:
		left, hasErr, err := Eval(node.Left, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}
		index, hasErr, err := Eval(node.Index, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}
		if left.Type() == object.RETURN_VALUE_OBJ {
			left = left.(*object.ReturnValue).Value[0]
		}
		if index.Type() == object.RETURN_VALUE_OBJ {
			index = index.(*object.ReturnValue).Value[0]
		}
		return evalIndexExpression(left, index)

	case *ast.ExpressionStatement:
		res, hasErr, err := Eval(node.Expression, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}
		if inTesting {
			return res, hasErr, err
		}
		return nil, false, nil
	case *ast.FunctionBody:
		return evalStatements(node.Statements, env)
	case *ast.Function:
		if node.Name.Value == "main" {
			for key, value := range parser.FunctionMap {
				newLocalEnv := object.NewEnclosedEnvironment(env)
				env.Set(key, &object.Function{Name: value.Name, Parameters: value.Parameters, ReturnType: value.ReturnType, Body: value.Body, Env: newLocalEnv}, object.FUNCTION)
			}

			funVar, _ := env.Get("main")
			return evalMainFunc(node, funVar.Value.(*object.Function).Env)
		}
		return nil, false, nil
	case *ast.VarStatement:
		return evalVarStatement(node, false, nil, env)
	case *ast.MultiValueAssignStmt:
		return evalMultiValueAssignStmt(node, env)
	case *ast.ReturnStatement:
		if node.Value == nil {
			mustReturn = true
			return &object.ReturnValue{Value: nil}, false, nil
		}
		var val []object.Object

		for i := 0; i < len(node.Value); i++ {
			rsObj, hasErr, err := evalReturnStatement(node, i, env)
			if err != nil {
				return NULL, hasErr, err
			}
			if rsObj.Type() == object.RETURN_VALUE_OBJ {
				rsObj = rsObj.(*object.ReturnValue).Value[0]
			}
			val = append(val, rsObj)
		}

		mustReturn = true
		return &object.ReturnValue{Value: val}, false, nil
	case *ast.IfStatement:
		localEnv := object.NewEnclosedEnvironment(env)
		res, hasErr, err := evalIfStatements(node, localEnv)
		if err != nil {
			return NULL, hasErr, err
		}
		if inTesting || ((mustReturn && res.Type() == object.RETURN_VALUE_OBJ) || res == BREAK || res == CONTINUE) {
			return res, hasErr, err
		}
		return nil, false, nil
	case *ast.ForLoopStatement:
		localEnv := object.NewEnclosedEnvironment(env)
		return evalForLoop(node, localEnv)
	case *ast.ContinueStatement:
		return CONTINUE, false, nil
	case *ast.BreakStatement:
		return BREAK, false, nil
	default:
		return NULL, true, fmt.Errorf("no eval function for given node type, got: %T", node)
	}
}

func evalStatements(stmts []ast.Statement, env *object.Environment) (object.Object, bool, error) {
	var result object.Object
	var hasErr bool
	var err error
	for _, statement := range stmts {
		result, hasErr, err = Eval(statement, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}

		if mustReturn && result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result, hasErr, err
		}

		if inForLoop && result == BREAK {
			return BREAK, false, nil
		} else if inForLoop && result == CONTINUE {
			return CONTINUE, false, nil
		}
	}

	if inTesting {
		return result, hasErr, err
	}
	return nil, false, nil
}

func evalMainFunc(node *ast.Function, env *object.Environment) (object.Object, bool, error) {
	_, hasErr, err := Eval(node.Body, env, inTesting)
	if err != nil {
		return NULL, hasErr, err
	}
	return nil, false, nil
}

func evalVarStatement(node *ast.VarStatement, injectObj bool, obj object.Object, env *object.Environment) (object.Object, bool, error) {
	var val object.Object
	var hasErr bool
	var err error

	if injectObj {
		val = obj
	} else if !injectObj {
		val, hasErr, err = Eval(node.Value, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}
	}

	for val.Type() == object.RETURN_VALUE_OBJ {
		val = val.(*object.ReturnValue).Value[0]
	}

	if node.Token.Kind == lexer.VAR {
		env.Set(node.Name.Value, val, object.VAR)
	} else if node.Token.Kind == lexer.CONST {
		env.Set(node.Name.Value, val, object.CONST)
	}

	return nil, false, nil
}

func evalMultiValueAssignStmt(node *ast.MultiValueAssignStmt, env *object.Environment) (object.Object, bool, error) {
	if !node.SingleCallExp {
		for _, element := range node.Objects {
			switch element := element.(type) {
			case *ast.VarStatement:
				_, hasErr, err := Eval(element, env, inTesting)
				if err != nil {
					return NULL, hasErr, err
				}
			case *ast.ExpressionStatement:
				_, hasErr, err := Eval(element.Expression.(*ast.AssignmentExpression), env, inTesting)
				if err != nil {
					return NULL, hasErr, err
				}
			}
		}
	} else {
		isVar := false
		varEntry, ok := node.Objects[0].(*ast.VarStatement)
		var expStmtEntry *ast.CallExpression
		if ok {
			isVar = true
		} else {
			isVar = false
			expStmtEntry, _ = node.Objects[0].(*ast.ExpressionStatement).Expression.(*ast.AssignmentExpression).Right.(*ast.CallExpression)
		}

		var returnObj object.Object
		if isVar {
			newObj, hasErr, err := Eval(varEntry.Value, env, inTesting)
			if err != nil {
				return NULL, hasErr, err
			}
			returnObj = newObj
		} else {
			newObj, hasErr, err := Eval(expStmtEntry, env, inTesting)
			if err != nil {
				return NULL, hasErr, err
			}
			returnObj = newObj
		}
		returnObjList := returnObj.(*object.ReturnValue).Value

		for i, element := range node.Objects {
			switch element := element.(type) {
			case *ast.VarStatement:
				_, hasErr, err := evalVarStatement(element, true, returnObjList[i], env)
				if err != nil {
					return NULL, hasErr, err
				}
			case *ast.ExpressionStatement:
				_, hasErr, err := evalAssignmentExpression(element.Expression.(*ast.AssignmentExpression), true, returnObjList[i], env)
				if err != nil {
					return NULL, hasErr, err
				}
			}
		}
	}

	return nil, false, nil
}

func evalReturnStatement(rs *ast.ReturnStatement, idx int, env *object.Environment) (object.Object, bool, error) {
	currNode := rs.Value[idx]
	return Eval(currNode, env, inTesting)
}

func evalIfStatements(node *ast.IfStatement, env *object.Environment) (object.Object, bool, error) {
	condition, hasErr, err := Eval(node.Value, env, inTesting)
	if err != nil {
		return NULL, hasErr, err
	}

	conditionRes, hasErr, err := execIf(condition)
	if err != nil {
		return NULL, hasErr, err
	}

	if conditionRes == TRUE {
		return evalStatements(node.Body.Statements, env)
	} else if node.MultiConseq != nil {
		for i := 0; i < len(node.MultiConseq); i++ {
			obj, success, hasErr, err := evalElseIfStatement(node.MultiConseq[i], env)
			if err != nil {
				return NULL, hasErr, err
			}
			if success {
				return obj, false, nil
			}
		}
	}
	if node.Consequence != nil {
		return evalStatements(node.Consequence.Body.Statements, env)
	}
	return nil, false, nil
}

func evalForLoop(node *ast.ForLoopStatement, env *object.Environment) (object.Object, bool, error) {
	inForLoop = true

	_, hasErr, err := Eval(node.Left, env, inTesting)
	if err != nil {
		return NULL, hasErr, err
	}

	infixObj, hasErr, err := Eval(node.Middle, env, inTesting)
	if err != nil {
		return NULL, hasErr, err
	}

	for infixObj == TRUE {
		resStmtObj, hasErr, err := evalStatements(node.Body.Statements, env)
		if err != nil {
			return NULL, hasErr, err
		}

		if resStmtObj == BREAK {
			break
		} else if resReturnObj, ok := resStmtObj.(*object.ReturnValue); ok {
			inForLoop = false
			return resReturnObj, false, nil
		}

		postfixObj, hasErr, err := Eval(node.Right, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}

		env.Update(node.Left.Name.Value, postfixObj, object.VAR)
		infixObj, hasErr, err = Eval(node.Middle, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}
		if infixObj == FALSE {
			break
		}
	}

	inForLoop = false
	return nil, false, nil
}

func evalElseIfStatement(node *ast.ElseIfStatement, env *object.Environment) (object.Object, bool, bool, error) {
	condition, hasErr, err := Eval(node.Value, env, inTesting)
	if err != nil {
		return NULL, false, hasErr, err
	}
	conditionRes, hasErr, err := execIf(condition)
	if err != nil {
		return NULL, false, hasErr, err
	}
	if conditionRes == TRUE {
		resObj, hasErr, err := Eval(node.Body, env, inTesting)
		if err != nil {
			return NULL, false, hasErr, err
		}
		return resObj, true, false, nil
	}
	return nil, false, false, nil
}

func execIf(obj object.Object) (object.Object, bool, error) {
	if obj.Type() == object.RETURN_VALUE_OBJ {
		return execIf(obj.(*object.ReturnValue).Value[0])
	}
	switch obj {
	case TRUE:
		return TRUE, false, nil
	case FALSE:
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("conditions for `if` and `else if` statements must result in a `bool`, got: " + string(obj.Type()))
	}
}
