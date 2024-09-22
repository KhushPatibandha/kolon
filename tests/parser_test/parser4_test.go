package parser_test

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func Test4(t *testing.T) {
	input := `
    fun: main() {
        var a: int = 10;
        var b: int = 20;
        var c: int = 30;

        if: (a > b): {
            var d: int = 40;
        } else if: (b > c): {
            var e: int = 50;
        } else: {
            var f: int = 60;
        }
    }
    fun: add(a: int, b: int): (int) {
        return a + b;
    }
    `

	tokens := lexer.Tokenizer(input)
	parser := parser.New(tokens)
	program := parser.ParseProgram()
	ce4(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	fun1, ok := program.Statements[0].(*ast.Function)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.Function. got=%T", program.Statements[0])
	}
	if len(fun1.Body.Statements) != 6 {
		t.Fatalf("expected 6 statements in function body, got %d", len(fun1.Body.Statements))
	}

	fun1VarStmt1, ok := fun1.Body.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not *ast.VarStatement. got=%T", fun1.Body.Statements[0])
	}
	if fun1VarStmt1.Name.Value != "a" {
		t.Fatalf("expected variable name 'a', got %s", fun1VarStmt1.Name.Value)
	}
	if fun1VarStmt1.Type.Value != "int" {
		t.Fatalf("expected variable type 'int', got %s", fun1VarStmt1.Type.Value)
	}

	fun1VarStmt2, ok := fun1.Body.Statements[1].(*ast.VarStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[1] is not *ast.VarStatement. got=%T", fun1.Body.Statements[1])
	}
	if fun1VarStmt2.Name.Value != "b" {
		t.Fatalf("expected variable name 'b', got %s", fun1VarStmt2.Name.Value)
	}
	if fun1VarStmt2.Type.Value != "int" {
		t.Fatalf("expected variable type 'int', got %s", fun1VarStmt2.Type.Value)
	}

	fun1VarStmt3, ok := fun1.Body.Statements[2].(*ast.VarStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[2] is not *ast.VarStatement. got=%T", fun1.Body.Statements[2])
	}
	if fun1VarStmt3.Name.Value != "c" {
		t.Fatalf("expected variable name 'c', got %s", fun1VarStmt3.Name.Value)
	}
	if fun1VarStmt3.Type.Value != "int" {
		t.Fatalf("expected variable type 'int', got %s", fun1VarStmt3.Type.Value)
	}

	fun1IfStmt, ok := fun1.Body.Statements[3].(*ast.IfStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[3] is not *ast.IfStatement. got=%T", fun1.Body.Statements[3])
	}
	if fun1IfStmt.Token.Value != "if" {
		t.Fatalf("expected if statement, got %s", fun1IfStmt.Token.Value)
	}
	if fun1IfStmt.Body.Statements[0].(*ast.VarStatement).Name.Value != "d" {
		t.Fatalf("expected variable name 'd', got %s", fun1IfStmt.Body.Statements[0].(*ast.VarStatement).Name.Value)
	}
	if fun1IfStmt.Body.Statements[0].(*ast.VarStatement).Type.Value != "int" {
		t.Fatalf("expected variable type 'int', got %s", fun1IfStmt.Body.Statements[0].(*ast.VarStatement).Type.Value)
	}

	fun1ElseIfStmt, ok := fun1.Body.Statements[4].(*ast.ElseIfStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[4] is not *ast.IfStatement. got=%T", fun1.Body.Statements[4])
	}
	if fun1ElseIfStmt.Token.Value != "else if" {
		t.Fatalf("expected else if statement, got %s", fun1ElseIfStmt.Token.Value)
	}
	if fun1ElseIfStmt.Body.Statements[0].(*ast.VarStatement).Name.Value != "e" {
		t.Fatalf("expected variable name 'e', got %s", fun1ElseIfStmt.Body.Statements[0].(*ast.VarStatement).Name.Value)
	}
	if fun1ElseIfStmt.Body.Statements[0].(*ast.VarStatement).Type.Value != "int" {
		t.Fatalf("expected variable type 'int', got %s", fun1ElseIfStmt.Body.Statements[0].(*ast.VarStatement).Type.Value)
	}

	fun1ElseStmt, ok := fun1.Body.Statements[5].(*ast.ElseStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[5] is not *ast.IfStatement. got=%T", fun1.Body.Statements[5])
	}
	if fun1ElseStmt.Token.Value != "else" {
		t.Fatalf("expected else statement, got %s", fun1ElseStmt.Token.Value)
	}
	if fun1ElseStmt.Body.Statements[0].(*ast.VarStatement).Name.Value != "f" {
		t.Fatalf("expected variable name 'f', got %s", fun1ElseStmt.Body.Statements[0].(*ast.VarStatement).Name.Value)
	}
	if fun1ElseStmt.Body.Statements[0].(*ast.VarStatement).Type.Value != "int" {
		t.Fatalf("expected variable type 'int', got %s", fun1ElseStmt.Body.Statements[0].(*ast.VarStatement).Type.Value)
	}

	fun2, ok := program.Statements[1].(*ast.Function)
	if !ok {
		t.Fatalf("program.Statements[1] is not *ast.Function. got=%T", program.Statements[1])
	}
	if len(fun2.Body.Statements) != 1 {
		t.Fatalf("expected 1 statement in function body, got %d", len(fun2.Body.Statements))
	}

	_, ok = fun2.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not *ast.ReturnStatement. got=%T", fun2.Body.Statements[0])
	}
}

func ce4(t *testing.T, p *parser.Parser) {
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
