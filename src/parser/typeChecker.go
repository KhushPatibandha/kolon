package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

var (
	currFunction *ast.Function
	builtinMap   = map[string]bool{
		"print":       true,
		"println":     true,
		"len":         true,
		"toString":    true,
		"push":        true,
		"pop":         true,
		"insert":      true,
		"remove":      true,
		"getIndex":    true,
		"keys":        true,
		"values":      true,
		"containsKey": true,
		"typeOf":      true,
		"slice":       true,
	}
	inTesting bool
)

type expType struct {
	Type    ast.Type
	CallExp bool
}

func TypeCheckProgram(program *ast.Program, env *Environment, inTest bool) error {
	inTesting = inTest
	for key, value := range FunctionMap {
		newLocalEnv := NewEnclosedEnvironment(env)
		for _, param := range value.Parameters {
			err := newLocalEnv.Set(param.ParameterName.Value, *param.ParameterName, *param.ParameterType, VAR, nil)
			if err != nil {
				return err
			}
		}
		err := env.Set(key, *value.Name, ast.Type{}, FUNCTION, newLocalEnv)
		if err != nil {
			return err
		}
	}
	return checkStmts(program.Statements, env)
}

func typeCheckStmts(stmtNode ast.Statement, env *Environment) error {
	switch node := stmtNode.(type) {
	case *ast.FunctionBody:
		return checkStmts(node.Statements, env)
	case *ast.MultiValueAssignStmt:
		return checkMultiAssignStmt(node, env)
	case *ast.VarStatement:
		return checkVarStmt(node, env)
	case *ast.ReturnStatement:
		return checkReturnStmt(node, env)
	case *ast.ContinueStatement:
		if !inForLoop {
			return errors.New("continue statement can only be used inside a `for loop`")
		}
		return nil
	case *ast.BreakStatement:
		if !inForLoop {
			return errors.New("break statement can only be used inside a `for loop`")
		}
		return nil
	case *ast.Function:
		currFunction = node
		if node.Name.Value == "main" {
			return checkMainFunStmt(node, env)
		}
		return checkFunStmt(node, env)
	case *ast.IfStatement:
		return checkIfStmt(node, env)
	case *ast.ForLoopStatement:
		localEnv := NewEnclosedEnvironment(env)
		return checkForLoopStmt(node, localEnv)
	case *ast.WhileLoopStatement:
		localEnv := NewEnclosedEnvironment(env)
		return checkWhileLoopStmt(node, localEnv)
	case *ast.ExpressionStatement:
		return checkExpStmt(node, env)
	default:
		msg := fmt.Errorf("unknown statement type, got: %T", node)
		return msg
	}
}

func getExpType(exp ast.Expression, env *Environment) (expType, error) {
	switch exp := exp.(type) {
	case *ast.HashMap:
		return checkHashmapExp(exp, env)
	case *ast.ArrayValue:
		return checkArrayExp(exp, env)
	case *ast.IntegerValue:
		return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.INT, Value: "int"}, Value: "int", IsArray: false, IsHash: false, SubTypes: nil}}, nil
	case *ast.FloatValue:
		return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.FLOAT, Value: "float"}, Value: "float", IsArray: false, IsHash: false, SubTypes: nil}}, nil
	case *ast.StringValue:
		return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.STRING, Value: "string"}, Value: "string", IsArray: false, IsHash: false, SubTypes: nil}}, nil
	case *ast.CharValue:
		return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.CHAR, Value: "char"}, Value: "char", IsArray: false, IsHash: false, SubTypes: nil}}, nil
	case *ast.BooleanValue:
		return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}}, nil
	case *ast.Identifier:
		return checkIdentExp(exp, env)
	case *ast.PrefixExpression:
		return checkPrefixExp(exp, env)
	case *ast.PostfixExpression:
		return checkPostfixExp(exp, env)
	case *ast.InfixExpression:
		return checkInfixExp(exp, env)
	case *ast.AssignmentExpression:
		return checkAssignExp(exp, env)
	case *ast.CallExpression:
		return checkCallExp(exp, env)
	case *ast.IndexExpression:
		return checkIndexExp(exp, env)
	default:
		msg := fmt.Errorf("unknown expression type, got: %T", exp)
		return expType{}, msg
	}
}

func checkStmts(stmts []ast.Statement, env *Environment) error {
	for _, stmt := range stmts {
		err := typeCheckStmts(stmt, env)
		if err != nil {
			return err
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// Statement Checking
// -----------------------------------------------------------------------------
func checkVarStmt(node *ast.VarStatement, env *Environment) error {
	definedType := node.Type

	expType, err := getExpType(node.Value, env)
	if err != nil {
		return err
	}

	if expType.CallExp {
		function := FunctionMap[node.Value.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return errors.New("call expression must return only 1 value for `var` statement, got: " + strconv.Itoa(len(function.ReturnType)))
		}
		expType.Type = *function.ReturnType[0].ReturnType
	}

	err = varTypeCheckerHelper(*definedType, expType.Type)
	if err != nil {
		return err
	}

	if node.Token.Kind == lexer.VAR {
		err := env.Set(node.Name.Value, *node.Name, *node.Type, VAR, nil)
		if err != nil {
			return err
		}
	} else if node.Token.Kind == lexer.CONST {
		err := env.Set(node.Name.Value, *node.Name, *node.Type, CONST, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkMultiAssignStmt(node *ast.MultiValueAssignStmt, env *Environment) error {
	if !node.SingleCallExp {
		for _, element := range node.Objects {
			switch element := element.(type) {
			case *ast.VarStatement:
				err := checkVarStmt(element, env)
				if err != nil {
					return err
				}
			case *ast.ExpressionStatement:
				_, err := getExpType(element.Expression.(*ast.AssignmentExpression), env)
				if err != nil {
					return err
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

		var function *ast.Function
		if isVar {
			function = FunctionMap[varEntry.Value.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		} else {
			function = FunctionMap[expStmtEntry.Name.(*ast.Identifier).Value]
		}
		returnTypes := function.ReturnType

		if len(returnTypes) != len(node.Objects) {
			return errors.New("number of expression on the right side of multi-assignment do not match the number of declarations on the left, left=" + fmt.Sprint(len(node.Objects)) + ", right=" + fmt.Sprint(len(returnTypes)))
		}

		for i, obj := range node.Objects {
			switch obj := obj.(type) {
			case *ast.VarStatement:
				err := varTypeCheckerHelper(*obj.Type, *returnTypes[i].ReturnType)
				if err != nil {
					return err
				}
				if obj.Token.Kind == lexer.VAR {
					err := env.Set(obj.Name.Value, *obj.Name, *obj.Type, VAR, nil)
					if err != nil {
						return err
					}
				} else if obj.Token.Kind == lexer.CONST {
					err := env.Set(obj.Name.Value, *obj.Name, *obj.Type, CONST, nil)
					if err != nil {
						return err
					}
				}
			case *ast.ExpressionStatement:
				variable, ok := env.Get(obj.Expression.(*ast.AssignmentExpression).Left.Value)
				if !ok {
					return errors.New("variable `" + obj.Expression.(*ast.AssignmentExpression).Left.Value + "` not found, can't assign value to undefined variables")
				}
				if variable.VarType == CONST {
					return errors.New("variable `" + obj.Expression.(*ast.AssignmentExpression).Left.Value + "` is a constant, can't assign value to a constant variable")
				} else if variable.VarType == FUNCTION {
					return errors.New("variable `" + obj.Expression.(*ast.AssignmentExpression).Left.Value + "` is a function, can't assign value to a function")
				}
				definedVarType := variable.Type
				err := varTypeCheckerHelper(definedVarType, *returnTypes[i].ReturnType)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func checkExpStmt(node *ast.ExpressionStatement, env *Environment) error {
	_, err := getExpType(node.Expression, env)
	if err != nil {
		return err
	}
	return nil
}

func checkIfStmt(node *ast.IfStatement, env *Environment) error {
	localEnvForIf := NewEnclosedEnvironment(env)

	expType, err := getExpType(node.Value, localEnvForIf)
	if err != nil {
		return err
	}

	if expType.CallExp {
		function := FunctionMap[node.Value.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return errors.New("call expression must return only 1 value for `if` statement, got: " + strconv.Itoa(len(function.ReturnType)))
		}
		expType.Type = *function.ReturnType[0].ReturnType
	}
	if expType.Type.Value != "bool" {
		return errors.New("condition for `if` statement must always result in a boolean value, got: " + expType.Type.Value)
	}

	err = checkStmts(node.Body.Statements, localEnvForIf)
	if err != nil {
		return err
	}

	if node.MultiConseq != nil {
		for _, elseIfStmt := range node.MultiConseq {
			localEnfForElseIf := NewEnclosedEnvironment(env)
			expType, err := getExpType(elseIfStmt.Value, localEnfForElseIf)
			if err != nil {
				return err
			}
			if expType.CallExp {
				function := FunctionMap[elseIfStmt.Value.(*ast.CallExpression).Name.(*ast.Identifier).Value]
				if len(function.ReturnType) != 1 {
					return errors.New("call expression must return only 1 value for `else if` statement, got: " + strconv.Itoa(len(function.ReturnType)))
				}
				expType.Type = *function.ReturnType[0].ReturnType
			}
			if expType.Type.Value != "bool" {
				return errors.New("condition for `else if` statement must always result in a boolean value, got: " + expType.Type.Value)
			}

			err = checkStmts(elseIfStmt.Body.Statements, localEnfForElseIf)
			if err != nil {
				return err
			}
		}
	}

	if node.Consequence != nil {
		localEnvForElse := NewEnclosedEnvironment(env)
		err = checkStmts(node.Consequence.Body.Statements, localEnvForElse)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkForLoopStmt(node *ast.ForLoopStatement, env *Environment) error {
	inForLoop = true

	if _, ok := node.Left.(*ast.VarStatement); ok {
		err := checkVarStmt(node.Left.(*ast.VarStatement), env)
		if err != nil {
			return err
		}
		varVariable, ok := env.Get(node.Left.(*ast.VarStatement).Name.Value)
		if !ok {
			return errors.New("`var` variable in `for loop` condition not found")
		}
		if varVariable.VarType != VAR {
			return errors.New("can't use `const` to define variable in `for loop` condition")
		}
		if varVariable.Type.Value != "int" {
			return errors.New("can only define variable in `for loop` condition as `int`, got: " + varVariable.Type.Value)
		}
	} else {
		err := checkExpStmt(node.Left.(*ast.ExpressionStatement), env)
		if err != nil {
			return err
		}
	}

	expType, err := getExpType(node.Middle, env)
	if err != nil {
		return err
	}
	if expType.Type.Value != "bool" {
		return errors.New("infix operation for `for loop` condition should always result in a `bool`, got: " + expType.Type.Value)
	}

	expType, err = getExpType(node.Right, env)
	if err != nil {
		return err
	}
	if expType.Type.Value != "int" {
		if _, ok := node.Right.(*ast.AssignmentExpression); ok {
			return errors.New("assignment operation for `for loop` condition should always result in an `int`, got: " + expType.Type.Value)
		}
		return errors.New("postfix operation for `for loop` condition should always result in an `int`, got: " + expType.Type.Value)
	}

	err = checkStmts(node.Body.Statements, env)
	if err != nil {
		return err
	}

	inForLoop = false
	return nil
}

func checkWhileLoopStmt(node *ast.WhileLoopStatement, env *Environment) error {
	inForLoop = true

	expType, err := getExpType(node.Condition, env)
	if err != nil {
		return err
	}
	if expType.CallExp {
		function := FunctionMap[node.Condition.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return errors.New("call expression must return only 1 value for `while` statement's break condition, got: " + strconv.Itoa(len(function.ReturnType)))
		}
		expType.Type = *function.ReturnType[0].ReturnType
	}
	if expType.Type.Value != "bool" {
		return errors.New("break condition for `while loop` should always result in a `bool`, got: " + expType.Type.Value)
	}

	err = checkStmts(node.Body.Statements, env)
	if err != nil {
		return err
	}

	inForLoop = false
	return nil
}

func checkFunStmt(node *ast.Function, env *Environment) error {
	funVariable, ok := env.Get(node.Name.Value)
	if !ok {
		return errors.New("function `" + node.Name.Value + "` not found in the environment")
	}
	if funVariable.VarType != FUNCTION {
		return errors.New("function `" + node.Name.Value + "` not found in the environment")
	}
	localFunEnv := funVariable.Env

	err := checkStmts(node.Body.Statements, localFunEnv)
	if err != nil {
		return err
	}

	if node.ReturnType != nil {
		err := checkAvailableReturnStmt(node.Body.Statements)
		if err != nil {
			return errors.New(err.Error() + " for function: " + node.Name.Value)
		}
	}
	return nil
}

func checkAvailableReturnStmt(stmts []ast.Statement) error {
	lastStmt := stmts[len(stmts)-1]
	switch node := lastStmt.(type) {
	case *ast.ReturnStatement:
		return nil
	case *ast.IfStatement:
		if node.Consequence == nil {
			return errors.New("missing `return` statement")
		}
		ifErr := checkAvailableReturnStmt(node.Body.Statements)
		if ifErr != nil {
			return ifErr
		}
		if node.MultiConseq != nil {
			for _, elseIfStmt := range node.MultiConseq {
				elseIfErr := checkAvailableReturnStmt(elseIfStmt.Body.Statements)
				if elseIfErr != nil {
					return elseIfErr
				}
			}
		}
		elseErr := checkAvailableReturnStmt(node.Consequence.Body.Statements)
		if elseErr != nil {
			return elseErr
		}
		return nil
	default:
		return errors.New("missing `return` statement")
	}
}

func checkMainFunStmt(node *ast.Function, env *Environment) error {
	funVariable, ok := env.Get(node.Name.Value)
	if !ok {
		return errors.New("`main` function not found in the environment")
	}
	if funVariable.VarType != FUNCTION {
		return errors.New("`main` function not found in the environment")
	}
	localMainEnv := funVariable.Env

	if len(node.Parameters) != 0 {
		return errors.New("`main` function must not have any parameters")
	}

	if node.ReturnType != nil {
		return errors.New("can't return anything from `main` function")
	}

	err := checkStmts(node.Body.Statements, localMainEnv)
	if err != nil {
		return err
	}
	return nil
}

func checkReturnStmt(node *ast.ReturnStatement, env *Environment) error {
	if inTesting {
		err := returnTypeCheckerHelper(node, env)
		if err != nil {
			return err
		}
		return nil
	}

	if currFunction.Name.Value == "main" {
		if node.Value != nil {
			return errors.New("can't return anything from `main` function")
		}
		return nil
	}

	if node.Value == nil && currFunction.ReturnType == nil {
		return nil
	}
	if node.Value != nil && currFunction.ReturnType != nil {
		if len(node.Value) != len(currFunction.ReturnType) {
			return errors.New("number of return values does not match the number of return types, got: " + strconv.Itoa(len(node.Value)) + ", expected: " + strconv.Itoa(len(currFunction.ReturnType)))
		}
	} else {
		return errors.New("number of return values does not match the number of return types")
	}

	err := returnTypeCheckerHelper(node, env)
	if err != nil {
		return err
	}
	return nil
}

// -----------------------------------------------------------------------------
// Expression Checking
// -----------------------------------------------------------------------------
func checkArrayExp(node *ast.ArrayValue, env *Environment) (expType, error) {
	if len(node.Values) == 0 {
		return expType{Type: ast.Type{Value: "", IsArray: true, IsHash: false, SubTypes: nil}, CallExp: false}, nil
	}
	typeStr := ""
	for _, element := range node.Values {
		exp, err := getExpType(element, env)
		if err != nil {
			return expType{}, err
		}
		if typeStr == "" {
			typeStr = exp.Type.Value
		} else {
			if typeStr != exp.Type.Value {
				return expType{}, errors.New("array can only have one type of elements, got: " + typeStr + " and " + exp.Type.Value)
			}
		}
	}

	return expType{Type: ast.Type{Value: typeStr, IsArray: true, IsHash: false, SubTypes: nil}, CallExp: false}, nil
}

func checkHashmapExp(hashmap *ast.HashMap, env *Environment) (expType, error) {
	if len(hashmap.Pairs) == 0 {
		return expType{Type: ast.Type{IsArray: false, IsHash: true, SubTypes: nil}, CallExp: false}, nil
	}
	keyTypeStr := ""
	valueTypeStr := ""
	for key, value := range hashmap.Pairs {
		keyExp, err := getExpType(key, env)
		if err != nil {
			return expType{}, err
		}
		valueExp, err := getExpType(value, env)
		if err != nil {
			return expType{}, err
		}
		if keyTypeStr == "" && valueTypeStr == "" {
			keyTypeStr = keyExp.Type.Value
			valueTypeStr = valueExp.Type.Value
		} else {
			if keyTypeStr != keyExp.Type.Value {
				return expType{}, errors.New("hashmap can only have one type of key, got: " + keyTypeStr + " and " + keyExp.Type.Value)
			}
			if valueTypeStr != valueExp.Type.Value {
				return expType{}, errors.New("hashmap can only have one type of value, got: " + valueTypeStr + " and " + valueExp.Type.Value)
			}
		}
	}
	var subTypes []*ast.Type
	subTypes = append(subTypes, &ast.Type{Value: keyTypeStr, IsArray: false, IsHash: false, SubTypes: nil})
	subTypes = append(subTypes, &ast.Type{Value: valueTypeStr, IsArray: false, IsHash: false, SubTypes: nil})
	return expType{Type: ast.Type{IsArray: false, IsHash: true, SubTypes: subTypes}, CallExp: false}, nil
}

func checkIdentExp(ident *ast.Identifier, env *Environment) (expType, error) {
	variable, ok := env.Get(ident.Value)
	if !ok {
		return expType{}, errors.New("variable `" + ident.Value + "` is undefined/not found")
	}
	return expType{Type: variable.Type, CallExp: false}, nil
}

func checkPrefixExp(prefix *ast.PrefixExpression, env *Environment) (expType, error) {
	rightExpType, err := getExpType(prefix.Right, env)
	if err != nil {
		return expType{}, err
	}
	if rightExpType.CallExp {
		function := FunctionMap[prefix.Right.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return expType{}, errors.New("call expression must return only 1 value for `prefix` expression, got: " + strconv.Itoa(len(function.ReturnType)))
		}
		rightExpType.Type = *function.ReturnType[0].ReturnType
	}
	if rightExpType.Type.IsArray || rightExpType.Type.IsHash {
		return expType{}, errors.New("prefix operator can't be used with array or hashmap")
	}
	switch prefix.Operator {
	case "!":
		if rightExpType.Type.Value != "bool" {
			return expType{}, errors.New("bang operator (`!`) can be only used with `bool` entities, got: " + rightExpType.Type.Value)
		}
	case "-":
		if rightExpType.Type.Value != "int" && rightExpType.Type.Value != "float" {
			return expType{}, errors.New("dash/minus (`-`) operator can be only used with `int` and `float` entities, got: " + rightExpType.Type.Value)
		}
	default:
		return expType{}, errors.New("only 2 `prefix` operator's supported, bang (`!`) and dash/minus (`-`), got: " + prefix.Operator)
	}

	return expType{Type: rightExpType.Type, CallExp: false}, nil
}

func checkPostfixExp(postfix *ast.PostfixExpression, env *Environment) (expType, error) {
	leftExpType, err := getExpType(postfix.Left, env)
	if err != nil {
		return expType{}, err
	}
	if leftExpType.CallExp {
		function := FunctionMap[postfix.Left.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return expType{}, errors.New("call expression must return only 1 value for `postfix` expression, got: " + strconv.Itoa(len(function.ReturnType)))
		}
		leftExpType.Type = *function.ReturnType[0].ReturnType
	}
	if leftExpType.Type.IsArray || leftExpType.Type.IsHash {
		return expType{}, errors.New("postfix operator can't be used with array or hashmap")
	}
	if leftExpType.Type.Value != "int" && leftExpType.Type.Value != "float" {
		return expType{}, errors.New("only `int` and `float` datatypes supported with `postfix` operation, got: " + leftExpType.Type.Value)
	}
	if postfix.Operator != "++" && postfix.Operator != "--" {
		return expType{}, errors.New("only 2 `postfix` operator's supported, increment (`++`) and decrement (`--`), got: " + postfix.Operator)
	}

	return expType{Type: leftExpType.Type, CallExp: false}, nil
}

func checkInfixExp(infix *ast.InfixExpression, env *Environment) (expType, error) {
	leftExpType, err := getExpType(infix.Left, env)
	if err != nil {
		return expType{}, err
	}
	rightExpType, err := getExpType(infix.Right, env)
	if err != nil {
		return expType{}, err
	}
	if leftExpType.CallExp {
		function := FunctionMap[infix.Left.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return expType{}, errors.New("call expression must return only 1 value for `infix` expression, got: " + strconv.Itoa(len(function.ReturnType)))
		}
		leftExpType.Type = *function.ReturnType[0].ReturnType
	}
	if rightExpType.CallExp {
		function := FunctionMap[infix.Right.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return expType{}, errors.New("call expression must return only 1 value for `infix` expression, got: " + strconv.Itoa(len(function.ReturnType)))
		}
		rightExpType.Type = *function.ReturnType[0].ReturnType
	}

	if leftExpType.Type.IsHash || rightExpType.Type.IsHash {
		return expType{}, errors.New("hashmap can't be used with infix operations")
	}

	switch {
	case leftExpType.Type.IsArray && rightExpType.Type.IsArray:
		if leftExpType.Type.Value != rightExpType.Type.Value {
			return expType{}, errors.New("can only add arrays of same type, got: `" + leftExpType.Type.Value + "` and `" + rightExpType.Type.Value + "`")
		}
		if infix.Operator == "+" {
			return expType{Type: ast.Type{Value: leftExpType.Type.Value, IsArray: true, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if infix.Operator == "==" || infix.Operator == "!=" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("can only use `+`, `==`, `!=` infix operators with 2 arrays, got: " + infix.Operator)
		}
	case leftExpType.Type.Value == "int" && rightExpType.Type.Value == "int":
		if infix.Operator == "+" || infix.Operator == "-" || infix.Operator == "*" || infix.Operator == "/" || infix.Operator == "%" || infix.Operator == "|" || infix.Operator == "&" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.INT, Value: "int"}, Value: "int", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if infix.Operator == ">" || infix.Operator == "<" || infix.Operator == "<=" || infix.Operator == ">=" || infix.Operator == "==" || infix.Operator == "!=" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("can only use `+`, `-`, `*`, `/`, `%`, `>`, `<`, `<=`, `>=`, `!=`, `==`, `|`, `&` infix operators with 2 `int`, got: " + infix.Operator)
		}
	case leftExpType.Type.Value == "float" && rightExpType.Type.Value == "float", (leftExpType.Type.Value == "int" && rightExpType.Type.Value == "float") || (leftExpType.Type.Value == "float" && rightExpType.Type.Value == "int"):
		if infix.Operator == "+" || infix.Operator == "-" || infix.Operator == "*" || infix.Operator == "/" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.FLOAT, Value: "float"}, Value: "float", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if infix.Operator == ">" || infix.Operator == "<" || infix.Operator == "<=" || infix.Operator == ">=" || infix.Operator == "==" || infix.Operator == "!=" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("can only use `+`, `-`, `*`, `/`, `>`, `<`, `<=`, `>=`, `!=`, `==` infix operators with 2 `float`, got: " + infix.Operator)
		}
	case leftExpType.Type.Value == "string" && rightExpType.Type.Value == "string":
		if infix.Operator == "+" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.STRING, Value: "string"}, Value: "string", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if infix.Operator == "==" || infix.Operator == "!=" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("can only use `+`, `==`, `!=` infix operators with 2 `string`, got: " + infix.Operator)
		}
	case leftExpType.Type.Value == "char" && rightExpType.Type.Value == "char":
		if infix.Operator == "+" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.STRING, Value: "string"}, Value: "string", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if infix.Operator == "==" || infix.Operator == "!=" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("can only use `+`, `==`, `!=` infix operators with 2 `char`, got: " + infix.Operator)
		}
	case leftExpType.Type.Value == "bool" && rightExpType.Type.Value == "bool":
		if infix.Operator == "==" || infix.Operator == "!=" || infix.Operator == "&&" || infix.Operator == "||" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("can only use `==`, `!=`, `&&`, `||` infix operators with 2 `bool`, got: " + infix.Operator)
		}
	default:
		return expType{}, errors.New("invalid `infix` operation with variable types on left and right, got: `" + leftExpType.Type.Value + "` and `" + rightExpType.Type.Value + "`")
	}
}

func checkAssignExp(assign *ast.AssignmentExpression, env *Environment) (expType, error) {
	leftVar, ok := env.Get(assign.Left.Value)
	if !ok {
		return expType{}, errors.New("variable `" + assign.Left.Value + "` not found, can't assign value to undefined variables")
	}
	if leftVar.VarType == CONST {
		return expType{}, errors.New("variable `" + leftVar.Ident.Value + "` is a constant, can't re-assign value to a constant variable")
	} else if leftVar.VarType == FUNCTION {
		return expType{}, errors.New("variable `" + leftVar.Ident.Value + "` is a function, can't assign value to a function")
	}
	leftType := leftVar.Type

	rightExpType, err := getExpType(assign.Right, env)
	if err != nil {
		return expType{}, err
	}
	if rightExpType.CallExp {
		function := FunctionMap[assign.Right.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return expType{}, errors.New("call expression must return only 1 value for `assignment` operator, got: " + strconv.Itoa(len(function.ReturnType)))
		}
		rightExpType.Type = *function.ReturnType[0].ReturnType
	}

	switch assign.Operator {
	case "=":
		err := varTypeCheckerHelper(leftType, rightExpType.Type)
		if err != nil {
			return expType{}, err
		}
		return expType{Type: leftType, CallExp: false}, nil
	case "+=", "-=", "*=", "/=", "%=":
		infixAst := &ast.InfixExpression{
			Left:     assign.Left,
			Operator: strings.TrimSuffix(assign.Operator, "="),
			Right:    assign.Right,
		}
		infixExpType, err := getExpType(infixAst, env)
		if err != nil {
			return expType{}, err
		}
		err = varTypeCheckerHelper(leftType, infixExpType.Type)
		if err != nil {
			return expType{}, err
		}
		return expType{Type: leftType, CallExp: false}, nil
	default:
		return expType{}, errors.New("only `=`, `+=`, `-=`, `*=`, `/=`, `%=` `assignment` operators supported, got: " + assign.Operator)
	}
}

func checkIndexExp(exp *ast.IndexExpression, env *Environment) (expType, error) {
	leftType, err := getExpType(exp.Left, env)
	if err != nil {
		return expType{}, err
	}
	indexType, err := getExpType(exp.Index, env)
	if err != nil {
		return expType{}, err
	}
	if leftType.CallExp {
		function := FunctionMap[exp.Left.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return expType{}, errors.New("call expression must return only 1 value for `index` operator, got: " + strconv.Itoa(len(function.ReturnType)))
		}
		leftType.Type = *function.ReturnType[0].ReturnType
	}
	if indexType.CallExp {
		function := FunctionMap[exp.Index.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return expType{}, errors.New("call expression must return only 1 value for `index` operator, got: " + strconv.Itoa(len(function.ReturnType)))
		}
		indexType.Type = *function.ReturnType[0].ReturnType
	}

	if leftType.Type.IsArray {
		if leftType.Type.Value == "" {
			return expType{}, errors.New("array is empty, can't index empty array")
		}
		if indexType.Type.Value != "int" {
			return expType{}, errors.New("array index must be an integer, got: " + indexType.Type.Value)
		}
		return expType{Type: ast.Type{Value: leftType.Type.Value, IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
	} else if leftType.Type.IsHash {
		if leftType.Type.SubTypes == nil {
			return expType{}, errors.New("hashmap is empty, can't index empty hashmap")
		}
		if indexType.Type.Value != leftType.Type.SubTypes[0].Value {
			return expType{}, errors.New("hashmap key type must be of datatype `" + leftType.Type.SubTypes[0].Value + "`, got: " + indexType.Type.Value)
		}
		return expType{Type: ast.Type{Value: leftType.Type.SubTypes[1].Value, IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
	} else if leftType.Type.Value == "string" {
		if indexType.Type.Value != "int" {
			return expType{}, errors.New("string index must be an integer, got: " + indexType.Type.Value)
		}
		return expType{Type: ast.Type{Value: "char", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
	} else {
		return expType{}, errors.New("index operation not supported for " + leftType.Type.Value + "[" + indexType.Type.Value + "]")
	}
}

func checkCallExp(callExp *ast.CallExpression, env *Environment) (expType, error) {
	nameVar, ok := env.Get(callExp.Name.(*ast.Identifier).Value)
	if !ok {
		if _, ok := builtinMap[callExp.Name.(*ast.Identifier).Value]; ok {
			resType, err := checkBuiltins(callExp, env)
			if err != nil {
				return expType{}, err
			}
			return resType, nil
		}
		return expType{}, errors.New("function `" + callExp.Name.(*ast.Identifier).Value + "` not found")
	}
	if nameVar.VarType != FUNCTION {
		return expType{}, errors.New("identifier `" + nameVar.Ident.Value + "` is not a function")
	}

	function := FunctionMap[callExp.Name.(*ast.Identifier).Value]
	if len(function.Parameters) != len(callExp.Args) {
		return expType{}, errors.New("number of arguments does not match the number of parameters for function `" + function.Name.Value + "`, got: " + strconv.Itoa(len(callExp.Args)) + ", expected: " + strconv.Itoa(len(function.Parameters)))
	}

	var subTypes []*ast.Type
	for i, arg := range callExp.Args {
		argType, err := getExpType(arg, env)
		if err != nil {
			return expType{}, err
		}
		if argType.CallExp {
			function := FunctionMap[arg.(*ast.CallExpression).Name.(*ast.Identifier).Value]
			if len(function.ReturnType) != 1 {
				return expType{}, errors.New("call expression must return only 1 value to be used as an argument for `call expression`, got: " + strconv.Itoa(len(function.ReturnType)))
			}
			argType.Type = *function.ReturnType[0].ReturnType
		}
		err = varTypeCheckerHelper(*function.Parameters[i].ParameterType, argType.Type)
		if err != nil {
			return expType{}, err
		}
		subTypes = append(subTypes, &argType.Type)
	}
	if function.ReturnType == nil {
		return expType{CallExp: true}, nil
	}
	return expType{Type: ast.Type{IsArray: false, IsHash: false, SubTypes: subTypes}, CallExp: true}, nil
}

func checkBuiltins(callExp *ast.CallExpression, env *Environment) (expType, error) {
	var argTypes []expType
	for _, arg := range callExp.Args {
		argType, err := getExpType(arg, env)
		if err != nil {
			return expType{}, err
		}
		if argType.CallExp {
			function := FunctionMap[arg.(*ast.CallExpression).Name.(*ast.Identifier).Value]
			if len(function.ReturnType) != 1 {
				return expType{}, errors.New("call expression must return only 1 value to be used as an argument for `built-in methods`, got: " + strconv.Itoa(len(function.ReturnType)))
			}
			argType.Type = *function.ReturnType[0].ReturnType
		}
		argTypes = append(argTypes, argType)
	}

	switch callExp.Name.(*ast.Identifier).Value {
	case "len":
		if len(callExp.Args) != 1 {
			return expType{}, errors.New("wrong number of arguments for `len`, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 1")
		}
		if argTypes[0].Type.IsArray || argTypes[0].Type.IsHash || argTypes[0].Type.Value == "string" {
			return expType{Type: ast.Type{Value: "int", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("argument for `len` not supported, got: " + argTypes[0].Type.Value + ", want: array, hashmap or `string`")
		}
	case "toString":
		if len(callExp.Args) != 1 {
			return expType{}, errors.New("wrong number of arguments for `toString`, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 1")
		}
		if argTypes[0].Type.IsArray || argTypes[0].Type.IsHash || argTypes[0].Type.Value == "int" || argTypes[0].Type.Value == "float" || argTypes[0].Type.Value == "bool" || argTypes[0].Type.Value == "char" || argTypes[0].Type.Value == "string" {
			return expType{Type: ast.Type{Value: "string", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("argument for `toString` not supported, got: " + argTypes[0].Type.Value + ", want: array, hashmap, `int`, `float`, `bool`, `char` or `string`")
		}
	case "print", "println":
		if len(callExp.Args) != 1 {
			return expType{}, errors.New("wrong number of arguments for `" + callExp.Name.(*ast.Identifier).Value + "`, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 1")
		}
		if argTypes[0].Type.IsArray || argTypes[0].Type.IsHash || argTypes[0].Type.Value == "int" || argTypes[0].Type.Value == "float" || argTypes[0].Type.Value == "bool" || argTypes[0].Type.Value == "char" || argTypes[0].Type.Value == "string" {
			return expType{Type: ast.Type{Value: callExp.Name.(*ast.Identifier).Value, IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("argument to `" + callExp.Name.(*ast.Identifier).Value + "` not supported, got: " + argTypes[0].Type.Value + ", want: array, hashmap, `int`, `float`, `bool`, `char` or `string`. use `toString` to convert to `string` in case of using other datatypes with `string`")
		}
	case "push":
		if len(callExp.Args) != 2 && len(callExp.Args) != 3 {
			return expType{}, errors.New("wrong number of arguments for `push`, got: " + strconv.Itoa(len(callExp.Args)) + ", want=2 or 3. for array: `push(array, element)` and for hashmap: `push(map, key, value)`")
		}
		if argTypes[0].Type.IsArray {
			if len(callExp.Args) != 2 {
				return expType{}, errors.New("wrong number of arguments for `push` for array, got: " + strconv.Itoa(len(callExp.Args)) + ", want=2. `push(array, element)`")
			}

			arrayType := argTypes[0].Type.Value
			if arrayType == "int" || arrayType == "float" || arrayType == "string" || arrayType == "char" || arrayType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("argument type mismatch for `push`, expected element to be " + arrayType + ", got: array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("argument type mismatch for `push`, expected element to be " + arrayType + ", got: hashmap")
				}
				if argTypes[1].Type.Value != arrayType {
					return expType{}, errors.New("argument type mismatch for `push`, expected element to be " + arrayType + ", got: " + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("array type not supported for `push`, got: " + arrayType + ", want: `int`, `float`, `string`, `char` or `bool`")
			}

			return argTypes[0], nil
		} else if argTypes[0].Type.IsHash {
			if len(callExp.Args) != 3 {
				return expType{}, errors.New("wrong number of arguments for `push` for hashmap, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 3. `push(map, key, value)`")
			}

			keyType := argTypes[0].Type.SubTypes[0].Value
			if keyType == "int" || keyType == "float" || keyType == "string" || keyType == "char" || keyType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("key type mismatch for `push`, expected key to be " + keyType + ", got: array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("key type mismatch for `push`, expected key to be " + keyType + ", got: hashmap")
				}
				if argTypes[1].Type.Value != keyType {
					return expType{}, errors.New("key type mismatch for `push`, expected key to be " + keyType + ", got: " + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("key type not supported, got: " + keyType + ", want: `int`, `float`, `string`, `char` or `bool`")
			}

			valueType := argTypes[0].Type.SubTypes[1].Value
			if valueType == "int" || valueType == "float" || valueType == "string" || valueType == "char" || valueType == "bool" {
				if argTypes[2].Type.IsArray {
					return expType{}, errors.New("value type mismatch for `push`, expected value to be " + valueType + ", got: array")
				} else if argTypes[2].Type.IsHash {
					return expType{}, errors.New("value type mismatch for `push`, expected value to be " + valueType + ", got: hashmap")
				}
				if argTypes[2].Type.Value != valueType {
					return expType{}, errors.New("value type mismatch for `push`, expected value to be " + valueType + ", got: " + argTypes[2].Type.Value)
				}
			} else {
				return expType{}, errors.New("value type not supported, got: " + valueType + ", want: `int`, `float`, `string`, `char` or `bool`")
			}

			return argTypes[0], nil
		} else {
			return expType{}, errors.New("data structure not supported by `push`, got: " + argTypes[0].Type.Value + ", want: array or hashmap")
		}
	case "pop":
		if len(callExp.Args) != 1 && len(callExp.Args) != 2 {
			return expType{}, errors.New("wrong number of arguments for `pop` for array, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 1 or 2. `pop(array)` or `pop(array, index)`")
		}
		if argTypes[0].Type.IsArray {
			if len(callExp.Args) == 1 {
				return expType{Type: ast.Type{Value: argTypes[0].Type.Value, IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
			} else {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("index must be an integer for `pop`, got: array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("index must be an integer for `pop`, got: hashmap")
				} else if argTypes[1].Type.Value != "int" {
					return expType{}, errors.New("index must be an integer for `pop`, got: " + argTypes[1].Type.Value)
				}
				return expType{Type: ast.Type{Value: argTypes[0].Type.Value, IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
			}
		} else if argTypes[0].Type.IsHash {
			return expType{}, errors.New("data structure not supported by `pop`, got: hashmap, want: array")
		} else {
			return expType{}, errors.New("data structure not supported by `pop`, got: " + argTypes[0].Type.Value + ", want: array")
		}
	case "insert":
		if len(callExp.Args) != 3 {
			return expType{}, errors.New("wrong number of arguments for `insert` for array, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 3. `insert(array, index, element)`")
		}
		if argTypes[0].Type.IsArray {
			if argTypes[1].Type.IsArray {
				return expType{}, errors.New("index must be an integer for `insert`, got: array")
			} else if argTypes[1].Type.IsHash {
				return expType{}, errors.New("index must be an integer for `insert`, got: hashmap")
			} else if argTypes[1].Type.Value != "int" {
				return expType{}, errors.New("index must be an integer for `insert`, got: " + argTypes[1].Type.Value)
			}

			arrayType := argTypes[0].Type.Value
			if arrayType == "int" || arrayType == "float" || arrayType == "string" || arrayType == "char" || arrayType == "bool" {
				if argTypes[2].Type.IsArray {
					return expType{}, errors.New("argument type mismatch for `insert`, expected element to be " + arrayType + ", got: array")
				} else if argTypes[2].Type.IsHash {
					return expType{}, errors.New("argument type mismatch for `insert`, expected element to be " + arrayType + ", got: hashmap")
				}
				if argTypes[2].Type.Value != arrayType {
					return expType{}, errors.New("argument type mismatch for `insert`, expected element to be " + arrayType + ", got: " + argTypes[2].Type.Value)
				}
			} else {
				return expType{}, errors.New("array type not supported for `insert`, got: " + arrayType + ", want: `int`, `float`, `string`, `char` or `bool`")
			}
			return argTypes[0], nil
		} else if argTypes[0].Type.IsHash {
			return expType{}, errors.New("data structure not supported by `insert`, got: hashmap, want: array")
		} else {
			return expType{}, errors.New("data structure not supported by `insert`, got: " + argTypes[0].Type.Value + ", want: array")
		}
	case "remove":
		if len(callExp.Args) != 2 {
			return expType{}, errors.New("wrong number of arguments for `remove` for array, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 2. for array: `remove(array, element)`, for hashmap: `remove(map, key)`")
		}
		if argTypes[0].Type.IsArray {
			if len(callExp.Args) != 2 {
				return expType{}, errors.New("wrong number of arguments for `remove` for array, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 2. `remove(array, element)`")
			}

			arrayType := argTypes[0].Type.Value
			if arrayType == "int" || arrayType == "float" || arrayType == "string" || arrayType == "char" || arrayType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("argument type mismatch for `remove`, expected element to be " + arrayType + ", got: array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("argument type mismatch for `remove`, expected element to be " + arrayType + ", got: hashmap")
				}
				if argTypes[1].Type.Value != arrayType {
					return expType{}, errors.New("argument type mismatch for `remove`, expected element to be " + arrayType + ", got: " + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("array type not supported for `remove`, got: " + arrayType + ", want: `int`, `float`, `string`, `char` or `bool`")
			}
			return argTypes[0], nil
		} else if argTypes[0].Type.IsHash {
			if len(callExp.Args) != 2 {
				return expType{}, errors.New("wrong number of arguments for `remove` for hashmap, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 2. `remove(map, key)`")
			}

			keyType := argTypes[0].Type.SubTypes[0].Value
			if keyType == "int" || keyType == "float" || keyType == "string" || keyType == "char" || keyType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("key type mismatch for `remove`, expected key to be " + keyType + ", got: array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("key type mismatch for `remove`, expected key to be " + keyType + ", got: hashmap")
				}
				if argTypes[1].Type.Value != keyType {
					return expType{}, errors.New("key type mismatch for `remove`, expected key to be " + keyType + ", got: " + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("key type not supported, got: " + keyType + ", want: `int`, `float`, `string`, `char` or `bool`")
			}
			return expType{Type: *argTypes[0].Type.SubTypes[1], CallExp: false}, nil
		} else {
			return expType{}, errors.New("data structure not supported by `remove`, got: " + argTypes[0].Type.Value + ", want: array or hashmap")
		}
	case "getIndex":
		if len(callExp.Args) != 2 {
			return expType{}, errors.New("wrong number of arguments for `getIndex` for array, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 2. `getIndex(array, element)`")
		}
		if argTypes[0].Type.IsArray {
			arrayType := argTypes[0].Type.Value
			if arrayType == "int" || arrayType == "float" || arrayType == "string" || arrayType == "char" || arrayType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("argument type mismatch for `getIndex`, expected element to be " + arrayType + ", got: array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("argument type mismatch for `getIndex`, expected element to be " + arrayType + ", got: hashmap")
				}
				if argTypes[1].Type.Value != arrayType {
					return expType{}, errors.New("argument type mismatch for `getIndex`, expected element to be " + arrayType + ", got: " + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("array type not supported for `getIndex`, got: " + arrayType + ", want: `int`, `float`, `string`, `char` or `bool`")
			}

			return expType{Type: ast.Type{Value: "int", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if argTypes[0].Type.IsHash {
			return expType{}, errors.New("data structure not supported by `getIndex`, got: hashmap, want: array")
		} else {
			return expType{}, errors.New("data structure not supported by `getIndex`, got: " + argTypes[0].Type.Value + ", want: array")
		}
	case "keys":
		if len(callExp.Args) != 1 {
			return expType{}, errors.New("wrong number of arguments for `keys` for hashmap, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 1. `keys(map)`")
		}
		if argTypes[0].Type.IsHash {
			return expType{Type: ast.Type{Value: argTypes[0].Type.SubTypes[0].Value, IsArray: true, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if argTypes[0].Type.IsArray {
			return expType{}, errors.New("data structure not supported by `keys`, got: array, want: hashmap")
		} else {
			return expType{}, errors.New("data structure not supported by `keys`, got: " + argTypes[0].Type.Value + ", want: hashmap")
		}
	case "values":
		if len(callExp.Args) != 1 {
			return expType{}, errors.New("wrong number of arguments for `values` for hashmap, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 1. `values(map)`")
		}
		if argTypes[0].Type.IsHash {
			return expType{Type: ast.Type{Value: argTypes[0].Type.SubTypes[1].Value, IsArray: true, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if argTypes[0].Type.IsArray {
			return expType{}, errors.New("data structure not supported by `keys`, got: array, want: hashmap")
		} else {
			return expType{}, errors.New("data structure not supported by `values`, got: " + argTypes[0].Type.Value + ", want: hashmap")
		}
	case "containsKey":
		if len(callExp.Args) != 2 {
			return expType{}, errors.New("wrong number of arguments for `containsKey` for hashmap, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 2. `containsKey(map, key)`")
		}
		if argTypes[0].Type.IsHash {
			keyType := argTypes[0].Type.SubTypes[0].Value
			if keyType == "int" || keyType == "float" || keyType == "string" || keyType == "char" || keyType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("key type mismatch for `containsKey`, expected key to be " + keyType + ", got: array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("key type mismatch for `containsKey`, expected key to be " + keyType + ", got: hashmap")
				}
				if argTypes[1].Type.Value != keyType {
					return expType{}, errors.New("key type mismatch for `containsKey`, expected key to be " + keyType + ", got: " + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("key type not supported, got: " + keyType + ", want: `int`, `float`, `string`, `char` or `bool`")
			}
			return expType{Type: ast.Type{Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if argTypes[0].Type.IsArray {
			return expType{}, errors.New("data structure not supported by `containsKey`, got: array, want: hashmap")
		} else {
			return expType{}, errors.New("data structure not supported by `containsKey`, got: " + argTypes[0].Type.Value + ", want: hashmap")
		}
	case "typeOf":
		if len(callExp.Args) != 1 {
			return expType{}, errors.New("wrong number of arguments for `typeOf`, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 1")
		}
		if argTypes[0].Type.IsArray || argTypes[0].Type.IsHash || argTypes[0].Type.Value == "int" || argTypes[0].Type.Value == "float" || argTypes[0].Type.Value == "char" || argTypes[0].Type.Value == "string" || argTypes[0].Type.Value == "bool" {
			return expType{Type: ast.Type{Value: "string", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("data structure not supported by `typeOf`, got: " + argTypes[0].Type.Value)
		}
	case "slice":
		if len(callExp.Args) != 3 && len(callExp.Args) != 4 {
			return expType{}, errors.New("wrong number of arguments for `slice`, got: " + strconv.Itoa(len(callExp.Args)) + ", want: 3 or 4. `slice(array/string, start, end)` or `slice(array/string, start, end, step)`")
		}

		if argTypes[1].Type.IsArray {
			return expType{}, errors.New("start index must be an `int` for `slice`, got: array")
		} else if argTypes[1].Type.IsHash {
			return expType{}, errors.New("start index must be an `int` for `slice`, got: hashmap")
		} else {
			if argTypes[1].Type.Value != "int" {
				return expType{}, errors.New("start index must be an `int` for `slice`, got: " + argTypes[1].Type.Value)
			}
		}

		if argTypes[2].Type.IsArray {
			return expType{}, errors.New("end index must be an `int` for `slice`, got: array")
		} else if argTypes[2].Type.IsHash {
			return expType{}, errors.New("end index must be an `int` for `slice`, got: hashmap")
		} else {
			if argTypes[2].Type.Value != "int" {
				return expType{}, errors.New("end index must be an `int` for `slice`, got: " + argTypes[2].Type.Value)
			}
		}

		if len(callExp.Args) == 4 {
			if argTypes[3].Type.IsArray {
				return expType{}, errors.New("step must be an `int` for `slice`, got: array")
			} else if argTypes[3].Type.IsHash {
				return expType{}, errors.New("step must be an `int` for `slice`, got: hashmap")
			} else {
				if argTypes[3].Type.Value != "int" {
					return expType{}, errors.New("step must be an `int` for `slice`, got: " + argTypes[3].Type.Value)
				}
			}
		}

		if argTypes[0].Type.IsArray {
			return argTypes[0], nil
		} else if argTypes[0].Type.IsHash {
			return expType{}, errors.New("data structure not supported by `slice`, got: hashmap, want: array or string")
		} else {
			if argTypes[0].Type.Value != "string" {
				return expType{}, errors.New("data structure not supported by `slice`, got: " + argTypes[0].Type.Value + ", want: array or string")
			}
			return argTypes[0], nil
		}
	default:
		return expType{}, errors.New("unknown builtin function `" + callExp.Name.(*ast.Identifier).Value + "`")
	}
}

// -----------------------------------------------------------------------------
// Helper Methods
// -----------------------------------------------------------------------------
func varTypeCheckerHelper(definedType ast.Type, expType ast.Type) error {
	if definedType.IsArray {
		if !expType.IsArray {
			if expType.IsHash {
				return errors.New("defined type is of array, but got: hashmap")
			}
			return errors.New("defined type is of array, but got: " + expType.Value)
		}
		if definedType.Value != expType.Value {
			if expType.Value == "" {
				return nil // means it is defined empty and can hold anytype
			}
			return errors.New("array declared as `" + definedType.Value + "`, got: " + expType.Value)
		}
	} else if definedType.IsHash {
		if !expType.IsHash {
			if expType.IsArray {
				return errors.New("defined type is of hashmap, but got: array")
			}
			return errors.New("defined type is of hashmap, got: " + expType.Value)
		}
		if expType.SubTypes == nil {
			return nil // means it is defined empty and can hold anytype
		}
		if definedType.SubTypes[0].Value != expType.SubTypes[0].Value {
			return errors.New("hashmap key type declared as `" + definedType.SubTypes[0].Value + "`, got: " + expType.SubTypes[0].Value)
		}
		if definedType.SubTypes[1].Value != expType.SubTypes[1].Value {
			return errors.New("hashmap value type declared as `" + definedType.SubTypes[1].Value + "`, got: " + expType.SubTypes[1].Value)
		}
	} else {
		if expType.IsArray {
			return errors.New("defined type is `" + definedType.Value + "`, got: array")
		} else if expType.IsHash {
			return errors.New("defined type is `" + definedType.Value + "`, got: hashmap")
		}
		if definedType.Value != expType.Value {
			return errors.New("defined type is `" + definedType.Value + "`, got: " + expType.Value)
		}
	}

	return nil
}

func returnTypeCheckerHelper(node *ast.ReturnStatement, env *Environment) error {
	for i, entry := range node.Value {
		switch entry.(type) {
		case *ast.Identifier, *ast.IntegerValue, *ast.FloatValue, *ast.StringValue, *ast.CharValue, *ast.BooleanValue, *ast.PrefixExpression, *ast.PostfixExpression, *ast.InfixExpression, *ast.ArrayValue, *ast.HashMap, *ast.CallExpression, *ast.IndexExpression:
			expType, err := getExpType(entry, env)
			if err != nil {
				return err
			}

			if !inTesting {
				if expType.CallExp {
					callFun := FunctionMap[entry.(*ast.CallExpression).Name.(*ast.Identifier).Value]
					if len(callFun.ReturnType) != 1 {
						return errors.New("call expression must return only 1 value for `return` statement, got: " + strconv.Itoa(len(callFun.ReturnType)))
					}
					expType.Type = *callFun.ReturnType[0].ReturnType
				}
				err = varTypeCheckerHelper(*currFunction.ReturnType[i].ReturnType, expType.Type)
				if err != nil {
					return err
				}
			}
		default:
			return errors.New("can only return expressions and datatypes, got: " + fmt.Sprintf("%T", entry))
		}
	}
	return nil
}
