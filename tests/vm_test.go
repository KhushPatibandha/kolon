package tests

import (
	"fmt"
	"testing"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/compiler/compiler"
	"github.com/KhushPatibandha/Kolon/src/compiler/vm"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/object"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func Test52(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"2 / 2", 1},
		{"0 / 1", 0},
		{"4 % 2", 0},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"5", 5},
		{"10", 10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"10 % 2", 0},
		{"10 % 3", 1},
		{"2 & 3", 2},
		{"2 | 3", 3},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"20 + 2 * -10", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"5.0", 5.0},
		{"10.0", 10.0},
		{"5.5", 5.5},
		{"-5.0", -5.0},
		{"-10.0", -10.0},
		{"4 - 5.5", -1.5},
		{"5.5 - 4", 1.5},
		{"5.5 + 5.5 + 5.5 + 5.5 - 10.0", 12.0},
		{"2.5 * 2.5 * 2.5 * 2.5 * 2.5", 97.65625},
		{"-50.0 + 100.0 + -50.0", 0.0},
		{"5.5 * 2.5 + 10.0", 23.75},
		{"5.5 + 2.5 * 10.0", 30.5},
		{"20.0 + 2.5 * -10.0", -5.0},
		{"50.0 / 2.5 * 2.5 + 10.0", 60.0},
		{"2.5 * (5.5 + 10.0)", 38.75},
		{"3.5 * 3.5 * 3.5 + 10.0", 52.875},
		{"3.5 * (3.5 * 3.5) + 10.0", 52.875},
		{"3.5 * (3.5 + 10.0)", 47.25},
		{"(5.5 + 10.0 * 2.5 + 15.0 / 3.0) * 2.5 + -10.0", 78.75},
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"10 > 5", true},
		{"5 < 10", true},
		{"10 < 5", false},
		{"5 > 10", false},
		{"5 > 5", false},
		{"5 < 5", false},
		{"5 == 5", true},
		{"5 != 5", false},
		{"5 == 10", false},
		{"5 != 10", true},
		{"5 >= 5", true},
		{"5 <= 5", true},
		{"5 >= 10", false},
		{"5 <= 10", true},
		{"5 >= 4", true},
		{"5 <= 4", false},
		{"true && true", true},
		{"1 == 1 && 2 == 2 && 3 == 3", true},
		{"true && false", false},
		{"false && true", false},
		{"false && false", false},
		{"true || true", true},
		{"true || false", true},
		{"false || true", true},
		{"false || false", false},
		{"true == false", false},
		{"false == true", false},
		{"true == true", true},
		{"false == false", true},
		{"true != false", true},
		{"false != true", true},
		{"true != true", false},
		{"false != false", false},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"(1 >= 2) == false", true},
		{"2.5 > 3.5", false},
		{"2.5 < 3.5", true},
		{"2.5 == 3.5", false},
		{"2.5 != 3.5", true},
		{"2.5 >= 3.5", false},
		{"2.5 <= 3.5", true},
		{"2.5 >= 3.5", false},
		{"2.5 <= 3.5", true},
		{"2.5 >= 4.5", false},
		{"2.5 <= 4.5", true},
		{"2.5 >= 2.5", true},
		{"2.5 <= 2.5", true},
		{"2.5 >= 1.5", true},
		{"2.5 <= 1.5", false},
		{"(1.5 >= 2.5) == true", false},
		{"\"Hello, World!\" == \"Hello, World!\"", true},
		{"\"Hello, World!\" != \"Hello, World!\"", false},
		{"\"Hello, World!\" == \"Hello, World\"", false},
		{"\"Hello, World!\" != \"Hello, World\"", true},
		{"'a' == 'a'", true},
		{"'a' != 'a'", false},
		{"'a' == 'b'", false},
		{"'a' != 'b'", true},
		{"\"Hello, World!\"", "\"Hello, World!\""},
		{"\"Hello, \" + \"World!\"", "\"Hello, World!\""},
		{"'a' + 'b'", "\"ab\""},
		{"'a'", "'a'"},
		{"5++", 6},
		{"5--", 4},
		{"-5--", -4},
		{"-5++", -6},
		{"(-5)++", -4},
		{"10.1++", 11.1},
		{"10.1--", 9.1},
		{"-10.1++", -11.1},
		{"-10.1--", -9.1},
		{`
		if: (true): {
		10;
		}
		`, 10},
		{`
		if: (true): {
		10;
		} else: {
		20;
		}
		`, 10},
		{`
		if: (false): {
		10;
		} else: {
		20;
		}
		`, 20},
		{`
		if: (1 < 2): {
		10;
		}
		`, 10},
		{`
		if: (1 < 2): {
		10;
		} else: {
		20;
		}
		`, 10},
		{`
		if: (1 > 2): {
		10;
		} else: {
		20;
		}
		`, 20},
		{"if: (true): { 10; }", 10},
		{"if: (1 < 2): { 10; }", 10},
		{"if: (1 > 2): { 10; } else: { 20; }", 20},
		{"if: (1 < 2): { 10; } else: { 20; }", 10},
		{"if: (1 > 2): { 10; } else if: (1 == 2): { 20; } else if: (1 < 2): { 30; }", 30},
		{"if: (1 > 2): { 10; } else if: (1 < 2): { 20; } else if: (1 == 2): { 30; }", 20},
		{"if: (1 > 2): { 10; } else if: (1 == 2): { 20; } else if: (1 >= 2): { 30; } else: { 40; }", 40},

		{"if: (true): {}", true},
		{"if: (false): { 10; }", false},
		{"if: (1 > 2): { 10; }", false},
		{"if: (false): { 10; } else if: (false): { 20; } else if: (true): {}", true},
		{"if: (false): { 10; } else: {}", false},

		{"var a: int = 1; a;", 1},
		{"var a: int = 1; var b: int = 2; a + b;", 3},
		{"var a: int = 1; var b: int = a + a; a + b;", 3},

		{"var a: int = 10; a = 20; a;", 20},
		{"var a: int = 10; a += 10; a;", 20},
		{"var a: int = 10; a -= 10; a;", 0},
		{"var a: int = 10; a *= 10; a;", 100},
		{"var a: int = 10; a /= 10; a;", 1},
		{"var a: int = 4; a %= 2; a;", 0},
	}
	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()
	for _, tt := range tests {
		// fmt.Println(tt.input)
		program := parseVM(tt.input)
		if program == nil {
			return
		}
		c := compiler.New()
		err := c.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := vm.New(c.Bytecode())
		err = vm.Run()
		if err != nil {
			fmt.Println(tt.input)
			t.Fatalf("vm error: %s", err)
		}
		stackElem := vm.LastPoppedStackEle()
		err = testExpectedObject(t, tt.expected, stackElem)
		if err != nil {
			fmt.Println(tt.input)
			fmt.Println(stackElem)
			t.Errorf("%s", err)
		}
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) error {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObjectVM(int64(expected), actual)
		if err != nil {
			return fmt.Errorf("testIntegerObject failed: %w", err)
		}
	case bool:
		err := testBooleanObjectVM(bool(expected), actual)
		if err != nil {
			return fmt.Errorf("testBooleanObject failed: %w", err)
		}
	case float64:
		err := testFloatObjectVM(float64(expected), actual)
		if err != nil {
			return fmt.Errorf("testFloatObject failed: %w", err)
		}
	case string:
		err := testStringObjectVM(expected, actual)
		if err != nil {
			return fmt.Errorf("testStringObject failed: %w", err)
		}
	case nil:
		if actual != nil {
			return fmt.Errorf("object is not nil. got=%T (%+v)", actual, actual)
		}
	default:
		return fmt.Errorf("type of expected not handled. got=%T", expected)
	}
	return nil
}

func testIntegerObjectVM(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func testFloatObjectVM(expected float64, actual object.Object) error {
	result, ok := actual.(*object.Float)
	if !ok {
		return fmt.Errorf("object is not Float. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
	}
	return nil
}

func testBooleanObjectVM(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}
	return nil
}

func testStringObjectVM(expected string, actual object.Object) error {
	if expected[0] == '\'' {
		return testCharObjectVM(expected, actual)
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

func testCharObjectVM(expected string, actual object.Object) error {
	result, ok := actual.(*object.Char)
	if !ok {
		return fmt.Errorf("object is not Char. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%s, want=%s", string(result.Value), string(expected[1]))
	}
	return nil
}

func parseVM(input string) *ast.Program {
	l := lexer.Tokenizer(input)
	p := parser.New(l, true)
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Println("Error parsing program: ", err)
		return nil
	}
	typeCheckerEnv := parser.NewEnvironment()
	err = parser.TypeCheckProgram(program, typeCheckerEnv, true)
	if err != nil {
		fmt.Println("Error type checking program: ", err)
		return nil
	}
	return program
}
