package ast

import (
	"bytes"
	"fmt"
	"strings"

	ktype "github.com/KhushPatibandha/Kolon/src/kType"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

// ------------------------------------------------------------------------------------------------------------------
// Expressions
// ------------------------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------------------------
// Identifier
// ------------------------------------------------------------------------------------------------------------------
type Identifier struct {
	Token lexer.Token
	Value string
	Type  *ktype.Type
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{i.Type},
		TypeLen: 1,
	}
}
func (i *Identifier) TokenValue() string { return i.Token.Value }
func (i *Identifier) String() string     { return i.TokenValue() }
func (i *Identifier) Equals(other *Identifier) bool {
	return i.Value == other.Value
}

// ------------------------------------------------------------------------------------------------------------------
// Integer
// ------------------------------------------------------------------------------------------------------------------
type Integer struct {
	Token lexer.Token
	Value int64
	Type  *ktype.Type
}

func (i *Integer) expressionNode() {}
func (i *Integer) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{i.Type},
		TypeLen: 1,
	}
}
func (i *Integer) baseTypeNode()      {}
func (i *Integer) TokenValue() string { return i.Token.Value }
func (i *Integer) String() string     { return i.TokenValue() }

// ------------------------------------------------------------------------------------------------------------------
// Float
// ------------------------------------------------------------------------------------------------------------------
type Float struct {
	Token lexer.Token
	Value float64
	Type  *ktype.Type
}

func (f *Float) expressionNode() {}
func (f *Float) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{f.Type},
		TypeLen: 1,
	}
}
func (f *Float) baseTypeNode()      {}
func (f *Float) TokenValue() string { return fmt.Sprintf("%g", f.Value) }
func (f *Float) String() string     { return f.Token.Value }

// ------------------------------------------------------------------------------------------------------------------
// Bool
// ------------------------------------------------------------------------------------------------------------------
type Bool struct {
	Token lexer.Token
	Value bool
	Type  *ktype.Type
}

func (b *Bool) expressionNode() {}
func (b *Bool) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{b.Type},
		TypeLen: 1,
	}
}
func (b *Bool) baseTypeNode()      {}
func (b *Bool) TokenValue() string { return b.Token.Value }
func (b *Bool) String() string     { return b.TokenValue() }

// ------------------------------------------------------------------------------------------------------------------
// String
// ------------------------------------------------------------------------------------------------------------------
type String struct {
	Token lexer.Token
	Value string
	Type  *ktype.Type
}

func (s *String) expressionNode() {}
func (s *String) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{s.Type},
		TypeLen: 1,
	}
}
func (s *String) baseTypeNode()      {}
func (s *String) TokenValue() string { return s.Token.Value }
func (s *String) String() string     { return s.TokenValue() }

// ------------------------------------------------------------------------------------------------------------------
// Char
// ------------------------------------------------------------------------------------------------------------------
type Char struct {
	Token lexer.Token
	Value string
	Type  *ktype.Type
}

func (c *Char) expressionNode() {}
func (c *Char) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{c.Type},
		TypeLen: 1,
	}
}
func (c *Char) baseTypeNode()      {}
func (c *Char) TokenValue() string { return c.Token.Value }
func (c *Char) String() string     { return c.TokenValue() }

// ------------------------------------------------------------------------------------------------------------------
// HashMap
// ------------------------------------------------------------------------------------------------------------------
type HashMap struct {
	Token     lexer.Token
	KeyType   *ktype.Type
	ValueType *ktype.Type
	Pairs     map[BaseType]Expression
}

func (hm *HashMap) expressionNode() {}
func (hm *HashMap) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types: []*ktype.Type{
			{
				Kind:      ktype.TypeHashMap,
				KeyType:   hm.KeyType,
				ValueType: hm.ValueType,
			},
		},
		TypeLen: 1,
	}
}
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
	Token  lexer.Token
	Type   *ktype.Type
	Values []Expression
}

func (a *Array) expressionNode() {}
func (a *Array) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types: []*ktype.Type{
			{
				Kind:        ktype.TypeArray,
				ElementType: a.Type,
			},
		},
		TypeLen: 1,
	}
}
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
	Token    lexer.Token
	Operator string
	Right    Expression
	Type     *ktype.Type
}

func (p *Prefix) expressionNode() {}
func (p *Prefix) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{p.Type},
		TypeLen: 1,
	}
}
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
	Token    lexer.Token
	Left     Expression
	Operator string
	Right    Expression
	Type     *ktype.Type
}

func (i *Infix) expressionNode() {}
func (i *Infix) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{i.Type},
		TypeLen: 1,
	}
}
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
	Token    lexer.Token
	Left     Expression
	Operator string
	Type     *ktype.Type
}

func (p *Postfix) canBeStatement() {}
func (p *Postfix) expressionNode() {}
func (p *Postfix) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{p.Type},
		TypeLen: 1,
	}
}
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
	Token    lexer.Token
	Left     *Identifier
	Operator string
	Right    Expression
	Type     *ktype.Type
}

func (a *Assignment) canBeStatement() {}
func (a *Assignment) expressionNode() {}
func (a *Assignment) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{a.Type},
		TypeLen: 1,
	}
}
func (a *Assignment) TokenValue() string { return a.Token.Value }
func (a *Assignment) String() string {
	return a.Left.String() + " " + a.Operator + " " + a.Right.String()
}

// ------------------------------------------------------------------------------------------------------------------
// CallExpression
// ------------------------------------------------------------------------------------------------------------------
type CallExpression struct {
	Token lexer.Token
	Name  *Identifier
	Args  []Expression
	Type  []*ktype.Type
}

func (ce *CallExpression) canBeStatement() {}
func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   ce.Type,
		TypeLen: len(ce.Type),
	}
}
func (ce *CallExpression) TokenValue() string { return ce.Token.Value }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Name.String())
	out.WriteString("(")

	if ce.Args != nil {
		args := []string{}
		for _, a := range ce.Args {
			args = append(args, a.String())
		}
		out.WriteString(strings.Join(args, ", "))
	}

	out.WriteString(")")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// IndexExpression
// ------------------------------------------------------------------------------------------------------------------
type IndexExpression struct {
	Token lexer.Token
	Left  Expression
	Index Expression
	Type  *ktype.Type
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) GetType() *ktype.TypeCheckResult {
	return &ktype.TypeCheckResult{
		Types:   []*ktype.Type{ie.Type},
		TypeLen: 1,
	}
}
func (ie *IndexExpression) TokenValue() string { return ie.Token.Value }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	return out.String()
}
