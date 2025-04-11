package parser

import (
	"errors"

	"github.com/KhushPatibandha/Kolon/src/ast"
)

type VariableType int

const (
	VAR VariableType = iota
	CONST
	FUNCTION
)

type Variable struct {
	Ident   ast.Identifier
	Type    ast.Type
	VarType VariableType
	Env     *Environment
}

type Environment struct {
	store map[string]*Variable
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]*Variable)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (*Variable, bool) {
	variable, ok := e.store[name]
	if !ok && e.outer != nil {
		variable, ok = e.outer.Get(name)
	}
	return variable, ok
}

func (e *Environment) Set(name string, ident ast.Identifier, identType ast.Type, identVarType VariableType, env *Environment) error {
	variable, ok := e.Get(ident.Value)
	if ok {
		if variable.VarType == CONST {
			return errors.New("variable `" + ident.Value + "` is a constant, can't re-declare const variables")
		}
		if variable.VarType != identVarType {
			return errors.New("variable `" + ident.Value + "` already declared as `" + varTypeToString(variable.VarType) + "` can't re-declare as `" + varTypeToString(identVarType) + "`")
		}
	}
	e.store[name] = &Variable{
		Ident:   ident,
		Type:    identType,
		VarType: identVarType,
		Env:     env,
	}
	return nil
}

func varTypeToString(varType VariableType) string {
	switch varType {
	case VAR:
		return "variable"
	case CONST:
		return "constant"
	case FUNCTION:
		return "function"
	default:
		return "unknown"
	}
}
