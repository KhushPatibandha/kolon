package evaluator

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/interpreter/object"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

var (
	NULL      = &object.Null{}
	CONTINUE  = &object.Continue{}
	BREAK     = &object.Break{}
	TRUE      = &object.Boolean{Value: true}
	FALSE     = &object.Boolean{Value: false}
	inForLoop = false
)

func Eval(node ast.Node, env *object.Environment) (object.Object, bool, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
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
	case *ast.IndexExpression:
		left, hasErr, err := Eval(node.Left, env)
		if err != nil {
			return left, hasErr, err
		}
		index, hasErr, err := Eval(node.Index, env)
		if err != nil {
			return index, hasErr, err
		}
		if left.Type() == object.RETURN_VALUE_OBJ {
			left = left.(*object.ReturnValue).Value[0]
		}
		if index.Type() == object.RETURN_VALUE_OBJ {
			index = index.(*object.ReturnValue).Value[0]
		}
		return evalIndexExpression(left, index)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.PrefixExpression:
		right, hasErr, err := Eval(node.Right, env)
		if err != nil {
			return right, hasErr, err
		}
		if right.Type() == object.RETURN_VALUE_OBJ {
			right = right.(*object.ReturnValue).Value[0]
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left, hasErr, err := Eval(node.Left, env)
		if err != nil {
			return left, hasErr, err
		}
		right, hasErr, err := Eval(node.Right, env)
		if err != nil {
			return right, hasErr, err
		}
		if left.Type() == object.RETURN_VALUE_OBJ {
			left = left.(*object.ReturnValue).Value[0]
		}
		if right.Type() == object.RETURN_VALUE_OBJ {
			right = right.(*object.ReturnValue).Value[0]
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PostfixExpression:
		left, hasErr, err := Eval(node.Left, env)
		if err != nil {
			return left, hasErr, err
		}
		if left.Type() == object.RETURN_VALUE_OBJ {
			left = left.(*object.ReturnValue).Value[0]
		}

		resObj, hasErr, err := evalPostfixExpression(node.Operator, left)
		if err != nil {
			return resObj, hasErr, err
		}

		if node.IsStmt {
			if id, ok := node.Left.(*ast.Identifier); ok {
				idVariable, hasErr, err := getIdentifierVariable(id, env)
				if err != nil {
					return idVariable.Value, hasErr, err
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
	case *ast.FunctionBody:
		return evalStatements(node.Statements, env)
	case *ast.IfStatement:
		localEnv := object.NewEnclosedEnvironment(env)
		return evalIfStatements(node, localEnv)
	case *ast.ReturnStatement:
		if node.Value == nil {
			return &object.ReturnValue{Value: nil}, false, nil
		}
		var val []object.Object

		for i := 0; i < len(node.Value); i++ {
			rsObj, hasErr, err := evalReturnValue(node, i, env)
			if err != nil {
				return NULL, hasErr, err
			}
			if rsObj.Type() == object.RETURN_VALUE_OBJ {
				rsObj = rsObj.(*object.ReturnValue).Value[0]
			}
			val = append(val, rsObj)
		}

		return &object.ReturnValue{Value: val}, false, nil
	case *ast.VarStatement:
		return evalVarStatement(node, false, nil, env)
	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, false, nil, env)
	case *ast.ForLoopStatement:
		localEnv := object.NewEnclosedEnvironment(env)
		return evalForLoop(node, localEnv)
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
	case *ast.CallExpression:
		function, hasErr, err := Eval(node.Name, env)
		if err != nil {
			return function, hasErr, err
		}
		args, hasErr, err := evalCallArgs(node.Args, env)
		if err != nil {
			return NULL, hasErr, err
		}
		return applyFunction(function, args)
	case *ast.MultiValueAssignStmt:
		return evalMultiValueAssignStmt(node, env)
	case *ast.ContinueStatement:
		return CONTINUE, false, nil
	case *ast.BreakStatement:
		return BREAK, false, nil
	default:
		return nil, true, fmt.Errorf("no eval function for given node type, got: %T", node)
	}
}

func evalStatements(stmts []ast.Statement, env *object.Environment) (object.Object, bool, error) {
	var result object.Object
	var hasErr bool
	var err error
	for _, statement := range stmts {
		result, hasErr, err = Eval(statement, env)
		if err != nil {
			return NULL, hasErr, err
		}

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result, hasErr, err
		}

		if inForLoop && result == BREAK {
			return BREAK, false, nil
		} else if inForLoop && result == CONTINUE {
			return CONTINUE, false, nil
		}
	}
	return result, hasErr, err
}

func evalMainFunc(node *ast.Function, env *object.Environment) (object.Object, bool, error) {
	resObj, hasErr, err := Eval(node.Body, env)
	if err != nil {
		return resObj, hasErr, err
	}
	return nil, false, nil
}

func evalIndexExpression(left object.Object, index object.Object) (object.Object, bool, error) {
	if left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ {
		return evalArrayIndexExpression(left, index)
	} else {
		return evalHashIndexExpression(left, index)
	}
}

func evalArrayIndexExpression(array object.Object, index object.Object) (object.Object, bool, error) {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	maxIdx := int64(len(arrayObj.Elements) - 1)
	if idx < 0 || idx > maxIdx {
		return NULL, true, errors.New("index out of range, index: " + strconv.FormatInt(idx, 10) + " max index: " + strconv.FormatInt(maxIdx, 10))
	}
	return arrayObj.Elements[idx], false, nil
}

func evalHashIndexExpression(hash object.Object, index object.Object) (object.Object, bool, error) {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return NULL, true, errors.New("unusable as hash key: " + string(index.Type()))
	}
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL, true, errors.New("key not found: " + index.Inspect())
	}
	return pair.Value, false, nil
}

func evalCallArgs(args []ast.Expression, env *object.Environment) ([]object.Object, bool, error) {
	var res []object.Object
	for _, e := range args {
		evaluated, hasErr, err := Eval(e, env)
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
	evaluated, hasErr, err := Eval(function.Body, function.Env)
	if err != nil {
		return evaluated, hasErr, err
	}

	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		if returnValue.Value == nil && function.ReturnType == nil {
			return nil, false, nil
		}
		return returnValue, false, nil
	}

	return evaluated, false, nil
}

func evalArrayValue(node *ast.ArrayValue, env *object.Environment) (object.Object, bool, error) {
	var res []object.Object
	for _, e := range node.Values {
		evaluated, hasErr, err := Eval(e, env)
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

func evalHashMap(node *ast.HashMap, env *object.Environment) (object.Object, bool, error) {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key, hasErr, err := Eval(keyNode, env)
		if err != nil {
			return NULL, hasErr, err
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return NULL, true, errors.New("unusable as hash key: " + string(key.Type()))
		}

		value, hasErr, err := Eval(valueNode, env)
		if err != nil {
			return NULL, hasErr, err
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	if node.KeyType != nil && node.ValueType != nil {
		return &object.Hash{Pairs: pairs, KeyType: node.KeyType.Value, ValueType: node.ValueType.Value}, false, nil
	}

	return &object.Hash{Pairs: pairs}, false, nil
}

func evalVarStatement(node *ast.VarStatement, injectObj bool, obj object.Object, env *object.Environment) (object.Object, bool, error) {
	var val object.Object
	var hasErr bool
	var err error

	if injectObj {
		val = obj
	} else if !injectObj {
		val, hasErr, err = Eval(node.Value, env)
		if err != nil {
			return val, hasErr, err
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

	return val, false, nil
}

func evalMultiValueAssignStmt(node *ast.MultiValueAssignStmt, env *object.Environment) (object.Object, bool, error) {
	if !node.SingleCallExp {
		for _, element := range node.Objects {
			switch element := element.(type) {
			case *ast.VarStatement:
				_, hasErr, err := Eval(element, env)
				if err != nil {
					return NULL, hasErr, err
				}
			case *ast.ExpressionStatement:
				_, hasErr, err := Eval(element.Expression.(*ast.AssignmentExpression), env)
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
			newObj, hasErr, err := Eval(varEntry.Value, env)
			if err != nil {
				return NULL, hasErr, err
			}
			returnObj = newObj
		} else {
			newObj, hasErr, err := Eval(expStmtEntry, env)
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

func evalForLoop(node *ast.ForLoopStatement, env *object.Environment) (object.Object, bool, error) {
	inForLoop = true

	varStmtObj, hasErr, err := Eval(node.Left, env)
	if err != nil {
		return varStmtObj, hasErr, err
	}

	infixObj, hasErr, err := Eval(node.Middle, env)
	if err != nil {
		return infixObj, hasErr, err
	}

	for infixObj == TRUE {
		resStmtObj, hasErr, err := evalStatements(node.Body.Statements, env)
		if err != nil {
			return resStmtObj, hasErr, err
		}

		if resStmtObj == BREAK {
			break
		}

		postfixObj, hasErr, err := Eval(node.Right, env)
		if err != nil {
			return postfixObj, hasErr, err
		}

		env.Update(node.Left.Name.Value, postfixObj, object.VAR)
		infixObj, hasErr, err = Eval(node.Middle, env)
		if err != nil {
			return infixObj, hasErr, err
		}
		if infixObj == FALSE {
			break
		}
	}

	inForLoop = false
	return nil, false, nil
}

func evalReturnValue(rs *ast.ReturnStatement, idx int, env *object.Environment) (object.Object, bool, error) {
	currNode := rs.Value[idx]
	return Eval(currNode, env)
}

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

func evalIfStatements(node *ast.IfStatement, env *object.Environment) (object.Object, bool, error) {
	condition, hasErr, err := Eval(node.Value, env)
	if err != nil {
		return condition, hasErr, err
	}

	conditionRes, hasErr, err := execIf(condition)
	if err != nil {
		return NULL, hasErr, err
	}

	if conditionRes {
		return Eval(node.Body, env)
	} else if node.MultiConseq != nil {
		for i := 0; i < len(node.MultiConseq); i++ {
			resObj, hasErr, err := evalElseIfStatement(node.MultiConseq[i], env)
			if err != nil {
				return resObj, hasErr, err
			}
			if resObj != nil {
				return resObj, hasErr, err
			}
		}
	}

	if node.Consequence != nil {
		return evalStatements(node.Consequence.Body.Statements, env)
	}
	return nil, false, nil
}

func evalElseIfStatement(node *ast.ElseIfStatement, env *object.Environment) (object.Object, bool, error) {
	condition, hasErr, err := Eval(node.Value, env)
	if err != nil {
		return condition, hasErr, err
	}
	conditionRes, hasErr, err := execIf(condition)
	if err != nil {
		return NULL, hasErr, err
	}
	if conditionRes {
		return Eval(node.Body, env)
	}
	return nil, false, nil
}

func execIf(obj object.Object) (bool, bool, error) {
	if obj.Type() == object.RETURN_VALUE_OBJ {
		return execIf(obj.(*object.ReturnValue).Value[0])
	}
	switch obj {
	case TRUE:
		return true, false, nil
	case FALSE:
		return false, false, nil
	default:
		return false, true, errors.New("conditions for `if` and `else if` statements must result in a `bool`, got: " + string(obj.Type()))
	}
}

// -----------------------------------------------------------------------------
// Assignment op
// -----------------------------------------------------------------------------
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

	rightSideObj, hasErr, err := Eval(node.Right, env)
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
		return resObj, hasErr, err
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

// -----------------------------------------------------------------------------
// Prefix op
// -----------------------------------------------------------------------------
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

// -----------------------------------------------------------------------------
// Postfix op
// -----------------------------------------------------------------------------
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

// -----------------------------------------------------------------------------
// Infix op
// -----------------------------------------------------------------------------
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
	case "%":
		return &object.Float{Value: float64(int(leftVal) % int(rightVal))}, false, nil
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
		return NULL, true, errors.New("can only use `+`, `-`, `*`, `/`, `%`, `>`, `<`, `<=`, `>=`, `!=`, `==` infix operators with 2 `float`, got: " + operator)
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
