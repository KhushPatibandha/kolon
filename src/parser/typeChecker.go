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
	}
	inTesting bool
)

type expType struct {
	Type    ast.Type
	CallExp bool
}

func TypeCheckProgram(program *ast.Program, env *Environment, inTest bool) error {
	inTesting = inTest

	// before starting add all the functions from FunctionMap to the environment.
	// because we might be calling functions that are written afterwards and because
	// of that there will not be in the global environment
	for key, value := range FunctionMap {
		newLocalEnv := NewEnclosedEnvironment(env)

		// loop throught all the function parameters and add them to the local env
		for _, param := range value.Parameters {
			newLocalEnv.Set(param.ParameterName.Value, *param.ParameterName, *param.ParameterType, VAR, nil)
		}

		env.Set(key, *value.Name, ast.Type{}, FUNCTION, newLocalEnv)
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
		// TODO: FIX the return statement bug
		return checkReturnStmt(node, env)
	case *ast.ContinueStatement:
		if !inForLoop {
			return errors.New("Continue statement can only be used inside a for loop")
		}
		return nil
	case *ast.BreakStatement:
		if !inForLoop {
			return errors.New("Break statement can only be used inside a for loop")
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
	case *ast.ExpressionStatement:
		return checkExpStmt(node, env)
	default:
		msg := fmt.Errorf("Unknown statement type. got: %T", node)
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
		msg := fmt.Errorf("Unknown expression type. got: %T", exp)
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
// Stmt checking
// -----------------------------------------------------------------------------
func checkVarStmt(node *ast.VarStatement, env *Environment) error {
	definedType := node.Type

	// get the type for the expression on the right
	expType, err := getExpType(node.Value, env)
	if err != nil {
		return err
	}

	// if the expType is a call expression, then the function must return only 1 value
	// we have to get the types from FunctionMap in parser
	if expType.CallExp {
		function := FunctionMap[node.Value.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return errors.New("Call expression must return only 1 value for var statement")
		}
		expType.Type = *function.ReturnType[0].ReturnType
	}

	// now we need to check the type of the variable with the expType.Type
	err = varTypeCheckerHelper(*definedType, expType.Type)
	if err != nil {
		return err
	}

	// If everything is correct, then add the variable to the environment
	if node.Token.Kind == lexer.VAR {
		env.Set(node.Name.Value, *node.Name, *node.Type, VAR, nil)
	} else if node.Token.Kind == lexer.CONST {
		env.Set(node.Name.Value, *node.Name, *node.Type, CONST, nil)
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
			return errors.New("Number of return values does not match the number of variables")
		}

		for i, obj := range node.Objects {
			switch obj := obj.(type) {
			case *ast.VarStatement:
				err := varTypeCheckerHelper(*obj.Type, *returnTypes[i].ReturnType)
				if err != nil {
					return nil
				}
				if obj.Token.Kind == lexer.VAR {
					env.Set(obj.Name.Value, *obj.Name, *obj.Type, VAR, nil)
				} else if obj.Token.Kind == lexer.CONST {
					env.Set(obj.Name.Value, *obj.Name, *obj.Type, CONST, nil)
				}
			case *ast.ExpressionStatement:
				variable, ok := env.Get(obj.Expression.(*ast.AssignmentExpression).Left.Value)
				if !ok {
					return errors.New("Variable not found, Can't assign value to undefined variables")
				}
				if variable.VarType == CONST || variable.VarType == FUNCTION {
					return errors.New("Can't assign value to a variable that is constant or a function")
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

	// Check the condition first. The condition must always result in a boolean value
	expType, err := getExpType(node.Value, localEnvForIf)
	if err != nil {
		return err
	}

	if expType.CallExp {
		function := FunctionMap[node.Value.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return errors.New("Call expression must return only 1 value for if statement")
		}
		expType.Type = *function.ReturnType[0].ReturnType
	}
	if expType.Type.Value != "bool" {
		return errors.New("Condition for IF statements must result in a boolean value")
	}

	// Check the statements inside the IF block
	err = checkStmts(node.Body.Statements, localEnvForIf)
	if err != nil {
		return err
	}

	// Check the else if blocks
	if node.MultiConseq != nil {
		for _, elseIfStmt := range node.MultiConseq {
			// Do this for every else if block because we are going through all if, else if and else blocks. hence we don't want variables in "if" to be accessible in "else if"
			localEnfForElseIf := NewEnclosedEnvironment(env)
			expType, err := getExpType(elseIfStmt.Value, localEnfForElseIf)
			if err != nil {
				return err
			}
			if expType.CallExp {
				function := FunctionMap[elseIfStmt.Value.(*ast.CallExpression).Name.(*ast.Identifier).Value]
				if len(function.ReturnType) != 1 {
					return errors.New("Call expression must return only 1 value for if statement")
				}
				expType.Type = *function.ReturnType[0].ReturnType
			}
			if expType.Type.Value != "bool" {
				return errors.New("Condition for ELSE IF statements must result in a boolean value")
			}

			// Check the statements inside the ELSE IF block
			err = checkStmts(elseIfStmt.Body.Statements, localEnfForElseIf)
			if err != nil {
				return err
			}
		}
	}

	// check the else block
	if node.Consequence != nil {
		localEnvForElse := NewEnclosedEnvironment(env)
		err = checkStmts(node.Consequence.Body.Statements, localEnvForElse)
		if err != nil {
			return nil
		}
	}

	return nil
}

func checkForLoopStmt(node *ast.ForLoopStatement, env *Environment) error {
	inForLoop = true

	// first check the var statement
	err := checkVarStmt(node.Left, env)
	if err != nil {
		return nil
	}
	// Get the variable type to check for VAR or CONST
	varVariable, ok := env.Get(node.Left.Name.Value)
	if !ok {
		return errors.New("Variable in FOR loop condition not found")
	}
	if varVariable.VarType != VAR {
		return errors.New("Can't use CONST to define variable in FOR loop condition")
	}
	// Check if the variable is INT or not
	if varVariable.Type.Value != "int" {
		return errors.New("Can only define variable in FOR loop condition as INT.")
	}

	// Check infix expression. The infix expression must result in a boolean value
	expType, err := getExpType(node.Middle, env)
	if err != nil {
		return nil
	}
	if expType.Type.Value != "bool" {
		return errors.New("Infix operation of FOR loop condition should always result in a BOOLEAN.")
	}

	// Check the postfix expression.
	expType, err = getExpType(node.Right, env)
	if err != nil {
		return err
	}
	if expType.Type.Value != "int" {
		return errors.New("Postfix operation of FOR loop condition should always result in an INT.")
	}

	// Check the statements inside the FOR block
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
		return errors.New("Function not found in the environment")
	}
	if funVariable.VarType != FUNCTION {
		return errors.New("Function not found in the environment")
	}
	localFunEnv := funVariable.Env

	// no need to check the parameters and return types. Because they are just bunch of types and also they are already checked in the parser

	// check the body
	err := checkStmts(node.Body.Statements, localFunEnv)
	if err != nil {
		return err
	}
	return nil
}

func checkMainFunStmt(node *ast.Function, env *Environment) error {
	funVariable, ok := env.Get(node.Name.Value)
	if !ok {
		return errors.New("Main function not found in the environment")
	}
	if funVariable.VarType != FUNCTION {
		return errors.New("Main function not found in the environment")
	}
	localMainEnv := funVariable.Env

	// First check the parameters. the list of parameters must be empty
	if len(node.Parameters) != 0 {
		return errors.New("Main function must not have any parameters")
	}

	// return type must be nil
	if node.ReturnType != nil {
		return errors.New("Can't return anything in main function")
	}

	// check the body
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
			return errors.New("Can't return anything in main function")
		}
		return nil
	}

	if node.Value == nil && currFunction.ReturnType == nil {
		return nil
	}
	if node.Value != nil && currFunction.ReturnType != nil {
		if len(node.Value) != len(currFunction.ReturnType) {
			return errors.New("Number of return values does not match the number of return types")
		}
	} else {
		return errors.New("Number of return values does not match the number of return types")
	}

	err := returnTypeCheckerHelper(node, env)
	if err != nil {
		return err
	}
	return nil
}

// -----------------------------------------------------------------------------
// Expression checking
// -----------------------------------------------------------------------------
func checkArrayExp(node *ast.ArrayValue, env *Environment) (expType, error) {
	if len(node.Values) == 0 {
		return expType{Type: ast.Type{IsArray: true, IsHash: false, SubTypes: nil}, CallExp: false}, nil
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
				return expType{}, errors.New("Array can only have one type of elements")
			}
		}
	}

	return expType{Type: ast.Type{Value: typeStr, IsArray: true, IsHash: false, SubTypes: nil}, CallExp: false}, nil
}

func checkHashmapExp(hashmap *ast.HashMap, env *Environment) (expType, error) {
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
				return expType{}, errors.New("Hashmap can only have one type of keys")
			}
			if valueTypeStr != valueExp.Type.Value {
				return expType{}, errors.New("Hashmap can only have one type of values")
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
		return expType{}, errors.New("Variable " + ident.Value + " is undefined/not found")
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
			return expType{}, errors.New("Call expression must return only 1 value for prefix expression")
		}
		rightExpType.Type = *function.ReturnType[0].ReturnType
	}
	if rightExpType.Type.IsArray || rightExpType.Type.IsHash {
		return expType{}, errors.New("Prefix operator can't be used with array or hashmap")
	}
	switch prefix.Operator {
	case "!":
		if rightExpType.Type.Value != "bool" {
			return expType{}, errors.New("Only Integer and Float datatypes supported with Postfix operation. got: " + rightExpType.Type.Value)
		}
	case "-":
		if rightExpType.Type.Value != "int" && rightExpType.Type.Value != "float" {
			return expType{}, errors.New("Dash/Minus(-) operator can be only used with Integers(int) and Floats(float) entities. got: " + rightExpType.Type.Value)
		}
	default:
		return expType{}, errors.New("Only 2 Prefix Operator's supported. !(Bang) and -(Dash/Minus). got: " + prefix.Operator)
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
			return expType{}, errors.New("Call expression must return only 1 value for postfix expression")
		}
		leftExpType.Type = *function.ReturnType[0].ReturnType
	}
	if leftExpType.Type.IsArray || leftExpType.Type.IsHash {
		return expType{}, errors.New("Postfix operator can't be used with array or hashmap")
	}
	if leftExpType.Type.Value != "int" && leftExpType.Type.Value != "float" {
		return expType{}, errors.New("Only Integer and Float datatypes supported with Postfix operation. got: " + leftExpType.Type.Value)
	}
	if postfix.Operator != "++" && postfix.Operator != "--" {
		return expType{}, errors.New("Only 2 Postfix Operator's supported. ++(Increment) and --(Decrement). got: " + postfix.Operator)
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
			return expType{}, errors.New("Call expression must return only 1 value for infix expression")
		}
		leftExpType.Type = *function.ReturnType[0].ReturnType
	}
	if rightExpType.CallExp {
		function := FunctionMap[infix.Right.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return expType{}, errors.New("Call expression must return only 1 value for infix expression")
		}
		rightExpType.Type = *function.ReturnType[0].ReturnType
	}

	if leftExpType.Type.IsHash || rightExpType.Type.IsHash {
		return expType{}, errors.New("Hashmap can't be used with infix operations")
	}

	switch {
	case leftExpType.Type.IsArray && rightExpType.Type.IsArray:
		if leftExpType.Type.Value != rightExpType.Type.Value {
			return expType{}, errors.New("Can only add arrays of same type. got: " + leftExpType.Type.Value + " and " + rightExpType.Type.Value)
		}
		if infix.Operator == "+" {
			return expType{Type: ast.Type{Value: leftExpType.Type.Value, IsArray: true, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if infix.Operator == "==" || infix.Operator == "!=" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("Can only perform \"+\", \"==\", \"!=\" with arrays. got: " + infix.Operator)
		}
	case leftExpType.Type.Value == "int" && rightExpType.Type.Value == "int":
		if infix.Operator == "+" || infix.Operator == "-" || infix.Operator == "*" || infix.Operator == "/" || infix.Operator == "%" || infix.Operator == "|" || infix.Operator == "&" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.INT, Value: "int"}, Value: "int", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if infix.Operator == ">" || infix.Operator == "<" || infix.Operator == "<=" || infix.Operator == ">=" || infix.Operator == "==" || infix.Operator == "!=" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("Can only perform \"+\", \"-\", \"/\", \"*\", \"%\", \"&\", \"|\", \">\", \"<\", \"<=\", \">=\", \"!=\", \"==\" with 2 Integers. got: " + infix.Operator)
		}
	case leftExpType.Type.Value == "float" && rightExpType.Type.Value == "float", (leftExpType.Type.Value == "int" && rightExpType.Type.Value == "float") || (leftExpType.Type.Value == "float" && rightExpType.Type.Value == "int"):
		if infix.Operator == "+" || infix.Operator == "-" || infix.Operator == "*" || infix.Operator == "/" || infix.Operator == "%" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.FLOAT, Value: "float"}, Value: "float", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if infix.Operator == ">" || infix.Operator == "<" || infix.Operator == "<=" || infix.Operator == ">=" || infix.Operator == "==" || infix.Operator == "!=" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("Can only perform \"+\", \"-\", \"/\", \"*\", \"%\", \">\", \"<\", \"<=\", \">=\", \"!=\", \"==\" with 2 Float var. got: " + infix.Operator)
		}
	case leftExpType.Type.Value == "string" && rightExpType.Type.Value == "string":
		if infix.Operator == "+" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.STRING, Value: "string"}, Value: "string", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if infix.Operator == "==" || infix.Operator == "!=" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("Can only perform \"+\", \"==\", \"!=\" with strings. got: " + infix.Operator)
		}
	case leftExpType.Type.Value == "char" && rightExpType.Type.Value == "char":
		if infix.Operator == "+" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.STRING, Value: "string"}, Value: "string", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if infix.Operator == "==" || infix.Operator == "!=" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("Can only perform \"+\", \"==\", \"!=\" with chars. got: " + infix.Operator)
		}
	case leftExpType.Type.Value == "bool" && rightExpType.Type.Value == "bool":
		if infix.Operator == "==" || infix.Operator == "!=" || infix.Operator == "&&" || infix.Operator == "||" {
			return expType{Type: ast.Type{Token: lexer.Token{Kind: lexer.BOOL, Value: "bool"}, Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("Can only perform \"&&\", \"||\", \"!=\", \"==\" with 2 Boolean values. got: " + infix.Operator)
		}
	default:
		return expType{}, errors.New("Invalid operation with variable types on left and right.")
	}
}

func checkAssignExp(assign *ast.AssignmentExpression, env *Environment) (expType, error) {
	leftVar, ok := env.Get(assign.Left.Value)
	if !ok {
		return expType{}, errors.New("Variable not found, Can't assign value to undefined variables")
	}
	if leftVar.VarType == CONST || leftVar.VarType == FUNCTION {
		return expType{}, errors.New("Can't assign value to a variable that is constant or a function")
	}
	leftType := leftVar.Type

	rightExpType, err := getExpType(assign.Right, env)
	if err != nil {
		return expType{}, err
	}
	if rightExpType.CallExp {
		function := FunctionMap[assign.Right.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return expType{}, errors.New("Call expression must return only 1 value for assignment expression")
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
		return expType{}, errors.New("Only (=, +=, -=, *=, /=, %=) assignment operators are supported. got: " + assign.Operator)
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
			return expType{}, errors.New("Call expression must return only 1 value for Index expression")
		}
		leftType.Type = *function.ReturnType[0].ReturnType
	}
	if indexType.CallExp {
		function := FunctionMap[exp.Index.(*ast.CallExpression).Name.(*ast.Identifier).Value]
		if len(function.ReturnType) != 1 {
			return expType{}, errors.New("Call expression must return only 1 value for Index expression")
		}
		indexType.Type = *function.ReturnType[0].ReturnType
	}

	if leftType.Type.IsArray {
		if indexType.Type.Value != "int" {
			return expType{}, errors.New("Array index must be an integer. got " + indexType.Type.Value)
		}
		return expType{Type: ast.Type{Value: leftType.Type.Value, IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
	} else if leftType.Type.IsHash {
		if indexType.Type.Value != leftType.Type.SubTypes[0].Value {
			return expType{}, errors.New("Hashmap key type must be " + leftType.Type.SubTypes[0].Value + ". got: " + indexType.Type.Value)
		}
		return expType{Type: ast.Type{Value: leftType.Type.SubTypes[1].Value, IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
	} else {
		return expType{}, errors.New("Index operator not supported: " + leftType.Type.Value + "[" + indexType.Type.Value + "]")
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
		return expType{}, errors.New("Function not found")
	}
	if nameVar.VarType != FUNCTION {
		return expType{}, errors.New("Identifier " + nameVar.Ident.Value + " is not a function")
	}

	function := FunctionMap[callExp.Name.(*ast.Identifier).Value]
	if len(function.Parameters) != len(callExp.Args) {
		return expType{}, errors.New("Number of arguments does not match the number of parameters")
	}
	for i, arg := range callExp.Args {
		argType, err := getExpType(arg, env)
		if err != nil {
			return expType{}, err
		}
		if argType.CallExp {
			function := FunctionMap[arg.(*ast.CallExpression).Name.(*ast.Identifier).Value]
			if len(function.ReturnType) != 1 {
				return expType{}, errors.New("Call expression must return only 1 value for call expression")
			}
			argType.Type = *function.ReturnType[0].ReturnType
		}
		err = varTypeCheckerHelper(*function.Parameters[i].ParameterType, argType.Type)
		if err != nil {
			return expType{}, err
		}
	}
	if len(function.ReturnType) != 1 {
		return expType{}, errors.New("Function must return only 1 value")
	}
	return expType{Type: *function.ReturnType[0].ReturnType, CallExp: true}, nil
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
				return expType{}, errors.New("Call expression must return only 1 value for call expression")
			}
			argType.Type = *function.ReturnType[0].ReturnType
		}
		argTypes = append(argTypes, argType)
	}

	switch callExp.Name.(*ast.Identifier).Value {
	case "len":
		if len(callExp.Args) != 1 {
			return expType{}, errors.New("Wrong number of arguments for `len`. got=" + strconv.Itoa(len(callExp.Args)) + ", want=1")
		}
		if argTypes[0].Type.IsArray || argTypes[0].Type.IsHash || argTypes[0].Type.Value == "string" {
			return expType{Type: ast.Type{Value: "int", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("Argument to `len` not supported, got=" + argTypes[0].Type.Value + ", want= array, hashmap or string")
		}
	case "toString":
		if len(callExp.Args) != 1 {
			return expType{}, errors.New("Wrong number of arguments for `toString`. got=" + strconv.Itoa(len(callExp.Args)) + ", want=1")
		}
		if argTypes[0].Type.IsArray || argTypes[0].Type.IsHash || argTypes[0].Type.Value == "int" || argTypes[0].Type.Value == "float" || argTypes[0].Type.Value == "bool" || argTypes[0].Type.Value == "char" || argTypes[0].Type.Value == "string" {
			return expType{Type: ast.Type{Value: "string", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else {
			return expType{}, errors.New("Argument to `toString` not supported, got=" + argTypes[0].Type.Value + ", want= array, hashmap, int, float, bool, char or string")
		}
	case "print", "println":
		if len(callExp.Args) != 1 {
			return expType{}, errors.New("Wrong number of arguments for `print`. got=" + strconv.Itoa(len(callExp.Args)) + ", want=1")
		}
		if argTypes[0].Type.IsArray || argTypes[0].Type.IsHash || argTypes[0].Type.Value == "int" || argTypes[0].Type.Value == "float" || argTypes[0].Type.Value == "bool" || argTypes[0].Type.Value == "char" || argTypes[0].Type.Value == "string" {
			if callExp.Name.(*ast.Identifier).Value == "print" {
				return expType{Type: ast.Type{Value: "print", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
			} else {
				return expType{Type: ast.Type{Value: "println", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
			}
		} else {
			if callExp.Name.(*ast.Identifier).Value == "print" {
				return expType{}, errors.New("Argument to `print` not supported, got=" + argTypes[0].Type.Value + ", want= array, hashmap, int, float, bool, char or string. Use `toString` to convert to string in case of using other types with string")
			} else {
				return expType{}, errors.New("Argument to `println` not supported, got=" + argTypes[0].Type.Value + ", want= array, hashmap, int, float, bool, char or string. Use `toString` to convert to string in case of using other types with string")
			}
		}
	case "push":
		if len(callExp.Args) != 2 && len(callExp.Args) != 3 {
			return expType{}, errors.New("Wrong number of arguments for `push`. got=" + strconv.Itoa(len(callExp.Args)) + ", want=2 or 3. For Array: `push(array, element)` and for Map: `push(map, key, value)`")
		}
		if argTypes[0].Type.IsArray {
			if len(callExp.Args) != 2 {
				return expType{}, errors.New("Wrong number of arguments for `push` on array. got=" + strconv.Itoa(len(callExp.Args)) + ", want=2. `push(array, element)`")
			}

			arrayType := argTypes[0].Type.Value
			if arrayType == "int" || arrayType == "float" || arrayType == "string" || arrayType == "char" || arrayType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=hashmap")
				}
				if argTypes[1].Type.Value != arrayType {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=" + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("Array type not supported. got=" + arrayType + ", want= int, float, string, char or bool")
			}

			return argTypes[0], nil
		} else if argTypes[0].Type.IsHash {
			if len(callExp.Args) != 3 {
				return expType{}, errors.New("Wrong number of arguments for `push` on hash. got=" + strconv.Itoa(len(callExp.Args)) + ", want=3. `push(map, key, value)`")
			}

			keyType := argTypes[0].Type.SubTypes[0].Value
			if keyType == "int" || keyType == "float" || keyType == "string" || keyType == "char" || keyType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("Key type mismatch. Expected " + keyType + ", got=array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("Key type mismatch. Expected " + keyType + ", got=hashmap")
				}
				if argTypes[1].Type.Value != keyType {
					return expType{}, errors.New("Key type mismatch. Expected " + keyType + ", got=" + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("Key type not supported. got=" + keyType)
			}

			valueType := argTypes[0].Type.SubTypes[1].Value
			if valueType == "int" || valueType == "float" || valueType == "string" || valueType == "char" || valueType == "bool" {
				if argTypes[2].Type.IsArray {
					return expType{}, errors.New("Value type mismatch. Expected " + valueType + ", got=array")
				} else if argTypes[2].Type.IsHash {
					return expType{}, errors.New("Value type mismatch. Expected " + valueType + ", got=hashmap")
				}
				if argTypes[2].Type.Value != valueType {
					return expType{}, errors.New("Value type mismatch. Expected " + valueType + ", got " + argTypes[2].Type.Value)
				}
			} else {
				return expType{}, errors.New("Value type not supported. got=" + valueType)
			}

			return argTypes[0], nil
		} else {
			return expType{}, errors.New("Data structure not supported by `push`. got=" + argTypes[0].Type.Value + ", want= array or hashmap")
		}
	case "pop":
		if argTypes[0].Type.IsArray {
			if len(callExp.Args) != 1 && len(callExp.Args) != 2 {
				return expType{}, errors.New("Wrong number of arguments for `pop` on array. got=" + strconv.Itoa(len(callExp.Args)) + ", want=1 or 2. `pop(array)` or `pop(array, index)`")
			}
			if len(callExp.Args) == 1 {
				return expType{Type: ast.Type{Value: argTypes[0].Type.Value, IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
			} else {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("Index must be an integer. got=array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("Index must be an integer. got=map")
				} else if argTypes[1].Type.Value != "int" {
					return expType{}, errors.New("Index must be an integer. got=" + argTypes[1].Type.Value)
				}
				return expType{Type: ast.Type{Value: argTypes[0].Type.Value, IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
			}
		} else if argTypes[0].Type.IsHash {
			return expType{}, errors.New("Data structure not supported by `pop`. got=map, want=array")
		} else {
			return expType{}, errors.New("Data structure not supported by `pop`. got=" + argTypes[0].Type.Value + ", want=array")
		}
	case "insert":
		if argTypes[0].Type.IsArray {
			if len(callExp.Args) != 3 {
				return expType{}, errors.New("Wrong number of arguments for `add` on array. got=" + strconv.Itoa(len(callExp.Args)) + ", want=3. `insert(array, index, element)`")
			}

			if argTypes[1].Type.IsArray {
				return expType{}, errors.New("Index must be an integer. got=array")
			} else if argTypes[1].Type.IsHash {
				return expType{}, errors.New("Index must be an integer. got=map")
			} else if argTypes[1].Type.Value != "int" {
				return expType{}, errors.New("Index must be an integer. got=" + argTypes[1].Type.Value)
			}

			arrayType := argTypes[0].Type.Value
			if arrayType == "int" || arrayType == "float" || arrayType == "string" || arrayType == "char" || arrayType == "bool" {
				if argTypes[2].Type.IsArray {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=array")
				} else if argTypes[2].Type.IsHash {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=hashmap")
				}
				if argTypes[2].Type.Value != arrayType {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=" + argTypes[2].Type.Value)
				}
			} else {
				return expType{}, errors.New("Array type not supported. got=" + argTypes[0].Type.Value + ", want= int, float, string, char or bool")
			}
			return argTypes[0], nil
		} else if argTypes[0].Type.IsHash {
			return expType{}, errors.New("Data structure not supported by `insert`. got=map, want=array")
		} else {
			return expType{}, errors.New("Data structure not supported by `insert`. got=" + argTypes[0].Type.Value + ", want=array")
		}
	case "remove":
		if len(callExp.Args) != 2 {
			return expType{}, errors.New("Wrong number of arguments for `remove` on array. got=" + strconv.Itoa(len(callExp.Args)) + ", want=2. For Array: `remove(array, element)`, for Hashmap: `remove(map, key)`")
		}
		if argTypes[0].Type.IsArray {
			if len(callExp.Args) != 2 {
				return expType{}, errors.New("Wrong number of arguments for `remove` on array. got=" + strconv.Itoa(len(callExp.Args)) + ", want=2. `remove(array, element)`")
			}

			arrayType := argTypes[0].Type.Value
			if arrayType == "int" || arrayType == "float" || arrayType == "string" || arrayType == "char" || arrayType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=hashmap")
				}
				if argTypes[1].Type.Value != arrayType {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=" + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("Array type not supported. got=" + arrayType + ", want= int, float, string, char or bool")
			}
			return argTypes[0], nil
		} else if argTypes[0].Type.IsHash {
			if len(callExp.Args) != 2 {
				return expType{}, errors.New("Wrong number of arguments for `remove` on hash. got=" + strconv.Itoa(len(callExp.Args)) + ", want=2. `remove(map, key)`")
			}

			keyType := argTypes[0].Type.SubTypes[0].Value
			if keyType == "int" || keyType == "float" || keyType == "string" || keyType == "char" || keyType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("Key type mismatch. Expected " + keyType + ", got=array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("Key type mismatch. Expected " + keyType + ", got=hashmap")
				}
				if argTypes[1].Type.Value != keyType {
					return expType{}, errors.New("Key type mismatch. Expected " + keyType + ", got=" + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("Key type not supported. got=" + keyType)
			}
			return expType{Type: *argTypes[0].Type.SubTypes[1], CallExp: false}, nil
		} else {
			return expType{}, errors.New("Data structure not supported by `remove`. got=" + argTypes[0].Type.Value)
		}
	case "getIndex":
		if argTypes[0].Type.IsArray {
			if len(callExp.Args) != 2 {
				return expType{}, errors.New("Wrong number of arguments for `getIndex` on array. got=" + strconv.Itoa(len(callExp.Args)) + ", want=2. `getIndex(array, element)`")
			}

			arrayType := argTypes[0].Type.Value
			if arrayType == "int" || arrayType == "float" || arrayType == "string" || arrayType == "char" || arrayType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=hashmap")
				}
				if argTypes[1].Type.Value != arrayType {
					return expType{}, errors.New("Argument type mismatch. Expected " + arrayType + ", got=" + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("Array type not supported. got=" + arrayType)
			}

			return expType{Type: ast.Type{Value: "int", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if argTypes[0].Type.IsHash {
			return expType{}, errors.New("Data structure not supported by `getIndex`. got=hashmap, want=array")
		} else {
			return expType{}, errors.New("Data structure not supported by `getIndex`. got=" + argTypes[0].Type.Value + ", want=array")
		}
	case "keys":
		if argTypes[0].Type.IsHash {
			if len(callExp.Args) != 1 {
				return expType{}, errors.New("Wrong number of arguments for `keys` on hash. got=" + strconv.Itoa(len(callExp.Args)) + ", want=1. `keys(map)`")
			}
			return expType{Type: ast.Type{Value: argTypes[0].Type.SubTypes[0].Value, IsArray: true, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if argTypes[0].Type.IsArray {
			return expType{}, errors.New("Data structure not supported by `keys`. got=array, want=hashmap")
		} else {
			return expType{}, errors.New("Data structure not supported by `keys`. got=" + argTypes[0].Type.Value + ", want=hashmap")
		}
	case "values":
		if argTypes[0].Type.IsHash {
			if len(callExp.Args) != 1 {
				return expType{}, errors.New("Wrong number of arguments for `values` on hash. got=" + strconv.Itoa(len(callExp.Args)) + ", want=1. `values(map)`")
			}
			return expType{Type: ast.Type{Value: argTypes[0].Type.SubTypes[1].Value, IsArray: true, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if argTypes[0].Type.IsArray {
			return expType{}, errors.New("Data structure not supported by `keys`. got=array, want=hashmap")
		} else {
			return expType{}, errors.New("Data structure not supported by `values`. got=" + argTypes[0].Type.Value + ", want=hashmap")
		}
	case "containsKey":
		if argTypes[0].Type.IsHash {
			if len(callExp.Args) != 2 {
				return expType{}, errors.New("Wrong number of arguments for `containsKey` on hash. got=" + strconv.Itoa(len(callExp.Args)) + ", want=2. `containsKey(map, key)`")
			}
			keyType := argTypes[0].Type.SubTypes[0].Value
			if keyType == "int" || keyType == "float" || keyType == "string" || keyType == "char" || keyType == "bool" {
				if argTypes[1].Type.IsArray {
					return expType{}, errors.New("Key type mismatch. Expected " + keyType + ", got=array")
				} else if argTypes[1].Type.IsHash {
					return expType{}, errors.New("Key type mismatch. Expected " + keyType + ", got=hashmap")
				}
				if argTypes[1].Type.Value != keyType {
					return expType{}, errors.New("Key type mismatch. Expected " + keyType + ", got=" + argTypes[1].Type.Value)
				}
			} else {
				return expType{}, errors.New("Key type not supported. got=" + keyType)
			}
			return expType{Type: ast.Type{Value: "bool", IsArray: false, IsHash: false, SubTypes: nil}, CallExp: false}, nil
		} else if argTypes[0].Type.IsArray {
			return expType{}, errors.New("Data structure not supported by `containsKey`. got=array, want=hashmap")
		} else {
			return expType{}, errors.New("Data structure not supported by `containsKey`. got=" + argTypes[0].Type.Value + ", want=hashmap")
		}
	default:
		return expType{}, errors.New("Unknown builtin function: " + callExp.Name.(*ast.Identifier).Value)
	}
}

// -----------------------------------------------------------------------------
// Helper Methods
// -----------------------------------------------------------------------------
func varTypeCheckerHelper(definedType ast.Type, expType ast.Type) error {
	if definedType.IsArray {
		if !expType.IsArray {
			if expType.IsHash {
				return errors.New("Defined type is of array, but got: hashmap")
			}
			return errors.New("Defined type is of array, but got: " + expType.Value)
		}
		if definedType.Value != expType.Value {
			if expType.Value == "" {
				return nil // means it is defined empty and can hold anytype
			}
			return errors.New("Array declared as " + definedType.Value + " but got: " + expType.Value)
		}
	} else if definedType.IsHash {
		if !expType.IsHash {
			if expType.IsArray {
				return errors.New("Defined type is of hashmap, but got: array")
			}
			return errors.New("Defined type is of hashmap, but got: " + expType.Value)
		}
		if definedType.SubTypes[0].Value != expType.SubTypes[0].Value {
			return errors.New("Hashmap Key type declared as " + definedType.SubTypes[0].Value + " but got: " + expType.SubTypes[0].Value)
		}
		if definedType.SubTypes[1].Value != expType.SubTypes[1].Value {
			return errors.New("Hashmap Value type declared as " + definedType.SubTypes[1].Value + " but got: " + expType.SubTypes[1].Value)
		}
	} else if definedType.Value != expType.Value {
		return errors.New("Defined type is " + definedType.Value + " but got: " + expType.Value)
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
						return errors.New("Call expression must return only 1 value for return statement")
					}
					expType.Type = *callFun.ReturnType[0].ReturnType
				}
				err = varTypeCheckerHelper(*currFunction.ReturnType[i].ReturnType, expType.Type)
				if err != nil {
					return err
				}
			}
		default:
			return errors.New("Can Only return expressions and datatypes. got: " + fmt.Sprintf("%T", entry))
		}
	}
	return nil
}
