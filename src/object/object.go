package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/KhushPatibandha/Kolon/src/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	FLOAT_OBJ        = "FLOAT"
	CHAR_OBJ         = "CHARACTER"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

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
