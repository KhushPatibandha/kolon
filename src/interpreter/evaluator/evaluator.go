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
	InForLoop = false
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
		return evalIndexExpression(left, index)
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
		return evalVarStatement(node, false, nil, env)
	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, false, nil, env)
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
		return nil, true, fmt.Errorf("No Eval function for given node type. got: %T", node)
	}
}

func evalStatements(stmts []ast.Statement, env *object.Environment) (object.Object, bool, error) {
	var result object.Object
	var hasErr bool
	var err error
	// fmt.Println(stmts)
	for _, statement := range stmts {
		result, hasErr, err = Eval(statement, env)
		if err != nil {
			return NULL, hasErr, err
		}

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result, hasErr, err
		}

		if InForLoop && result == BREAK {
			return BREAK, false, nil
		} else if InForLoop && result == CONTINUE {
			return CONTINUE, false, nil
		}
	}
	return result, hasErr, err
}

func evalMainFunc(node *ast.Function, env *object.Environment) (object.Object, bool, error) {
	resObj, hasErr, err := evalStatements(node.Body.Statements, env)
	if err != nil {
		return resObj, hasErr, err
	}
	if resObj != nil && resObj.Type() == object.RETURN_VALUE_OBJ {
		// It must be nil, because you cant return anything in main function.
		if resObj.(*object.ReturnValue).Value != nil {
			return NULL, true, errors.New("Can't return anything in main function.")
		}
	}
	return nil, false, nil
}

func evalIndexExpression(left object.Object, index object.Object) (object.Object, bool, error) {
	if left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ {
		return evalArrayIndexExpression(left, index)
	} else if left.Type() == object.HASH_OBJ {
		return evalHashIndexExpression(left, index)
	}
	return NULL, true, errors.New("Index operator not supported: " + string(left.Type()) + "[" + string(index.Type()) + "]")
}

func evalArrayIndexExpression(array object.Object, index object.Object) (object.Object, bool, error) {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	maxIdx := int64(len(arrayObj.Elements) - 1)
	if idx < 0 || idx > maxIdx {
		return NULL, true, errors.New("Index out of range. Index: " + strconv.FormatInt(idx, 10) + " Max Index: " + strconv.FormatInt(maxIdx, 10))
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
			return NULL, true, errors.New("Not a function: " + string(fn.Type()))
		}
	}

	// check if the type of the args and the type of the params are same or not.
	if len(function.Parameters) != len(args) {
		return NULL, true, errors.New("Number of arguments doesn't match.")
	}
	for i, param := range function.Parameters {
		if args[i].Type() == object.RETURN_VALUE_OBJ {
			args[i] = args[i].(*object.ReturnValue).Value[0]
		}
		if param.ParameterType.IsArray {
			if args[i].Type() != object.ARRAY_OBJ {
				return NULL, true, errors.New("Parameter type doesn't match. Expected: array at: " + strconv.Itoa(i+1) + "got: " + string(args[i].Type()))
			}
			arrayType := param.ParameterType.Value

			// validate if the array has the object of type return type or not
			_, err := validateArrayTypes(*args[i].(*object.Array), arrayType)
			if err != nil {
				return NULL, true, err
			}
			args[i].(*object.Array).TypeOf = arrayType
		} else if param.ParameterType.IsHash {
			if args[i].Type() != object.HASH_OBJ {
				return NULL, true, errors.New("Parameter type doesn't match. Expected: hash at: " + strconv.Itoa(i+1) + "got: " + string(args[i].Type()))
			}
			keyType := param.ParameterType.SubTypes[0].Value
			valueType := param.ParameterType.SubTypes[1].Value
			for _, pair := range args[i].(*object.Hash).Pairs {
				err := validateMapTypes(pair.Key, pair.Value, keyType, valueType)
				if err != nil {
					return NULL, true, err
				}
			}
			args[i].(*object.Hash).KeyType = keyType
			args[i].(*object.Hash).ValueType = valueType
		} else if param.ParameterType.Value == "int" && args[i].Type() != object.INTEGER_OBJ {
			return NULL, true, errors.New("Parameter type doesn't match. Expected: int at: " + strconv.Itoa(i+1) + "got: " + string(args[i].Type()))
		} else if param.ParameterType.Value == "string" && args[i].Type() != object.STRING_OBJ {
			return NULL, true, errors.New("Parameter type doesn't match. Expected: string at: " + strconv.Itoa(i+1) + "got: " + string(args[i].Type()))
		} else if param.ParameterType.Value == "float" && args[i].Type() != object.FLOAT_OBJ {
			return NULL, true, errors.New("Parameter type doesn't match. Expected: float at: " + strconv.Itoa(i+1) + "got: " + string(args[i].Type()))
		} else if param.ParameterType.Value == "char" && args[i].Type() != object.CHAR_OBJ {
			return NULL, true, errors.New("Parameter type doesn't match. Expected: char at: " + strconv.Itoa(i+1) + "got: " + string(args[i].Type()))
		} else if param.ParameterType.Value == "bool" && args[i].Type() != object.BOOLEAN_OBJ {
			return NULL, true, errors.New("Parameter type doesn't match. Expected: bool at: " + strconv.Itoa(i+1) + "got: " + string(args[i].Type()))
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
			if ret.Type() == object.RETURN_VALUE_OBJ {
				ret = ret.(*object.ReturnValue).Value[0]
			}
			if function.ReturnType[i].ReturnType.IsArray {
				if ret.Type() != object.ARRAY_OBJ {
					return NULL, true, errors.New("Return type doesn't match. Expected: array at: " + strconv.Itoa(i+1) + "got: " + string(ret.Type()))
				}
				arrayReturnType := function.ReturnType[i].ReturnType.Value

				// validate if the array has the object of type return type or not
				_, err := validateArrayTypes(*ret.(*object.Array), arrayReturnType)
				if err != nil {
					return NULL, true, err
				}
				ret.(*object.Array).TypeOf = arrayReturnType
			} else if function.ReturnType[i].ReturnType.IsHash {
				if ret.Type() != object.HASH_OBJ {
					return NULL, true, errors.New("Return type doesn't match. Expected: hash at: " + strconv.Itoa(i+1) + "got: " + string(ret.Type()))
				}
				keyType := function.ReturnType[i].ReturnType.SubTypes[0].Value
				valueType := function.ReturnType[i].ReturnType.SubTypes[1].Value
				for _, pair := range ret.(*object.Hash).Pairs {
					err := validateMapTypes(pair.Key, pair.Value, keyType, valueType)
					if err != nil {
						return NULL, true, err
					}
				}
				ret.(*object.Hash).KeyType = keyType
				ret.(*object.Hash).ValueType = valueType
			} else if function.ReturnType[i].ReturnType.Value == "int" && ret.Type() != object.INTEGER_OBJ {
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

	arrayObj := object.Array{Elements: res}

	// we have the context.
	if node.Type != nil {
		_, err := validateArrayTypes(arrayObj, node.Type.Value)
		if err != nil {
			return NULL, true, err
		}
		arrayObj.TypeOf = node.Type.Value
		return &arrayObj, false, nil
	}

	return &arrayObj, false, nil
}

func validateArrayTypes(arrayObj object.Array, typeOf string) (bool, error) {
	for _, element := range arrayObj.Elements {
		if typeOf == "int" && element.Type() != object.INTEGER_OBJ {
			return false, errors.New("Array declared as int but got an element of type: " + string(element.Type()))
		} else if typeOf == "string" && element.Type() != object.STRING_OBJ {
			return false, errors.New("Array declared as string but got an element of type: " + string(element.Type()))
		} else if typeOf == "float" && element.Type() != object.FLOAT_OBJ {
			return false, errors.New("Array declared as float but got an element of type: " + string(element.Type()))
		} else if typeOf == "char" && element.Type() != object.CHAR_OBJ {
			return false, errors.New("Array declared as char but got an element of type: " + string(element.Type()))
		} else if typeOf == "bool" && element.Type() != object.BOOLEAN_OBJ {
			return false, errors.New("Array declared as bool but got an element of type: " + string(element.Type()))
		}
	}
	return true, nil
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
		// we are given the types so we should validate them.
		for _, pair := range pairs {
			err := validateMapTypes(pair.Key, pair.Value, node.KeyType.Value, node.ValueType.Value)
			if err != nil {
				return NULL, true, err
			}
		}
		return &object.Hash{Pairs: pairs, KeyType: node.KeyType.Value, ValueType: node.ValueType.Value}, false, nil
	}

	return &object.Hash{Pairs: pairs}, false, nil
}

func validateMapTypes(key object.Object, value object.Object, keyType string, valueType string) error {
	if keyType == "int" && key.Type() != object.INTEGER_OBJ {
		return errors.New("Map declared as int but got key of type: " + string(key.Type()))
	} else if keyType == "string" && key.Type() != object.STRING_OBJ {
		return errors.New("Map declared as string but got key of type: " + string(key.Type()))
	} else if keyType == "float" && key.Type() != object.FLOAT_OBJ {
		return errors.New("Map declared as float but got key of type: " + string(key.Type()))
	} else if keyType == "char" && key.Type() != object.CHAR_OBJ {
		return errors.New("Map declared as char but got key of type: " + string(key.Type()))
	} else if keyType == "bool" && key.Type() != object.BOOLEAN_OBJ {
		return errors.New("Map declared as bool but got key of type: " + string(key.Type()))
	}
	if valueType == "int" && value.Type() != object.INTEGER_OBJ {
		return errors.New("Map declared as int but got value of type: " + string(value.Type()))
	} else if valueType == "string" && value.Type() != object.STRING_OBJ {
		return errors.New("Map declared as string but got value of type: " + string(value.Type()))
	} else if valueType == "float" && value.Type() != object.FLOAT_OBJ {
		return errors.New("Map declared as float but got value of type: " + string(value.Type()))
	} else if valueType == "char" && value.Type() != object.CHAR_OBJ {
		return errors.New("Map declared as char but got value of type: " + string(value.Type()))
	} else if valueType == "bool" && value.Type() != object.BOOLEAN_OBJ {
		return errors.New("Map declared as bool but got value of type: " + string(value.Type()))
	}
	return nil
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

	// The idea is, if we only have a single value on the left, then we can directly assign the 0th element of the return object to the val.
	for val.Type() == object.RETURN_VALUE_OBJ {
		val = val.(*object.ReturnValue).Value[0]
	}

	if node.Type.IsArray {
		if val.Type() != object.ARRAY_OBJ {
			return NULL, true, errors.New("Identifier declared as array but got: " + string(val.Type()))
		}
		arrayType := node.Type.Value
		if val.(*object.Array).TypeOf != arrayType {
			return NULL, true, errors.New("Array type doesn't match. Expected: " + arrayType + " got: " + val.(*object.Array).TypeOf)
		}
	} else if node.Type.IsHash {
		if val.Type() != object.HASH_OBJ {
			return NULL, true, errors.New("Identifier declared as hash but got: " + string(val.Type()))
		}
		keyType := node.Type.SubTypes[0].Value
		valueType := node.Type.SubTypes[1].Value
		if val.(*object.Hash).KeyType != keyType || val.(*object.Hash).ValueType != valueType {
			return NULL, true, errors.New("Hash type doesn't match. Expected: " + keyType + " -> " + valueType + " got: " + val.(*object.Hash).KeyType + " -> " + val.(*object.Hash).ValueType)
		}
	} else if node.Type.Value == "int" && val.Type() != object.INTEGER_OBJ {
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
}

func evalMultiValueAssignStmt(node *ast.MultiValueAssignStmt, env *object.Environment) (object.Object, bool, error) {
	// Check if the right side is bunch of expressions or a call expression.
	isCallExpr := false
	isVar := false
	varEntry, ok := node.Objects[0].(*ast.VarStatement)
	var expStmtEntry *ast.CallExpression
	if ok {
		isVar = true
		_, ok := varEntry.Value.(*ast.CallExpression)
		if ok {
			isCallExpr = true
		}
	} else {
		isVar = false
		expStmtEntry, ok = node.Objects[0].(*ast.ExpressionStatement).Expression.(*ast.AssignmentExpression).Right.(*ast.CallExpression)
		if ok {
			isCallExpr = true
		}
	}

	if !isCallExpr {
		// loop throught all the elements in the mvas object param and execute them.
		for _, element := range node.Objects {
			switch element := element.(type) {
			case *ast.VarStatement:
				_, hasErr, err := Eval(element, env)
				if err != nil {
					return NULL, hasErr, err
				}
			case *ast.ExpressionStatement:
				// simply update the value of the variable.
				_, hasErr, err := Eval(element.Expression.(*ast.AssignmentExpression), env)
				if err != nil {
					return NULL, hasErr, err
				}
			}
		}
	} else {
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

		// return object will always have the same number of elements that are required(i.e on the left side of the assignment)
		// because we are checking that at the time of parsing.
		// we just have to update the values of the variables.
		for i, element := range node.Objects {
			switch element := element.(type) {
			case *ast.VarStatement:
				_, hasErr, err := evalVarStatement(element, true, returnObjList[i], env)
				if err != nil {
					return NULL, hasErr, err
				}
			case *ast.ExpressionStatement:
				// simply update the value of the variable.
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
	InForLoop = true

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

	InForLoop = false
	return nil, false, nil
}

func evalReturnValue(rs *ast.ReturnStatement, idx int, env *object.Environment) (object.Object, bool, error) {
	currNode := rs.Value[idx]
	switch currNode.(type) {
	case *ast.Identifier, *ast.IntegerValue, *ast.FloatValue, *ast.BooleanValue, *ast.StringValue, *ast.CharValue, *ast.PrefixExpression, *ast.PostfixExpression, *ast.InfixExpression, *ast.ArrayValue, *ast.HashMap, *ast.CallExpression, *ast.IndexExpression:
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
		// if not found, check if builtin method exists.
		builtin, ok := builtins[node.Value]
		if ok {
			return &object.Variable{Type: object.FUNCTION, Value: builtin}, false, nil
		} else {
			return &object.Variable{Type: object.VAR, Value: NULL}, true, errors.New("Identifier not found: " + node.Value)
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
	switch obj {
	case TRUE:
		return true, false, nil
	case FALSE:
		return false, false, nil
	default:
		return false, true, errors.New("Conditions for 'if' and 'else if' statements must result in a boolean. got: " + string(obj.Type()))
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
		return NULL, true, errors.New("Only (=, +=, -=, *=, /=, %=) assignment operators are supported. got: " + node.Operator)
	}
}

func assignOpHelper(node *ast.AssignmentExpression, injectObj bool, rightVal object.Object, env *object.Environment) (object.Object, object.Object, bool, bool, error) {
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

	// If we are injecting than no need to evaluate the right side
	if injectObj {
		return leftSideVariable.Value, rightVal, isVar, false, nil
	}

	// If the variable is not constant, then evaluate the expression on the right.
	rightSideObj, hasErr, err := Eval(node.Right, env)
	if err != nil {
		return NULL, NULL, isVar, hasErr, err
	}

	return leftSideVariable.Value, rightSideObj, isVar, false, nil
}

func evalSymbolAssignOp(operator string, node *ast.AssignmentExpression, env *object.Environment) (object.Object, bool, error) {
	leftSideObj, rightSideObj, isVar, hasErr, err := assignOpHelper(node, false, nil, env)
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

func evalAssignOp(node *ast.AssignmentExpression, injectObj bool, rightVal object.Object, env *object.Environment) (object.Object, bool, error) {
	leftSideObj, rightSideObj, leftIsVar, hasErr, err := assignOpHelper(node, injectObj, rightVal, env)
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

// -----------------------------------------------------------------------------
// Prefix op
// -----------------------------------------------------------------------------
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

// -----------------------------------------------------------------------------
// Postfix op
// -----------------------------------------------------------------------------
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

// -----------------------------------------------------------------------------
// Infix op
// -----------------------------------------------------------------------------
func evalInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	// if left or right is return value object, then convert it into the data type object first.
	// but if the return value object has multiple values, then return error.
	if left.Type() == object.RETURN_VALUE_OBJ {
		if len(left.(*object.ReturnValue).Value) > 1 {
			return NULL, true, errors.New("Can't use multiple return values in infix operation.")
		}
		left = left.(*object.ReturnValue).Value[0]
	}
	if right.Type() == object.RETURN_VALUE_OBJ {
		if len(right.(*object.ReturnValue).Value) > 1 {
			return NULL, true, errors.New("Can't use multiple return values in infix operation.")
		}
		right = right.(*object.ReturnValue).Value[0]
	}
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
		if left.(*object.Array).TypeOf != right.(*object.Array).TypeOf {
			return NULL, true, errors.New("Array types don't match. Expected: " + left.(*object.Array).TypeOf + " got: " + right.(*object.Array).TypeOf)
		}
		return evalArrayInfixExpression(operator, left, right)
	default:
		return NULL, true, errors.New("Invalid operation with variable types on left and right.")
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
		return NULL, true, errors.New("Can only perform \"+\", \"==\", \"!=\" with arrays. got: " + operator)
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
