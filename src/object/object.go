package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

type (
	ObjectType string
	SignalType int
)

const (
	ARRAY_OBJ   = "ARRAY"
	HASHMAP_OBJ = "HASHMAP"
	INTEGER_OBJ = "INT"
	FLOAT_OBJ   = "FLOAT"
	BOOLEAN_OBJ = "BOOL"
	STRING_OBJ  = "STRING"
	CHAR_OBJ    = "CHAR"
	MULTI_OBJ   = "MULTI"
)

const (
	SIGNAL_NONE SignalType = iota
	SIGNAL_RETURN
	SIGNAL_BREAK
	SIGNAL_CONTINUE
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type EvalResult struct {
	Value  Object
	Signal SignalType
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

// ------------------------------------------------------------------------------------------------------------------
// Integer
// ------------------------------------------------------------------------------------------------------------------
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// ------------------------------------------------------------------------------------------------------------------
// Float
// ------------------------------------------------------------------------------------------------------------------
type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: uint64(f.Value)}
}

// ------------------------------------------------------------------------------------------------------------------
// Bool
// ------------------------------------------------------------------------------------------------------------------
type Bool struct {
	Value bool
}

func (b *Bool) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Bool) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Bool) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

// ------------------------------------------------------------------------------------------------------------------
// String
// ------------------------------------------------------------------------------------------------------------------
type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// ------------------------------------------------------------------------------------------------------------------
// Char
// ------------------------------------------------------------------------------------------------------------------
type Char struct {
	Value string
}

func (c *Char) Inspect() string  { return c.Value }
func (c *Char) Type() ObjectType { return CHAR_OBJ }
func (c *Char) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(c.Value))
	return HashKey{Type: c.Type(), Value: h.Sum64()}
}

// ------------------------------------------------------------------------------------------------------------------
// Array
// ------------------------------------------------------------------------------------------------------------------
type Array struct {
	Elements []Object
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

// ------------------------------------------------------------------------------------------------------------------
// HashMap
// ------------------------------------------------------------------------------------------------------------------
type HashPair struct {
	Key   Object
	Value Object
}
type Hashable interface {
	HashKey() HashKey
}
type HashMap struct {
	Pairs map[HashKey]HashPair
}

func (h *HashMap) Inspect() string {
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
func (h *HashMap) Type() ObjectType { return HASHMAP_OBJ }

// ------------------------------------------------------------------------------------------------------------------
// MultiValue
// ------------------------------------------------------------------------------------------------------------------
type MultiValue struct {
	Values []Object
}

func (mv *MultiValue) Inspect() string {
	var values []string
	for _, v := range mv.Values {
		values = append(values, v.Inspect())
	}
	return fmt.Sprintf("(%s)", strings.Join(values, ", "))
}
func (mv *MultiValue) Type() ObjectType { return MULTI_OBJ }
