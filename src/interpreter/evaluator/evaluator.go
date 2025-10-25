package evaluator

import (
	"fmt"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/environment"
	"github.com/KhushPatibandha/Kolon/src/object"
)

var (
	TRUE  = &object.Bool{Value: true}
	FALSE = &object.Bool{Value: false}
)

// ------------------------------------------------------------------------------------------------------------------
// Evaluator
// ------------------------------------------------------------------------------------------------------------------
func Evaluate(node ast.Node, env *environment.Environment) (*object.EvalResult, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalStmts(node.Statements, env)
	case *ast.Identifier:
		return nil, nil
	case *ast.Integer:
		return nil, nil
	case *ast.Float:
		return nil, nil
	case *ast.Bool:
		return nil, nil
	case *ast.String:
		return nil, nil
	case *ast.Char:
		return nil, nil
	case *ast.HashMap:
		return nil, nil
	case *ast.Array:
		return nil, nil
	case *ast.Prefix:
		return nil, nil
	case *ast.Infix:
		return nil, nil
	case *ast.Postfix:
		return nil, nil
	case *ast.Assignment:
		return nil, nil
	case *ast.IndexExpression:
		return nil, nil
	case *ast.CallExpression:
		return nil, nil
	case *ast.ExpressionStatement:
		return nil, nil
	case *ast.Body:
		return evalStmts(node.Statements, env)
	case *ast.Function:
		if node.Name.Value == "main" {
			if sym, ok := env.GetFunc("main"); ok {
				return evalStmts(node.Body.Statements, sym.Env)
			}
		}
		return nil, nil
	case *ast.VarAndConst:
		return nil, nil
	case *ast.MultiAssignment:
		return nil, nil
	case *ast.Return:
		return nil, nil
	case *ast.If:
		return nil, nil
	case *ast.ForLoop:
		return nil, nil
	case *ast.WhileLoop:
		return nil, nil
	case *ast.Continue:
		return nil, nil
	case *ast.Break:
		return nil, nil
	default:
		return nil, fmt.Errorf("no eval function for given node type, got: %T", node)
	}
}
