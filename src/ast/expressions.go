package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/KhushPatibandha/Kolon/src/lexer"
)

// ------------------------------------------------------------------------------------------------------------------
// Expressions
// ------------------------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------------------------
// Identifier
// ------------------------------------------------------------------------------------------------------------------
type Identifier struct {
	Token *lexer.Token
	value string
}

func (i *Identifier) expressionNode()    {}
func (i *Identifier) TokenValue() string { return i.Token.Value }
func (i *Identifier) String() string     { return i.TokenValue() }

// ------------------------------------------------------------------------------------------------------------------
// Integer
// ------------------------------------------------------------------------------------------------------------------
type Integer struct {
	Token *lexer.Token
	value int64
}

func (i *Integer) expressionNode()    {}
func (i *Integer) baseTypeNode()      {}
func (i *Integer) TokenValue() string { return i.Token.Value }
func (i *Integer) String() string     { return i.TokenValue() }

// ------------------------------------------------------------------------------------------------------------------
// Float
// ------------------------------------------------------------------------------------------------------------------
type Float struct {
	Token *lexer.Token
	value float64
}

func (f *Float) expressionNode()    {}
func (f *Float) baseTypeNode()      {}
func (f *Float) TokenValue() string { return fmt.Sprintf("%g", f.value) }
func (f *Float) String() string     { return f.TokenValue() }

// ------------------------------------------------------------------------------------------------------------------
// Bool
// ------------------------------------------------------------------------------------------------------------------
type Bool struct {
	Token *lexer.Token
	value bool
}

func (b *Bool) expressionNode()    {}
func (b *Bool) baseTypeNode()      {}
func (b *Bool) TokenValue() string { return b.Token.Value }
func (b *Bool) String() string     { return b.TokenValue() }

// ------------------------------------------------------------------------------------------------------------------
// String
// ------------------------------------------------------------------------------------------------------------------
type String struct {
	Token *lexer.Token
	value string
}

func (s *String) expressionNode()    {}
func (s *String) baseTypeNode()      {}
func (s *String) TokenValue() string { return s.Token.Value }
func (s *String) String() string     { return s.TokenValue() }

// ------------------------------------------------------------------------------------------------------------------
// Char
// ------------------------------------------------------------------------------------------------------------------
type Char struct {
	Token *lexer.Token
	value string
}

func (c *Char) expressionNode()    {}
func (c *Char) baseTypeNode()      {}
func (c *Char) TokenValue() string { return c.Token.Value }
func (c *Char) String() string     { return c.TokenValue() }

// ------------------------------------------------------------------------------------------------------------------
// HashMap
// ------------------------------------------------------------------------------------------------------------------
type HashMap struct {
	Token     *lexer.Token
	KeyType   *Type
	ValueType *Type
	Pairs     map[BaseType]Expression
}

func (hm *HashMap) expressionNode()    {}
func (hm *HashMap) TokenValue() string { return hm.Token.Value }
func (hm *HashMap) String() string {
	var out bytes.Buffer
	pair := []string{}
	for key, val := range hm.Pairs {
		pair = append(pair, key.String()+": "+val.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pair, ", "))
	out.WriteString("}")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// Array
// ------------------------------------------------------------------------------------------------------------------
type Array struct {
	Token  *lexer.Token
	Type   *Type
	Values []Expression
}

func (a *Array) expressionNode()    {}
func (a *Array) TokenValue() string { return a.Token.Value }
func (a *Array) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range a.Values {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// Prefix
// ------------------------------------------------------------------------------------------------------------------
type Prefix struct {
	Token    *lexer.Token
	Operator string
	Right    Expression
}

func (p *Prefix) expressionNode()    {}
func (p *Prefix) TokenValue() string { return p.Token.Value }
func (p *Prefix) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// Infix
// ------------------------------------------------------------------------------------------------------------------
type Infix struct {
	Token    *lexer.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *Infix) expressionNode()    {}
func (i *Infix) TokenValue() string { return i.Token.Value }
func (i *Infix) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// Postfix
// ------------------------------------------------------------------------------------------------------------------
type Postfix struct {
	Token    *lexer.Token
	Left     Expression
	Operator string
}

func (p *Postfix) canBeStatement()    {}
func (p *Postfix) expressionNode()    {}
func (p *Postfix) TokenValue() string { return p.Token.Value }
func (p *Postfix) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Left.String())
	out.WriteString(p.Operator)
	out.WriteString(")")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// Assignment
// ------------------------------------------------------------------------------------------------------------------
type Assignment struct {
	Token    *lexer.Token
	Left     *Identifier
	Operator string
	Right    Expression
}

func (a *Assignment) canBeStatement()    {}
func (a *Assignment) expressionNode()    {}
func (a *Assignment) TokenValue() string { return a.Token.Value }
func (a *Assignment) String() string {
	return a.Left.String() + " " + a.Operator + " " + a.Right.String()
}

// ------------------------------------------------------------------------------------------------------------------
// CallExpression
// ------------------------------------------------------------------------------------------------------------------
type CallExpression struct {
	Token *lexer.Token
	Name  *Identifier
	Args  []Expression
}

func (ce *CallExpression) canBeStatement()    {}
func (ce *CallExpression) expressionNode()    {}
func (ce *CallExpression) TokenValue() string { return ce.Token.Value }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Args {
		args = append(args, a.String())
	}
	out.WriteString(ce.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// IndexExpression
// ------------------------------------------------------------------------------------------------------------------
type IndexExpression struct {
	Token *lexer.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()    {}
func (ie *IndexExpression) TokenValue() string { return ie.Token.Value }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	return out.String()
}
