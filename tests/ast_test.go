package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

func Test4(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "myVar"},
					Value: "myVar",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Name:        "int",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "anotherVar"},
					Value: "anotherVar",
				},
			},
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.CONST, Value: "const"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "someVar"},
					Value: "someVar",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Name:        "int",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: &ast.Integer{
					Token: lexer.Token{Kind: lexer.INT, Value: "1"},
					Value: 1,
				},
			},
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "someString"},
					Value: "someString",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "string"},
					Name:        "string",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: &ast.String{
					Token: lexer.Token{Kind: lexer.STRING, Value: "\"someString\""},
					Value: "\"someString\"",
				},
			},
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "someFloat"},
					Value: "someFloat",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "float"},
					Name:        "float",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: &ast.Float{
					Token: lexer.Token{Kind: lexer.FLOAT, Value: "1.0"},
					Value: 1.0,
				},
			},
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.CONST, Value: "const"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "someBool"},
					Value: "someBool",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "bool"},
					Name:        "bool",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: &ast.Bool{
					Token: lexer.Token{Kind: lexer.BOOL, Value: "true"},
					Value: true,
				},
			},
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "someChar"},
					Value: "someChar",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "char"},
					Name:        "char",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: &ast.Char{
					Token: lexer.Token{Kind: lexer.CHAR, Value: "'a'"},
					Value: "'a'",
				},
			},
		},
	}
	assert.Equal(t, program.String(), "var myVar: int = anotherVar;const "+
		"someVar: int = 1;var someString: string = \"someString\";var "+
		"someFloat: float = 1.0;const someBool: bool = true;var someChar: char = 'a';")
}

func Test5(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "main"},
					Value: "main",
				},
				Parameters:  nil,
				ReturnTypes: nil,
				Body: &ast.Body{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					Statements: []ast.Statement{
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Integer{
								Token: lexer.Token{Kind: lexer.INT, Value: "10"},
								Value: 10,
							},
						},
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "bool"},
								Name:        "bool",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Bool{
								Token: lexer.Token{Kind: lexer.BOOL, Value: "true"},
								Value: true,
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, program.String(), "fun: main() {var a: int = 10;var b: bool = true;}")
}

func Test6(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "main"},
					Value: "main",
				},
				Parameters:  nil,
				ReturnTypes: nil,
				Body: &ast.Body{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					Statements: []ast.Statement{
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Integer{
								Token: lexer.Token{Kind: lexer.INT, Value: "10"},
								Value: 10,
							},
						},
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Integer{
								Token: lexer.Token{Kind: lexer.INT, Value: "20"},
								Value: 20,
							},
						},
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
								Value: "c",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Integer{
								Token: lexer.Token{Kind: lexer.INT, Value: "30"},
								Value: 30,
							},
						},
						&ast.If{
							Token: lexer.Token{Kind: lexer.IF, Value: "if"},
							Condition: &ast.Infix{
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
							Body: &ast.Body{
								Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
								Statements: []ast.Statement{
									&ast.VarAndConst{
										Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
										Name: &ast.Identifier{
											Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "d"},
											Value: "d",
										},
										Type: &ast.Type{
											Kind:        ast.TypeBase,
											Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
											Name:        "int",
											ElementType: nil,
											KeyType:     nil,
											ValueType:   nil,
										},
										Value: &ast.Integer{
											Token: lexer.Token{Kind: lexer.INT, Value: "40"},
											Value: 40,
										},
									},
								},
							},
							MultiConditionals: []*ast.ElseIf{
								{
									Token: lexer.Token{Kind: lexer.ELSE_IF, Value: "else if"},
									Condition: &ast.Infix{
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
									Body: &ast.Body{
										Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
										Statements: []ast.Statement{
											&ast.VarAndConst{
												Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
												Name: &ast.Identifier{
													Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "e"},
													Value: "e",
												},
												Type: &ast.Type{
													Kind:        ast.TypeBase,
													Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
													Name:        "int",
													ElementType: nil,
													KeyType:     nil,
													ValueType:   nil,
												},
												Value: &ast.Integer{
													Token: lexer.Token{Kind: lexer.INT, Value: "50"},
													Value: 50,
												},
											},
										},
									},
								},
							},
							Alternate: &ast.Else{
								Token: lexer.Token{Kind: lexer.ELSE, Value: "else"},
								Body: &ast.Body{
									Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
									Statements: []ast.Statement{
										&ast.VarAndConst{
											Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
											Name: &ast.Identifier{
												Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "f"},
												Value: "f",
											},
											Type: &ast.Type{
												Kind:        ast.TypeBase,
												Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
												Name:        "int",
												ElementType: nil,
												KeyType:     nil,
												ValueType:   nil,
											},
											Value: &ast.Integer{
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
			},
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "add"},
					Value: "add",
				},
				Parameters: []*ast.FunctionParameter{
					{
						ParameterName: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
							Value: "a",
						},
						ParameterType: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
					},
					{
						ParameterName: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
							Value: "b",
						},
						ParameterType: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
					},
				},
				ReturnTypes: []*ast.Type{
					{
						Kind:        ast.TypeBase,
						Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
						Name:        "int",
						ElementType: nil,
						KeyType:     nil,
						ValueType:   nil,
					},
					{
						Kind:        ast.TypeBase,
						Token:       lexer.Token{Kind: lexer.TYPE, Value: "bool"},
						Name:        "bool",
						ElementType: nil,
						KeyType:     nil,
						ValueType:   nil,
					},
				},
				Body: &ast.Body{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					Statements: []ast.Statement{
						&ast.Return{
							Token: lexer.Token{Kind: lexer.RETURN, Value: "return"},
							Value: []ast.Expression{
								&ast.Infix{
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
								&ast.Bool{
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

	assert.Equal(t, program.String(), "fun: main() {var a: int = 10;var "+
		"b: int = 20;var c: int = 30;if: ((a > b)): {var d: "+
		"int = 40;}else if: ((b > c)): {var e: int = 50;}"+
		"else: {var f: int = 60;}}fun: add(a: int, b: int): (int, bool) {return: ((a + b), true);}")
}

func Test7(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "main"},
					Value: "main",
				},
				Parameters:  nil,
				ReturnTypes: nil,
				Body: &ast.Body{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					Statements: []ast.Statement{
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Integer{
								Token: lexer.Token{Kind: lexer.INT, Value: "10"},
								Value: 10,
							},
						},
						&ast.ExpressionStatement{
							Expression: &ast.Postfix{
								Token:    lexer.Token{Kind: lexer.PLUS_PLUS, Value: "++"},
								Operator: "++",
								Left: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
									Value: "a",
								},
							},
						},
						&ast.ExpressionStatement{
							Expression: &ast.Postfix{
								Token:    lexer.Token{Kind: lexer.MINUS_MINUS, Value: "--"},
								Operator: "--",
								Left: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
									Value: "a",
								},
							},
						},
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Integer{
								Token: lexer.Token{Kind: lexer.INT, Value: "20"},
								Value: 20,
							},
						},
						&ast.If{
							Token: lexer.Token{Kind: lexer.IF, Value: "if"},
							Condition: &ast.Infix{
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
							Body: &ast.Body{
								Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
								Statements: []ast.Statement{
									&ast.VarAndConst{
										Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
										Name: &ast.Identifier{
											Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
											Value: "c",
										},
										Type: &ast.Type{
											Kind:        ast.TypeBase,
											Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
											Name:        "int",
											ElementType: nil,
											KeyType:     nil,
											ValueType:   nil,
										},
										Value: &ast.Infix{
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
							Alternate: &ast.Else{
								Token: lexer.Token{Kind: lexer.ELSE, Value: "else"},
								Body: &ast.Body{
									Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
									Statements: []ast.Statement{
										&ast.VarAndConst{
											Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
											Name: &ast.Identifier{
												Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "d"},
												Value: "d",
											},
											Type: &ast.Type{
												Kind:        ast.TypeBase,
												Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
												Name:        "int",
												ElementType: nil,
												KeyType:     nil,
												ValueType:   nil,
											},
											Value: &ast.Infix{
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
		},
	}

	assert.Equal(t, "fun: main() {var a: int = 10;(a++);(a--);var "+
		"b: int = 20;if: ((a > b)): {var c: int = (a + b);}else: {var d: int = (a - b);}}", program.String())
}

func Test8(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "main"},
					Value: "main",
				},
				Parameters:  nil,
				ReturnTypes: nil,
				Body: &ast.Body{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					Statements: []ast.Statement{
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Integer{
								Token: lexer.Token{Kind: lexer.INT, Value: "10"},
								Value: 10,
							},
						},
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
						},
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
								Value: "c",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: nil,
						},
					},
				},
			},
		},
	}

	assert.Equal(t, "fun: main() {var a: int = 10;var b: int = a;var c: int;}", program.String())
}

func Test9(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.Function{
				Token: lexer.Token{Kind: lexer.FUN, Value: "fun"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "main"},
					Value: "main",
				},
				Parameters:  nil,
				ReturnTypes: nil,
				Body: &ast.Body{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					Statements: []ast.Statement{
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Integer{
								Token: lexer.Token{Kind: lexer.INT, Value: "10"},
								Value: 10,
							},
						},
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
						},
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
								Value: "c",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.CallExpression{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "add"},
								Name: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "add"},
									Value: "add",
								},
								Args: []ast.Expression{
									&ast.Infix{
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
									&ast.Infix{
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
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "d"},
								Value: "d",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "string"},
								Name:        "string",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.String{
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
				Parameters: []*ast.FunctionParameter{
					{
						ParameterName: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
							Value: "a",
						},
						ParameterType: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
					},
					{
						ParameterName: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
							Value: "b",
						},
						ParameterType: &ast.Type{
							Kind:  ast.TypeBase,
							Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},

							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
					},
				},
				ReturnTypes: []*ast.Type{
					{
						Kind:        ast.TypeBase,
						Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
						Name:        "int",
						ElementType: nil,
						KeyType:     nil,
						ValueType:   nil,
					},
				},
				Body: &ast.Body{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					Statements: []ast.Statement{
						&ast.Return{
							Token: lexer.Token{Kind: lexer.RETURN, Value: "return"},
							Value: []ast.Expression{
								&ast.Infix{
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

	assert.Equal(t, "fun: main() {var a: int = "+
		"10;var b: int = a;var c: int = add((a + b), (b * b));var "+
		"d: string = \"Hello\";}fun: add(a: int, b: "+
		"int): (int) {return: (a + b);}", program.String())
}

func Test10(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "x"},
					Value: "x",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Name:        "int",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: &ast.Integer{
					Token: lexer.Token{Kind: lexer.INT, Value: "5"},
					Value: 5,
				},
			},
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "s"},
					Value: "s",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Name:        "int",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: &ast.Postfix{
					Token:    lexer.Token{Kind: lexer.PLUS_PLUS, Value: "++"},
					Operator: "++",
					Left: &ast.Identifier{
						Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "x"},
						Value: "x",
					},
				},
			},
			&ast.ForLoop{
				Token: lexer.Token{Kind: lexer.FOR, Value: "for"},
				Left: &ast.VarAndConst{
					Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
					Name: &ast.Identifier{
						Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "i"},
						Value: "i",
					},
					Type: &ast.Type{
						Kind:        ast.TypeBase,
						Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
						Name:        "int",
						ElementType: nil,
						KeyType:     nil,
						ValueType:   nil,
					},
					Value: &ast.Integer{
						Token: lexer.Token{Kind: lexer.INT, Value: "0"},
						Value: 0,
					},
				},
				Middle: &ast.Infix{
					Token: lexer.Token{Kind: lexer.LESS_THAN, Value: "<"},
					Left: &ast.Identifier{
						Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "i"},
						Value: "i",
					},
					Operator: "<",
					Right: &ast.Integer{
						Token: lexer.Token{Kind: lexer.INT, Value: "10"},
						Value: 10,
					},
				},
				Right: &ast.Postfix{
					Token:    lexer.Token{Kind: lexer.PLUS_PLUS, Value: "++"},
					Operator: "++",
					Left: &ast.Identifier{
						Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "i"},
						Value: "i",
					},
				},
				Body: &ast.Body{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					Statements: []ast.Statement{
						&ast.VarAndConst{
							Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Type: &ast.Type{
								Kind:        ast.TypeBase,
								Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
								Name:        "int",
								ElementType: nil,
								KeyType:     nil,
								ValueType:   nil,
							},
							Value: &ast.Integer{
								Token: lexer.Token{Kind: lexer.INT, Value: "5"},
								Value: 5,
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, "var x: int = 5;var s: int = (x++);for: (var i: int = 0; (i < 10); (i++)): {var a: int = 5;}", program.String())
}

func Test11(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "s"},
					Value: "s",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Name:        "int",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: nil,
			},
			&ast.ExpressionStatement{
				Expression: &ast.Assignment{
					Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
					Left: &ast.Identifier{
						Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "s"},
						Value: "s",
					},
					Operator: "=",
					Right: &ast.Integer{
						Token: lexer.Token{Kind: lexer.INT, Value: "10"},
						Value: 10,
					},
				},
			},
		},
	}

	assert.Equal(t, "var s: int;s = 10;", program.String())
}

func Test12(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.MultiAssignment{
				Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
				Objects: []ast.Statement{
					&ast.VarAndConst{
						Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
						Name: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
							Value: "a",
						},
						Type: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Value: &ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "10"},
							Value: 10,
						},
					},
					&ast.VarAndConst{
						Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
						Name: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
							Value: "b",
						},
						Type: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Value: &ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "20"},
							Value: 20,
						},
					},
				},
			},
		},
	}

	// var a: int, var b: int = 10, 20;
	assert.Equal(t, "var a: int = 10;var b: int = 20;", program.String())
}

func Test13(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "d"},
					Value: "d",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Name:        "int",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: nil,
			},
			&ast.MultiAssignment{
				Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
				Objects: []ast.Statement{
					&ast.VarAndConst{
						Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
						Name: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
							Value: "a",
						},
						Type: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Value: &ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "10"},
							Value: 10,
						},
					},
					&ast.VarAndConst{
						Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
						Name: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
							Value: "b",
						},
						Type: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Value: &ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "20"},
							Value: 20,
						},
					},
					&ast.VarAndConst{
						Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
						Name: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
							Value: "c",
						},
						Type: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "string"},
							Name:        "string",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Value: &ast.String{
							Token: lexer.Token{Kind: lexer.STRING, Value: "\"hello\""},
							Value: "\"hello\"",
						},
					},
					&ast.ExpressionStatement{
						Expression: &ast.Assignment{
							Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
							Left: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "d"},
								Value: "d",
							},
							Operator: "=",
							Right: &ast.Integer{
								Token: lexer.Token{Kind: lexer.INT, Value: "100"},
								Value: 100,
							},
						},
					},
				},
			},
		},
	}

	// var d: int; var a: int, var b: int, var c: string, d = 10, 20, "hello", 100;
	assert.Equal(t, "var d: int;var a: int = 10;var b: int = 20;var c: string = \"hello\";d = 100;", program.String())
}

func Test14(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
					Value: "a",
				},
				Type: &ast.Type{
					Kind:  ast.TypeBase,
					Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},

					Name:        "int",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: nil,
			},
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
					Value: "b",
				},
				Type: &ast.Type{
					Kind:  ast.TypeBase,
					Token: lexer.Token{Kind: lexer.TYPE, Value: "string"},

					Name:        "string",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: nil,
			},
			&ast.MultiAssignment{
				Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
				Objects: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.Assignment{
							Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
							Left: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Operator: "=",
							Right: &ast.Integer{
								Token: lexer.Token{Kind: lexer.INT, Value: "10"},
								Value: 10,
							},
						},
					},
					&ast.ExpressionStatement{
						Expression: &ast.Assignment{
							Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
							Left: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Operator: "=",
							Right: &ast.String{
								Token: lexer.Token{Kind: lexer.STRING, Value: "\"hello\""},
								Value: "\"hello\"",
							},
						},
					},
					&ast.VarAndConst{
						Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
						Name: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
							Value: "c",
						},
						Type: &ast.Type{
							Kind:  ast.TypeBase,
							Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},

							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Value: &ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "20"},
							Value: 20,
						},
					},
				},
			},
		},
	}

	// var a: int; var b: string; a, b, var c: int = 10, "hello", 20;
	assert.Equal(t, "var a: int;var b: string;a = 10;b = \"hello\";var c: int = 20;", program.String())
}

func Test15(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
					Value: "c",
				},
				Type: &ast.Type{
					Kind:  ast.TypeBase,
					Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},

					Name:        "int",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: nil,
			},
			&ast.MultiAssignment{
				Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
				Objects: []ast.Statement{
					&ast.VarAndConst{
						Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
						Name: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
							Value: "a",
						},
						Type: &ast.Type{
							Kind:  ast.TypeBase,
							Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},

							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Value: &ast.CallExpression{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
								Value: "getValues",
							},
							Args: nil,
						},
					},
					&ast.VarAndConst{
						Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
						Name: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
							Value: "b",
						},
						Type: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Value: &ast.CallExpression{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
								Value: "getValues",
							},
							Args: nil,
						},
					},
					&ast.ExpressionStatement{
						Expression: &ast.Assignment{
							Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
							Left: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
								Value: "c",
							},
							Operator: "=",
							Right: &ast.CallExpression{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
								Name: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
									Value: "getValues",
								},
								Args: nil,
							},
						},
					},
				},
			},
		},
	}

	// var c: int; var a: int, var: b: int, c = getValues();
	assert.Equal(t, "var c: int;var a: int = getValues();var b: int = getValues();c = getValues();", program.String())
}

func Test16(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
					Value: "a",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Name:        "int",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: nil,
			},
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
					Value: "b",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Name:        "int",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: nil,
			},
			&ast.MultiAssignment{
				Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
				Objects: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.Assignment{
							Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
							Left: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
								Value: "a",
							},
							Operator: "=",
							Right: &ast.CallExpression{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
								Name: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
									Value: "getValues",
								},
								Args: nil,
							},
						},
					},
					&ast.ExpressionStatement{
						Expression: &ast.Assignment{
							Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
							Left: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
								Value: "b",
							},
							Operator: "=",
							Right: &ast.CallExpression{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
								Name: &ast.Identifier{
									Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
									Value: "getValues",
								},
								Args: nil,
							},
						},
					},
					&ast.VarAndConst{
						Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
						Name: &ast.Identifier{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "c"},
							Value: "c",
						},
						Type: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Value: &ast.CallExpression{
							Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
							Name: &ast.Identifier{
								Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "getValues"},
								Value: "getValues",
							},
							Args: nil,
						},
					},
				},
			},
		},
	}

	// var a: int; var b: int; a, b, var c: int = getValues();
	assert.Equal(t, "var a: int;var b: int;a = getValues();b = getValues();var c: int = getValues();", program.String())
}

func Test17(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
					Value: "a",
				},
				Type: &ast.Type{
					Kind:  ast.TypeArray,
					Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
					ElementType: &ast.Type{
						Kind:        ast.TypeBase,
						Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
						Name:        "int",
						ElementType: nil,
						KeyType:     nil,
						ValueType:   nil,
					},
					Name:      "",
					KeyType:   nil,
					ValueType: nil,
				},
				Value: &ast.Array{
					Token: lexer.Token{Kind: lexer.OPEN_SQUARE_BRACKET, Value: "["},
					Values: []ast.Expression{
						&ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "1"},
							Value: 1,
						},
						&ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "2"},
							Value: 2,
						},
						&ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "3"},
							Value: 3,
						},
						&ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "4"},
							Value: 4,
						},
						&ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "5"},
							Value: 5,
						},
					},
				},
			},
		},
	}

	assert.Equal(t, "var a: int[] = [1, 2, 3, 4, 5];", program.String())
}

func Test18(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
					Value: "a",
				},
				Type: &ast.Type{
					Kind:  ast.TypeArray,
					Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
					ElementType: &ast.Type{
						Kind:        ast.TypeBase,
						Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
						Name:        "int",
						ElementType: nil,
						KeyType:     nil,
						ValueType:   nil,
					},
					Name:      "",
					KeyType:   nil,
					ValueType: nil,
				},
				Value: &ast.Array{
					Token: lexer.Token{Kind: lexer.OPEN_SQUARE_BRACKET, Value: "["},
					Values: []ast.Expression{
						&ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "1"},
							Value: 1,
						},
						&ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "2"},
							Value: 2,
						},
						&ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "3"},
							Value: 3,
						},
					},
				},
			},
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "b"},
					Value: "b",
				},
				Type: &ast.Type{
					Kind:        ast.TypeBase,
					Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
					Name:        "int",
					ElementType: nil,
					KeyType:     nil,
					ValueType:   nil,
				},
				Value: &ast.IndexExpression{
					Token: lexer.Token{Kind: lexer.OPEN_SQUARE_BRACKET, Value: "["},
					Left: &ast.Identifier{
						Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
						Value: "a",
					},
					Index: &ast.Integer{
						Token: lexer.Token{Kind: lexer.INT, Value: "0"},
						Value: 0,
					},
				},
			},
		},
	}

	assert.Equal(t, "var a: int[] = [1, 2, 3];var b: int = a[0];", program.String())
}

func Test19(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
					Value: "a",
				},
				Type: &ast.Type{
					Kind:  ast.TypeHashMap,
					Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
					KeyType: &ast.Type{
						Kind:        ast.TypeBase,
						Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
						Name:        "int",
						ElementType: nil,
						KeyType:     nil,
						ValueType:   nil,
					},
					ValueType: &ast.Type{
						Kind:        ast.TypeBase,
						Token:       lexer.Token{Kind: lexer.TYPE, Value: "bool"},
						Name:        "bool",
						ElementType: nil,
						KeyType:     nil,
						ValueType:   nil,
					},
					ElementType: nil,
					Name:        "",
				},
				Value: &ast.HashMap{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					KeyType: &ast.Type{
						Kind:        ast.TypeBase,
						Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
						Name:        "int",
						ElementType: nil,
						KeyType:     nil,
						ValueType:   nil,
					},
					ValueType: &ast.Type{
						Kind:        ast.TypeBase,
						Token:       lexer.Token{Kind: lexer.TYPE, Value: "bool"},
						Name:        "bool",
						ElementType: nil,
						KeyType:     nil,
						ValueType:   nil,
					},
					Pairs: map[ast.BaseType]ast.Expression{
						&ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "1"},
							Value: 1,
						}: &ast.Bool{
							Token: lexer.Token{Kind: lexer.BOOL, Value: "true"},
							Value: true,
						},
						&ast.Integer{
							Token: lexer.Token{Kind: lexer.INT, Value: "2"},
							Value: 2,
						}: &ast.Bool{
							Token: lexer.Token{Kind: lexer.BOOL, Value: "false"},
							Value: false,
						},
					},
				},
			},
			&ast.WhileLoop{
				Token: lexer.Token{Kind: lexer.WHILE, Value: "while"},
				Condition: &ast.Bool{
					Token: lexer.Token{Kind: lexer.BOOL, Value: "true"},
					Value: true,
				},
				Body: &ast.Body{
					Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
					Statements: []ast.Statement{
						&ast.If{
							Token: lexer.Token{Kind: lexer.IF, Value: "if"},
							Condition: &ast.Bool{
								Token: lexer.Token{Kind: lexer.BOOL, Value: "true"},
								Value: true,
							},
							Body: &ast.Body{
								Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
								Statements: []ast.Statement{
									&ast.Break{
										Token: lexer.Token{Kind: lexer.BREAK, Value: "break"},
									},
								},
							},
							MultiConditionals: nil,
							Alternate: &ast.Else{
								Token: lexer.Token{Kind: lexer.ELSE, Value: "else"},
								Body: &ast.Body{
									Token: lexer.Token{Kind: lexer.OPEN_CURLY_BRACKET, Value: "{"},
									Statements: []ast.Statement{
										&ast.Continue{
											Token: lexer.Token{Kind: lexer.CONTINUE, Value: "continue"},
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
	assert.Equal(t, "var a: int[bool] = {1: true, 2: false};while: (true): {if: (true): {break;}else: {continue;}}", program.String())
}

func Test20(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
					Value: "a",
				},
				Type: &ast.Type{
					Kind:  ast.TypeHashMap,
					Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
					KeyType: &ast.Type{
						Kind:        ast.TypeBase,
						Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
						Name:        "int",
						ElementType: nil,
						KeyType:     nil,
						ValueType:   nil,
					},
					ValueType: &ast.Type{
						Kind:  ast.TypeArray,
						Token: lexer.Token{Kind: lexer.TYPE, Value: "string"},
						ElementType: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "string"},
							Name:        "string",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Name:      "",
						KeyType:   nil,
						ValueType: nil,
					},
					Name:        "",
					ElementType: nil,
				},
				Value: nil,
			},
		},
	}
	assert.Equal(t, "var a: int[string[]];", program.String())
}

func Test21(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarAndConst{
				Token: lexer.Token{Kind: lexer.VAR, Value: "var"},
				Name: &ast.Identifier{
					Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: "a"},
					Value: "a",
				},
				Type: &ast.Type{
					Kind:  ast.TypeHashMap,
					Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
					KeyType: &ast.Type{
						Kind:  ast.TypeArray,
						Token: lexer.Token{Kind: lexer.TYPE, Value: "int"},
						ElementType: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "int"},
							Name:        "int",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Name:      "",
						KeyType:   nil,
						ValueType: nil,
					},
					ValueType: &ast.Type{
						Kind:  ast.TypeArray,
						Token: lexer.Token{Kind: lexer.TYPE, Value: "string"},
						ElementType: &ast.Type{
							Kind:        ast.TypeBase,
							Token:       lexer.Token{Kind: lexer.TYPE, Value: "string"},
							Name:        "string",
							ElementType: nil,
							KeyType:     nil,
							ValueType:   nil,
						},
						Name:      "",
						KeyType:   nil,
						ValueType: nil,
					},
					Name:        "",
					ElementType: nil,
				},
				Value: nil,
			},
		},
	}
	assert.Equal(t, "var a: int[][string[]];", program.String())
}

func Test22(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.Return{
				Token: lexer.Token{Kind: lexer.RETURN, Value: "return"},
				Value: nil,
			},
		},
	}
	assert.Equal(t, "return;", program.String())
}
