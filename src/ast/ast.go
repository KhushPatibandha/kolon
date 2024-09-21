package ast

import "github.com/KhushPatibandha/Kolon/src/lexer"

type Node interface {
	TokenValue() string
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

// -----------------------------------------------------------------------------
// For Type (int, string, char, long, float, double, bool)
// -----------------------------------------------------------------------------
type Type struct {
	Token lexer.Token
	Value string
}

func (t *Type) TokenValue() string { return t.Token.Value }

// -----------------------------------------------------------------------------
// For Function Parameters
// -----------------------------------------------------------------------------
type FunctionParameters struct {
	// Token         lexer.Token
	ParameterName *Identifier
	ParameterType *Type
}

func (fp *FunctionParameters) TokenValue() string { return fp.ParameterName.Token.Value }

// -----------------------------------------------------------------------------
// For Functions Return Type
// -----------------------------------------------------------------------------
type FunctionReturnType struct {
	// Token      lexer.Token
	ReturnType *Type
}

func (frt *FunctionReturnType) TokenValue() string { return frt.ReturnType.Token.Value }

// -----------------------------------------------------------------------------
// For Function Body
// -----------------------------------------------------------------------------
type FunctionBody struct {
	Token      lexer.Token // '{' token
	Statements []Statement
}

func (fb *FunctionBody) statementNode()     {}
func (fb *FunctionBody) TokenValue() string { return fb.Token.Value }

// -----------------------------------------------------------------------------
// For Identifier
// -----------------------------------------------------------------------------
type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode()    {}
func (i *Identifier) TokenValue() string { return i.Token.Value }

// -----------------------------------------------------------------------------
// For Var Statement
// -----------------------------------------------------------------------------
type VarStatement struct {
	Token lexer.Token
	Name  *Identifier
	Type  *Type
	Value Expression
}

func (vr *VarStatement) statementNode()     {}
func (vr *VarStatement) TokenValue() string { return vr.Token.Value }

// -----------------------------------------------------------------------------
// For Return Statement
// -----------------------------------------------------------------------------
type ReturnStatement struct {
	Token lexer.Token
	Value Expression
}

func (r *ReturnStatement) statementNode()     {}
func (r *ReturnStatement) TokenValue() string { return r.Token.Value }

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

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
