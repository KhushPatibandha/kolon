package parser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/environment"
	ktype "github.com/KhushPatibandha/Kolon/src/kType"
)

// ------------------------------------------------------------------------------------------------------------------
// Expressions
// ------------------------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------------------------
// Identifier
// ------------------------------------------------------------------------------------------------------------------
func typeCheckIdent(ident *ast.Identifier,
	env *environment.Environment,
) (*ktype.TypeCheckResult, error) {
	if sym, ok := env.GetVar(ident.Value); ok {
		t := ktype.InternType(sym.Type)
		return &ktype.TypeCheckResult{Types: []*ktype.Type{t}, TypeLen: 1}, nil
	}
	if _, ok := env.GetFunc(ident.Value); ok {
		return nil,
			errors.New(
				"`" + ident.Value + "` is a function, did you mean to " +
					"call it as `" + ident.Value + "(...)`?",
			)
	}
	return nil, errors.New("variable `" + ident.Value + "` is undefined/not found")
}

// ------------------------------------------------------------------------------------------------------------------
// Integer
// ------------------------------------------------------------------------------------------------------------------
func typeCheckInteger() (*ktype.TypeCheckResult, error) {
	it := ktype.NewBaseType("int")
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{it},
		TypeLen: 1,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Float
// ------------------------------------------------------------------------------------------------------------------
func typeCheckFloat() (*ktype.TypeCheckResult, error) {
	ft := ktype.NewBaseType("float")
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{ft},
		TypeLen: 1,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Bool
// ------------------------------------------------------------------------------------------------------------------
func typeCheckBool() (*ktype.TypeCheckResult, error) {
	bt := ktype.NewBaseType("bool")
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{bt},
		TypeLen: 1,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// String
// ------------------------------------------------------------------------------------------------------------------
func typeCheckString() (*ktype.TypeCheckResult, error) {
	st := ktype.NewBaseType("string")
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{st},
		TypeLen: 1,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Char
// ------------------------------------------------------------------------------------------------------------------
func typeCheckChar() (*ktype.TypeCheckResult, error) {
	ct := ktype.NewBaseType("char")
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{ct},
		TypeLen: 1,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// HashMap
// ------------------------------------------------------------------------------------------------------------------
func typeCheckHashMap(exp *ast.HashMap,
	env *environment.Environment,
) (*ktype.TypeCheckResult, error) {
	if len(exp.Pairs) == 0 {
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewHashMapType(nil, nil)},
			TypeLen: 1,
		}, nil
	}

	var keyType *ktype.Type = nil
	var valueType *ktype.Type = nil

	for k, v := range exp.Pairs {
		key, err := typeCheckExp(k, env)
		if err != nil {
			return nil, err
		}
		value, err := typeCheckExp(v, env)
		if err != nil {
			return nil, err
		}
		if key.TypeLen != 1 {
			return nil,
				errors.New(
					"hashmap keys must be of a single type, got: " +
						strconv.Itoa(key.TypeLen) +
						". in case of call expression, it must return a single value",
				)
		}
		if key.Types[0].Kind != ktype.TypeBase {
			return nil,
				errors.New(
					"key in a hashmap can only be of `BaseType`, got: " +
						key.Types[0].TypeKindToString(),
				)
		}

		if value.TypeLen != 1 {
			return nil,
				errors.New(
					"hashmap values must be of a single type, got: " +
						strconv.Itoa(value.TypeLen) +
						". in case of call expression, it must return a single value",
				)
		}
		if keyType == nil && valueType == nil {
			keyType = key.Types[0]
			valueType = value.Types[0]
		} else {
			if !keyType.Equals(key.Types[0]) {
				return nil,
					errors.New(
						"hashmap can only have one type of key, got: " +
							keyType.String() + " and " + key.Types[0].String(),
					)
			}
			if !valueType.Equals(value.Types[0]) {
				return nil,
					errors.New(
						"hashmap can only have one type of value, got: " +
							valueType.String() + " and " + value.Types[0].String(),
					)
			}
		}
	}

	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{ktype.NewHashMapType(keyType, valueType)},
		TypeLen: 1,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Array
// ------------------------------------------------------------------------------------------------------------------
func typeCheckArray(exp *ast.Array,
	env *environment.Environment,
) (*ktype.TypeCheckResult, error) {
	if len(exp.Values) == 0 {
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewArrayType(nil)},
			TypeLen: 1,
		}, nil
	}

	var arrayType *ktype.Type = nil

	for _, ele := range exp.Values {
		e, err := typeCheckExp(ele, env)
		if err != nil {
			return nil, err
		}
		if e.TypeLen != 1 {
			return nil,
				errors.New(
					"array elements must be of a single type, got: " +
						strconv.Itoa(e.TypeLen) +
						". in case of call expression, it must return a single value",
				)
		}
		if arrayType == nil {
			arrayType = e.Types[0]
		} else {
			if !arrayType.Equals(e.Types[0]) {
				return nil,
					errors.New(
						"array can only have one type of elements, got: " +
							arrayType.String() + " and " + e.Types[0].String(),
					)
			}
		}
	}

	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{ktype.NewArrayType(arrayType)},
		TypeLen: 1,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Prefix
// ------------------------------------------------------------------------------------------------------------------
func typeCheckPrefix(exp *ast.Prefix,
	env *environment.Environment,
) (*ktype.TypeCheckResult, error) {
	right, err := typeCheckExp(exp.Right, env)
	if err != nil {
		return nil, err
	}
	if right.TypeLen != 1 {
		return nil,
			errors.New(
				"prefix operator can only be applied to a single type, got: " +
					strconv.Itoa(right.TypeLen) +
					". in case of call expression, it must return a single value",
			)
	}

	if right.Types[0].Kind != ktype.TypeBase {
		return nil,
			errors.New(
				"prefix operator can't be used with array or hashmap",
			)
	}
	switch exp.Operator {
	case "!":
		if right.Types[0].Name != "bool" {
			return nil,
				errors.New(
					"bang operator (`!`) can be only used with `bool` entities, got: " +
						right.Types[0].String(),
				)
		}
	case "-":
		if right.Types[0].Name != "int" && right.Types[0].Name != "float" {
			return nil,
				errors.New(
					"dash/minus (`-`) operator can be only used " +
						"with `int` and `float` entities, got: " +
						right.Types[0].String(),
				)
		}
	default:
		return nil,
			errors.New(
				"only 2 `prefix` operator's supported, bang (`!`) and " +
					"dash/minus (`-`), got: " + exp.Operator,
			)
	}
	return right, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Infix
// ------------------------------------------------------------------------------------------------------------------
func typeCheckInfix(exp *ast.Infix,
	env *environment.Environment,
) (*ktype.TypeCheckResult, error) {
	right, err := typeCheckExp(exp.Right, env)
	if err != nil {
		return nil, err
	}
	if right.TypeLen != 1 {
		return nil, errors.New(
			"infix operator can only be applied to a single type on the right side, got: " +
				strconv.Itoa(right.TypeLen) +
				". in case of call expression, it must return a single value",
		)
	}
	return typeCheckInfixWithRightType(exp, right.Types[0], env)
}

func typeCheckInfixWithRightType(exp *ast.Infix,
	right *ktype.Type,
	env *environment.Environment,
) (*ktype.TypeCheckResult, error) {
	left, err := typeCheckExp(exp.Left, env)
	if err != nil {
		return nil, err
	}
	if left.TypeLen != 1 {
		return nil, errors.New(
			"infix operator can only be applied to a single type on the left side, got: " +
				strconv.Itoa(left.TypeLen) +
				". in case of call expression, it must return a single value",
		)
	}

	if left.Types[0].Kind == ktype.TypeHashMap || right.Kind == ktype.TypeHashMap {
		return nil, errors.New("hashmap can't be used with infix operations")
	}

	switch {
	case left.Types[0].Kind == ktype.TypeArray && right.Kind == ktype.TypeArray:
		if !left.Types[0].ElementType.Equals(right.ElementType) {
			return nil,
				errors.New(
					"can only add arrays of same type, got: `" +
						left.Types[0].String() + "` and `" +
						right.String() + "`",
				)
		}
		switch exp.Operator {
		case "+":
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewArrayType(left.Types[0].ElementType)},
				TypeLen: 1,
			}, nil
		case "==", "!=":
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewBaseType("bool")},
				TypeLen: 1,
			}, nil
		default:
			return nil,
				errors.New(
					"can only use `+`, `==`, `!=` infix operators with 2 arrays, got: " +
						exp.Operator,
				)
		}
	case left.Types[0].Name == "int" && right.Name == "int":
		switch exp.Operator {
		case "+", "-", "*", "/", "%", "|", "&":
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewBaseType("int")},
				TypeLen: 1,
			}, nil
		case ">", "<", "<=", ">=", "==", "!=":
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewBaseType("bool")},
				TypeLen: 1,
			}, nil
		default:
			return nil,
				errors.New(
					"can only use `+`, `-`, `*`, `/`, `%`, `>`, `<`, " +
						"`<=`, `>=`, `!=`, `==`, `|`, `&` " +
						"infix operators with 2 `int`, got: " +
						exp.Operator,
				)
		}
	case left.Types[0].Name == "float" && right.Name == "float",
		((left.Types[0].Name == "int" && right.Name == "float") ||
			(left.Types[0].Name == "float" && right.Name == "int")):
		switch exp.Operator {
		case "+", "-", "*", "/":
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewBaseType("float")},
				TypeLen: 1,
			}, nil
		case ">", "<", "<=", ">=", "==", "!=":
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewBaseType("bool")},
				TypeLen: 1,
			}, nil
		default:
			return nil,
				errors.New(
					"can only use `+`, `-`, `*`, `/`, `>`, `<`, `<=`, `>=`, `!=`, `==` " +
						"infix operators with 2 `float`, got: " +
						exp.Operator,
				)
		}
	case left.Types[0].Name == "string" && right.Name == "string":
		switch exp.Operator {
		case "+":
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewBaseType("string")},
				TypeLen: 1,
			}, nil
		case "==", "!=":
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewBaseType("bool")},
				TypeLen: 1,
			}, nil
		default:
			return nil,
				errors.New(
					"can only use `+`, `==`, `!=` infix operators with 2 `string`, got: " +
						exp.Operator,
				)
		}
	case left.Types[0].Name == "char" && right.Name == "char":
		switch exp.Operator {
		case "+":
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewBaseType("string")},
				TypeLen: 1,
			}, nil
		case "==", "!=":
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewBaseType("bool")},
				TypeLen: 1,
			}, nil
		default:
			return nil,
				errors.New(
					"can only use `+`, `==`, `!=` infix operators with 2 `char`, got: " +
						exp.Operator,
				)
		}
	case left.Types[0].Name == "bool" && right.Name == "bool":
		switch exp.Operator {
		case "==", "!=", "&&", "||":
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewBaseType("bool")},
				TypeLen: 1,
			}, nil
		default:
			return nil,
				errors.New(
					"can only use `==`, `!=`, `&&`, `||` infix operators with 2 `bool`, got: " +
						exp.Operator,
				)
		}
	default:
		return nil,
			errors.New(
				"invalid `infix` operation with variable types on left and right, got: `" +
					left.Types[0].String() + "` and `" +
					right.String() + "`",
			)
	}
}

// ------------------------------------------------------------------------------------------------------------------
// Postfix
// ------------------------------------------------------------------------------------------------------------------
func typeCheckPostfix(exp *ast.Postfix,
	env *environment.Environment,
) (*ktype.TypeCheckResult, error) {
	left, err := typeCheckExp(exp.Left, env)
	if err != nil {
		return nil, err
	}
	if left.TypeLen != 1 {
		return nil,
			errors.New(
				"postfix operator can only be applied to a single type, got: " +
					strconv.Itoa(left.TypeLen) +
					". in case of call expression, it must return a single value",
			)
	}
	if left.Types[0].Kind != ktype.TypeBase {
		return nil, errors.New("postfix operator can't be used with array or hashmap")
	}
	if left.Types[0].Name != "int" && left.Types[0].Name != "float" {
		return nil,
			errors.New(
				"only `int` and `float` datatypes supported with `postfix` operation, got: " +
					left.Types[0].String(),
			)
	}
	if exp.Operator != "++" && exp.Operator != "--" {
		return nil,
			errors.New(
				"only 2 `postfix` operator's supported, increment (`++`) and decrement (`--`), got: " +
					exp.Operator,
			)
	}
	return left, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Assignment
// ------------------------------------------------------------------------------------------------------------------
func typeCheckAssignment(exp *ast.Assignment,
	env *environment.Environment,
) (*ktype.TypeCheckResult, error) {
	right, err := typeCheckExp(exp.Right, env)
	if err != nil {
		return nil, err
	}
	if right.TypeLen != 1 {
		return nil,
			errors.New(
				"assignment can only be applied to a single type on the right side, got: " +
					strconv.Itoa(right.TypeLen) +
					". in case of call expression, it must return a single value",
			)
	}

	return typeCheckAssignmentWithRightType(exp, right.Types[0], env, true)
}

func typeCheckAssignmentWithRightType(exp *ast.Assignment, right *ktype.Type,
	env *environment.Environment,
	passInfixExp bool,
) (*ktype.TypeCheckResult, error) {
	left, err := typeCheckIdent(exp.Left, env)
	if err != nil {
		return nil, err
	}
	if left.TypeLen != 1 {
		return nil,
			errors.New(
				"assignment can only be applied to a single type on the left side, got: " +
					strconv.Itoa(left.TypeLen) +
					". in case of call expression, it must return a single value",
			)
	}
	leftSym, _ := env.GetVar(exp.Left.Value)
	if leftSym.IdentType == environment.CONST {
		return nil,
			errors.New(
				"variable `" + leftSym.Ident.Value +
					"` is a constant, can't re-assign value to a constant variable",
			)
	}

	switch exp.Operator {
	case "=":
		if !left.Types[0].Equals(right) {
			return nil,
				errors.New(
					"type mismatch at the time of assignment, got: `" +
						left.Types[0].String() + "` on left and `" +
						right.String() + "` on right",
				)
		}
	case "+=", "-=", "*=", "/=", "%=":
		infixExp := &ast.Infix{
			Left:     exp.Left,
			Operator: strings.TrimSuffix(exp.Operator, "="),
		}

		var iT *ktype.TypeCheckResult
		if passInfixExp {
			infixExp.Right = exp.Right
			iT, err = typeCheckInfix(infixExp, env)
		} else {
			iT, err = typeCheckInfixWithRightType(infixExp, right, env)
		}
		if err != nil {
			return nil, err
		}

		if !left.Types[0].Equals(iT.Types[0]) {
			return nil,
				errors.New(
					"type mismatch at the time of assignment, got: `" +
						left.Types[0].String() + "` on left and `" +
						iT.Types[0].String() + "` on right",
				)
		}
	default:
		return nil,
			errors.New(
				"only `=`, `+=`, `-=`, `*=`, `/=`, `%=` " +
					"`assignment` operators supported, got: " +
					exp.Operator,
			)
	}
	return left, nil
}

// ------------------------------------------------------------------------------------------------------------------
// IndexExpression
// ------------------------------------------------------------------------------------------------------------------
func typeCheckIndexExp(exp *ast.IndexExpression,
	env *environment.Environment,
) (*ktype.TypeCheckResult, error) {
	left, err := typeCheckExp(exp.Left, env)
	if err != nil {
		return nil, err
	}
	if left.TypeLen != 1 {
		return nil,
			errors.New(
				"index expression can only be applied to a single type on the left side, got: " +
					strconv.Itoa(left.TypeLen) +
					". in case of call expression, it must return a single value",
			)
	}

	index, err := typeCheckExp(exp.Index, env)
	if err != nil {
		return nil, err
	}
	if index.TypeLen != 1 {
		return nil,
			errors.New(
				"index expression can only be applied to a single type on the index side, got: " +
					strconv.Itoa(index.TypeLen) +
					". in case of call expression, it must return a single value",
			)
	}

	switch left.Types[0].Kind {
	case ktype.TypeArray:
		if left.Types[0].ElementType == nil {
			return nil, errors.New("array is empty, can't index empty array")
		}
		if index.Types[0].Kind != ktype.TypeBase || index.Types[0].Name != "int" {
			return nil,
				errors.New(
					"array index must be an `int`, got: " +
						index.Types[0].String(),
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{left.Types[0].ElementType},
			TypeLen: 1,
		}, nil
	case ktype.TypeHashMap:
		if left.Types[0].KeyType == nil || left.Types[0].ValueType == nil {
			return nil, errors.New("hashmap is empty, can't index empty hashmap")
		}
		if !left.Types[0].KeyType.Equals(index.Types[0]) {
			return nil,
				errors.New(
					"hashmap index type must be of datatype `" +
						left.Types[0].KeyType.String() + "`, got: " +
						index.Types[0].String(),
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{left.Types[0].ValueType},
			TypeLen: 1,
		}, nil
	default:
		if left.Types[0].Kind == ktype.TypeBase && left.Types[0].Name == "string" &&
			index.Types[0].Kind == ktype.TypeBase && index.Types[0].Name == "int" {
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{ktype.NewBaseType("char")},
				TypeLen: 1,
			}, nil
		} else {
			return nil,
				errors.New(
					"index operation not supported for " + left.Types[0].String() +
						"[" + index.Types[0].String() + "]",
				)
		}
	}
}

// ------------------------------------------------------------------------------------------------------------------
// CallExpression
// ------------------------------------------------------------------------------------------------------------------
func typeCheckCallExp(exp *ast.CallExpression,
	env *environment.Environment,
) (*ktype.TypeCheckResult, error) {
	funcSym, ok := env.GetFunc(exp.Name.Value)
	if !ok {
		return nil,
			errors.New(
				"function `" + exp.Name.Value +
					"` not found",
			)
	}
	if funcSym.Func.Builtin {
		return typeCheckBuiltin(exp, env)
	}
	if len(funcSym.Func.Function.Parameters) != len(exp.Args) {
		return nil,
			errors.New(
				"number of arguments does not match the number of parameters for function `" +
					funcSym.Ident.Value + "`, got: " + strconv.Itoa(len(exp.Args)) +
					", expected: " + strconv.Itoa(len(funcSym.Func.Function.Parameters)),
			)
	}

	if exp.Args != nil {
		for i, arg := range exp.Args {
			argType, err := typeCheckExp(arg, env)
			if err != nil {
				return nil, err
			}
			if argType.TypeLen != 1 {
				return nil,
					errors.New(
						"argument at position " + strconv.Itoa(i+1) +
							" must be of a single type, got: " +
							strconv.Itoa(argType.TypeLen) +
							". in case of call expression, it must return a single value",
					)
			}
			paramType := funcSym.Func.Function.Parameters[i].ParameterType
			if !paramType.Equals(argType.Types[0]) {
				return nil,
					errors.New(
						"type mismatch for argument at position " + strconv.Itoa(i+1) +
							" for function call `" + exp.Name.Value + "`, expected: `" +
							paramType.String() + "`, got: `" +
							argType.Types[0].String() + "`",
					)
			}
		}
	}
	if funcSym.Func.Function.ReturnTypes == nil {
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{},
			TypeLen: 0,
		}, nil
	}
	return &ktype.TypeCheckResult{
		Types:   funcSym.Func.Function.ReturnTypes,
		TypeLen: len(funcSym.Func.Function.ReturnTypes),
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Builtin
// ------------------------------------------------------------------------------------------------------------------
func typeCheckBuiltin(exp *ast.CallExpression,
	env *environment.Environment,
) (*ktype.TypeCheckResult, error) {
	var argTypes []*ktype.Type
	for _, arg := range exp.Args {
		t, err := typeCheckExp(arg, env)
		if err != nil {
			return nil, err
		}
		if t.TypeLen != 1 {
			return nil,
				errors.New(
					"builtin function arguments must be of a single type, got: " +
						strconv.Itoa(t.TypeLen) +
						". in case of call expression, it must return a single value",
				)
		}
		argTypes = append(argTypes, t.Types[0])
	}

	switch exp.Name.Value {
	case "len":
		if exp.Args == nil || len(exp.Args) != 1 {
			return nil,
				errors.New(
					"wrong number of arguments for `len`, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 1",
				)
		}
		if argTypes[0].Kind != ktype.TypeArray &&
			argTypes[0].Kind != ktype.TypeHashMap &&
			(argTypes[0].Kind != ktype.TypeBase || argTypes[0].Name != "string") {
			return nil,
				errors.New(
					"argument for `len` not supported, got: " +
						argTypes[0].String() + ", want: array, hashmap or `string`",
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewBaseType("int")},
			TypeLen: 1,
		}, nil
	case "toString":
		if exp.Args == nil || len(exp.Args) != 1 {
			return nil,
				errors.New(
					"wrong number of arguments for `toString`, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 1",
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewBaseType("string")},
			TypeLen: 1,
		}, nil
	case "toInt":
		if exp.Args == nil || len(exp.Args) != 1 {
			return nil,
				errors.New(
					"wrong number of arguments for `toInt`, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 1",
				)
		}
		switch argTypes[0].Kind {
		case ktype.TypeArray:
			return nil,
				errors.New(
					"argument for `toInt` not supported, got: array, " +
						"want: `int`, `float`, `string`, `char`",
				)
		case ktype.TypeHashMap:
			return nil,
				errors.New(
					"argument for `toInt` not supported, " +
						"got: hashmap, want: `int`, `float`, `string`, `char`",
				)
		case ktype.TypeBase:
			if argTypes[0].Name == "bool" {
				return nil,
					errors.New(
						"argument for `toInt` not supported, got: bool" +
							", want: `int`, `float`, `string`, `char`",
					)
			}
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewBaseType("int")},
			TypeLen: 1,
		}, nil
	case "toFloat":
		if exp.Args == nil || len(exp.Args) != 1 {
			return nil,
				errors.New(
					"wrong number of arguments for `toFloat`, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 1",
				)
		}
		switch argTypes[0].Kind {
		case ktype.TypeArray:
			return nil,
				errors.New(
					"argument for `toFloat` not supported, " +
						"got: array, want: `int` or `float` or `string`",
				)
		case ktype.TypeHashMap:
			return nil,
				errors.New(
					"argument for `toFloat` not supported, " +
						"got: hashmap, want: `int` or `float` or `string`",
				)
		case ktype.TypeBase:
			if argTypes[0].Name == "bool" || argTypes[0].Name == "char" {
				return nil,
					errors.New(
						"argument for `toFloat` not supported, got: " +
							argTypes[0].String() + ", want: `int` or `float` or `string`",
					)
			}
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewBaseType("float")},
			TypeLen: 1,
		}, nil
	case "print":
		if exp.Args == nil || len(exp.Args) != 1 {
			return nil,
				errors.New(
					"wrong number of arguments for `print`, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 1",
				)
		}
		return &ktype.TypeCheckResult{Types: []*ktype.Type{}, TypeLen: 0}, nil
	case "println":
		if exp.Args != nil && len(exp.Args) != 1 {
			return nil,
				errors.New(
					"wrong number of arguments for `println`, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 0 or 1",
				)
		}
		return &ktype.TypeCheckResult{Types: []*ktype.Type{}, TypeLen: 0}, nil
	case "scan", "scanln":
		fncName := exp.Name.Value
		if exp.Args != nil && len(exp.Args) != 1 && len(exp.Args) != 2 {
			return nil,
				errors.New(
					"wrong number of arguments for `" +
						fncName + "`, got: " + strconv.Itoa(len(exp.Args)) +
						", want: 0 or 1 or 2. " + fncName +
						"() or " + fncName + "(prompt) or " +
						fncName + "(prompt, newline[true/false])",
				)
		}
		if exp.Args != nil {
			if argTypes[0].Kind != ktype.TypeBase || argTypes[0].Name != "string" {
				return nil,
					errors.New(
						"argument for `" + fncName +
							"` not supported, got: " + argTypes[0].String() + ", want: `string`",
					)
			}
			if len(exp.Args) == 2 {
				if argTypes[1].Kind != ktype.TypeBase || argTypes[1].Name != "bool" {
					return nil,
						errors.New(
							"2nd argument for `" + fncName +
								"` not supported, got: " + argTypes[1].String() + ", want: `bool`",
						)
				}
			}
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewBaseType("string")},
			TypeLen: 1,
		}, nil
	case "push":
		if exp.Args == nil || (len(exp.Args) != 2 && len(exp.Args) != 3) {
			return nil,
				errors.New(
					"wrong number of arguments for `push`, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 2 or 3. for array: " +
						"`push(array, element)` and for hashmap: `push(map, key, value)`",
				)
		}
		if argTypes[0].Kind != ktype.TypeArray && argTypes[0].Kind != ktype.TypeHashMap {
			return nil,
				errors.New(
					"data structure not supported by `push`, got: " +
						argTypes[0].String() + ", want: array or hashmap",
				)
		}
		if argTypes[0].Kind == ktype.TypeArray {
			if len(exp.Args) != 2 {
				return nil,
					errors.New(
						"wrong number of arguments for `push` for array, got: " +
							strconv.Itoa(len(exp.Args)) + ", want: 2. `push(array, element)`",
					)
			}
			if !argTypes[0].ElementType.Equals(argTypes[1]) {
				return nil,
					errors.New(
						"argument type mismatch for `push`, expected element to be " +
							argTypes[0].ElementType.String() + ", got: " + argTypes[1].String(),
					)
			}
		} else {
			if len(exp.Args) != 3 {
				return nil,
					errors.New(
						"wrong number of arguments for `push` for hashmap, got: " +
							strconv.Itoa(len(exp.Args)) + ", want: 3. `push(map, key, value)`",
					)
			}
			if !argTypes[0].KeyType.Equals(argTypes[1]) {
				return nil,
					errors.New(
						"key type mismatch for `push`, expected key to be " +
							argTypes[0].KeyType.String() + ", got: " + argTypes[1].String(),
					)
			}
			if !argTypes[0].ValueType.Equals(argTypes[2]) {
				return nil,
					errors.New(
						"value type mismatch for `push`, expected value to be " +
							argTypes[0].ValueType.String() + ", got: " + argTypes[2].String(),
					)
			}
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{argTypes[0]},
			TypeLen: 1,
		}, nil
	case "pop":
		if exp.Args == nil || (len(exp.Args) != 1 && len(exp.Args) != 2) {
			return nil,
				errors.New(
					"wrong number of arguments for `pop` for array, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 1 or 2. `pop(array)` or `pop(array, index)`",
				)
		}
		if argTypes[0].Kind != ktype.TypeArray {
			return nil,
				errors.New(
					"data structure not supported by `pop`, got: " +
						argTypes[0].String() + ", want: array",
				)
		}
		if len(exp.Args) == 2 {
			if argTypes[1].Kind != ktype.TypeBase || argTypes[1].Name != "int" {
				return nil,
					errors.New(
						"index must be an integer for `pop`, got: " + argTypes[1].String(),
					)
			}
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{argTypes[0].ElementType},
			TypeLen: 1,
		}, nil
	case "insert":
		if exp.Args == nil || len(exp.Args) != 3 {
			return nil,
				errors.New(
					"wrong number of arguments for `insert` for array, got: " +
						strconv.Itoa(len(exp.Args)) +
						", want: 3. `insert(array, index, element)`",
				)
		}
		if argTypes[0].Kind != ktype.TypeArray {
			return nil,
				errors.New(
					"data structure not supported by `insert`, got: " +
						argTypes[0].String() + ", want: array",
				)
		}
		if argTypes[1].Kind != ktype.TypeBase || argTypes[1].Name != "int" {
			return nil,
				errors.New(
					"index must be an `int` for `insert`, got: " + argTypes[1].String(),
				)
		}
		if !argTypes[0].ElementType.Equals(argTypes[2]) {
			return nil,
				errors.New(
					"argument type mismatch for `insert`, expected element to be " +
						argTypes[0].ElementType.String() + ", got: " + argTypes[2].String(),
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{argTypes[0]},
			TypeLen: 1,
		}, nil
	case "delete":
		if exp.Args == nil || len(exp.Args) != 2 {
			return nil,
				errors.New(
					"wrong number of arguments for `delete`, got: " +
						strconv.Itoa(len(exp.Args)) +
						", want: 2. for array: `remove(array, element)`, " +
						"for hashmap: `remove(map, key)`",
				)
		}
		if argTypes[0].Kind != ktype.TypeArray && argTypes[0].Kind != ktype.TypeHashMap {
			return nil,
				errors.New(
					"data structure not supported by `delete`, got: " +
						argTypes[0].String() + ", want: array or hashmap",
				)
		}
		if argTypes[0].Kind == ktype.TypeArray {
			if !argTypes[0].ElementType.Equals(argTypes[1]) {
				return nil,
					errors.New(
						"argument type mismatch for `delete`, expected element to be " +
							argTypes[0].ElementType.String() + ", got: " + argTypes[1].String(),
					)
			}
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{argTypes[0].ElementType},
				TypeLen: 1,
			}, nil
		} else {
			if !argTypes[0].KeyType.Equals(argTypes[1]) {
				return nil,
					errors.New(
						"key type mismatch for `delete`, expected key to be " +
							argTypes[0].KeyType.String() + ", got: " + argTypes[1].String(),
					)
			}
			return &ktype.TypeCheckResult{
				Types:   []*ktype.Type{argTypes[0].ValueType},
				TypeLen: 1,
			}, nil
		}
	case "remove":
		if exp.Args == nil || len(exp.Args) != 2 {
			return nil,
				errors.New(
					"wrong number of arguments for `remove`, got: " +
						strconv.Itoa(len(exp.Args)) +
						", want: 2. for array: `remove(array, element)`, " +
						"for hashmap: `remove(map, key)`",
				)
		}
		if argTypes[0].Kind != ktype.TypeArray && argTypes[0].Kind != ktype.TypeHashMap {
			return nil,
				errors.New(
					"data structure not supported by `remove`, got: " +
						argTypes[0].String() + ", want: array or hashmap",
				)
		}
		if argTypes[0].Kind == ktype.TypeArray {
			if !argTypes[0].ElementType.Equals(argTypes[1]) {
				return nil,
					errors.New(
						"argument type mismatch for `remove`, expected element to be " +
							argTypes[0].ElementType.String() + ", got: " + argTypes[1].String(),
					)
			}
		} else {
			if !argTypes[0].KeyType.Equals(argTypes[1]) {
				return nil,
					errors.New(
						"key type mismatch for `remove`, expected key to be " +
							argTypes[0].KeyType.String() + ", got: " + argTypes[1].String(),
					)
			}
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{argTypes[0]},
			TypeLen: 1,
		}, nil
	case "getIndex":
		if exp.Args == nil || len(exp.Args) != 2 {
			return nil,
				errors.New(
					"wrong number of arguments for `getIndex` for array, got: " +
						strconv.Itoa(len(exp.Args)) +
						", want: 2. `getIndex(array, element)`",
				)
		}
		if argTypes[0].Kind != ktype.TypeArray {
			return nil,
				errors.New(
					"data structure not supported by `getIndex`, got: " +
						argTypes[0].String() + ", want: array",
				)
		}
		if !argTypes[0].ElementType.Equals(argTypes[1]) {
			return nil,
				errors.New(
					"argument type mismatch for `getIndex`, expected element to be " +
						argTypes[0].ElementType.String() + ", got: " + argTypes[1].String(),
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewBaseType("int")},
			TypeLen: 1,
		}, nil
	case "keys":
		if exp.Args == nil || len(exp.Args) != 1 {
			return nil,
				errors.New(
					"wrong number of arguments for `keys` for hashmap, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 1. `keys(map)`",
				)
		}
		if argTypes[0].Kind != ktype.TypeHashMap {
			return nil,
				errors.New(
					"data structure not supported by `keys`, got: " +
						argTypes[0].String() + ", want: hashmap",
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewArrayType(argTypes[0].KeyType)},
			TypeLen: 1,
		}, nil
	case "values":
		if exp.Args == nil || len(exp.Args) != 1 {
			return nil,
				errors.New(
					"wrong number of arguments for `values` for hashmap, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 1. `values(map)`",
				)
		}
		if argTypes[0].Kind != ktype.TypeHashMap {
			return nil,
				errors.New(
					"data structure not supported by `values`, got: " +
						argTypes[0].String() + ", want: hashmap",
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewArrayType(argTypes[0].ValueType)},
			TypeLen: 1,
		}, nil
	case "containsKey":
		if exp.Args == nil || len(exp.Args) != 2 {
			return nil,
				errors.New(
					"wrong number of arguments for `containsKey` for hashmap, got: " +
						strconv.Itoa(len(exp.Args)) +
						", want: 2. `containsKey(map, key)`",
				)
		}
		if argTypes[0].Kind != ktype.TypeHashMap {
			return nil,
				errors.New(
					"data structure not supported by `containsKey`, got: " +
						argTypes[0].String() + ", want: hashmap",
				)
		}
		if !argTypes[0].KeyType.Equals(argTypes[1]) {
			return nil,
				errors.New(
					"key type mismatch for `containsKey`, expected key to be " +
						argTypes[0].KeyType.String() + ", got: " + argTypes[1].String(),
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewBaseType("bool")},
			TypeLen: 1,
		}, nil
	case "typeOf":
		if exp.Args == nil || len(exp.Args) != 1 {
			return nil,
				errors.New(
					"wrong number of arguments for `typeOf`, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 1",
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewBaseType("string")},
			TypeLen: 1,
		}, nil
	case "slice":
		if exp.Args == nil || (len(exp.Args) != 3 && len(exp.Args) != 4) {
			return nil,
				errors.New(
					"wrong number of arguments for `slice`, got: " +
						strconv.Itoa(len(exp.Args)) +
						", want: 3 or 4. `slice(array/string, " +
						"start, end)` or `slice(array/string, start, end, step)`",
				)
		}
		if argTypes[0].Kind != ktype.TypeArray &&
			(argTypes[0].Kind != ktype.TypeBase || argTypes[0].Name != "string") {
			return nil,
				errors.New(
					"data structure not supported by `slice`, got: " +
						argTypes[0].String() + ", want: array or string",
				)
		}
		if argTypes[1].Kind != ktype.TypeBase || argTypes[1].Name != "int" {
			return nil,
				errors.New(
					"start index must be an `int` for `slice`, got: " +
						argTypes[1].String(),
				)
		}
		if argTypes[2].Kind != ktype.TypeBase || argTypes[2].Name != "int" {
			return nil,
				errors.New(
					"end index must be an `int` for `slice`, got: " +
						argTypes[2].String(),
				)
		}
		if len(exp.Args) == 4 {
			if argTypes[3].Kind != ktype.TypeBase || argTypes[3].Name != "int" {
				return nil,
					errors.New(
						"step must be an `int` for `slice`, got: " +
							argTypes[3].String(),
					)
			}
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{argTypes[0]},
			TypeLen: 1,
		}, nil
	case "copy":
		if exp.Args == nil || len(exp.Args) != 1 {
			return nil,
				errors.New(
					"wrong number of arguments for `copy`, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 1",
				)
		}
		if argTypes[0].Kind != ktype.TypeArray &&
			argTypes[0].Kind != ktype.TypeHashMap {
			return nil,
				errors.New(
					"data structure not supported by `copy`, got: " +
						argTypes[0].String() + ", want: array or hashmap",
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{argTypes[0]},
			TypeLen: 1,
		}, nil
	case "equals":
		if exp.Args == nil || len(exp.Args) != 2 {
			return nil,
				errors.New(
					"wrong number of arguments for `equals`, got: " +
						strconv.Itoa(len(exp.Args)) + ", want: 2",
				)
		}
		if !argTypes[0].Equals(argTypes[1]) {
			return nil,
				errors.New(
					"type mismatch for arguments of `equals`, got: `" +
						argTypes[0].String() + "` and `" +
						argTypes[1].String() + "`",
				)
		}
		return &ktype.TypeCheckResult{
			Types:   []*ktype.Type{ktype.NewBaseType("bool")},
			TypeLen: 1,
		}, nil
	default:
		return nil,
			errors.New(
				"unknown builtin function `" +
					exp.Name.Value + "`",
			)
	}
}
