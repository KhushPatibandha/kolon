package parser_test

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func Test1(t *testing.T) {
	input := `
        var x: int = 10;
        var y: int = 100;
        var foobar: int = 10000;
    `

	tokens := lexer.Tokenizer(input)
	parser := parser.New(tokens)

	program := parser.ParseProgram()
	ce1(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       string
	}{
		{"x", "int"},
		{"y", "int"},
		{"foobar", "int"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, identifier string, typeOfvar string) bool {
	if s.TokenValue() != "var" {
		t.Errorf("s.TokenValue not 'var'. got=%q", s.TokenValue())
		return false
	}

	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}
	if varStmt.Name.Value != identifier {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", identifier, varStmt.Name.Value)
		return false
	}
	if varStmt.Name.TokenValue() != identifier {
		t.Errorf("varStmt.Name.TokenValue() not '%s'. got=%s", identifier, varStmt.Name.TokenValue())
		return false
	}
	if varStmt.Type.Value != typeOfvar {
		t.Errorf("varStmt.Type.Value not '%s'. got=%s", typeOfvar, varStmt.Type.Value)
		return false
	}
	if varStmt.Type.TokenValue() != typeOfvar {
		t.Errorf("varStmt.Type.TokenValue() not '%s'. got=%s", typeOfvar, varStmt.Type.TokenValue())
		return false
	}
	return true
}

func ce1(t *testing.T, p *parser.Parser) {
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
