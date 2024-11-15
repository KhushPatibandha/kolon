package evaluator

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/object"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
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
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.PrefixExpression:
		if postfix, ok := node.Right.(*ast.PostfixExpression); ok {
			left, hasErr, err := Eval(postfix.Left, env)
			if err != nil {
				return left, hasErr, err
			}
			operator := postfix.Operator
			res, hasErr, err := evalPostfixExpression(operator, left)
			if err != nil {
				return res, hasErr, err
			}

			var resVal interface{}

			switch {
			case res.Type() == object.INTEGER_OBJ:
				resVal = res.(*object.Integer).Value
			case res.Type() == object.FLOAT_OBJ:
				resVal = res.(*object.Float).Value
			default:
				return NULL, true, errors.New("Only Integer and Float datatypes supported with Postfix operation. got: " + string(res.Type()))
			}

			switch val := resVal.(type) {
			case int64:
				if operator == "++" {
					return evalPrefixExpression(node.Operator, &object.Integer{Value: val - 2})
				} else {
					return evalPrefixExpression(node.Operator, &object.Integer{Value: val + 2})
				}
			case float64:
				if operator == "++" {
					return evalPrefixExpression(node.Operator, &object.Float{Value: val - 2})
				} else {
					return evalPrefixExpression(node.Operator, &object.Float{Value: val + 2})
				}
			}
		}

		right, hasErr, err := Eval(node.Right, env)
		if err != nil {
			return right, hasErr, err
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
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PostfixExpression:
		left, hasErr, err := Eval(node.Left, env)
		if err != nil {
			return left, hasErr, err
		}

		resObj, hasErr, err := evalPostfixExpression(node.Operator, left)
		if err != nil {
			return resObj, hasErr, err
		}

		// will only update the identifier if the postfix is a stmt.
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
		localEnv := object.NewEnclosedEnvironment(env)
		return evalStatements(node.Statements, localEnv)
	case *ast.IfStatement:
		localEnv := object.NewEnclosedEnvironment(env)
		return evalIfStatements(node, localEnv)
	case *ast.ReturnStatement:
		if node.Value == nil {
			return &object.ReturnValue{Value: nil}, false, nil
		}
		var val []object.Object

		for i := 0; i < len(node.Value); i++ {
			// fmt.Println(node.Value[i])
			rsObj, hasErr, err := evalReturnValue(node, i, env)
			if err != nil {
				return NULL, hasErr, err
			}
			val = append(val, rsObj)
		}

		return &object.ReturnValue{Value: val}, false, nil
	case *ast.VarStatement:
		val, hasErr, err := Eval(node.Value, env)
		if err != nil {
			return val, hasErr, err
		}

		if node.Type.Value == "int" && val.Type() != object.INTEGER_OBJ {
			return NULL, true, errors.New("Identifier declared as int but got: " + string(val.Type()))
		} else if node.Type.Value == "string" && val.Type() != object.STRING_OBJ {
			return NULL, true, errors.New("Identifier declared as string but got: " + string(val.Type()))
		} else if node.Type.Value == "float" && val.Type() != object.FLOAT_OBJ {
			return NULL, true, errors.New("Identifier declared as float but got: " + string(val.Type()))
		} else if node.Type.Value == "char" && val.Type() != object.CHAR_OBJ {
			return NULL, true, errors.New("Identifier declared as char but got: " + string(val.Type()))
		} else if node.Type.Value == "bool" && val.Type() != object.BOOLEAN_OBJ {
			return NULL, true, errors.New("Identifier declared as bool but got: " + string(val.Type()))
		}

		if node.Token.Kind == lexer.VAR {
			env.Set(node.Name.Value, val, object.VAR)
		} else if node.Token.Kind == lexer.CONST {
			env.Set(node.Name.Value, val, object.CONST)
		}

		return val, false, nil
	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, env)
	case *ast.ForLoopStatement:
		localEnv := object.NewEnclosedEnvironment(env)
		return evalForLoop(node, localEnv)
	case *ast.Function:
		// skip all the function execpt main
		if node.Name.Value == "main" {
			// add all the functions in the code from function map to the environment.
			for key, value := range parser.FunctionMap {
				newLocalEnv := object.NewEnclosedEnvironment(env)
				env.Set(key.Value, &object.Function{Name: value.Name, Parameters: value.Parameters, ReturnType: value.ReturnType, Body: value.Body, Env: newLocalEnv}, object.FUNCTION)
			}

			// evaluate main function
			return evalMainFunc(node, env)
		}
		return nil, false, nil
	case *ast.CallExpression:
		function, hasErr, err := Eval(node.Name, env)
		if err != nil {
			return function, hasErr, err
		}
		args, hasErr, err := evalCallArgs(node.Args, env)
		if args == nil || err != nil {
			return NULL, hasErr, err
		}
		return applyFunction(function, args)
	default:
		return nil, true, fmt.Errorf("No Eval function for given node type. got: %T", node)
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
	}
	return result, hasErr, err
}

func evalMainFunc(node *ast.Function, env *object.Environment) (object.Object, bool, error) {
	resObj, hasErr, err := evalStatements(node.Body.Statements, env)
	if err != nil {
		return resObj, hasErr, err
	}
	if resObj.Type() == object.RETURN_VALUE_OBJ {
		// It must be nil, because you cant return anything in main function.
		if resObj.(*object.ReturnValue).Value != nil {
			return NULL, true, errors.New("Can't return anything in main function.")
		}
	}
	return nil, false, nil
}

func evalCallArgs(args []ast.Expression, env *object.Environment) ([]object.Object, bool, error) {
	var res []object.Object
	for _, e := range args {
		evaluated, hasErr, err := Eval(e, env)
		if err != nil {
			return nil, hasErr, err
		}
		res = append(res, evaluated)
	}
	return res, false, nil
}

func applyFunction(fn object.Object, args []object.Object) (object.Object, bool, error) {
	function, ok := fn.(*object.Function)
	if !ok {
		return NULL, true, errors.New("Not a function: " + string(fn.Type()))
	}

	for i, param := range function.Parameters {
		function.Env.Set(param.ParameterName.Value, args[i], object.VAR)
	}
	evaluated, hasErr, err := Eval(function.Body, function.Env)
	if err != nil {
		return evaluated, hasErr, err
	}

	if returnValue, ok := evaluated.(*object.ReturnValue); ok {

		// first check if the return type is nil or not. because you can still use return stmt without any return types defined in the function.
		// just like how we can use return stmt in main function.
		// in this case the return type param in function will be nil and the returnObj.Value will also be nil
		if returnValue.Value == nil && function.ReturnType == nil {
			return nil, false, nil
		}

		// check if you are returning the correct number of values.
		if len(returnValue.Value) != len(function.ReturnType) {
			return NULL, true, errors.New("Number of return values doesn't match.")
		}

		// check if the return types are correct.
		for i, ret := range returnValue.Value {
			if function.ReturnType[i].ReturnType.Value == "int" && ret.Type() != object.INTEGER_OBJ {
				return NULL, true, errors.New("Return type doesn't match. Expected: int at: " + strconv.Itoa(i+1) + "got: " + string(ret.Type()))
			} else if function.ReturnType[i].ReturnType.Value == "string" && ret.Type() != object.STRING_OBJ {
				return NULL, true, errors.New("Return type doesn't match. Expected: string at: " + strconv.Itoa(i+1) + "got: " + string(ret.Type()))
			} else if function.ReturnType[i].ReturnType.Value == "float" && ret.Type() != object.FLOAT_OBJ {
				return NULL, true, errors.New("Return type doesn't match. Expected: float at: " + strconv.Itoa(i+1) + "got: " + string(ret.Type()))
			} else if function.ReturnType[i].ReturnType.Value == "char" && ret.Type() != object.CHAR_OBJ {
				return NULL, true, errors.New("Return type doesn't match. Expected: char at: " + strconv.Itoa(i+1) + "got: " + string(ret.Type()))
			} else if function.ReturnType[i].ReturnType.Value == "bool" && ret.Type() != object.BOOLEAN_OBJ {
				return NULL, true, errors.New("Return type doesn't match. Expected: bool at: " + strconv.Itoa(i+1) + "got: " + string(ret.Type()))
			}
		}

		// TODO: fix this
		return returnValue.Value[0], false, nil
	}

	return evaluated, false, nil
}

func evalForLoop(node *ast.ForLoopStatement, env *object.Environment) (object.Object, bool, error) {
	// Evaluate the VAR stmt in the for loop(.)
	varStmtObj, hasErr, err := Eval(node.Left, env)
	if err != nil {
		return varStmtObj, hasErr, err
	}

	// Get the variable type to check for VAR or CONST
	varVariable, hasErr, err := getIdentifierVariable(node.Left.Name, env)
	if err != nil {
		return NULL, hasErr, err
	}
	if varVariable.Type != object.VAR {
		return NULL, true, errors.New("Can't use CONST to define variable in FOR loop condition")
	}

	// Check if the variable is INT or not
	if varStmtObj.Type() != object.INTEGER_OBJ {
		return NULL, true, errors.New("Can only define variable in FOR loop condition as INT.")
	}

	// Eval infix operation
	infixObj, hasErr, err := Eval(node.Middle, env)
	if err != nil {
		return infixObj, hasErr, err
	}
	// this infix obj should always result in a boolean.
	if infixObj.Type() != object.BOOLEAN_OBJ {
		return NULL, true, errors.New("Infix operation of FOR loop condition should always result in a BOOLEAN.")
	}

	for infixObj == TRUE {
		resStmtObj, hasErr, err := evalStatements(node.Body.Statements, env)
		if err != nil {
			return resStmtObj, hasErr, err
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

	return nil, false, nil
}

func evalReturnValue(rs *ast.ReturnStatement, idx int, env *object.Environment) (object.Object, bool, error) {
	currNode := rs.Value[idx]
	switch currNode.(type) {
	case *ast.Identifier, *ast.IntegerValue, *ast.FloatValue, *ast.BooleanValue, *ast.StringValue, *ast.CharValue, *ast.PrefixExpression, *ast.PostfixExpression, *ast.InfixExpression:
		return Eval(currNode, env)
	default:
		return NULL, true, errors.New("Can Only return expressions and datatypes. got: " + fmt.Sprintf("%T", currNode))
	}
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
		return &object.Variable{Type: object.VAR, Value: NULL}, true, errors.New("Identifier not found: " + node.Value)
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
	switch obj {
	case TRUE:
		return true, false, nil
	case FALSE:
		return false, false, nil
	default:
		return false, true, errors.New("Conditions for 'if' and 'else if' statements must result in a boolean. got: " + string(obj.Type()))
	}
}

func evalAssignmentExpression(node *ast.AssignmentExpression, env *object.Environment) (object.Object, bool, error) {
	switch node.Operator {
	case "=":
		return evalAssignOp(node, env)
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
		return NULL, true, errors.New("Only (=, +=, -=, *=, /=, %=) assignment operators are supported. got: " + node.Operator)
	}
}

func assignOpHelper(node *ast.AssignmentExpression, env *object.Environment) (object.Object, object.Object, bool, bool, error) {
	leftSideVariable, hasErr, err := getIdentifierVariable(node.Left, env)
	if err != nil {
		return NULL, NULL, false, hasErr, err
	}

	// Check if the variable is a constant. if so, return error.
	if leftSideVariable.Type == object.CONST {
		return NULL, NULL, false, true, errors.New("Can't re-assign CONST variables. variable '" + node.Left.Value + "' is a constant.")
	}

	var isVar bool
	if leftSideVariable.Type == object.VAR {
		isVar = true
	}

	// If the variable is not constant, then evaluate the expression on the right.
	rightSideObj, hasErr, err := Eval(node.Right, env)
	if err != nil {
		return NULL, NULL, isVar, hasErr, err
	}

	return leftSideVariable.Value, rightSideObj, isVar, false, nil
}

func evalSymbolAssignOp(operator string, node *ast.AssignmentExpression, env *object.Environment) (object.Object, bool, error) {
	leftSideObj, rightSideObj, isVar, hasErr, err := assignOpHelper(node, env)
	if err != nil {
		return NULL, hasErr, err
	}

	// Only exception.
	if leftSideObj.Type() == object.INTEGER_OBJ && rightSideObj.Type() == object.FLOAT_OBJ {
		return NULL, true, errors.New("Can't convert types. Variable on the left is of type INT and variable on the right is of type FLOAT.")
	}

	if isVar {
		resObj, hasErr, err := evalInfixExpression(operator, leftSideObj, rightSideObj)
		if err != nil {
			return resObj, hasErr, err
		}
		env.Update(node.Left.Value, resObj, object.VAR)
		return resObj, false, nil
	}
	return NULL, true, errors.New("Something went wrong, in evalSymbolAssign.")
}

func evalAssignOp(node *ast.AssignmentExpression, env *object.Environment) (object.Object, bool, error) {
	leftSideObj, rightSideObj, leftIsVar, hasErr, err := assignOpHelper(node, env)
	if err != nil {
		return NULL, hasErr, err
	}

	leftSideObjType := leftSideObj.Type()
	rightSideObjType := rightSideObj.Type()
	if leftIsVar {
		if leftSideObjType == rightSideObjType {
			env.Update(node.Left.Value, rightSideObj, object.VAR)
			return rightSideObj, false, nil
		}
	}
	return NULL, true, errors.New("Can't convert types. Either re-assign the variable or keep the current type. Original type: " + string(leftSideObjType) + " new assigned value's type: " + string(rightSideObjType))
}

func evalPrefixExpression(operator string, right object.Object) (object.Object, bool, error) {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return NULL, true, errors.New("Only 2 Prefix Operator's supported. !(Bang) and -(Dash/Minus). got: " + operator)
	}
}

func evalMinusOperatorExpression(right object.Object) (object.Object, bool, error) {
	if right.Type() != object.INTEGER_OBJ && right.Type() != object.FLOAT_OBJ {
		return NULL, true, errors.New("Dash/Minus(-) operator can be only used with Integers(int) and Floats(float) entities. got: " + string(right.Type()))
	}

	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}, false, nil
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}, false, nil
}

func evalBangOperatorExpression(right object.Object) (object.Object, bool, error) {
	switch right := right.(type) {
	case *object.Boolean:
		if right == TRUE {
			return FALSE, false, nil
		} else {
			return TRUE, false, nil
		}
	default:
		return NULL, true, errors.New("Bang Operator(!) can be only used with boolean(bool) entities. got: " + string(right.Type()))
	}
}

func evalPostfixExpression(operator string, left object.Object) (object.Object, bool, error) {
	switch {
	case left.Type() == object.INTEGER_OBJ:
		return evalIntegerPostfixExpression(operator, left)
	case left.Type() == object.FLOAT_OBJ:
		return evalFloatPostfixExpression(operator, left)
	default:
		return NULL, true, errors.New("Only Integer and Float datatypes supported with Postfix operation. got: " + string(left.Type()))
	}
}

func evalFloatPostfixExpression(operator string, left object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Float).Value

	switch operator {
	case "++":
		return &object.Float{Value: leftVal + 1}, false, nil
	case "--":
		return &object.Float{Value: leftVal - 1}, false, nil
	default:
		return NULL, true, errors.New("Only ++ and -- Postfix operation supported. got: " + operator)
	}
}

func evalIntegerPostfixExpression(operator string, left object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Integer).Value

	switch operator {
	case "++":
		return &object.Integer{Value: leftVal + 1}, false, nil
	case "--":
		return &object.Integer{Value: leftVal - 1}, false, nil
	default:
		return NULL, true, errors.New("Only ++ and -- Postfix operation supported. got: " + operator)
	}
}

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
	default:
		return NULL, true, errors.New("Invalid operation with variable types on left and right.")
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
		return NULL, true, errors.New("Can only perform \"+\", \"==\", \"!=\" with chars. got: " + operator)
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
		return NULL, true, errors.New("Can only perform \"+\", \"==\", \"!=\" with strings. got: " + operator)
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
		return NULL, true, errors.New("Can only perform \"+\", \"-\", \"/\", \"*\", \"%\", \">\", \"<\", \"<=\", \">=\", \"!=\", \"==\" with 2 Float var. got: " + operator)
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
		return NULL, true, errors.New("Can only perform \"+\", \"-\", \"/\", \"*\", \"%\", \"&\", \"|\", \">\", \"<\", \"<=\", \">=\", \"!=\", \"==\" with 2 Integers. got: " + operator)
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
		return NULL, true, errors.New("Can only perform \"&&\", \"||\", \"!=\", \"==\" with 2 Boolean values. got: " + operator)
	}
}
