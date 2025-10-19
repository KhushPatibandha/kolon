package ast

import (
	"bytes"

	ktype "github.com/KhushPatibandha/Kolon/src/kType"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

// ------------------------------------------------------------------------------------------------------------------
// Statements
// ------------------------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------------------------
// Expression
// ------------------------------------------------------------------------------------------------------------------
type ExpressionStatement struct {
	Token      lexer.Token
	Expression StatementableExpression
}

func (es *ExpressionStatement) statementNode()     {}
func (es *ExpressionStatement) TokenValue() string { return es.Token.Value }
func (es *ExpressionStatement) String() string     { return es.Expression.String() + ";" }

// ------------------------------------------------------------------------------------------------------------------
// Body
// ------------------------------------------------------------------------------------------------------------------
type Body struct {
	Token      lexer.Token
	Statements []Statement
}

func (b *Body) statementNode()     {}
func (b *Body) TokenValue() string { return b.Token.Value }
func (b *Body) String() string {
	var out bytes.Buffer
	for _, s := range b.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// Function
// ------------------------------------------------------------------------------------------------------------------
type Function struct {
	Token       lexer.Token
	Name        *Identifier
	Parameters  []*FunctionParameter
	ReturnTypes []*ktype.Type
	Body        *Body
}

func (f *Function) statementNode()     {}
func (f *Function) TokenValue() string { return f.Token.Value }
func (f *Function) String() string {
	var out bytes.Buffer

	out.WriteString(f.TokenValue() + ": ")
	out.WriteString(f.Name.String() + "(")
	if f.Parameters != nil {
		for i, param := range f.Parameters {
			out.WriteString(param.String())
			if i != len(f.Parameters)-1 {
				out.WriteString(", ")
			}
		}
	}
	out.WriteString(")")

	if f.ReturnTypes != nil {
		out.WriteString(": (")
		for i, param := range f.ReturnTypes {
			out.WriteString(param.String())
			if i != len(f.ReturnTypes)-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString(")")
	}

	if f.Body != nil {
		out.WriteString(" {")
		out.WriteString(f.Body.String())
		out.WriteString("}")
	} else {
		out.WriteString(";")
	}

	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// Var and Const
// ------------------------------------------------------------------------------------------------------------------
type VarAndConst struct {
	Token lexer.Token
	Name  *Identifier
	Type  *ktype.Type
	Value Expression
}

func (vac *VarAndConst) statementNode()     {}
func (vac *VarAndConst) TokenValue() string { return vac.Token.Value }
func (vac *VarAndConst) String() string {
	var out bytes.Buffer
	out.WriteString(vac.TokenValue() + " ")
	out.WriteString(vac.Name.String() + ": ")
	out.WriteString(vac.Type.String())

	if vac.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vac.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// Multi-Assignment
// ------------------------------------------------------------------------------------------------------------------
type MultiAssignment struct {
	Token              lexer.Token
	Objects            []Statement
	SingleFunctionCall bool
}

func (ma *MultiAssignment) statementNode()     {}
func (ma *MultiAssignment) TokenValue() string { return ma.Token.Value }
func (ma *MultiAssignment) String() string {
	var out bytes.Buffer
	for _, obj := range ma.Objects {
		out.WriteString(obj.String())
	}
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// Return
// ------------------------------------------------------------------------------------------------------------------
type Return struct {
	Token lexer.Token
	Value []Expression
}

func (r *Return) statementNode()     {}
func (r *Return) TokenValue() string { return r.Token.Value }
func (r *Return) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenValue())
	if r.Value != nil {
		out.WriteString(": ")
		if len(r.Value) > 1 {
			out.WriteString("(")
		}
		for i, val := range r.Value {
			out.WriteString(val.String())
			if i != len(r.Value)-1 {
				out.WriteString(", ")
			}
		}
		if len(r.Value) > 1 {
			out.WriteString(")")
		}
	}
	out.WriteString(";")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// Continue
// ------------------------------------------------------------------------------------------------------------------
type Continue struct {
	Token lexer.Token
}

func (c *Continue) statementNode()     {}
func (c *Continue) TokenValue() string { return c.Token.Value }
func (c *Continue) String() string     { return c.TokenValue() + ";" }

// ------------------------------------------------------------------------------------------------------------------
// Break
// ------------------------------------------------------------------------------------------------------------------
type Break struct {
	Token lexer.Token
}

func (b *Break) statementNode()     {}
func (b *Break) TokenValue() string { return b.Token.Value }
func (b *Break) String() string     { return b.TokenValue() + ";" }

// ------------------------------------------------------------------------------------------------------------------
// If
// ------------------------------------------------------------------------------------------------------------------
type If struct {
	Token             lexer.Token
	Condition         Expression
	Body              *Body
	MultiConditionals []*ElseIf
	Alternate         *Else
}

func (i *If) statementNode()     {}
func (i *If) TokenValue() string { return i.Token.Value }
func (i *If) String() string {
	var out bytes.Buffer
	out.WriteString(i.TokenValue() + ": (")
	out.WriteString(i.Condition.String() + "): {")
	out.WriteString(i.Body.String() + "}")

	if i.MultiConditionals != nil {
		for x := range i.MultiConditionals {
			out.WriteString(i.MultiConditionals[x].String())
		}
	}
	if i.Alternate != nil {
		out.WriteString(i.Alternate.String())
	}
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// ElseIf
// ------------------------------------------------------------------------------------------------------------------
type ElseIf struct {
	Token     lexer.Token
	Condition Expression
	Body      *Body
}

func (ei *ElseIf) statementNode()     {}
func (ei *ElseIf) TokenValue() string { return ei.Token.Value }
func (ei *ElseIf) String() string {
	var out bytes.Buffer
	out.WriteString(ei.TokenValue() + ": (")
	out.WriteString(ei.Condition.String() + "): {")
	out.WriteString(ei.Body.String() + "}")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// Else
// ------------------------------------------------------------------------------------------------------------------
type Else struct {
	Token lexer.Token
	Body  *Body
}

func (e *Else) statementNode()     {}
func (e *Else) TokenValue() string { return e.Token.Value }
func (e *Else) String() string {
	var out bytes.Buffer
	out.WriteString(e.TokenValue() + ": {")
	out.WriteString(e.Body.String() + "}")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// ForLoop
// ------------------------------------------------------------------------------------------------------------------
type ForLoop struct {
	Token  lexer.Token
	Left   Statement
	Middle *Infix
	Right  StatementableExpression
	Body   *Body
}

func (f *ForLoop) statementNode()     {}
func (f *ForLoop) TokenValue() string { return f.Token.Value }
func (f *ForLoop) String() string {
	var out bytes.Buffer
	out.WriteString(f.TokenValue() + ": (")
	out.WriteString(f.Left.String() + " ")
	out.WriteString(f.Middle.String() + "; ")
	out.WriteString(f.Right.String() + "): {")
	out.WriteString(f.Body.String() + "}")
	return out.String()
}

// ------------------------------------------------------------------------------------------------------------------
// WhileLoop
// ------------------------------------------------------------------------------------------------------------------
type WhileLoop struct {
	Token     lexer.Token
	Condition Expression
	Body      *Body
}

func (w *WhileLoop) statementNode()     {}
func (w *WhileLoop) TokenValue() string { return w.Token.Value }
func (w *WhileLoop) String() string {
	var out bytes.Buffer
	out.WriteString(w.TokenValue() + ": (")
	out.WriteString(w.Condition.String() + "): {")
	out.WriteString(w.Body.String() + "}")
	return out.String()
}
