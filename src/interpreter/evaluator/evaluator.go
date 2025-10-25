package evaluator

import (
	"fmt"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/environment"
	"github.com/KhushPatibandha/Kolon/src/object"
)

var (
	TRUE = &object.EvalResult{
		Value:  &object.Bool{Value: true},
		Signal: object.SIGNAL_NONE,
	}
	FALSE = &object.EvalResult{
		Value:  &object.Bool{Value: false},
		Signal: object.SIGNAL_NONE,
	}
	CONTINUE = &object.EvalResult{
		Value:  nil,
		Signal: object.SIGNAL_CONTINUE,
	}
	BREAK = &object.EvalResult{
		Value:  nil,
		Signal: object.SIGNAL_BREAK,
	}
)

type Evaluator struct {
	inTesting bool
	env       *environment.Environment
	stack     *environment.Stack
}

// ------------------------------------------------------------------------------------------------------------------
// Evaluator
// ------------------------------------------------------------------------------------------------------------------
func New(inTesting bool) *Evaluator {
	e := &Evaluator{
		inTesting: inTesting,
		env:       environment.NewEnvironment(),
		stack:     environment.NewStack(),
	}

	e.stack.Push(e.env)
	environment.LoadBuiltins(e.env)

	return e
}

func (e *Evaluator) Evaluate(node ast.Node) (*object.EvalResult, error) {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalStmts(node.Statements)
	case *ast.Identifier:
		return e.evalIdentifier(node)
	case *ast.Integer:
		return e.evalInteger(node)
	case *ast.Float:
		return e.evalFloat(node)
	case *ast.Bool:
		return e.evalBoolean(node)
	case *ast.String:
		return e.evalString(node)
	case *ast.Char:
		return e.evalChar(node)
	case *ast.HashMap:
		return e.evalHashMap(node)
	case *ast.Array:
		return e.evalArray(node)
	case *ast.Prefix:
		return e.evalPrefix(node)
	case *ast.Infix:
		return e.evalInfix(node)
	case *ast.Postfix:
		return e.evalPostfix(node)
	case *ast.Assignment:
		return e.evalAssignment(node, false, nil)
	case *ast.IndexExpression:
		return e.evalIndex(node)
	case *ast.CallExpression:
		return e.evalCall(node)
	case *ast.ExpressionStatement:
		return e.evalExpressionStatement(node)
	case *ast.Body:
		return e.evalStmts(node.Statements)
	case *ast.Function:
		return e.evalFunc(node)
	case *ast.VarAndConst:
		return e.evalVarConst(node, false, nil)
	case *ast.MultiAssignment:
		return e.evalMultiAssign(node)
	case *ast.Return:
		return e.evalReturn(node)
	case *ast.If:
		return e.evalIf(node)
	case *ast.ForLoop:
		return e.evalForLoop(node)
	case *ast.WhileLoop:
		return e.evalWhileLoop(node)
	case *ast.Continue:
		return CONTINUE, nil
	case *ast.Break:
		return BREAK, nil
	default:
		return nil, fmt.Errorf("no eval function for given node type, got: %T", node)
	}
}
