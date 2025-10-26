package evaluator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/object"
)

// ------------------------------------------------------------------------------------------------------------------
// Expressions
// ------------------------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------------------------
// Identifier
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalIdentifier(i *ast.Identifier) (*object.EvalResult, error) {
	if sym, ok := e.stack.Top().GetVar(i.Value); ok {
		return &object.EvalResult{
			Value:  sym.ValueObject,
			Signal: object.SIGNAL_NONE,
		}, nil
	}
	return nil, fmt.Errorf("identifier not found: %s", i.Value)
}

// ------------------------------------------------------------------------------------------------------------------
// Integer
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalInteger(i *ast.Integer) (*object.EvalResult, error) {
	return &object.EvalResult{
		Value:  &object.Integer{Value: i.Value},
		Signal: object.SIGNAL_NONE,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Float
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalFloat(f *ast.Float) (*object.EvalResult, error) {
	return &object.EvalResult{
		Value:  &object.Float{Value: f.Value},
		Signal: object.SIGNAL_NONE,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Boolean
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalBoolean(b *ast.Bool) (*object.EvalResult, error) {
	if b.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

// ------------------------------------------------------------------------------------------------------------------
// String
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalString(s *ast.String) (*object.EvalResult, error) {
	return &object.EvalResult{
		Value:  &object.String{Value: s.Value},
		Signal: object.SIGNAL_NONE,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Char
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalChar(c *ast.Char) (*object.EvalResult, error) {
	return &object.EvalResult{
		Value:  &object.Char{Value: c.Value},
		Signal: object.SIGNAL_NONE,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// HashMap
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalHashMap(h *ast.HashMap) (*object.EvalResult, error) {
	pairs := make(map[object.HashKey]object.HashPair)

	for k, v := range h.Pairs {
		key, err := e.Evaluate(k)
		if err != nil {
			return nil, err
		}

		hashKey, ok := key.Value.(object.Hashable)
		if !ok {
			return nil, errors.New("unusable as hash key: " + string(key.Value.Type()))
		}

		value, err := e.Evaluate(v)
		if err != nil {
			return nil, err
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key.Value, Value: value.Value}
	}
	return &object.EvalResult{
		Value:  &object.HashMap{Pairs: pairs},
		Signal: object.SIGNAL_NONE,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Array
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalArray(a *ast.Array) (*object.EvalResult, error) {
	var res []object.Object
	for _, ele := range a.Values {
		r, err := e.Evaluate(ele)
		if err != nil {
			return nil, err
		}
		res = append(res, r.Value)
	}
	return &object.EvalResult{
		Value:  &object.Array{Elements: res},
		Signal: object.SIGNAL_NONE,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Prefix
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalPrefix(p *ast.Prefix) (*object.EvalResult, error) {
	right, err := e.Evaluate(p.Right)
	if err != nil {
		return nil, err
	}
	if p.Operator == "!" {
		return e.evalPrefixBang(right.Value)
	} else {
		return e.evalPrefixMinus(right.Value)
	}
}

func (e *Evaluator) evalPrefixBang(right object.Object) (*object.EvalResult, error) {
	if right == TRUE.Value {
		return FALSE, nil
	}
	return TRUE, nil
}

func (e *Evaluator) evalPrefixMinus(right object.Object) (*object.EvalResult, error) {
	if right.Type() == object.FLOAT_OBJ {
		return &object.EvalResult{
			Value:  &object.Float{Value: -right.(*object.Float).Value},
			Signal: object.SIGNAL_NONE,
		}, nil
	}
	return &object.EvalResult{
		Value:  &object.Integer{Value: -right.(*object.Integer).Value},
		Signal: object.SIGNAL_NONE,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Infix
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalInfix(i *ast.Infix) (*object.EvalResult, error) {
	left, err := e.Evaluate(i.Left)
	if err != nil {
		return nil, err
	}
	right, err := e.Evaluate(i.Right)
	if err != nil {
		return nil, err
	}
	switch {
	case left.Value.Type() == object.INTEGER_OBJ && right.Value.Type() == object.INTEGER_OBJ:
		return e.evalInfixInteger(i.Operator, left.Value, right.Value)
	case left.Value.Type() == object.FLOAT_OBJ && right.Value.Type() == object.FLOAT_OBJ:
		return e.evalInfixFloat(i.Operator, left.Value, right.Value)
	case left.Value.Type() == object.BOOLEAN_OBJ && right.Value.Type() == object.BOOLEAN_OBJ:
		return e.evalInfixBool(i.Operator, left.Value, right.Value)
	case left.Value.Type() == object.STRING_OBJ && right.Value.Type() == object.STRING_OBJ:
		return e.evalInfixString(i.Operator, left.Value, right.Value)
	case left.Value.Type() == object.CHAR_OBJ && right.Value.Type() == object.CHAR_OBJ:
		return e.evalInfixChar(i.Operator, left.Value, right.Value)
	case (left.Value.Type() == object.INTEGER_OBJ && right.Value.Type() == object.FLOAT_OBJ) ||
		(left.Value.Type() == object.FLOAT_OBJ && right.Value.Type() == object.INTEGER_OBJ):
		l := 0.0
		r := 0.0

		if left.Value.Type() == object.INTEGER_OBJ {
			l = float64(left.Value.(*object.Integer).Value)
			r = right.Value.(*object.Float).Value
		} else {
			l = left.Value.(*object.Float).Value
			r = float64(right.Value.(*object.Integer).Value)
		}

		return e.evalInfixFloat(i.Operator, &object.Float{Value: l}, &object.Float{Value: r})
	default:
		return e.evalInfixArray(i.Operator, left.Value, right.Value)
	}
}

func (e *Evaluator) evalInfixChar(operator string,
	left, right object.Object,
) (*object.EvalResult, error) {
	l := left.(*object.Char).Value
	r := right.(*object.Char).Value
	l = l[1 : len(l)-1]
	r = r[1 : len(r)-1]
	switch operator {
	case "+":
		return &object.EvalResult{
			Value:  &object.String{Value: "\"" + l + r + "\""},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "==":
		if l == r {
			return TRUE, nil
		}
		return FALSE, nil
	default:
		if l != r {
			return TRUE, nil
		}
		return FALSE, nil
	}
}

func (e *Evaluator) evalInfixString(operator string,
	left, right object.Object,
) (*object.EvalResult, error) {
	l := left.(*object.String).Value
	r := right.(*object.String).Value
	l = l[1 : len(l)-1]
	r = r[1 : len(r)-1]
	switch operator {
	case "+":
		return &object.EvalResult{
			Value:  &object.String{Value: "\"" + l + r + "\""},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "==":
		if l == r {
			return TRUE, nil
		}
		return FALSE, nil
	default:
		if l != r {
			return TRUE, nil
		}
		return FALSE, nil
	}
}

func (e *Evaluator) evalInfixBool(operator string,
	left, right object.Object,
) (*object.EvalResult, error) {
	l := left.(*object.Bool).Value
	r := right.(*object.Bool).Value
	switch operator {
	case "&&":
		if l && r {
			return TRUE, nil
		}
		return FALSE, nil
	case "||":
		if l || r {
			return TRUE, nil
		}
		return FALSE, nil
	case "==":
		if l == r {
			return TRUE, nil
		}
		return FALSE, nil
	default:
		if l != r {
			return TRUE, nil
		}
		return FALSE, nil
	}
}

func (e *Evaluator) evalInfixFloat(operator string,
	left, right object.Object,
) (*object.EvalResult, error) {
	l := left.(*object.Float).Value
	r := right.(*object.Float).Value
	switch operator {
	case "+":
		return &object.EvalResult{
			Value:  &object.Float{Value: l + r},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "-":
		return &object.EvalResult{
			Value:  &object.Float{Value: l - r},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "*":
		return &object.EvalResult{
			Value:  &object.Float{Value: l * r},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "/":
		if r == 0.0 {
			return nil, errors.New("float division by zero")
		}
		return &object.EvalResult{
			Value:  &object.Float{Value: l / r},
			Signal: object.SIGNAL_NONE,
		}, nil
	case ">":
		if l > r {
			return TRUE, nil
		}
		return FALSE, nil
	case "<":
		if l < r {
			return TRUE, nil
		}
		return FALSE, nil
	case "<=":
		if l <= r {
			return TRUE, nil
		}
		return FALSE, nil
	case ">=":
		if l >= r {
			return TRUE, nil
		}
		return FALSE, nil
	case "==":
		if l == r {
			return TRUE, nil
		}
		return FALSE, nil
	default:
		if l != r {
			return TRUE, nil
		}
		return FALSE, nil
	}
}

func (e *Evaluator) evalInfixInteger(operator string,
	left, right object.Object,
) (*object.EvalResult, error) {
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.EvalResult{
			Value:  &object.Integer{Value: l + r},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "-":
		return &object.EvalResult{
			Value:  &object.Integer{Value: l - r},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "*":
		return &object.EvalResult{
			Value:  &object.Integer{Value: l * r},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "/":
		if r == 0 {
			return nil, errors.New("integer division by zero")
		}
		return &object.EvalResult{
			Value:  &object.Integer{Value: l / r},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "%":
		if r == 0 {
			return nil, errors.New("modulo by zero")
		}
		return &object.EvalResult{
			Value:  &object.Integer{Value: l % r},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "&":
		return &object.EvalResult{
			Value:  &object.Integer{Value: l & r},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "|":
		return &object.EvalResult{
			Value:  &object.Integer{Value: l | r},
			Signal: object.SIGNAL_NONE,
		}, nil
	case ">":
		if l > r {
			return TRUE, nil
		}
		return FALSE, nil
	case "<":
		if l < r {
			return TRUE, nil
		}
		return FALSE, nil
	case "<=":
		if l <= r {
			return TRUE, nil
		}
		return FALSE, nil
	case ">=":
		if l >= r {
			return TRUE, nil
		}
		return FALSE, nil
	case "==":
		if l == r {
			return TRUE, nil
		}
		return FALSE, nil
	default:
		if l != r {
			return TRUE, nil
		}
		return FALSE, nil
	}
}

func (e *Evaluator) evalInfixArray(operator string,
	left, right object.Object,
) (*object.EvalResult, error) {
	l := left.(*object.Array)
	r := right.(*object.Array)
	switch operator {
	case "+":
		a := append(append([]object.Object{}, l.Elements...), r.Elements...)
		return &object.EvalResult{
			Value:  &object.Array{Elements: a},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "==":
		// TODO: write test for this
		if l == r {
			return TRUE, nil
		}
		return FALSE, nil
	default:
		if l == r {
			return FALSE, nil
		}
		return TRUE, nil
	}
}

// ------------------------------------------------------------------------------------------------------------------
// Postfix
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalPostfix(p *ast.Postfix) (*object.EvalResult, error) {
	left, err := e.Evaluate(p.Left)
	if err != nil {
		return nil, err
	}
	if left.Value.Type() == object.FLOAT_OBJ {
		return e.evalPostfixFloat(left.Value, p.Operator)
	} else {
		return e.evalPostfixInteger(left.Value, p.Operator)
	}
}

func (e *Evaluator) evalPostfixFloat(left object.Object, operator string) (*object.EvalResult, error) {
	l := left.(*object.Float).Value
	if operator == "++" {
		return &object.EvalResult{
			Value:  &object.Float{Value: l + 1},
			Signal: object.SIGNAL_NONE,
		}, nil
	}
	return &object.EvalResult{
		Value:  &object.Float{Value: l - 1},
		Signal: object.SIGNAL_NONE,
	}, nil
}

func (e *Evaluator) evalPostfixInteger(left object.Object, operator string) (*object.EvalResult, error) {
	l := left.(*object.Integer).Value
	if operator == "++" {
		return &object.EvalResult{
			Value:  &object.Integer{Value: l + 1},
			Signal: object.SIGNAL_NONE,
		}, nil
	}
	return &object.EvalResult{
		Value:  &object.Integer{Value: l - 1},
		Signal: object.SIGNAL_NONE,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Assignment
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalAssignment(a *ast.Assignment,
	injectValue bool,
	val object.Object,
) (*object.EvalResult, error) {
	switch a.Operator {
	case "=":
		return e.evalAssignmentEqual(a, injectValue, val)
	default:
		return e.evalAssignmentSymbol(a)
	}
}

func (e *Evaluator) evalAssignmentSymbol(a *ast.Assignment) (*object.EvalResult, error) {
	in := &ast.Infix{
		Left:     a.Left,
		Operator: strings.TrimSuffix(a.Operator, "="),
		Right:    a.Right,
	}
	right, err := e.evalInfix(in)
	if err != nil {
		return nil, err
	}
	return e.evalAssignmentEqual(a, true, right.Value)
}

func (e *Evaluator) evalAssignmentEqual(a *ast.Assignment,
	injectValue bool,
	right object.Object,
) (*object.EvalResult, error) {
	var r object.Object

	if injectValue {
		r = right
	} else {
		res, err := e.Evaluate(a.Right)
		if err != nil {
			return nil, err
		}
		r = res.Value
	}

	e.stack.Top().SetValue(a.Left.Value, r)
	return &object.EvalResult{
		Value:  r,
		Signal: object.SIGNAL_NONE,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// CallExpression
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalCall(c *ast.CallExpression) (*object.EvalResult, error) {
	args, err := e.evalCallArgs(c)
	if err != nil {
		return nil, err
	}

	sym, _ := e.env.GetFunc(c.Name.Value)
	if sym.Func.Builtin {
		return e.evalBuiltin(c.Name.Value, args)
	}

	localEnv := sym.Env
	e.stack.Push(localEnv)
	r, err := e.evalStmts(sym.Func.Function.Body.Statements)

	return r, err
}

func (e *Evaluator) evalCallArgs(c *ast.CallExpression) ([]object.Object, error) {
	var res []object.Object
	if c.Args != nil {
		for _, ex := range c.Args {
			r, err := e.Evaluate(ex)
			if err != nil {
				return nil, err
			}
			res = append(res, r.Value)
		}
	}
	return res, nil
}

// ------------------------------------------------------------------------------------------------------------------
// IndexExpression
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalIndex(i *ast.IndexExpression) (*object.EvalResult, error) {
	left, err := e.Evaluate(i.Left)
	if err != nil {
		return nil, err
	}
	index, err := e.Evaluate(i.Index)
	if err != nil {
		return nil, err
	}
	if left.Value.Type() == object.ARRAY_OBJ {
		return e.evalIndexArray(left.Value, index.Value)
	} else if left.Value.Type() == object.STRING_OBJ {
		return e.evalIndexString(left.Value, index.Value)
	} else {
		return e.evalIndexHashMap(left.Value, index.Value)
	}
}

func (e *Evaluator) evalIndexArray(left, index object.Object) (*object.EvalResult, error) {
	a := left.(*object.Array)
	i := index.(*object.Integer).Value
	maxIdx := int64(len(a.Elements) - 1)
	if i < 0 || i > maxIdx {
		return nil,
			errors.New(
				"index out of range, index: " +
					strconv.FormatInt(i, 10) + ", max index: " +
					strconv.FormatInt(maxIdx, 10) + ", min index: 0",
			)
	}
	return &object.EvalResult{
		Value:  a.Elements[i],
		Signal: object.SIGNAL_NONE,
	}, nil
}

func (e *Evaluator) evalIndexString(left, index object.Object) (*object.EvalResult, error) {
	s := left.(*object.String).Value
	s = s[1 : len(s)-1]

	i := index.(*object.Integer).Value
	maxIdx := int64(len(s) - 1)
	if i < 0 || i > maxIdx {
		return nil,
			errors.New(
				"index out of range, index: " +
					strconv.FormatInt(i, 10) + ", max index: " +
					strconv.FormatInt(maxIdx, 10) + ", min index: 0",
			)
	}
	return &object.EvalResult{
		Value:  &object.Char{Value: "'" + string([]rune(s)[i]) + "'"},
		Signal: object.SIGNAL_NONE,
	}, nil
}

func (e *Evaluator) evalIndexHashMap(left, index object.Object) (*object.EvalResult, error) {
	h := left.(*object.HashMap)
	k, ok := index.(object.Hashable)
	if !ok {
		return nil, errors.New("unusable as hash key: " + string(index.Type()))
	}
	pair, ok := h.Pairs[k.HashKey()]
	if !ok {
		return nil, errors.New("key not found: " + index.Inspect())
	}
	return &object.EvalResult{
		Value:  pair.Value,
		Signal: object.SIGNAL_NONE,
	}, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Builtin
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalBuiltin(name string, args []object.Object) (*object.EvalResult, error) {
	return nil, nil
}

// TODO: add equals() builtin function
