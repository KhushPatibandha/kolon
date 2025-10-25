package evaluator

import (
	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/environment"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/object"
)

// ------------------------------------------------------------------------------------------------------------------
// Statements
// ------------------------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------------------------
// ExpressionStatement
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalExpressionStatement(es *ast.ExpressionStatement) (*object.EvalResult, error) {
	_, err := e.Evaluate(es.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Body
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalStmts(stmts []ast.Statement) (*object.EvalResult, error) {
	defer e.stack.Pop()
	for _, stmt := range stmts {
		r, err := e.Evaluate(stmt)
		if err != nil {
			return nil, err
		}
		if r.Signal != object.SIGNAL_NONE {
			return r, nil
		}
	}
	return nil, nil
}

// ------------------------------------------------------------------------------------------------------------------
// VarAndConst
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalVarConst(vc *ast.VarAndConst,
	injectValue bool,
	val object.Object,
) (*object.EvalResult, error) {
	var r object.Object

	if injectValue {
		r = val
	} else {
		res, err := e.Evaluate(vc.Value)
		if err != nil {
			return nil, err
		}
		r = res.Value
	}

	sym := &environment.Symbol{
		Ident:       vc.Name,
		ValueObject: r,
		Type:        nil,
		Func:        nil,
		Env:         nil,
	}
	switch vc.Token.Kind {
	case lexer.VAR:
		sym.IdentType = environment.VAR
	case lexer.CONST:
		sym.IdentType = environment.CONST
	}
	e.stack.Top().Set(sym)

	return nil, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Multi-Assignment
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalMultiAssign(ma *ast.MultiAssignment) (*object.EvalResult, error) {
	if !ma.SingleFunctionCall {
		for _, ele := range ma.Objects {
			var err error
			switch ele := ele.(type) {
			case *ast.VarAndConst:
				_, err = e.evalVarConst(ele, false, nil)
			case *ast.ExpressionStatement:
				_, err = e.evalAssignment(ele.Expression.(*ast.Assignment), false, nil)
			}
			if err != nil {
				return nil, err
			}
		}
	} else {
		var call *ast.CallExpression

		switch v := ma.Objects[0].(type) {
		case *ast.VarAndConst:
			call = v.Value.(*ast.CallExpression)
		case *ast.ExpressionStatement:
			call = v.Expression.(*ast.Assignment).Right.(*ast.CallExpression)
		}

		r, err := e.evalCall(call)
		if err != nil {
			return nil, err
		}

		rList := r.Value.(*object.MultiValue).Values

		for i, ele := range ma.Objects {
			var err error
			switch ele := ele.(type) {
			case *ast.VarAndConst:
				_, err = e.evalVarConst(ele, true, rList[i])
			case *ast.ExpressionStatement:
				_, err = e.evalAssignment(ele.Expression.(*ast.Assignment), true, rList[i])
			}
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Function
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalFunc(f *ast.Function) (*object.EvalResult, error) {
	funcLocalEnv := environment.BootstrapFuncEnv(f, e.env)
	e.env.Set(&environment.Symbol{
		IdentType: environment.FUNCTION,
		Ident:     f.Name,
		Func: &environment.FuncInfo{
			Function: f,
			Builtin:  false,
		},
		Env:  funcLocalEnv,
		Type: nil,
	})
	if f.Name.Value == "main" {
		if sym, ok := e.env.GetFunc("main"); ok {
			e.stack.Push(sym.Env)
			return e.evalStmts(f.Body.Statements)
		}
	}
	return nil, nil
}

// ------------------------------------------------------------------------------------------------------------------
// If
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalIf(i *ast.If) (*object.EvalResult, error) {
	condition, err := e.Evaluate(i.Condition)
	if err != nil {
		return nil, err
	}

	if condition.Value == TRUE.Value {
		localEnv := environment.NewEnclosedEnvironment(e.stack.Top())
		e.stack.Push(localEnv)
		return e.evalStmts(i.Body.Statements)
	} else if i.MultiConditionals != nil {
		for _, mc := range i.MultiConditionals {
			r, success, err := e.evalElseIf(mc)
			if err != nil {
				return nil, err
			}
			if success {
				return r, nil
			}
		}
	}
	if i.Alternate != nil {
		localEnv := environment.NewEnclosedEnvironment(e.stack.Top())
		e.stack.Push(localEnv)
		return e.evalStmts(i.Alternate.Body.Statements)
	}

	return nil, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Else If
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalElseIf(ei *ast.ElseIf) (*object.EvalResult, bool, error) {
	condition, err := e.Evaluate(ei.Condition)
	if err != nil {
		return nil, false, err
	}
	if condition.Value == TRUE.Value {
		localEnv := environment.NewEnclosedEnvironment(e.stack.Top())
		e.stack.Push(localEnv)
		r, err := e.evalStmts(ei.Body.Statements)
		return r, true, err
	}
	return nil, false, nil
}

// ------------------------------------------------------------------------------------------------------------------
// ForLoop
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalForLoop(f *ast.ForLoop) (*object.EvalResult, error) {
	localEnv := environment.NewEnclosedEnvironment(e.stack.Top())
	e.stack.Push(localEnv)

	_, err := e.Evaluate(f.Left)
	if err != nil {
		return nil, err
	}

	condition, err := e.Evaluate(f.Middle)
	if err != nil {
		return nil, err
	}

	for condition.Value == TRUE.Value {
		bodyLocalEnv := environment.NewEnclosedEnvironment(e.stack.Top())
		e.stack.Push(bodyLocalEnv)
		r, err := e.evalStmts(f.Body.Statements)
		if err != nil {
			return nil, err
		}

		if r.Signal == object.SIGNAL_BREAK {
			break
		} else if r.Signal == object.SIGNAL_RETURN {
			return r, nil
		}

		_, err = e.Evaluate(f.Right)
		if err != nil {
			return nil, err
		}

		condition, err = e.Evaluate(f.Middle)
		if err != nil {
			return nil, err
		}
		if condition.Value == FALSE.Value {
			break
		}
	}
	e.stack.Pop()

	return nil, nil
}

// ------------------------------------------------------------------------------------------------------------------
// WhileLoop
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalWhileLoop(w *ast.WhileLoop) (*object.EvalResult, error) {
	condition, err := e.Evaluate(w.Condition)
	if err != nil {
		return nil, err
	}

	for condition.Value == TRUE.Value {
		localEnv := environment.NewEnclosedEnvironment(e.stack.Top())
		e.stack.Push(localEnv)
		r, err := e.evalStmts(w.Body.Statements)
		if err != nil {
			return nil, err
		}

		if r.Signal == object.SIGNAL_BREAK {
			break
		} else if r.Signal == object.SIGNAL_RETURN {
			return r, nil
		}

		condition, err = e.Evaluate(w.Condition)
		if err != nil {
			return nil, err
		}
		if condition.Value == FALSE.Value {
			break
		}
	}

	return nil, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Return
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalReturn(r *ast.Return) (*object.EvalResult, error) {
	if r.Value == nil {
		return &object.EvalResult{Value: nil, Signal: object.SIGNAL_RETURN}, nil
	}

	if len(r.Value) == 1 {
		r, err := e.Evaluate(r.Value[0])
		if err != nil {
			return nil, err
		}
		return &object.EvalResult{Value: r.Value, Signal: object.SIGNAL_RETURN}, nil
	}

	var val []object.Object

	for _, ex := range r.Value {
		r, err := e.Evaluate(ex)
		if err != nil {
			return nil, err
		}
		val = append(val, r.Value)
	}

	return &object.EvalResult{
		Value:  &object.MultiValue{Values: val},
		Signal: object.SIGNAL_RETURN,
	}, nil
}
