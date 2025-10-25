package evaluator

import (
	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/environment"
	"github.com/KhushPatibandha/Kolon/src/object"
)

// ------------------------------------------------------------------------------------------------------------------
// Statements
// ------------------------------------------------------------------------------------------------------------------
func evalStmts(stmts []ast.Statement, env *environment.Environment) (*object.EvalResult, error) {
	for _, stmt := range stmts {
		r, err := Evaluate(stmt, env)
		if err != nil {
			return nil, err
		}
		if r.Signal != object.SIGNAL_NONE {
			return r, nil
		}
	}
	return nil, nil
}
