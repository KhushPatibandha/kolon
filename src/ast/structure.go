package ast

import (
	"bytes"

	ktype "github.com/KhushPatibandha/Kolon/src/kType"
)

type Program struct{ Statements []Statement }

// ------------------------------------------------------------------------------------------------------------------
// Statement: A piece of code that performs an action and does not evaluate to a value
// ------------------------------------------------------------------------------------------------------------------
type Statement interface {
	Node
	statementNode()
}

// ------------------------------------------------------------------------------------------------------------------
// Expression: A piece of code that will generate or evaluate to a value
// ------------------------------------------------------------------------------------------------------------------
type Expression interface {
	Node
	GetType() *ktype.TypeCheckResult
	expressionNode()
}

// ------------------------------------------------------------------------------------------------------------------
// BaseType: Smallest unit of data in the language
// eg: Integer, Float, Bool, String, Char
// ------------------------------------------------------------------------------------------------------------------
type BaseType interface {
	Expression
	baseTypeNode()
}

// ------------------------------------------------------------------------------------------------------------------
// StatementableExpression: An expression that can also be used as a statement
// eg: Postfix expression like i++ or i--, function calls that returns no value
// ------------------------------------------------------------------------------------------------------------------
type StatementableExpression interface {
	Expression
	canBeStatement()
}

// ------------------------------------------------------------------------------------------------------------------
// Node: Any construct in the AST
// ------------------------------------------------------------------------------------------------------------------
type Node interface {
	TokenValue() string
	String() string
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
