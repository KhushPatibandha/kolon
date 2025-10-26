package environment

import (
	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

func LoadBuiltins(env *Environment) {
	builtins := []string{
		"print", "println", "scan", "scanln", "len", "toString", "toFloat", "toInt",
		"push", "pop", "insert", "remove", "getIndex", "keys", "values", "containsKey",
		"typeOf", "slice", "delete", "equals", "copy",
	}
	for _, name := range builtins {
		env.FuncNameSpace[name] = &Symbol{
			IdentType: FUNCTION,
			Ident: &ast.Identifier{
				Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: name},
				Value: name,
			},
			Func: &FuncInfo{
				Builtin:  true,
				Function: nil,
			},
			Type: nil,
			Env:  nil,
		}
	}
}

func BootstrapFuncEnv(stmt *ast.Function, env *Environment) *Environment {
	funcLocalEnv := NewEnclosedEnvironment(env)
	for _, param := range stmt.Parameters {
		funcLocalEnv.Set(&Symbol{
			IdentType: VAR,
			Ident:     param.ParameterName,
			Type:      param.ParameterType,
			Func:      nil,
			Env:       nil,
		})
	}
	return funcLocalEnv
}
