package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/KhushPatibandha/Kolon/src/lexer"
)

type Node interface {
	TokenValue() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenValue() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenValue()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// -----------------------------------------------------------------------------
// For HashMap
// -----------------------------------------------------------------------------
type HashMap struct {
	Token     lexer.Token // {
	KeyType   *Type
	ValueType *Type
	Pairs     map[Expression]Expression
}

func (hm *HashMap) expressionNode()    {}
func (hm *HashMap) TokenValue() string { return hm.Token.Value }
func (hm *HashMap) String() string {
	var out bytes.Buffer
	pair := []string{}
	for key, value := range hm.Pairs {
		pair = append(pair, key.String()+": "+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pair, ", "))
	out.WriteString("}")
	return out.String()
}

// -----------------------------------------------------------------------------
// For Array
// -----------------------------------------------------------------------------
type ArrayValue struct {
	Token  lexer.Token // {
	Type   *Type
	Values []Expression
}

func (av *ArrayValue) expressionNode()    {}
func (av *ArrayValue) TokenValue() string { return av.Token.Value }
func (av *ArrayValue) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range av.Values {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// -----------------------------------------------------------------------------
// For Integer
// -----------------------------------------------------------------------------
type IntegerValue struct {
	Token lexer.Token
	Value int64
}

func (iv *IntegerValue) expressionNode()    {}
func (iv *IntegerValue) TokenValue() string { return iv.Token.Value }
func (iv *IntegerValue) String() string     { return iv.Token.Value }

// -----------------------------------------------------------------------------
// For Float
// -----------------------------------------------------------------------------
type FloatValue struct {
	Token lexer.Token
	Value float64
}

func (fv *FloatValue) expressionNode() {}

func (fv *FloatValue) TokenValue() string { return fmt.Sprintf("%g", fv.Value) }
func (fv *FloatValue) String() string     { return fv.Token.Value }

// -----------------------------------------------------------------------------
// For String
// -----------------------------------------------------------------------------
type StringValue struct {
	Token lexer.Token
	Value string
}

func (sv *StringValue) expressionNode()    {}
func (sv *StringValue) TokenValue() string { return sv.Token.Value }
func (sv *StringValue) String() string     { return sv.Token.Value }

// -----------------------------------------------------------------------------
// For Boolean
// -----------------------------------------------------------------------------
type BooleanValue struct {
	Token lexer.Token
	Value bool
}

func (bv *BooleanValue) expressionNode()    {}
func (bv *BooleanValue) TokenValue() string { return bv.Token.Value }
func (bv *BooleanValue) String() string     { return bv.Token.Value }

// -----------------------------------------------------------------------------
// For char
// -----------------------------------------------------------------------------
type CharValue struct {
	Token lexer.Token
	Value string
}

func (cv *CharValue) expressionNode()    {}
func (cv *CharValue) TokenValue() string { return cv.Token.Value }
func (cv *CharValue) String() string     { return cv.Token.Value }

// -----------------------------------------------------------------------------
// For Identifier
// -----------------------------------------------------------------------------
type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode()    {}
func (i *Identifier) TokenValue() string { return i.Token.Value }
func (i *Identifier) String() string     { return i.Value }

// -----------------------------------------------------------------------------
// For Type (int, string, char, float, bool)
// -----------------------------------------------------------------------------
type Type struct {
	Token    lexer.Token
	Value    string
	IsArray  bool
	IsHash   bool
	SubTypes []*Type
}

func (t *Type) TokenValue() string { return t.Token.Value }
func (t *Type) String() string {
	var out bytes.Buffer
	if t.IsHash {
		// key[value]
		out.WriteString(t.SubTypes[0].String() + "[" + t.SubTypes[1].String() + "]")
	} else {
		out.WriteString(t.TokenValue())
		if t.IsArray {
			out.WriteString("[]")
		}
	}
	return out.String()
}

// -----------------------------------------------------------------------------
// For Function Parameters
// -----------------------------------------------------------------------------
type FunctionParameters struct {
	// Token         lexer.Token
	ParameterName *Identifier
	ParameterType *Type
}

func (fp *FunctionParameters) TokenValue() string { return fp.ParameterName.Token.Value }
func (fp *FunctionParameters) String() string {
	var out bytes.Buffer
	out.WriteString(fp.ParameterName.String())
	out.WriteString(": ")
	out.WriteString(fp.ParameterType.String())
	return out.String()
}

// -----------------------------------------------------------------------------
// For Functions Return Type
// -----------------------------------------------------------------------------
type FunctionReturnType struct {
	// Token      lexer.Token
	ReturnType *Type
}

func (frt *FunctionReturnType) TokenValue() string { return frt.ReturnType.Token.Value }
func (frt *FunctionReturnType) String() string     { return frt.ReturnType.String() }

// -----------------------------------------------------------------------------
// For Body (function, if, else, else if, for) -- Basically everything with {...}
// -----------------------------------------------------------------------------
type FunctionBody struct {
	Token      lexer.Token // '{' token
	Statements []Statement
}

func (fb *FunctionBody) statementNode()     {}
func (fb *FunctionBody) TokenValue() string { return fb.Token.Value }
func (fb *FunctionBody) String() string {
	var out bytes.Buffer

	for _, s := range fb.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// -----------------------------------------------------------------------------
// For Multiple declaration of var stmt or variable assignment
// eg: var a: int, var b: int, c, d = 1, 2, 3, 4
// Here a and b are declared and c and d are already declared, but assigned newly
// -----------------------------------------------------------------------------
type MultiValueAssignStmt struct {
	Token         lexer.Token // "=" (Equal Assign) Token
	Objects       []Statement // it will store objects of type VarStatement and ExpressionStatement(Expression in es will always be AssignemntExpression with operator "=")
	SingleCallExp bool
}

func (mvas *MultiValueAssignStmt) statementNode()     {}
func (mvas *MultiValueAssignStmt) TokenValue() string { return mvas.Token.Value }
func (mvas *MultiValueAssignStmt) String() string {
	var out bytes.Buffer

	for _, obj := range mvas.Objects {
		out.WriteString(obj.String())
	}

	return out.String()
}

// -----------------------------------------------------------------------------
// For Var AND Const Statement
// -----------------------------------------------------------------------------
type VarStatement struct {
	Token lexer.Token
	Name  *Identifier
	Type  *Type
	Value Expression
}

func (vr *VarStatement) statementNode()     {}
func (vr *VarStatement) TokenValue() string { return vr.Token.Value }
func (vr *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vr.TokenValue() + " ")
	out.WriteString(vr.Name.String())
	out.WriteString(": ")
	out.WriteString(vr.Type.String())

	if vr.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vr.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// -----------------------------------------------------------------------------
// For Return Statement
// -----------------------------------------------------------------------------
type ReturnStatement struct {
	Token lexer.Token
	Value []Expression
}

func (r *ReturnStatement) statementNode()     {}
func (r *ReturnStatement) TokenValue() string { return r.Token.Value }
func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenValue())

	if r.Value != nil {
		out.WriteString(": ")
	}

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

	out.WriteString(";")
	return out.String()
}

// -----------------------------------------------------------------------------
// For Continue Statement
// -----------------------------------------------------------------------------
type ContinueStatement struct {
	Token lexer.Token
}

func (cs *ContinueStatement) statementNode()     {}
func (cs *ContinueStatement) TokenValue() string { return cs.Token.Value }
func (cs *ContinueStatement) String() string     { return cs.TokenValue() + ";" }

// -----------------------------------------------------------------------------
// For Break Statement
// -----------------------------------------------------------------------------
type BreakStatement struct {
	Token lexer.Token
}

func (bs *BreakStatement) statementNode()     {}
func (bs *BreakStatement) TokenValue() string { return bs.Token.Value }
func (bs *BreakStatement) String() string     { return bs.TokenValue() + ";" }

// -----------------------------------------------------------------------------
// For Functions
// -----------------------------------------------------------------------------
type Function struct {
	Token      lexer.Token
	Name       *Identifier
	Parameters []*FunctionParameters
	ReturnType []*FunctionReturnType
	Body       *FunctionBody
}

func (f *Function) statementNode()     {}
func (f *Function) TokenValue() string { return f.Token.Value }
func (f *Function) String() string {
	var out bytes.Buffer

	out.WriteString(f.TokenValue() + ": ")
	out.WriteString(f.Name.String() + "(")
	for i, param := range f.Parameters {
		out.WriteString(param.String())
		if i != len(f.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")

	if f.ReturnType != nil {
		out.WriteString(": (")
		for i, param := range f.ReturnType {
			out.WriteString(param.String())
			if i != len(f.ReturnType)-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString(")")
	}

	out.WriteString(" {")
	out.WriteString(f.Body.String())
	out.WriteString("}")

	return out.String()
}

// -----------------------------------------------------------------------------
// For If Statement
// -----------------------------------------------------------------------------
type IfStatement struct {
	Token       lexer.Token
	Value       Expression
	Body        *FunctionBody
	MultiConseq []*ElseIfStatement
	Consequence *ElseStatement
}

func (ifs *IfStatement) statementNode()     {}
func (ifs *IfStatement) TokenValue() string { return ifs.Token.Value }
func (ifs *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ifs.TokenValue() + ": ")
	out.WriteString("(")
	out.WriteString(ifs.Value.String())
	out.WriteString("): ")
	out.WriteString("{")
	out.WriteString(ifs.Body.String())
	out.WriteString("}")

	if ifs.MultiConseq != nil {
		for i := 0; i < len(ifs.MultiConseq); i++ {
			out.WriteString(ifs.MultiConseq[i].String())
		}
	}

	if ifs.Consequence != nil {
		out.WriteString(ifs.Consequence.String())
	}

	return out.String()
}

// -----------------------------------------------------------------------------
// For Else Statement
// -----------------------------------------------------------------------------
type ElseStatement struct {
	Token lexer.Token
	Body  *FunctionBody
}

func (el *ElseStatement) statementNode()     {}
func (el *ElseStatement) TokenValue() string { return el.Token.Value }
func (el *ElseStatement) String() string {
	var out bytes.Buffer

	out.WriteString(el.TokenValue() + ": {")
	out.WriteString(el.Body.String())
	out.WriteString("}")
	return out.String()
}

// -----------------------------------------------------------------------------
// For Else If Statement
// -----------------------------------------------------------------------------
type ElseIfStatement struct {
	Token lexer.Token
	Value Expression
	Body  *FunctionBody
}

func (eis *ElseIfStatement) statementNode()     {}
func (eis *ElseIfStatement) TokenValue() string { return eis.Token.Value }
func (eis *ElseIfStatement) String() string {
	var out bytes.Buffer

	out.WriteString(eis.TokenValue() + ": ")
	out.WriteString("(")
	out.WriteString(eis.Value.String())
	out.WriteString("): {")
	out.WriteString(eis.Body.String())
	out.WriteString("}")
	return out.String()
}

// -----------------------------------------------------------------------------
// for loop statement
// -----------------------------------------------------------------------------
type ForLoopStatement struct {
	Token  lexer.Token
	Left   Statement
	Middle *InfixExpression
	Right  Expression
	Body   *FunctionBody
}

func (fls *ForLoopStatement) statementNode()     {}
func (fls *ForLoopStatement) TokenValue() string { return fls.Token.Value }
func (fls *ForLoopStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fls.TokenValue() + ": (")
	out.WriteString(fls.Left.String() + " ")
	out.WriteString(fls.Middle.String())
	out.WriteString("; ")
	out.WriteString(fls.Right.String())
	out.WriteString("): {")
	out.WriteString(fls.Body.String())
	out.WriteString("}")

	return out.String()
}

// -----------------------------------------------------------------------------
// For While loop Statement
// -----------------------------------------------------------------------------
type WhileLoopStatement struct {
	Token     lexer.Token
	Condition Expression
	Body      *FunctionBody
}

func (wls *WhileLoopStatement) statementNode()     {}
func (wls *WhileLoopStatement) TokenValue() string { return wls.Token.Value }
func (wls *WhileLoopStatement) String() string {
	var out bytes.Buffer

	out.WriteString(wls.TokenValue() + ": (" + wls.Condition.String() + "): {")
	out.WriteString(wls.Body.String())
	out.WriteString("}")

	return out.String()
}

// -----------------------------------------------------------------------------
// For Expression Statement
// -----------------------------------------------------------------------------
type ExpressionStatement struct {
	Token      lexer.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()     {}
func (es *ExpressionStatement) TokenValue() string { return es.Token.Value }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// -----------------------------------------------------------------------------
// For Prefix Expression
// -----------------------------------------------------------------------------
type PrefixExpression struct {
	Token    lexer.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()    {}
func (pe *PrefixExpression) TokenValue() string { return pe.Token.Value }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// -----------------------------------------------------------------------------
// For Infix Expression
// -----------------------------------------------------------------------------
type InfixExpression struct {
	Token    lexer.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()    {}
func (ie *InfixExpression) TokenValue() string { return ie.Token.Value }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// -----------------------------------------------------------------------------
// For Postfix Expression
// -----------------------------------------------------------------------------
type PostfixExpression struct {
	Token    lexer.Token
	Operator string
	Left     Expression
	IsStmt   bool
}

func (poe *PostfixExpression) statementNode()     {}
func (poe *PostfixExpression) expressionNode()    {}
func (poe *PostfixExpression) TokenValue() string { return poe.Token.Value }
func (poe *PostfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(poe.Left.String())
	out.WriteString(poe.Operator)
	out.WriteString(")")

	if poe.IsStmt {
		out.WriteString(";")
	}

	return out.String()
}

// -----------------------------------------------------------------------------
// For Assignment Expression, Operator (=, +=, -=, *=, /=, %=)
// -----------------------------------------------------------------------------
type AssignmentExpression struct {
	Token    lexer.Token
	Left     *Identifier
	Operator string
	Right    Expression
}

func (ae *AssignmentExpression) expressionNode()    {}
func (ae *AssignmentExpression) TokenValue() string { return ae.Token.Value }
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Left.String() + " " + ae.Operator + " " + ae.Right.String() + ";")
	return out.String()
}

// -----------------------------------------------------------------------------
// For Call Expression
// -----------------------------------------------------------------------------
type CallExpression struct {
	Token lexer.Token
	Name  Expression
	Args  []Expression
}

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

// -----------------------------------------------------------------------------
// Index Expression
// -----------------------------------------------------------------------------
type IndexExpression struct {
	Token lexer.Token // [
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()    {}
func (ie *IndexExpression) TokenValue() string { return ie.Token.Value }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ie.Left.String() + "[" + ie.Index.String() + "]")
	return out.String()
}
