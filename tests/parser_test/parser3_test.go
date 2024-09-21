package parser_test

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func Test3(t *testing.T) {
	input := `
    fun: hehe(name: string, age: int): (bool, int) {
        var a: int = 5;
        return true, 5;
    }
    `

	tokens := lexer.Tokenizer(input)
	parser := parser.New(tokens)

	program := parser.ParseProgram()
	ce3(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	functionStmt, ok := program.Statements[0].(*ast.Function)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.Function. got=%T", program.Statements[0])
	}

	if len(functionStmt.Body.Statements) != 2 {
		t.Fatalf("expected 2 statements in function body, got %d", len(functionStmt.Body.Statements))
	}

	varStmt, ok := functionStmt.Body.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not *ast.VarStatement. got=%T", functionStmt.Body.Statements[0])
	}
	if varStmt.Name.Value != "a" {
		t.Fatalf("expected variable name 'a', got %s", varStmt.Name.Value)
	}
	if varStmt.Type.Value != "int" {
		t.Fatalf("expected variable type 'int', got %s", varStmt.Type.Value)
	}

	_, ok = functionStmt.Body.Statements[1].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[1] is not *ast.ReturnStatement. got=%T", functionStmt.Body.Statements[1])
	}
}

func ce3(t *testing.T, p *parser.Parser) {
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
