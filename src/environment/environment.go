package environment

import (
	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/kType"
)

type IdentType int

const (
	VAR IdentType = iota
	CONST
	FUNCTION
)

type Symbol struct {
	IdentType IdentType
	Ident     *ast.Identifier
	Type      *ktype.Type
	Func      *FuncInfo
	Env       *Environment
}

type FuncInfo struct {
	Function *ast.Function
	Builtin  bool
}

type Environment struct {
	VariableNameSpace map[string]*Symbol
	FuncNameSpace     map[string]*Symbol
	Outer             *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		VariableNameSpace: make(map[string]*Symbol),
		FuncNameSpace:     make(map[string]*Symbol),
		Outer:             nil,
	}
}

func NewEnclosedEnvironment(Outer *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = Outer
	return env
}

func (e *Environment) GetVar(name string) (*Symbol, bool) {
	sym, ok := e.VariableNameSpace[name]
	if !ok && e.Outer != nil {
		sym, ok = e.Outer.GetVar(name)
	}
	if !ok {
		return nil, false
	}
	return sym, true
}

func (e *Environment) GetFunc(name string) (*Symbol, bool) {
	sym, ok := e.FuncNameSpace[name]
	if !ok && e.Outer != nil {
		sym, ok = e.Outer.GetFunc(name)
	}
	if !ok {
		return nil, false
	}
	return sym, true
}

func (e *Environment) Set(sym *Symbol) {
	switch sym.IdentType {
	case VAR, CONST:
		e.VariableNameSpace[sym.Ident.Value] = sym
	case FUNCTION:
		e.FuncNameSpace[sym.Ident.Value] = sym
	}
}
