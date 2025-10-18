package parser

import (
	"fmt"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/environment"
	"github.com/KhushPatibandha/Kolon/src/kType"
)

// ------------------------------------------------------------------------------------------------------------------
// TypeChecker
// ------------------------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------------------------
// Expressions
// ------------------------------------------------------------------------------------------------------------------
func typeCheckExp(exp ast.Expression, env *environment.Environment) (*ktype.TypeCheckResult, error) {
	switch exp := exp.(type) {
	case *ast.Identifier:
		return typeCheckIdent(exp, env)
	case *ast.Integer:
		return typeCheckInteger()
	case *ast.Float:
		return typeCheckFloat()
	case *ast.String:
		return typeCheckString()
	case *ast.Char:
		return typeCheckChar()
	case *ast.Bool:
		return typeCheckBool()
	case *ast.HashMap:
		return typeCheckHashMap(exp, env)
	case *ast.Array:
		return typeCheckArray(exp, env)
	case *ast.Prefix:
		return typeCheckPrefix(exp, env)
	case *ast.Postfix:
		return typeCheckPostfix(exp, env)
	case *ast.Infix:
		return typeCheckInfix(exp, env)
	case *ast.Assignment:
		return typeCheckAssignment(exp, env)
	case *ast.IndexExpression:
		return typeCheckIndexExp(exp, env)
	case *ast.CallExpression:
		return typeCheckCallExp(exp, env)
	default:
		return nil, fmt.Errorf("unknown expression type, got: %T", exp)
	}
}
