package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/KhushPatibandha/Kolon/src/ast"
)

type (
	ObjectType      string
	BuiltinFunction func(args ...Object) (Object, bool, error)
)

const (
	ARRAY_OBJ        = "ARRAY"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	FLOAT_OBJ        = "FLOAT"
	CHAR_OBJ         = "CHARACTER"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	HASH_OBJ         = "HASH"
	CONTINUE_OBJ     = "CONTINUE"
	BREAK_OBJ        = "BREAK"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: uint64(f.Value)}
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (c *Char) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(c.Value))
	return HashKey{Type: c.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs     map[HashKey]HashPair
	KeyType   string
	ValueType string
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

type Array struct {
	Elements []Object
	TypeOf   string
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (a *Array) Type() ObjectType { return ARRAY_OBJ }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }

type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }

type Char struct {
	Value string
}

func (c *Char) Inspect() string  { return c.Value }
func (c *Char) Type() ObjectType { return CHAR_OBJ }

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type Continue struct{}

func (c *Continue) Inspect() string  { return "continue" }
func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }

type Break struct{}

func (b *Break) Inspect() string  { return "break" }
func (b *Break) Type() ObjectType { return BREAK_OBJ }

type ReturnValue struct {
	Value []Object
}

func (rv *ReturnValue) Inspect() string {
	var out bytes.Buffer

	for i := 0; i < len(rv.Value); i++ {
		out.WriteString(rv.Value[i].Inspect())
		if i != len(rv.Value)-1 {
			out.WriteString(", ")
		}
	}

	return out.String()
}
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Function struct {
	Name       *ast.Identifier
	Parameters []*ast.FunctionParameters
	ReturnType []*ast.FunctionReturnType
	Body       *ast.FunctionBody
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	returnTypes := []string{}
	for _, r := range f.ReturnType {
		returnTypes = append(returnTypes, r.String())
	}

	out.WriteString("fun: " + f.Name.Value + "(")

	if len(f.Parameters) > 0 {
		out.WriteString(strings.Join(params, ", "))
	}

	out.WriteString(")")

	if f.ReturnType != nil {
		out.WriteString(": (")
		out.WriteString(strings.Join(returnTypes, ", "))
		out.WriteString(")")
	}

	out.WriteString(" { " + f.Body.String() + " }")
	return out.String()
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
