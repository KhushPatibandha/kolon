package tests

import (
	"fmt"
	"testing"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func Test4(t *testing.T) {
	input := `
        var x: int = 10;
        var y: int = 100;
        var foobar: int = 10000;
        const age: int = 100;
        const heh: string = "hello";
    `

	tokens := lexer.Tokenizer(input)
	parser := parser.New(tokens)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 5 {
		t.Fatalf("program.Statements does not contain 5 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       string
	}{
		{"x", "int"},
		{"y", "int"},
		{"foobar", "int"},
		{"age", "int"},
		{"heh", "string"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
}

func Test5(t *testing.T) {
	input := `
        return: 5;
        return: 100;
        return: 312413;
        return: (5 + 1, true, "hello");
    `

	tokens := lexer.Tokenizer(input)
	parser := parser.New(tokens)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenValue() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenValue())
		}
	}
}

func Test6(t *testing.T) {
	input := `
    fun: hehe(name: string, age: int): (bool, int) {
        var a: int = 5;
        return: (true, 5);
    }
    `

	tokens := lexer.Tokenizer(input)
	parser := parser.New(tokens)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)
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

func Test7(t *testing.T) {
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
        return: a + b;
    }
    `

	tokens := lexer.Tokenizer(input)
	parser := parser.New(tokens)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
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

func Test8(t *testing.T) {
	input := `
    if: (a < b): {
        var a: int = 10;
    } else: {
        var b: int = 20;
    }
    `

	tokens := lexer.Tokenizer(input)
	p := parser.New(tokens)
	program := p.ParseProgram()
	checkParserErrors(t, p)
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

func Test9(t *testing.T) {
	input := `foobar;`

	tokens := lexer.Tokenizer(input)
	p := parser.New(tokens)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenValue() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenValue())
	}
}

func Test10(t *testing.T) {
	input := "5;"

	tokens := lexer.Tokenizer(input)
	p := parser.New(tokens)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not parser.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerValue)
	if !ok {
		t.Fatalf("exp not *parser.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Fatalf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenValue() != "5" {
		t.Errorf("literal.TokenValue not %s. got=%s", "5", literal.TokenValue())
	}
}

func Test11(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-5;", "-", 5},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		tokens := lexer.Tokenizer(tt.input)
		p := parser.New(tokens)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testValueExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func Test12(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 % 5;", 5, "%", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 && 5;", 5, "&&", 5},
		{"5 || 5;", 5, "||", 5},
		{"5 & 5;", 5, "&", 5},
		{"5 | 5;", 5, "|", 5},
		{"5 ++ 5", 5, "++", 5}, // should fail, but test will pass. Just to check if it gets confused with postfix operators
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
		// {"5 = 5;", 5, "=", 5},
		// {"5 += 5;", 5, "+=", 5},
		// {"5 -= 5;", 5, "-=", 5},
		// {"5 *= 5;", 5, "*=", 5},
		// {"5 /= 5;", 5, "/=", 5},
		// {"5 %= 5;", 5, "%=", 5},
	}

	for _, tt := range infixTests {
		tokens := lexer.Tokenizer(tt.input)
		p := parser.New(tokens)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			if tt.operator == "++" {
				continue
			}
			t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testValueExpression(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testValueExpression(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func Test13(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.Tokenizer(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func Test17(t *testing.T) {
	postfixTests := []struct {
		input        string
		integerValue int64
		operator     string
	}{
		{"5++", 5, "++"},
		{"5--", 5, "--"},
	}

	for _, tt := range postfixTests {
		l := lexer.Tokenizer(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PostfixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PostfixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Left, tt.integerValue) {
			return
		}
	}
}

func Test19(t *testing.T) {
	input := `
    true;
    false;
    var a: bool = true;
    var b: bool = false;
    `
	tokens := lexer.Tokenizer(input)
	p := parser.New(tokens)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
	}

	firstStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	if !testValueExpression(t, firstStmt.Expression, true) {
		return
	}

	secondStmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.ExpressionStatement. got=%T", program.Statements[1])
	}
	if !testValueExpression(t, secondStmt.Expression, false) {
		return
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       string
		expectedValue      bool
	}{
		{"a", "bool", true},
		{"b", "bool", false},
	}

	for i, tt := range tests {
		stmt := program.Statements[i+2]
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}

}

func testVarStatement(t *testing.T, s ast.Statement, identifier string, typeOfvar string) bool {
	if s.TokenValue() != "var" && s.TokenValue() != "const" {
		t.Errorf("s.TokenValue not 'var' || 'const'. got=%q", s.TokenValue())
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

func testValueExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanValue)
	if !ok {
		t.Errorf("exp not *ast.BooleanValue. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenValue() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenValue() not %t. got=%s", value, bo.TokenValue())
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerValue)
	if !ok {
		t.Errorf("il not *ast.IntegerValue. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenValue() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenValue() not %d. got=%s", value, integ.TokenValue())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenValue() != value {
		t.Errorf("ident.TokenValue() not %s. got=%s", value, ident.TokenValue())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testValueExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testValueExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
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
