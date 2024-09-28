package tests

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

func Test14(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarStatement{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "myVar"},
					Value: "myVar",
				},
				Type: &ast.Type{
					Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Value: "int",
				},
				Value: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "anotherVar"},
					Value: "anotherVar",
				},
			},
			&ast.VarStatement{
				Token: lexer.Token{Kind: lexer.CONST, Value: "const"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "someVar"},
					Value: "someVar",
				},
				Type: &ast.Type{
					Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Value: "int",
				},
				Value: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.INT, Value: "1"},
					Value: "1",
				},
			},
			&ast.VarStatement{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "someString"},
					Value: "someString",
				},
				Type: &ast.Type{
					Token: lexer.Token{Kind: lexer.TYPE, Value: "string"},
					Value: "string",
				},
				Value: &ast.StringValue{
					Token: lexer.Token{Kind: lexer.STRING, Value: "\"someString\""},
					Value: "\"someString\"",
				},
			},
			&ast.VarStatement{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "someFloat"},
					Value: "someFloat",
				},
				Type: &ast.Type{
					Token: lexer.Token{Kind: lexer.TYPE, Value: "float"},
					Value: "float",
				},
				Value: &ast.FloatValue{
					Token: lexer.Token{Kind: lexer.FLOAT, Value: "1.0"},
					Value: 1.0,
				},
			},
			&ast.VarStatement{
				Token: lexer.Token{Kind: lexer.CONST, Value: "const"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "someBool"},
					Value: "someBool",
				},
				Type: &ast.Type{
					Token: lexer.Token{Kind: lexer.TYPE, Value: "bool"},
					Value: "bool",
				},
				Value: &ast.BooleanValue{
					Token: lexer.Token{Kind: lexer.BOOL, Value: "true"},
					Value: true,
				},
			},
			&ast.VarStatement{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "someChar"},
					Value: "someChar",
				},
				Type: &ast.Type{
					Token: lexer.Token{Kind: lexer.TYPE, Value: "char"},
					Value: "char",
				},
				Value: &ast.CharValue{
					Token: lexer.Token{Kind: lexer.CHAR, Value: "'a'"},
					Value: "'a'",
				},
			},
		},
	}

	if program.String() != "var myVar: int = anotherVar;const someVar: int = 1;var someString: string = \"someString\";var someFloat: float = 1.0;const someBool: bool = true;var someChar: char = 'a';" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func Test15(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "main"},
					Value: "main",
				},
				Parameters: []*ast.FunctionParameters{},
				ReturnType: nil,
				Body: &ast.FunctionBody{
					Statements: []ast.Statement{
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.IntegerValue{
								Token: lexer.Token{Kind: lexer.INT, Value: "10"},
								Value: 10,
							},
						},
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "bool"},
								Value: "bool",
							},
							Value: &ast.BooleanValue{
								Token: lexer.Token{Kind: lexer.BOOL, Value: "true"},
								Value: true,
							},
						},
					},
				},
			},
		},
	}

	if program.String() != "fun: main() {var a: int = 10;var b: bool = true;}" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func Test16(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "main"},
					Value: "main",
				},
				Parameters: []*ast.FunctionParameters{},
				ReturnType: nil,
				Body: &ast.FunctionBody{
					Statements: []ast.Statement{
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.IntegerValue{
								Token: lexer.Token{Kind: lexer.INT, Value: "10"},
								Value: 10,
							},
						},
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.IntegerValue{
								Token: lexer.Token{Kind: lexer.INT, Value: "20"},
								Value: 20,
							},
						},
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
								Value: "c",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.IntegerValue{
								Token: lexer.Token{Kind: lexer.INT, Value: "30"},
								Value: 30,
							},
						},
						&ast.IfStatement{
							Token: lexer.Token{Kind: lexer.IF, Value: "if"},
							Value: &ast.InfixExpression{
								Token: lexer.Token{Kind: lexer.GREATER_THAN, Value: ">"},
								Left: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
									Value: "a",
								},
								Operator: ">",
								Right: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
									Value: "b",
								},
							},
							Body: &ast.FunctionBody{
								Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
								Statements: []ast.Statement{
									&ast.VarStatement{
										Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
										Name: &ast.Identifier{
											Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "d"},
											Value: "d",
										},
										Type: &ast.Type{
											Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
											Value: "int",
										},
										Value: &ast.IntegerValue{
											Token: lexer.Token{Kind: lexer.INT, Value: "40"},
											Value: 40,
										},
									},
								},
							},
						},
						&ast.ElseIfStatement{
							Token: lexer.Token{Kind: lexer.ELSE_IF, Value: "else if"},
							Value: &ast.InfixExpression{
								Token: lexer.Token{Kind: lexer.GREATER_THAN, Value: ">"},
								Left: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
									Value: "b",
								},
								Operator: ">",
								Right: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
									Value: "c",
								},
							},
							Body: &ast.FunctionBody{
								Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
								Statements: []ast.Statement{
									&ast.VarStatement{
										Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
										Name: &ast.Identifier{
											Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "e"},
											Value: "e",
										},
										Type: &ast.Type{
											Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
											Value: "int",
										},
										Value: &ast.IntegerValue{
											Token: lexer.Token{Kind: lexer.INT, Value: "50"},
											Value: 50,
										},
									},
								},
							},
						},
						&ast.ElseStatement{
							Token: lexer.Token{Kind: lexer.ELSE, Value: "else"},
							Body: &ast.FunctionBody{
								Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
								Statements: []ast.Statement{
									&ast.VarStatement{
										Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
										Name: &ast.Identifier{
											Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "f"},
											Value: "f",
										},
										Type: &ast.Type{
											Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
											Value: "int",
										},
										Value: &ast.IntegerValue{
											Token: lexer.Token{Kind: lexer.INT, Value: "60"},
											Value: 60,
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "add"},
					Value: "add",
				},
				Parameters: []*ast.FunctionParameters{
					{
						ParameterName: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
							Value: "a",
						},
						ParameterType: &ast.Type{
							Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Value: "int",
						},
					},
					{
						ParameterName: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
							Value: "b",
						},
						ParameterType: &ast.Type{
							Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Value: "int",
						},
					},
				},
				ReturnType: []*ast.FunctionReturnType{
					{
						ReturnType: &ast.Type{
							Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Value: "int",
						},
					},
					{
						ReturnType: &ast.Type{
							Token: lexer.Token{Kind: lexer.TYPE, Value: "bool"},
							Value: "bool",
						},
					},
				},
				Body: &ast.FunctionBody{
					Statements: []ast.Statement{
						&ast.ReturnStatement{
							Token: lexer.Token{Kind: lexer.RETURN, Value: "return"},
							Value: []ast.Expression{
								&ast.InfixExpression{
									Token: lexer.Token{Kind: lexer.PLUS, Value: "+"},
									Left: &ast.Identifier{
										Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
										Value: "a",
									},
									Operator: "+",
									Right: &ast.Identifier{
										Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
										Value: "b",
									},
								},
								&ast.BooleanValue{
									Token: lexer.Token{Kind: lexer.BOOL, Value: "true"},
									Value: true,
								},
							},
						},
					},
				},
			},
		},
	}

	if program.String() != "fun: main() {var a: int = 10;var b: int = 20;var c: int = 30;if: ((a > b)): {var d: int = 40;}else if: ((b > c)): {var e: int = 50;}else: {var f: int = 60;}}fun: add(a: int, b: int): (int, bool) {return: ((a + b), true);}" {
		t.Errorf("program.String() wrong. got=%q", program.String())
		t.Errorf("Expected: fun: main() {var a: int = 10;var b: int = 20;var c: int = 30;if: ((a > b)): {var d: int = 40;} else if: ((b > c)): {var e: int = 50;} else: {var f: int = 60;}}fun: add(a: int, b: int): (int, bool) {return: ((a + b), true);}")
	}
}

func Test18(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "main"},
					Value: "main",
				},
				Parameters: []*ast.FunctionParameters{},
				ReturnType: nil,
				Body: &ast.FunctionBody{
					Statements: []ast.Statement{
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.IntegerValue{
								Token: lexer.Token{Kind: lexer.INT, Value: "10"},
								Value: 10,
							},
						},
						&ast.PostfixExpression{
							Token:    lexer.Token{Kind: lexer.PLUS_PLUS, Value: "++"},
							Operator: "++",
							Left: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							IsStmt: true,
						},
						&ast.PostfixExpression{
							Token:    lexer.Token{Kind: lexer.MINUS_MINUS, Value: "--"},
							Operator: "--",
							Left: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							IsStmt: true,
						},
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.IntegerValue{
								Token: lexer.Token{Kind: lexer.INT, Value: "20"},
								Value: 20,
							},
						},
						&ast.IfStatement{
							Token: lexer.Token{Kind: lexer.IF, Value: "if"},
							Value: &ast.InfixExpression{
								Token: lexer.Token{Kind: lexer.GREATER_THAN, Value: ">"},
								Left: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
									Value: "a",
								},
								Operator: ">",
								Right: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
									Value: "b",
								},
							},
							Body: &ast.FunctionBody{
								Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
								Statements: []ast.Statement{
									&ast.VarStatement{
										Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
										Name: &ast.Identifier{
											Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
											Value: "c",
										},
										Type: &ast.Type{
											Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
											Value: "int",
										},
										Value: &ast.InfixExpression{
											Token: lexer.Token{Kind: lexer.PLUS, Value: "+"},
											Left: &ast.Identifier{
												Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
												Value: "a",
											},
											Operator: "+",
											Right: &ast.Identifier{
												Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
												Value: "b",
											},
										},
									},
								},
							},
						},
						&ast.ElseStatement{
							Token: lexer.Token{Kind: lexer.ELSE, Value: "else"},
							Body: &ast.FunctionBody{
								Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
								Statements: []ast.Statement{
									&ast.VarStatement{
										Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
										Name: &ast.Identifier{
											Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "d"},
											Value: "d",
										},
										Type: &ast.Type{
											Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
											Value: "int",
										},
										Value: &ast.InfixExpression{
											Token: lexer.Token{Kind: lexer.DASH, Value: "-"},
											Left: &ast.Identifier{
												Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
												Value: "a",
											},
											Operator: "-",
											Right: &ast.Identifier{
												Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
												Value: "b",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if program.String() != "fun: main() {var a: int = 10;(a++);(a--);var b: int = 20;if: ((a > b)): {var c: int = (a + b);}else: {var d: int = (a - b);}}" {
		t.Errorf("program.String() wrong. got=%q", program.String())
		t.Errorf("Expected: fun: main() {var a: int = 10;(a++);(a--);var b: int = 20;if: ((a > b)): {var c: int = (a + b);}else: {var d: int = (a - b);}}")
	}
}

func Test21(t *testing.T) {

	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "main"},
					Value: "main",
				},
				Parameters: []*ast.FunctionParameters{},
				ReturnType: nil,
				Body: &ast.FunctionBody{
					Statements: []ast.Statement{
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.IntegerValue{
								Token: lexer.Token{Kind: lexer.INT, Value: "10"},
								Value: 10,
							},
						},
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
						},
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
								Value: "c",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
						},
					},
				},
			},
		},
	}

	if program.String() != "fun: main() {var a: int = 10;var b: int = a;var c: int;}" {
		t.Errorf("program.String() wrong. got=%q", program.String())
		t.Errorf("Expected: fun: main() {var a: int = 10;var b: int = a;var c: int;}")
	}
}

func Test22(t *testing.T) {

	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "main"},
					Value: "main",
				},
				Parameters: []*ast.FunctionParameters{},
				ReturnType: nil,
				Body: &ast.FunctionBody{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					Statements: []ast.Statement{
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.IntegerValue{
								Token: lexer.Token{Kind: lexer.INT, Value: "10"},
								Value: 10,
							},
						},
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
						},
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
								Value: "c",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.CallExpression{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "add"},
								Name: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "add"},
									Value: "add",
								},
								Args: []ast.Expression{
									&ast.InfixExpression{
										Token: lexer.Token{Kind: lexer.PLUS, Value: "+"},
										Left: &ast.Identifier{
											Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
											Value: "a",
										},
										Operator: "+",
										Right: &ast.Identifier{
											Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
											Value: "b",
										},
									},
									&ast.InfixExpression{
										Token: lexer.Token{Kind: lexer.STAR, Value: "*"},
										Left: &ast.Identifier{
											Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
											Value: "b",
										},
										Operator: "*",
										Right: &ast.Identifier{
											Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
											Value: "b",
										},
									},
								},
							},
						},
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "d"},
								Value: "d",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "string"},
								Value: "string",
							},
							Value: &ast.StringValue{
								Token: lexer.Token{Kind: lexer.STRING, Value: "\"Hello\""},
								Value: "\"Hello\"",
							},
						},
					},
				},
			},
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "add"},
					Value: "add",
				},
				Parameters: []*ast.FunctionParameters{
					{
						ParameterName: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
							Value: "a",
						},
						ParameterType: &ast.Type{
							Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Value: "int",
						},
					},
					{
						ParameterName: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
							Value: "b",
						},
						ParameterType: &ast.Type{
							Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Value: "int",
						},
					},
				},
				ReturnType: []*ast.FunctionReturnType{
					{
						ReturnType: &ast.Type{
							Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Value: "int",
						},
					},
				},
				Body: &ast.FunctionBody{
					Statements: []ast.Statement{
						&ast.ReturnStatement{
							Token: lexer.Token{Kind: lexer.RETURN, Value: "return"},
							Value: []ast.Expression{
								&ast.InfixExpression{
									Token: lexer.Token{Kind: lexer.PLUS, Value: "+"},
									Left: &ast.Identifier{
										Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
										Value: "a",
									},
									Operator: "+",
									Right: &ast.Identifier{
										Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
										Value: "b",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if program.String() != "fun: main() {var a: int = 10;var b: int = a;var c: int = add((a + b), (b * b));var d: string = \"Hello\";}fun: add(a: int, b: int): (int) {return: (a + b);}" {
		t.Errorf("program.String() wrong. got=%q", program.String())
		t.Errorf("Expected: fun: main() {var a: int = 10;var b: int = a;var c: int = add((a + b), (b * b));var d: string = \"Hello\";}fun: add(a: int, b: int): (int) {return: (a + b);}")
	}
}

func Test23(t *testing.T) {

	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarStatement{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "x"},
					Value: "x",
				},
				Type: &ast.Type{
					Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Value: "int",
				},
				Value: &ast.IntegerValue{
					Token: lexer.Token{Kind: lexer.INT, Value: "5"},
					Value: 5,
				},
			},
			&ast.VarStatement{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "s"},
					Value: "s",
				},
				Type: &ast.Type{
					Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Value: "int",
				},
				Value: &ast.PostfixExpression{
					Token:    lexer.Token{Kind: lexer.PLUS_PLUS, Value: "++"},
					Operator: "++",
					Left: &ast.Identifier{
						Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "x"},
						Value: "x",
					},
					IsStmt: false,
				},
			},
			&ast.ForLoopStatement{
				Token: lexer.Token{Kind: lexer.FOR, Value: "for"},
				Left: &ast.VarStatement{
					Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
					Name: &ast.Identifier{
						Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "i"},
						Value: "i",
					},
					Type: &ast.Type{
						Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
						Value: "int",
					},
					Value: &ast.IntegerValue{
						Token: lexer.Token{Kind: lexer.INT, Value: "0"},
						Value: 0,
					},
				},
				Middle: &ast.InfixExpression{
					Token: lexer.Token{Kind: lexer.LESS_THAN, Value: "<"},
					Left: &ast.Identifier{
						Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "i"},
						Value: "i",
					},
					Operator: "<",
					Right: &ast.IntegerValue{
						Token: lexer.Token{Kind: lexer.INT, Value: "10"},
						Value: 10,
					},
				},
				Right: &ast.PostfixExpression{
					Token:    lexer.Token{Kind: lexer.PLUS_PLUS, Value: "++"},
					Operator: "++",
					Left: &ast.Identifier{
						Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "i"},
						Value: "i",
					},
					IsStmt: false,
				},
				Body: &ast.FunctionBody{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					Statements: []ast.Statement{
						&ast.VarStatement{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Value: "int",
							},
							Value: &ast.IntegerValue{
								Token: lexer.Token{Kind: lexer.INT, Value: "5"},
								Value: 5,
							},
						},
					},
				},
			},
		},
	}

	if program.String() != "var x: int = 5;var s: int = (x++);for: (var i: int = 0; (i < 10); (i++)): {var a: int = 5;}" {
		t.Errorf("program.String() wrong. got=%q", program.String())
		t.Errorf("Expected: for: (var i: int = 0; (i < 10); (i++)): {var a: int = 5;}")
	}
}
