package evaluator

import (
	"errors"
	"fmt"
	"strconv"

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
	return nil, nil
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
	return nil, nil
}

// ------------------------------------------------------------------------------------------------------------------
// CallExpression
// ------------------------------------------------------------------------------------------------------------------
func (e *Evaluator) evalCall(c *ast.CallExpression) (*object.EvalResult, error) {
	return nil, nil
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
