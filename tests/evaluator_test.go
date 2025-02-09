package tests

import (
	"fmt"
	"testing"

	"github.com/KhushPatibandha/Kolon/src/interpreter/evaluator"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/object"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func Test27(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		hasErr   bool
	}{
		{"5;", 5, false},
		{"10;", 10, false},
		{"-5;", -5, false},
		{"-10;", -10, false},
		{"5 + 5 + 5 + 5 - 10;", 10, false},
		{"2 * 2 * 2 * 2 * 2;", 32, false},
		{"-50 + 100 + -50;", 0, false},
		{"5 * 2 + 10;", 20, false},
		{"5 + 2 * 10;", 25, false},
		{"20 + 2 * -10;", 0, false},
		{"50 / 2 * 2 + 10;", 60, false},
		{"2 * (5 + 10);", 30, false},
		{"3 * 3 * 3 + 10;", 37, false},
		{"3 * (3 * 3) + 10;", 37, false},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10;", 50, false},
		{"10 % 2;", 0, false},
		{"10 % 3;", 1, false},
		{"2 & 3;", 2, false},
		{"2 | 3;", 3, false},
		{"5++;", 6, false},
		{"5--;", 4, false},
	}

	for _, tt := range tests {
		evaluated, hasErr, err := testEval(tt.input)
		if err != nil {
			t.Errorf("Error evaluating %s", tt.input)
		}
		if hasErr || hasErr != tt.hasErr {
			t.Errorf("Error evaluating %s", tt.input)
		}
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func Test28(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		hasErr   bool
	}{
		{"true;", true, false},
		{"false;", false, false},
		{"!true;", false, false},
		{"!false;", true, false},
		{"!!true;", true, false},
		{"!!false;", false, false},
		{"10 > 5;", true, false},
		{"5 < 10;", true, false},
		{"10 < 5;", false, false},
		{"5 > 10;", false, false},
		{"5 > 5;", false, false},
		{"5 < 5;", false, false},
		{"5 == 5;", true, false},
		{"5 != 5;", false, false},
		{"5 == 10;", false, false},
		{"5 != 10;", true, false},
		{"5 >= 5;", true, false},
		{"5 <= 5;", true, false},
		{"5 >= 10;", false, false},
		{"5 <= 10;", true, false},
		{"5 >= 4;", true, false},
		{"5 <= 4;", false, false},
		{"true && true;", true, false},
		{"1 == 1 && 2 == 2 && 3 == 3;", true, false},
		{"true && false;", false, false},
		{"false && true;", false, false},
		{"false && false;", false, false},
		{"true || true;", true, false},
		{"true || false;", true, false},
		{"false || true;", true, false},
		{"false || false;", false, false},
		{"true == false;", false, false},
		{"false == true;", false, false},
		{"true == true;", true, false},
		{"false == false;", true, false},
		{"true != false;", true, false},
		{"false != true;", true, false},
		{"true != true;", false, false},
		{"false != false;", false, false},
		{"(1 < 2) == true;", true, false},
		{"(1 < 2) == false;", false, false},
		{"(1 > 2) == true;", false, false},
		{"(1 > 2) == false;", true, false},
		{"2.5 > 3.5;", false, false},
		{"2.5 < 3.5;", true, false},
		{"2.5 == 3.5;", false, false},
		{"2.5 != 3.5;", true, false},
		{"2.5 >= 3.5;", false, false},
		{"2.5 <= 3.5;", true, false},
		{"2.5 >= 3.5;", false, false},
		{"2.5 <= 3.5;", true, false},
		{"2.5 >= 4.5;", false, false},
		{"2.5 <= 4.5;", true, false},
		{"2.5 >= 2.5;", true, false},
		{"2.5 <= 2.5;", true, false},
		{"2.5 >= 1.5;", true, false},
		{"2.5 <= 1.5;", false, false},
		{"(1 >= 2) == false;", true, false},
		{"(1.5 >= 2.5) == true;", false, false},
		{"\"Hello, World!\" == \"Hello, World!\";", true, false},
		{"\"Hello, World!\" != \"Hello, World!\";", false, false},
		{"\"Hello, World!\" == \"Hello, World\";", false, false},
		{"\"Hello, World!\" != \"Hello, World\";", true, false},
		{"'a' == 'a';", true, false},
		{"'a' != 'a';", false, false},
		{"'a' == 'b';", false, false},
		{"'a' != 'b';", true, false},
	}

	for _, tt := range tests {
		evaluated, hasErr, err := testEval(tt.input)
		if err != nil {
			t.Errorf("Error evaluating %s", tt.input)
		}
		if hasErr || hasErr != tt.hasErr {
			t.Errorf("Error evaluating %s", tt.input)
		}
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func Test29(t *testing.T) {
	test := []struct {
		input    string
		expected float64
		hasErr   bool
	}{
		{"5.0;", 5.0, false},
		{"10.0;", 10.0, false},
		{"5.5;", 5.5, false},
		{"-5.0;", -5.0, false},
		{"-10.0;", -10.0, false},
		{"4 - 5.5;", -1.5, false},
		{"5.5 - 4;", 1.5, false},
		{"5.5 + 5.5 + 5.5 + 5.5 - 10.0;", 12.0, false},
		{"2.5 * 2.5 * 2.5 * 2.5 * 2.5;", 97.65625, false},
		{"-50.0 + 100.0 + -50.0;", 0.0, false},
		{"5.5 * 2.5 + 10.0;", 23.75, false},
		{"5.5 + 2.5 * 10.0;", 30.5, false},
		{"20.0 + 2.5 * -10.0;", -5.0, false},
		{"50.0 / 2.5 * 2.5 + 10.0;", 60.0, false},
		{"2.5 * (5.5 + 10.0);", 38.75, false},
		{"3.5 * 3.5 * 3.5 + 10.0;", 52.875, false},
		{"3.5 * (3.5 * 3.5) + 10.0;", 52.875, false},
		{"3.5 * (3.5 + 10.0);", 47.25, false},
		{"(5.5 + 10.0 * 2.5 + 15.0 / 3.0) * 2.5 + -10.0;", 78.75, false},
	}

	for _, tt := range test {
		evaluated, hasErr, err := testEval(tt.input)
		if err != nil {
			t.Errorf("Error evaluating %s", tt.input)
		}
		if hasErr || hasErr != tt.hasErr {
			t.Errorf("Error evaluating %s", tt.input)
		}
		testFloatObject(t, evaluated, tt.expected)
	}
}

func Test30(t *testing.T) {
	test := []struct {
		input    string
		expected string
		hasErr   bool
	}{
		{"\"Hello, World!\";", "\"Hello, World!\"", false},
		{"\"Hello, \" + \"World!\";", "\"Hello, World!\"", false},
		{"'a' + 'b';", "\"ab\"", false},
	}

	for _, tt := range test {
		evaluated, hasErr, err := testEval(tt.input)
		if err != nil {
			t.Errorf("Error evaluating %s", tt.input)
		}
		if hasErr || hasErr != tt.hasErr {
			t.Errorf("Error evaluating %s", tt.input)
		}
		result, ok := evaluated.(*object.String)
		if !ok {
			t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if result.Value != tt.expected {
			t.Errorf("object has wrong value. got=%s, want=%s", result.Value, tt.expected)
		}
	}
}

func Test31(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{"'a';", "'a'"},
	}

	for _, tt := range test {
		evaluated, hasErr, err := testEval(tt.input)
		if err != nil {
			t.Errorf("Error evaluating %s", tt.input)
		}
		if hasErr {
			t.Errorf("Error evaluating %s", tt.input)
		}
		result, ok := evaluated.(*object.Char)
		if !ok {
			t.Errorf("object is not Char. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if result.Value != tt.expected {
			t.Errorf("object has wrong value. got=%s, want=%s", result.Value, tt.expected)
		}
	}
}

func Test32(t *testing.T) {
	test := []struct {
		input    string
		expected bool
		hasErr   bool
	}{
		{"!true;", false, false},
		{"!false;", true, false},
		{"!!false;", false, false},
		{"!!true;", true, false},
		{"!5;", false, true},
		{"!5.5;", false, true},
		{"!\"Hello, World!\";", false, true},
		{"!'a';", false, true},
	}

	for i, tt := range test {
		evaluated, hasErr, err := testEval(tt.input)
		// fmt.Printf("evaluated: %v\n", evaluated)
		// fmt.Printf("err: %v\n", err)
		if i == 0 || i == 1 || i == 2 || i == 3 {
			if !testBooleanObject(t, evaluated, tt.expected) {
				t.Errorf("Error evaluating %s", tt.input)
			}
		}
		if i == 4 || i == 5 || i == 6 || i == 7 {
			if !hasErr || hasErr != tt.hasErr {
				t.Errorf("Error not raised for %s", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("Error evaluating %s", tt.input)
		}
	}
}

func Test33(t *testing.T) {
	tests := []struct {
		input  string
		object object.Object
		hasErr bool
	}{
		{"-5;", &object.Integer{Value: -5}, false},
		{"-5.5;", &object.Float{Value: -5.5}, false},
		{"-\"Hello, World!\";", evaluator.NULL, true},
		{"-'a';", evaluator.NULL, true},
		{"-true;", evaluator.NULL, true},
		{"-false;", evaluator.NULL, true},
	}

	for _, tt := range tests {
		evaluated, hasErr, err := testEval(tt.input)
		if tt.hasErr {
			if !hasErr {
				t.Errorf("Error not raised for %s", tt.input)
			}
			if tt.object != evaluator.NULL && evaluated != tt.object {
				t.Errorf("Got different object. got: %s", string(tt.object.Type()))
			}
			continue
		}
		if err != nil {
			t.Errorf("Error evaluating %s", tt.input)
		}
		if tt.object.Type() == object.INTEGER_OBJ {
			testIntegerObject(t, evaluated, tt.object.(*object.Integer).Value)
		}
		if tt.object.Type() == object.FLOAT_OBJ {
			testFloatObject(t, evaluated, tt.object.(*object.Float).Value)
		}
	}
}

func Test34(t *testing.T) {
	postfixTest := []struct {
		input  string
		object object.Object
		hasErr bool
	}{
		{"5++;", &object.Integer{Value: 6}, false},
		{"5--;", &object.Integer{Value: 4}, false},
		{"-5--;", &object.Integer{Value: -4}, false},
		{"-5++;", &object.Integer{Value: -6}, false},
		{"(-5)++;", &object.Integer{Value: -4}, false},
		{"10.1++;", &object.Float{Value: 11.1}, false},
		{"10.1--;", &object.Float{Value: 9.1}, false},
		{"-10.1--;", &object.Float{Value: -9.1}, false},
		{"-10.1++;", &object.Float{Value: -11.1}, false},
		{"true++;", &object.Null{}, true},
	}

	for i, tt := range postfixTest {
		evaluated, hasErr, err := testEval(tt.input)
		if tt.hasErr {
			if !hasErr {
				t.Errorf("Error not raised for %s", tt.input)
			}
			if err == nil {
				t.Errorf("Expected error for %s, but got nil", tt.input)
			}
			if err != nil && i != 9 {
				t.Error(err.Error())
			}
			if tt.object != evaluator.NULL && evaluated != tt.object {
				t.Errorf("Expected to get NULL Object in times of error. got: %s", string(tt.object.Type()))
			}
			continue
		}
		if err != nil {
			t.Errorf("Error evaluating %s", tt.input)
		}
		if tt.object.Type() == object.INTEGER_OBJ {
			testIntegerObject(t, evaluated, tt.object.(*object.Integer).Value)
		} else if tt.object.Type() == object.FLOAT_OBJ {
			testFloatObject(t, evaluated, tt.object.(*object.Float).Value)
		}
	}
}

func Test35(t *testing.T) {
	ifElseTests := []struct {
		input    string
		expected interface{}
		hasErr   bool
	}{
		{"if: (true): { 10; }", 10, false},
		{"if: (false): { 10; }", nil, false},
		{"if: (1): { 10; }", evaluator.NULL, true}, // should throw error
		{"if: (1 < 2): { 10; }", 10, false},
		{"if: (1 > 2): { 10; }", nil, false},
		{"if: (1 > 2): { 10; } else: { 20; }", 20, false},
		{"if: (1 < 2): { 10; } else: { 20; }", 10, false},
		{"if: (1 > 2): { 10; } else if: (1 == 2): { 20; } else if: (1 < 2): { 30; }", 30, false},
		{"if: (1 > 2): { 10; } else if: (1 < 2): { 20; } else if: (1 == 2): { 30; }", 20, false},
		{"if: (1 > 2): { 10; } else if: (1 == 2): { 20; } else if: (1 >= 2): { 30; } else: { 40; }", 40, false},
		{"if: (1 > 2): { 10; } else if: (1): { 20; } else if: (1 >= 2): { 30; } else: { 40; }", evaluator.NULL, true}, // should throw error
	}

	for i, tt := range ifElseTests {
		evaluated, hasErr, err := testEval(tt.input)
		if hasErr && i != 2 && i != 10 {
			t.Errorf("Error evaluating %s", tt.input)
			if err != nil {
				t.Error(err.Error())
			}
		}
		if (i == 2 || i == 10) && !hasErr {
			t.Errorf("This Test should throw an error. %s", tt.input)
		}
		if i == 2 || i == 10 {
			continue
		}
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNil(t, evaluated)
		}
	}
}

func Test36(t *testing.T) {
	returnStmtTests := []struct {
		input          string
		expectedOutput []object.Object
		hasErr         bool
	}{
		{"return: true;", []object.Object{evaluator.TRUE}, false},
		{"return: (true, false, true);", []object.Object{evaluator.TRUE, evaluator.FALSE, evaluator.TRUE}, false},
		{"return: 10;", []object.Object{&object.Integer{Value: 10}}, false},
		{"return: (10, 20, 1 < 2, true != true, 1 + 1, 4.1);", []object.Object{&object.Integer{Value: 10}, &object.Integer{Value: 20}, evaluator.TRUE, evaluator.FALSE, &object.Integer{Value: 2}, &object.Float{Value: 4.1}}, false},
		{"return: \"Hello\";", []object.Object{&object.String{Value: "\"Hello\""}}, false},
		{"return: 'c';", []object.Object{&object.Char{Value: "'c'"}}, false},
		{"return: (\"Hello\", 'h', true, \"World\");", []object.Object{&object.String{Value: "\"Hello\""}, &object.Char{Value: "'h'"}, evaluator.TRUE, &object.String{Value: "\"World\""}}, false},
		{"return: 10; 9;", []object.Object{&object.Integer{Value: 10}}, false},
		{"return: 2 * 5; 9;", []object.Object{&object.Integer{Value: 10}}, false},
		{"9; return: 2 * 5; 9;", []object.Object{&object.Integer{Value: 10}}, false},
		{
			`if: (10 > 1): {
                if: (10 > 1): {
                    return: 10;
                }
                return: 1;
            }`, []object.Object{&object.Integer{Value: 10}}, false,
		},
		{
			`if: (10 > 1): {
                if: (10 > 1): {
                    return: (10, true, "Hello", 'w', 1.1);
                }
                return: (1, false, "World", 'h', 10.1);
            }`, []object.Object{&object.Integer{Value: 10}, evaluator.TRUE, &object.String{Value: "\"Hello\""}, &object.Char{Value: "'w'"}, &object.Float{Value: 1.1}}, false,
		},
	}

	for _, tt := range returnStmtTests {
		evaluated, hasErr, err := testEval(tt.input)
		if hasErr != tt.hasErr {
			t.Error("expected error and recived error not matching")
		}
		if err != nil {
			t.Error(err.Error())
		}

		if !hasErr {
			returnValue, ok := evaluated.(*object.ReturnValue)
			if !ok {
				t.Errorf("Expected Return Object. got=%T", evaluated)
			}

			if len(returnValue.Value) != len(tt.expectedOutput) {
				t.Fatalf("Expected %d values in ReturnValue, got=%d", len(tt.expectedOutput), len(returnValue.Value))
			}

			for i, expected := range tt.expectedOutput {
				actual := returnValue.Value[i]

				switch expected := expected.(type) {
				case *object.Boolean:
					testBooleanObject(t, actual, expected.Value)
				case *object.Integer:
					testIntegerObject(t, actual, expected.Value)
				case *object.Float:
					testFloatObject(t, actual, expected.Value)
				case *object.String:
					testStringObject(t, actual, expected.Value)
				case *object.Char:
					testCharObject(t, actual, expected.Value)
				}
			}
		}
	}
}

func Test37(t *testing.T) {
	varStmtTest := []struct {
		input          string
		expectedOutput []object.Object
		hasErr         bool
	}{
		{"var a: int; return: a;", []object.Object{&object.Integer{Value: 0}}, false},
		{"var a: string; return: a;", []object.Object{&object.String{Value: ""}}, false},
		{"var a: float; return: a;", []object.Object{&object.Float{Value: 0.0}}, false},
		{"var a: char; return: a;", []object.Object{&object.Char{Value: ""}}, false},
		{"var a: bool; return: a;", []object.Object{&object.Boolean{Value: false}}, false},
		{"var a: int = 10; return: a;", []object.Object{&object.Integer{Value: 10}}, false},
		{"var a: int = 10; const b: int = a++; return: (a, b);", []object.Object{&object.Integer{Value: 10}, &object.Integer{Value: 11}}, false},
		{"var a: int = 10; var b: int = a++; return: (a, b);", []object.Object{&object.Integer{Value: 10}, &object.Integer{Value: 11}}, false},
		{"var a: int = 10; return: a + a;", []object.Object{&object.Integer{Value: 20}}, false},
		{"var a: int = 10; a++; return: a;", []object.Object{&object.Integer{Value: 11}}, false},
		{"var a: int = 10; return: (a++, a, true);", []object.Object{&object.Integer{Value: 11}, &object.Integer{Value: 10}, evaluator.TRUE}, false},
		{"var a: int = 10; return: a++;", []object.Object{&object.Integer{Value: 11}}, false},
		{"var a: int = 10; var a: bool = true; return: a;", []object.Object{evaluator.TRUE}, false},
		{"var a: int = 5 * 5; return: a;", []object.Object{&object.Integer{Value: 25}}, false},
		{"var a: int = 10; var b: int = a; return: b;", []object.Object{&object.Integer{Value: 10}}, false},
		{"var a: int = 10; var b: int = a; b = 11; return: b;", []object.Object{&object.Integer{Value: 11}}, false},
		{"var a: int = 10; var b: int = a; b = 11; b = b + 1; return: b;", []object.Object{&object.Integer{Value: 12}}, false},
		{"var a: int = 10; var b: int = a; b = 11; a += b; return: a;", []object.Object{&object.Integer{Value: 21}}, false},
		{"var a: int = 10; var b: int = a; b = 11; a -= b; return: a;", []object.Object{&object.Integer{Value: -1}}, false},
		{"var a: int = 10; var b: int = a; var c: int = a + b + 5; return: c;", []object.Object{&object.Integer{Value: 25}}, false},
		{"var a: int = 10; var b: float = 1.1; var c: float = a + b; return: c;", []object.Object{&object.Float{Value: 11.1}}, false},
		{"var a: int = 10; var b: float = 1.1; b += a; return: b;", []object.Object{&object.Float{Value: 11.1}}, false},
		{"var a: bool = true; return: a;", []object.Object{evaluator.TRUE}, false},
		{"var a: string = \"Hello \"; var b: string = \"World!\"; var c: string = a + b; return: c;", []object.Object{&object.String{Value: "\"Hello World!\""}}, false},
		{"var a: string = \"Hello \"; var b: string = \"World!\"; a += b; return: a;", []object.Object{&object.String{Value: "\"Hello World!\""}}, false},
		{"var a: char = 'c'; var b: char = 'c'; var c: string = a + b; return: c;", []object.Object{&object.String{Value: "\"cc\""}}, false},
		{"const a: int = 1; return: a;", []object.Object{&object.Integer{Value: 1}}, false},
		{"const a: float = 1.1; return: a;", []object.Object{&object.Float{Value: 1.1}}, false},
		{"const a: bool = true; return: a;", []object.Object{evaluator.TRUE}, false},
		{"const a: string = \"a\"; return: a;", []object.Object{&object.String{Value: "\"a\""}}, false},
		{"const a: char = 'a'; return: a;", []object.Object{&object.Char{Value: "'a'"}}, false},
		{"var a: int = 0; for: (var i: int = 0; i < 2; i++): { a++; } return: a;", []object.Object{&object.Integer{Value: 2}}, false},
		{`
		    var a: int = 1;
		    var b: int = 1;
		    for: (var i: int = 0; i < 3; i++): {
		        a++;
		        b += 2;
		    }
		    return: (a, b, a++, b++, a + b);
		`, []object.Object{&object.Integer{Value: 4}, &object.Integer{Value: 7}, &object.Integer{Value: 5}, &object.Integer{Value: 8}, &object.Integer{Value: 11}}, false},
		{
			`
                var a: int = 1;
                var b: int = 1;
                for: (var i: int = 0; i < 3; i++): {
                    for: (var j: int = 0; j < 3; j++): {
                        a++;
                        b += 2;
                    }
                }
                return: (a, b, a++, b++, a + b);
            `, []object.Object{&object.Integer{Value: 10}, &object.Integer{Value: 19}, &object.Integer{Value: 11}, &object.Integer{Value: 20}, &object.Integer{Value: 29}}, false,
		},
		{"var a: int; var b: int, a, var c: int = 1, 2, 3; return: (a, b, c);", []object.Object{&object.Integer{Value: 2}, &object.Integer{Value: 1}, &object.Integer{Value: 3}}, false},
		{"var a: int; var b: int; a, b, var c: int = 1, 2, 3; return: (a, b, c);", []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}, &object.Integer{Value: 3}}, false},
		{"var a: int = len(\"Hello\"); return: a;", []object.Object{&object.Integer{Value: 5}}, false},
		{"var b: string = \"hello\"; var a: int = len(b); return: a;", []object.Object{&object.Integer{Value: 5}}, false},
		{"var a: int = len(\"\"); return: a;", []object.Object{&object.Integer{Value: 0}}, false},
		{"var a: int[] = [1, 2, 3, 4]; return: a;", []object.Object{&object.Array{Elements: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}, &object.Integer{Value: 3}, &object.Integer{Value: 4}}, TypeOf: "int"}}, false},
	}

	for _, tt := range varStmtTest {
		// if i == 37 {
		// 	fmt.Println(i)
		evaluated, hasErr, err := testEval(tt.input)
		if hasErr != tt.hasErr {
			fmt.Println(tt.input)
			t.Errorf("expected error an recived error not matching. got=%v expected=%v", hasErr, tt.hasErr)
		}
		if err != nil {
			t.Error(err.Error())
		}
		if !hasErr {
			returnValue, ok := evaluated.(*object.ReturnValue)
			if !ok {
				t.Errorf("Expected Return Object. got=%T", evaluated)
			}
			if len(returnValue.Value) != len(tt.expectedOutput) {
				t.Fatalf("Expected %d values in ReturnValue, got=%d", len(tt.expectedOutput), len(returnValue.Value))
			}

			for i, expected := range tt.expectedOutput {
				actual := returnValue.Value[i]

				switch expected := expected.(type) {
				case *object.Boolean:
					testBooleanObject(t, actual, expected.Value)
				case *object.Integer:
					testIntegerObject(t, actual, expected.Value)
				case *object.Float:
					testFloatObject(t, actual, expected.Value)
				case *object.String:
					testStringObject(t, actual, expected.Value)
				case *object.Char:
					testCharObject(t, actual, expected.Value)
				case *object.Null:
					testNullObject(t, actual)
				case *object.Array:
					testArrayObject(t, actual, expected.Elements)
				default:
					testNil(t, actual)
				}

			}
		}
		// fmt.Println("Pass")
		// }
	}
}

func Test38(t *testing.T) {
	varStmtErrTest := []struct {
		input          string
		expectedOutput object.Object
		hasErr         bool
	}{
		{"var a: string = 10; var b: int = a; var c: int = a + b + 5; return: c;", evaluator.NULL, true},
		{"a = 10;", evaluator.NULL, true},
		{"var a: int = 10; var b: float = 1.1; a += b; return: a;", &object.Float{Value: 11.1}, true},
		{"const a: int = 1; a += 1; return: a;", &object.Integer{Value: 2}, true},
		{"var a: int = 10; r++; return: a;", &object.Integer{Value: 11}, true},
		{"const a: int = 10; a++; return: a;", &object.Integer{Value: 11}, true},
		{`
            if: (10 > 1): {
                var a: int = 10;
            }
            return: a;
        `, &object.Integer{Value: 10}, true},
		{
			`
                if: (10 < 1): {
                    var a: int = 10;
                } else: {
                    var b: int = 20;
                }
                return: b;
            `, &object.Integer{Value: 20}, true,
		},
		{
			`
                if: (10 < 1): {
                    var a: int = 10;
                } else if: (1 < 2): {
                    var c: int = 30;
                } else: {
                    var b: int = 20;
                }
                return: c;
            `, &object.Integer{Value: 30}, true,
		},
		{`
		    for: (var i: int = 0; i < 3; i++): {
                var a: int = 10;
		    }
            return: a;
        `, &object.Integer{Value: 10}, true},
		{"var a: int = len(); return: a;", &object.Integer{Value: 0}, true},
		{"var a: int[] = {1, \"hello\", 3, 4}; return: a;", &object.Array{Elements: []object.Object{&object.Integer{Value: 1}, &object.String{Value: "\"hello\""}, &object.Integer{Value: 3}, &object.Integer{Value: 4}}}, true},
		{"var a: int[]; return: a;", &object.Array{Elements: []object.Object{}}, true},
		{"var a: char = 'c'; var b: char = 'c'; a += b; return: a;", &object.String{Value: "\"cc\""}, true},
		{"const a: int = 1; const a: bool = true; return: a;", &object.ReturnValue{Value: []object.Object{evaluator.TRUE}}, true},
	}

	for _, tt := range varStmtErrTest {
		_, hasErr, err := testEval(tt.input)
		if hasErr != tt.hasErr {
			t.Error("expected error an recived error not matching")
		}
		if err == nil {
			t.Error("All these test must throw an error")
		}
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
	}
}

func Test39(t *testing.T) {
	testFunctions := []struct {
		input    string
		expected object.Object
		hasErr   bool
	}{
		{"fun: main() { var a: int = 10; }", nil, false},
		{"fun: main() { var a: int = 10; if: (a == 10): { return; } else: { a++; }}", nil, false},
		{"fun: add(a: int, b: int): (int) { return: a + b; } fun: main() { var a: int = 10; if: (a == 10): { return; } else: { a++; }}", nil, false},
		{"fun: add(a: int, b: int): (int) { return: a + b; } fun: main() { var a: int = add(5, 5); if: (a == 10): { return; } else: { a++; }}", nil, false},
		{"fun: add(a: int, b: int): (int) { return: a + b; } fun: main() { var a: int = add(1, 5); if: (a == 10): { return; } else: { a++; } if: (a == 7): { return; }}", nil, false},
	}

	for _, tt := range testFunctions {
		evaluated, hasErr, err := testEval(tt.input)
		if hasErr != tt.hasErr {
			t.Error("expected error an recived error not matching")
		}
		if err != nil {
			t.Error(err.Error())
			return
		}
		if !hasErr {
			switch expected := tt.expected.(type) {
			case *object.Boolean:
				testBooleanObject(t, evaluated, expected.Value)
			case *object.Integer:
				testIntegerObject(t, evaluated, expected.Value)
			case *object.Float:
				testFloatObject(t, evaluated, expected.Value)
			case *object.String:
				testStringObject(t, evaluated, expected.Value)
			case *object.Char:
				testCharObject(t, evaluated, expected.Value)
			case *object.Null:
				testNullObject(t, evaluated)
			default:
				testNil(t, evaluated)
			}
		}
	}
}

func testEval(input string) (object.Object, bool, error) {
	l := lexer.Tokenizer(input)
	p := parser.New(l, true)
	program, err := p.ParseProgram()
	if err != nil {
		return nil, true, err
	}
	typeCheckEnv := parser.NewEnvironment()
	err = parser.TypeCheckProgram(program, typeCheckEnv, true)
	if err != nil {
		return evaluator.NULL, true, err
	}
	env := object.NewEnvironment()

	return evaluator.Eval(program, env, true)
}

func testNil(t *testing.T, obj object.Object) bool {
	if obj != nil {
		t.Errorf("Object is not nil. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != evaluator.NULL {
		t.Errorf("Object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testArrayObject(t *testing.T, obj object.Object, expected []object.Object) bool {
	result, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
		return false
	}
	if len(result.Elements) != len(expected) {
		t.Errorf("Array has wrong number of elements. got=%d, want=%d", len(result.Elements), len(expected))
		return false
	}
	for i, expectedElem := range expected {
		switch result.TypeOf {
		case "int":
			testIntegerObject(t, result.Elements[i], expectedElem.(*object.Integer).Value)
		case "float":
			testFloatObject(t, result.Elements[i], expectedElem.(*object.Float).Value)
		case "string":
			testStringObject(t, result.Elements[i], expectedElem.(*object.String).Value)
		case "char":
			testCharObject(t, result.Elements[i], expectedElem.(*object.Char).Value)
		case "bool":
			testBooleanObject(t, result.Elements[i], expectedElem.(*object.Boolean).Value)
		}
	}
	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("Object is not a String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Object has worng value. got=%s, want=%s", result.Value, expected)
		return false
	}
	return true
}

func testCharObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Char)
	if !ok {
		t.Errorf("Object is not a Char. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Object has worng value. got=%s, want=%s", result.Value, expected)
		return false
	}
	return true
}
