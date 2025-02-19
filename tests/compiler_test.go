package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/compiler/code"
	"github.com/KhushPatibandha/Kolon/src/compiler/compiler"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/object"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

type compileTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func Test49(t *testing.T) {
	tests := []compileTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1.1 + 2.1",
			expectedConstants: []interface{}{1.1, 2.1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 / 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 % 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMod),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 & 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAnd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 | 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpOr),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1; 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 >= 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThanEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1.1 < 2.1",
			expectedConstants: []interface{}{2.1, 1.1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1.1 <= 2.1",
			expectedConstants: []interface{}{2.1, 1.1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThanEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1.1 == 2.1",
			expectedConstants: []interface{}{1.1, 2.1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqualEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1.1 != 2.1",
			expectedConstants: []interface{}{1.1, 2.1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true && false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpAndAnd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true || false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpOrOr),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "\"someString\"",
			expectedConstants: []interface{}{"\"someString\""},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `"some" + "thing"`,
			expectedConstants: []interface{}{"\"some\"", "\"thing\""},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "'c'",
			expectedConstants: []interface{}{"'c'"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},

		{
			input:             "-1.12",
			expectedConstants: []interface{}{1.12},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpNot),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "10++",
			expectedConstants: []interface{}{10},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPlusPlus),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "10.1--",
			expectedConstants: []interface{}{10.1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinusMinus),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
            if: (true): {
                10;
            }
            1000;
            `,
			expectedConstants: []interface{}{10, 1000},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),           // 0000
				code.Make(code.OpJumpNotTrue, 8), // 0001
				code.Make(code.OpConstant, 0),    // 0004
				code.Make(code.OpPop),            // 0007
				code.Make(code.OpConstant, 1),    // 0008
				code.Make(code.OpPop),            // 0011
			},
		},
		{
			input: `
            if: (true): {
                10;
            } else: {
                20;
            }
            1000;
            `,
			expectedConstants: []interface{}{10, 20, 1000},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),            // 0000
				code.Make(code.OpJumpNotTrue, 11), // 0001
				code.Make(code.OpConstant, 0),     // 0004
				code.Make(code.OpPop),             // 0007
				code.Make(code.OpJump, 15),        // 0008
				code.Make(code.OpConstant, 1),     // 0011
				code.Make(code.OpPop),             // 0014
				code.Make(code.OpConstant, 2),     // 0015
				code.Make(code.OpPop),             // 0018
			},
		},
		{
			input: `
            if: (true): {
                10;
            } else if: (false):  {
                20;
            } else if: (false): {
                30;
            } else: {
                40;
            }
            1000;
            `,
			expectedConstants: []interface{}{10, 20, 30, 40, 1000},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),            // 0000
				code.Make(code.OpJumpNotTrue, 11), // 0001
				code.Make(code.OpConstant, 0),     // 0004
				code.Make(code.OpPop),             // 0007
				code.Make(code.OpJump, 37),        // 0008

				code.Make(code.OpFalse),           // 0011
				code.Make(code.OpJumpNotTrue, 22), // 0012
				code.Make(code.OpConstant, 1),     // 0015
				code.Make(code.OpPop),             // 0018
				code.Make(code.OpJump, 37),        // 0019

				code.Make(code.OpFalse),           // 0022
				code.Make(code.OpJumpNotTrue, 33), // 0023
				code.Make(code.OpConstant, 2),     // 0026
				code.Make(code.OpPop),             // 0029
				code.Make(code.OpJump, 37),        // 0030

				code.Make(code.OpConstant, 3), // 0033
				code.Make(code.OpPop),         // 0036

				code.Make(code.OpConstant, 4), // 0037
				code.Make(code.OpPop),         // 0040
			},
		},
		{
			input: `
            if: (true): {
                10;
            } else if: (true): {
                20;
            } else if: (true): {
                30;
            }
            1000;
            `,
			expectedConstants: []interface{}{10, 20, 30, 1000},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),            // 0000
				code.Make(code.OpJumpNotTrue, 11), // 0001
				code.Make(code.OpConstant, 0),     // 0004
				code.Make(code.OpPop),             // 0007
				code.Make(code.OpJump, 33),        // 0008

				code.Make(code.OpTrue),            // 0011
				code.Make(code.OpJumpNotTrue, 22), // 0012
				code.Make(code.OpConstant, 1),     // 0015
				code.Make(code.OpPop),             // 0018
				code.Make(code.OpJump, 33),        // 0019

				code.Make(code.OpTrue),            // 0022
				code.Make(code.OpJumpNotTrue, 33), // 0023
				code.Make(code.OpConstant, 2),     // 0026
				code.Make(code.OpPop),             // 0029
				code.Make(code.OpJump, 33),        // 0030

				code.Make(code.OpConstant, 3), // 0033
				code.Make(code.OpPop),         // 0036
			},
		},
		{
			input: `
            if: (false): {
                10;
            }
            `,
			expectedConstants: []interface{}{10},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),          // 0000
				code.Make(code.OpJumpNotTrue, 8), // 0001
				code.Make(code.OpConstant, 0),    // 0004
				code.Make(code.OpPop),            // 0007
			},
		},
		{
			input: `
            var a: int = 10;
            var b: int = 20;
            `,
			expectedConstants: []interface{}{10, 20},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
            var a: int = 10;
            a;
            `,
			expectedConstants: []interface{}{10},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
            var a: int = 10;
            var b: int = a;
            b;
            `,
			expectedConstants: []interface{}{10},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
                var a: int = 10;
                a += 20;
                a;
            `,
			expectedConstants: []interface{}{10, 20},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),

				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpSetGlobal, 0),

				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
                var a: int = 20;
                a -= 10;
                a;
            `,
			expectedConstants: []interface{}{20, 10},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),

				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpSetGlobal, 0),

				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
                var a: int = 10;
                a = 20;
                a;
            `,
			expectedConstants: []interface{}{10, 20},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),

				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 0),

				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
                var a: int = 10;
                var b: int = 20;
                a = b;
                a;
            `,
			expectedConstants: []interface{}{10, 20},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),

				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),

				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpSetGlobal, 0),

				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "[]",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "[1, 2, 3]",
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),

				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),

				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),

				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{}",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpHash, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{1: 2, 3: 4, 5: 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpHash, 6),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{1: 2 + 3, 4: 5 * 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),

				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpAdd),

				code.Make(code.OpConstant, 3),

				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),

				code.Make(code.OpHash, 4),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "[1, 2, 3][1]",
			expectedConstants: []interface{}{1, 2, 3, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "[1, 2, 3][1 + 1]",
			expectedConstants: []interface{}{1, 2, 3, 1, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpAdd),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{1: 2}[2 - 1]",
			expectedConstants: []interface{}{1, 2, 2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpHash, 2),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
	}
	err := runCompilerTests(t, tests)
	if err != nil {
		t.Fatal(err)
	}
}

func runCompilerTests(t *testing.T, tests []compileTestCase) error {
	t.Helper()
	for _, tt := range tests {
		program, err := parse(tt.input)
		if err != nil {
			return err
		}
		c := compiler.New()
		c.InTesting = true
		err = c.Compile(program)
		if err != nil {
			return errors.New("compiler error: " + err.Error())
		}
		bytecode := c.Bytecode()
		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			fmt.Println("Input: ", tt.input)
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			fmt.Println("Input: ", tt.input)
			t.Fatalf("testConstants failed: %s", err)
		}
	}
	return nil
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := cancatInstructions(expected)
	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length. \nwant=%q\ngot=%q", concatted, actual)
	}
	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d. \nwant=%q\ngot=%q", i, concatted, actual)
		}
	}
	return nil
}

func testConstants(_ *testing.T, expected []interface{}, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d", len(actual), len(expected))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObjectC(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		case float64:
			err := testFloatObjectC(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testFloatObject failed: %s", i, err)
			}
		case string:
			err := testStringObjectC(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s", i, err)
			}
		default:
			return fmt.Errorf("type of constant not handled. got=%T", constant)
		}
	}
	return nil
}

func cancatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

func testIntegerObjectC(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func testFloatObjectC(expected float64, actual object.Object) error {
	result, ok := actual.(*object.Float)
	if !ok {
		return fmt.Errorf("object is not Float. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
	}
	return nil
}

func testStringObjectC(expected string, actual object.Object) error {
	if expected[0] == '\'' {
		return testCharObjectC(expected, actual)
	}
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
	}
	return nil
}

func testCharObjectC(expected string, actual object.Object) error {
	result, ok := actual.(*object.Char)
	if !ok {
		return fmt.Errorf("object is not Char. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%s, want=%s", string(result.Value), string(expected[1]))
	}
	return nil
}

func parse(input string) (*ast.Program, error) {
	l := lexer.Tokenizer(input)
	p := parser.New(l, true)
	program, err := p.ParseProgram()
	if err != nil {
		return nil, errors.New("Error parsing program: " + err.Error())
	}
	typeCheckerEnv := parser.NewEnvironment()
	err = parser.TypeCheckProgram(program, typeCheckerEnv, true)
	if err != nil {
		return nil, errors.New("Error type checking program: " + err.Error())
	}
	return program, nil
}
