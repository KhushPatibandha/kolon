package parser_test

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func Test5(t *testing.T) {
	input := `
    if: (some condition): {
        var a: int = 10;
    } else: {
        var b: int = 20;
    }
    `

	tokens := lexer.Tokenizer(input)
	p := parser.New(tokens)
	program := p.ParseProgram()
	ce5(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	ifStmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
	}
	if len(ifStmt.Body.Statements) != 1 {
		t.Fatalf("consequence is not 1 statements. got=%d\n", len(ifStmt.Body.Statements))
	}
	ifStmtVarStmt, ok := ifStmt.Body.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.VarStatement. got=%T", ifStmt.Body.Statements[0])
	}
	if ifStmtVarStmt.Name.Value != "a" {
		t.Fatalf("varStmt.Name.Value not 'a'. got=%q", ifStmtVarStmt.Name.Value)
	}
	if ifStmtVarStmt.Type.Value != "int" {
		t.Fatalf("varStmt.Type.Value not 'int'. got=%q", ifStmtVarStmt.Type.Value)
	}

	elseStmt, ok := program.Statements[1].(*ast.ElseStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not *ast.ElseStatement. got=%T", program.Statements[1])
	}
	if len(elseStmt.Body.Statements) != 1 {
		t.Fatalf("consequence is not 1 statements. got=%d\n", len(elseStmt.Body.Statements))
	}
	elseStmtVarStmt, ok := elseStmt.Body.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.VarStatement. got=%T", elseStmt.Body.Statements[0])
	}
	if elseStmtVarStmt.Name.Value != "b" {
		t.Fatalf("varStmt.Name.Value not 'b'. got=%q", elseStmtVarStmt.Name.Value)
	}
	if elseStmtVarStmt.Type.Value != "int" {
		t.Fatalf("varStmt.Type.Value not 'int'. got=%q", elseStmtVarStmt.Type.Value)
	}
}

func ce5(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
