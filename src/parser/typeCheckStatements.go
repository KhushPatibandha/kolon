package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/environment"
	"github.com/KhushPatibandha/Kolon/src/kType"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

// ------------------------------------------------------------------------------------------------------------------
// Statements
// ------------------------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------------------------
// VarAndConst
// ------------------------------------------------------------------------------------------------------------------
func typeCheckVarAndConst(stmt *ast.VarAndConst, env *environment.Environment) error {
	if stmt.Value == nil {
		switch stmt.Type.Kind {
		case ktype.TypeArray:
			return errors.New(
				"array `" + stmt.Name.Value + "` must always be " +
					"initialized while declaring, for empty array use `[]`")
		case ktype.TypeHashMap:
			return errors.New(
				"hashmap `" + stmt.Name.Value +
					"` must always be initialized while declaring, for empty hashmap use `{}`")
		default:
			if stmt.Token.Kind == lexer.CONST {
				return errors.New(
					"const variable `" + stmt.Name.Value +
						"` must be initialized while declaring")
			}
		}
	}

	right := stmt.Value.GetType()
	if right.TypeLen != 1 {
		return errors.New(
			"variable (`var`) and constant (`const`) declarations " +
				"must be assigned a single value, got: " +
				strconv.Itoa(right.TypeLen) +
				". in case of call expression, it must return a single value",
		)
	}

	return typeCheckVarAndConstWithRightType(stmt, right.Types[0], env)
}

func typeCheckVarAndConstWithRightType(stmt *ast.VarAndConst,
	right *ktype.Type,
	env *environment.Environment,
) error {
	if sym, ok := env.GetVar(stmt.Name.Value); ok {
		if sym.IdentType == environment.CONST {
			return errors.New(
				"variable `" + sym.Ident.Value + "` is a constant," +
					" can't re-declare const variables",
			)
		}
		if !stmt.Type.Equals(sym.Type) {
			return errors.New(
				"variable `" + sym.Ident.Value + "` already declared as `" +
					sym.Type.String() + "` can't re-declare as `" + stmt.Type.String() + "`",
			)
		}
	}

	if !stmt.Type.Equals(right) {
		return errors.New(
			"type mismatch in variable/constant declaration, expected: " +
				stmt.Type.String() + ", got: " + right.String(),
		)
	}

	switch stmt.Token.Kind {
	case lexer.VAR:
		env.Set(&environment.Symbol{
			IdentType: environment.VAR,
			Ident:     stmt.Name,
			Type:      stmt.Type,
			Func:      nil,
			Env:       nil,
		})
	case lexer.CONST:
		env.Set(&environment.Symbol{
			IdentType: environment.CONST,
			Ident:     stmt.Name,
			Type:      stmt.Type,
			Func:      nil,
			Env:       nil,
		})
	}
	return nil
}

// ------------------------------------------------------------------------------------------------------------------
// Multi-Assignment
// ------------------------------------------------------------------------------------------------------------------
func typeCheckMultiAssign(stmt *ast.MultiAssignment, env *environment.Environment) error {
	if stmt.SingleFunctionCall {
		var call *ast.CallExpression
		var isCall bool

		varConst, ok := stmt.Objects[0].(*ast.VarAndConst)
		if ok {
			call, isCall = varConst.Value.(*ast.CallExpression)
		} else {
			call, isCall = stmt.Objects[0].(*ast.ExpressionStatement).
				Expression.(*ast.CallExpression)
		}

		if !isCall {
			return errors.New(
				"number of expressions on the right side of multi-assignment = 1, " +
					"expected a function call, got: " +
					fmt.Sprintf("%T", call),
			)
		}

		if len(call.Type) != len(stmt.Objects) {
			return errors.New(
				"number of return values from function call does not match " +
					"the number of variables in multi-assignment, expected: " +
					strconv.Itoa(len(stmt.Objects)) + ", got: " +
					strconv.Itoa(len(call.Type)),
			)
		}

		for i, obj := range stmt.Objects {
			var err error
			switch obj := obj.(type) {
			case *ast.VarAndConst:
				err = typeCheckVarAndConstWithRightType(obj, call.Type[i], env)
			case *ast.ExpressionStatement:
				if exp, ok := obj.Expression.(*ast.Assignment); ok {
					_, err = typeCheckAssignmentWithRightType(exp, call.Type[i], env, false)
				}
			}
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		for _, obj := range stmt.Objects {
			switch obj := obj.(type) {
			case *ast.VarAndConst:
				err := typeCheckVarAndConst(obj, env)
				if err != nil {
					return err
				}
			case *ast.ExpressionStatement:
				if exp, ok := obj.Expression.(*ast.Assignment); ok {
					_, err := typeCheckAssignment(exp, env)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

// ------------------------------------------------------------------------------------------------------------------
// Return
// ------------------------------------------------------------------------------------------------------------------
func typeCheckReturn(stmt *ast.Return, fun *ast.Function) error {
	if stmt.Value == nil && fun.ReturnTypes == nil {
		return nil
	}
	if stmt.Value == nil && fun.ReturnTypes != nil {
		return errors.New("not enough return values for function `" +
			fun.Name.Value + "`")
	}
	if stmt.Value != nil && fun.ReturnTypes == nil {
		return errors.New("too many return values for function `" +
			fun.Name.Value + "`, expected none")
	}
	if len(stmt.Value) != len(fun.ReturnTypes) {
		return errors.New(
			"number of return values does not match the number of return types, got: " +
				strconv.Itoa(len(stmt.Value)) +
				", expected: " + strconv.Itoa(len(fun.ReturnTypes)),
		)
	}
	for i, retExp := range stmt.Value {
		right := retExp.GetType()
		if right.TypeLen != 1 {
			return errors.New(
				"return expression must evaluate to a single value, got: " +
					strconv.Itoa(right.TypeLen) +
					". in case of call expression, it must return a single value",
			)
		}
		expectedType := fun.ReturnTypes[i]
		if !expectedType.Equals(right.Types[0]) {
			return errors.New(
				"type mismatch in return statement, expected: " +
					expectedType.String() + ", got: " + right.Types[0].String(),
			)
		}
	}
	return nil
}

// ------------------------------------------------------------------------------------------------------------------
// ForLoop
// ------------------------------------------------------------------------------------------------------------------
func typeCheckForLoop(stmt *ast.ForLoop, env *environment.Environment) error {
	left := stmt.Left
	var sym *environment.Symbol

	if _, ok := left.(*ast.VarAndConst); !ok {
		if _, ok := left.(*ast.ExpressionStatement); !ok {
			return errors.New(
				"expected a `var` statement or an assignment expression " +
					"after an open bracket (`(`) in `for loop` statement, got: " +
					fmt.Sprintf("%T", left),
			)
		} else {
			assign, ok := left.(*ast.ExpressionStatement).Expression.(*ast.Assignment)
			if !ok {
				return errors.New(
					"expected a `var` statement or an assignment expression " +
						"after an open bracket (`(`) in `for loop` statement, got: " +
						fmt.Sprintf("%T", left.(*ast.ExpressionStatement).Expression),
				)
			} else {
				if assign.Operator != "=" {
					return errors.New(
						"expected assignment-equal operator (`=`) in assignment " +
							"expression in `for loop` statement, got: " +
							assign.Operator,
					)
				}
				sym, _ = env.GetVar(assign.Left.Value)
			}
		}
	} else {
		sym, _ = env.GetVar(stmt.Left.(*ast.VarAndConst).Name.Value)
	}
	if sym.IdentType != environment.VAR {
		return errors.New("can't use `const` to define variable in `for loop` condition")
	}
	i, _ := typeCheckInteger()
	if !sym.Type.Equals(i.Types[0]) {
		return errors.New(
			"can only define variable in `for loop` condition as `int`, got: " +
				sym.Type.String(),
		)
	}

	err := typeCheckBoolCon(stmt.Middle, "for")
	if err != nil {
		return err
	}

	right := stmt.Right.GetType()
	if right.TypeLen != 1 {
		return errors.New(
			"`for` loop increment/decrement expression must evaluate to a single value, got: " +
				strconv.Itoa(right.TypeLen) +
				". in case of call expression, it must return a single value",
		)
	}
	if !right.Types[0].Equals(i.Types[0]) {
		return errors.New(
			"`for` loop increment/decrement expression must be of type `int`, got: " +
				right.Types[0].String(),
		)
	}
	return nil
}
