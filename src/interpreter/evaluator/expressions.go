package evaluator

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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
	switch name {
	case "len":
		var r object.Object
		switch arg := args[0].(type) {
		case *object.String:
			r = &object.Integer{Value: int64(len(arg.Value) - 2)}
		case *object.Array:
			r = &object.Integer{Value: int64(len(arg.Elements))}
		case *object.HashMap:
			r = &object.Integer{Value: int64(len(arg.Pairs))}
		}
		return &object.EvalResult{
			Value:  r,
			Signal: object.SIGNAL_NONE,
		}, nil
	case "toString":
		var r object.Object
		switch arg := args[0].(type) {
		case *object.Integer:
			s := strconv.FormatInt(arg.Value, 10)
			s = "\"" + s + "\""
			r = &object.String{Value: s}
		case *object.Float:
			s := strconv.FormatFloat(arg.Value, 'f', -1, 64)
			s = "\"" + s + "\""
			r = &object.String{Value: s}
		case *object.Bool:
			s := strconv.FormatBool(arg.Value)
			s = "\"" + s + "\""
			r = &object.String{Value: s}
		case *object.Char:
			s := arg.Value[1 : len(arg.Value)-1]
			s = "\"" + s + "\""
			r = &object.String{Value: s}
		case *object.String:
			r = arg
		case *object.Array:
			r = &object.String{Value: "\"" + arg.Inspect() + "\""}
		case *object.HashMap:
			r = &object.String{Value: "\"" + arg.Inspect() + "\""}
		}
		return &object.EvalResult{
			Value:  r,
			Signal: object.SIGNAL_NONE,
		}, nil
	case "toInt":
		var r object.Object
		switch arg := args[0].(type) {
		case *object.Integer:
			r = arg
		case *object.Float:
			r = &object.Integer{Value: int64(arg.Value)}
		case *object.Char:
			s := arg.Value[1 : len(arg.Value)-1]
			r = &object.Integer{Value: int64(s[0])}
		case *object.String:
			s := arg.Value[1 : len(arg.Value)-1]
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, errors.New("Error converting string to int, can't convert: " + s)
			}
			r = &object.Integer{Value: i}
		}
		return &object.EvalResult{
			Value:  r,
			Signal: object.SIGNAL_NONE,
		}, nil
	case "toFloat":
		var r object.Object
		switch arg := args[0].(type) {
		case *object.Integer:
			r = &object.Float{Value: float64(arg.Value)}
		case *object.Float:
			r = arg
		case *object.String:
			s := arg.Value[1 : len(arg.Value)-1]
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, errors.New("Error converting string to float, can't convert: " + s)
			}
			r = &object.Float{Value: f}
		}
		return &object.EvalResult{
			Value:  r,
			Signal: object.SIGNAL_NONE,
		}, nil
	case "print":
		switch arg := args[0].(type) {
		case *object.String:
			fmt.Print(arg.Value[1 : len(arg.Value)-1])
		case *object.Char:
			fmt.Print(arg.Value[1 : len(arg.Value)-1])
		case *object.Integer:
			fmt.Print(strconv.FormatInt(arg.Value, 10))
		case *object.Float:
			fmt.Print(strconv.FormatFloat(arg.Value, 'f', -1, 64))
		case *object.Bool:
			fmt.Print(strconv.FormatBool(arg.Value))
		default:
			fmt.Print(arg.Inspect())
		}
		return &object.EvalResult{
			Value:  nil,
			Signal: object.SIGNAL_NONE,
		}, nil
	case "println":
		if len(args) == 0 {
			fmt.Println()
			return nil, nil
		}
		switch arg := args[0].(type) {
		case *object.String:
			fmt.Println(arg.Value[1 : len(arg.Value)-1])
		case *object.Char:
			fmt.Println(arg.Value[1 : len(arg.Value)-1])
		case *object.Integer:
			fmt.Println(strconv.FormatInt(arg.Value, 10))
		case *object.Float:
			fmt.Println(strconv.FormatFloat(arg.Value, 'f', -1, 64))
		case *object.Bool:
			fmt.Println(strconv.FormatBool(arg.Value))
		default:
			fmt.Println(arg.Inspect())
		}
		return &object.EvalResult{
			Value:  nil,
			Signal: object.SIGNAL_NONE,
		}, nil
	case "scan":
		if len(args) != 0 {
			strToPrint := args[0].(*object.String).
				Value[1 : len(args[0].(*object.String).Value)-1]
			if len(args) == 2 && args[1].Inspect() == "true" {
				fmt.Println(strToPrint)
			} else {
				fmt.Print(strToPrint)
			}
		}
		reader := bufio.NewReader(os.Stdin)
		var input []string
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, errors.New("error reading input: " + err.Error())
			}
			line = strings.TrimSpace(line)
			line = strings.TrimSuffix(line, "\n")
			if line == "" {
				break
			}
			input = append(input, line)
		}
		res := strings.Join(input, " ")
		return &object.EvalResult{
			Value:  &object.String{Value: "\"" + res + "\""},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "scanln":
		if len(args) != 0 {
			strToPrint := args[0].(*object.String).
				Value[1 : len(args[0].(*object.String).Value)-1]
			if len(args) == 2 && args[1].Inspect() == "true" {
				fmt.Println(strToPrint)
			} else {
				fmt.Print(strToPrint)
			}
		}
		reader := bufio.NewReader(os.Stdin)
		var input string
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, errors.New("error reading input: " + err.Error())
		}
		input = strings.TrimSpace(input)
		input = strings.TrimSuffix(input, "\n")
		return &object.EvalResult{
			Value:  &object.String{Value: "\"" + input + "\""},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "getIndex":
		switch arg := args[0].(type) {
		case *object.Array:
			t := args[1].Inspect()
			for i, ele := range arg.Elements {
				if ele.Inspect() == t {
					return &object.EvalResult{
						Value:  &object.Integer{Value: int64(i)},
						Signal: object.SIGNAL_NONE,
					}, nil
				}
			}
		}
		return &object.EvalResult{
			Value:  &object.Integer{Value: -1},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "keys":
		var keys []object.Object
		for _, pair := range args[0].(*object.HashMap).Pairs {
			keys = append(keys, pair.Key)
		}
		return &object.EvalResult{
			Value:  &object.Array{Elements: keys},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "values":
		var values []object.Object
		for _, pair := range args[0].(*object.HashMap).Pairs {
			values = append(values, pair.Value)
		}
		return &object.EvalResult{
			Value:  &object.Array{Elements: values},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "containsKey":
		h := args[0].(*object.HashMap)
		k, ok := args[1].(object.Hashable)
		if !ok {
			return nil, errors.New("unusable as hash key: " + string(args[1].Type()))
		}
		_, ok = h.Pairs[k.HashKey()]
		if ok {
			return TRUE, nil
		}
		return FALSE, nil
	case "typeOf":
		return &object.EvalResult{
			Value:  &object.String{Value: "\"" + getType(args[0]) + "\""},
			Signal: object.SIGNAL_NONE,
		}, nil
	case "push":
		switch arg := args[0].(type) {
		case *object.Array:
			arg.Elements = append(arg.Elements, args[1])
			return &object.EvalResult{
				Value:  arg,
				Signal: object.SIGNAL_NONE,
			}, nil
		default:
			hashKey, ok := args[1].(object.Hashable)
			if !ok {
				return nil, errors.New("unusable as hash key: " + string(args[1].Type()))
			}
			arg.(*object.HashMap).Pairs[hashKey.HashKey()] = object.HashPair{
				Key:   args[1],
				Value: args[2],
			}
			return &object.EvalResult{
				Value:  arg,
				Signal: object.SIGNAL_NONE,
			}, nil
		}
	case "pop":
		a := args[0].(*object.Array)
		var popped object.Object
		if len(args) == 1 {
			if len(a.Elements) == 0 {
				return nil, errors.New("array is empty, can't pop any elements")
			}
			popped = a.Elements[len(a.Elements)-1]
			a.Elements = a.Elements[:len(a.Elements)-1]
		} else {
			idx := args[1].(*object.Integer).Value
			if idx < 0 || idx >= int64(len(a.Elements)) {
				return nil,
					errors.New("index out of range, can't pop element at index: " +
						strconv.FormatInt(idx, 10),
					)
			}
			popped = a.Elements[idx]
			a.Elements = append(a.Elements[:idx], a.Elements[idx+1:]...)
		}
		return &object.EvalResult{
			Value:  popped,
			Signal: object.SIGNAL_NONE,
		}, nil
	case "insert":
		a := args[0].(*object.Array)
		idx := args[1].(*object.Integer).Value
		if idx < 0 || idx > int64(len(a.Elements)) {
			return nil,
				errors.New("index out of range, can't insert element at index: " +
					strconv.FormatInt(idx, 10),
				)
		}
		value := args[2]
		a.Elements = append(a.Elements[:idx], append([]object.Object{value}, a.Elements[idx:]...)...)
		return &object.EvalResult{
			Value:  a,
			Signal: object.SIGNAL_NONE,
		}, nil
	case "remove":
		switch arg := args[0].(type) {
		case *object.Array:
			eleToRemove := args[1].Inspect()
			for i, ele := range arg.Elements {
				if ele.Inspect() == eleToRemove {
					arg.Elements = append(arg.Elements[:i], arg.Elements[i+1:]...)
					break
				}
			}
			return &object.EvalResult{
				Value:  arg,
				Signal: object.SIGNAL_NONE,
			}, nil
		default:
			hashKey, ok := args[1].(object.Hashable)
			if !ok {
				return nil, errors.New("unusable as hash key: " + string(args[1].Type()))
			}
			delete(arg.(*object.HashMap).Pairs, hashKey.HashKey())
			return &object.EvalResult{
				Value:  arg,
				Signal: object.SIGNAL_NONE,
			}, nil
		}
	case "delete":
		switch arg := args[0].(type) {
		case *object.Array:
			eleToDelete := args[1].Inspect()
			for i, ele := range arg.Elements {
				if ele.Inspect() == eleToDelete {
					arg.Elements = append(arg.Elements[:i], arg.Elements[i+1:]...)
					return &object.EvalResult{
						Value:  args[1],
						Signal: object.SIGNAL_NONE,
					}, nil
				}
			}
			return &object.EvalResult{
				Value:  nil,
				Signal: object.SIGNAL_NONE,
			}, nil
		default:
			hashKey, ok := args[1].(object.Hashable)
			if !ok {
				return nil, errors.New("unusable as hash key: " + string(args[1].Type()))
			}
			hMap := arg.(*object.HashMap)
			pair, ok := hMap.Pairs[hashKey.HashKey()]
			if !ok {
				return &object.EvalResult{
					Value:  nil,
					Signal: object.SIGNAL_NONE,
				}, nil
			}
			delete(hMap.Pairs, hashKey.HashKey())
			return &object.EvalResult{
				Value:  pair.Value,
				Signal: object.SIGNAL_NONE,
			}, nil
		}
	case "slice":
		start := args[1].(*object.Integer).Value
		end := args[2].(*object.Integer).Value
		switch arg := args[0].(type) {
		case *object.Array:
			if start < 0 || start >= int64(len(arg.Elements)) ||
				end < 0 || end > int64(len(arg.Elements)) || start > end {
				return nil,
					errors.New("index out of range, can't slice array from " +
						strconv.FormatInt(start, 10) + " to " +
						strconv.FormatInt(end, 10),
					)
			}
			newArr := &object.Array{}

			if len(args) == 3 {
				sliced := arg.Elements[start:end]
				newArr.Elements = sliced
			} else {
				step := args[3].(*object.Integer).Value
				if step <= 0 {
					return nil,
						errors.New("step must be a positive integer, got: " +
							strconv.FormatInt(step, 10),
						)
				}
				var sliced []object.Object
				for i := start; i < end; i += step {
					sliced = append(sliced, arg.Elements[i])
				}
				newArr.Elements = sliced
			}

			return &object.EvalResult{
				Value:  newArr,
				Signal: object.SIGNAL_NONE,
			}, nil
		default:
			s := arg.(*object.String).Value
			s = s[1 : len(s)-1]
			if start < 0 || start >= int64(len(s)) ||
				end < 0 || end > int64(len(s)) || start > end {
				return nil,
					errors.New("index out of range, can't slice string from " +
						strconv.FormatInt(start, 10) + " to " +
						strconv.FormatInt(end, 10),
					)
			}
			newStr := &object.String{}

			if len(args) == 3 {
				sliced := s[start:end]
				newStr.Value = "\"" + sliced + "\""
			} else {
				step := args[3].(*object.Integer).Value
				if step <= 0 {
					return nil,
						errors.New("step must be a positive integer, got: " +
							strconv.FormatInt(step, 10),
						)
				}
				var sliced strings.Builder
				for i := start; i < end; i += step {
					sliced.WriteByte(s[i])
				}
				newStr.Value = "\"" + sliced.String() + "\""
			}

			return &object.EvalResult{
				Value:  newStr,
				Signal: object.SIGNAL_NONE,
			}, nil
		}
	case "equals":
		switch arg := args[0].(type) {
		case *object.Integer:
			other := args[1].(*object.Integer)
			if arg.Value == other.Value {
				return TRUE, nil
			}
			return FALSE, nil
		case *object.Float:
			other := args[1].(*object.Float)
			if arg.Value == other.Value {
				return TRUE, nil
			}
			return FALSE, nil
		case *object.Bool:
			other := args[1].(*object.Bool)
			if arg.Value == other.Value {
				return TRUE, nil
			}
			return FALSE, nil
		case *object.String:
			other := args[1].(*object.String)
			if arg.Value == other.Value {
				return TRUE, nil
			}
			return FALSE, nil
		case *object.Char:
			other := args[1].(*object.Char)
			if arg.Value == other.Value {
				return TRUE, nil
			}
			return FALSE, nil
		case *object.Array:
			other := args[1].(*object.Array)
			if len(arg.Elements) != len(other.Elements) {
				return FALSE, nil
			}
			if arg.Inspect() != other.Inspect() {
				return FALSE, nil
			}
			return TRUE, nil
		default:
			h1 := args[0].(*object.HashMap)
			h2 := args[1].(*object.HashMap)
			if len(h1.Pairs) != len(h2.Pairs) {
				return FALSE, nil
			}
			if h1.Inspect() != h2.Inspect() {
				return FALSE, nil
			}
			return TRUE, nil
		}
	case "copy":
		return &object.EvalResult{
			Value:  deepCopy(args[0]),
			Signal: object.SIGNAL_NONE,
		}, nil
	default:
		return nil, nil
	}
}
