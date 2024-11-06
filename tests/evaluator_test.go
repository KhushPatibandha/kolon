package tests

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/evaluator"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/object"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func Test25(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		hasErr   bool
	}{
		{"5", 5, false},
		{"10", 10, false},
		{"-5", -5, false},
		{"-10", -10, false},
		{"5 + 5 + 5 + 5 - 10", 10, false},
		{"2 * 2 * 2 * 2 * 2", 32, false},
		{"-50 + 100 + -50", 0, false},
		{"5 * 2 + 10", 20, false},
		{"5 + 2 * 10", 25, false},
		{"20 + 2 * -10", 0, false},
		{"50 / 2 * 2 + 10", 60, false},
		{"2 * (5 + 10)", 30, false},
		{"3 * 3 * 3 + 10", 37, false},
		{"3 * (3 * 3) + 10", 37, false},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50, false},
		{"10 % 2", 0, false},
		{"10 % 3", 1, false},
		{"2 & 3", 2, false},
		{"2 | 3", 3, false},
		{"5++", 6, false},
		{"5--", 4, false},
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

func Test26(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		hasErr   bool
	}{
		{"true", true, false},
		{"false", false, false},
		{"!true", false, false},
		{"!false", true, false},
		{"!!true", true, false},
		{"!!false", false, false},
		{"10 > 5", true, false},
		{"5 < 10", true, false},
		{"10 < 5", false, false},
		{"5 > 10", false, false},
		{"5 > 5", false, false},
		{"5 < 5", false, false},
		{"5 == 5", true, false},
		{"5 != 5", false, false},
		{"5 == 10", false, false},
		{"5 != 10", true, false},
		{"5 >= 5", true, false},
		{"5 <= 5", true, false},
		{"5 >= 10", false, false},
		{"5 <= 10", true, false},
		{"5 >= 4", true, false},
		{"5 <= 4", false, false},
		{"true && true", true, false},
		{"true && false", false, false},
		{"false && true", false, false},
		{"false && false", false, false},
		{"true || true", true, false},
		{"true || false", true, false},
		{"false || true", true, false},
		{"false || false", false, false},
		{"true == false", false, false},
		{"false == true", false, false},
		{"true == true", true, false},
		{"false == false", true, false},
		{"true != false", true, false},
		{"false != true", true, false},
		{"true != true", false, false},
		{"false != false", false, false},
		{"(1 < 2) == true", true, false},
		{"(1 < 2) == false", false, false},
		{"(1 > 2) == true", false, false},
		{"(1 > 2) == false", true, false},
		{"2.5 > 3.5", false, false},
		{"2.5 < 3.5", true, false},
		{"2.5 == 3.5", false, false},
		{"2.5 != 3.5", true, false},
		{"2.5 >= 3.5", false, false},
		{"2.5 <= 3.5", true, false},
		{"2.5 >= 3.5", false, false},
		{"2.5 <= 3.5", true, false},
		{"2.5 >= 4.5", false, false},
		{"2.5 <= 4.5", true, false},
		{"2.5 >= 2.5", true, false},
		{"2.5 <= 2.5", true, false},
		{"2.5 >= 1.5", true, false},
		{"2.5 <= 1.5", false, false},
		{"(1 >= 2) == false", true, false},
		{"(1.5 >= 2.5) == true", false, false},
		{"\"Hello, World!\" == \"Hello, World!\"", true, false},
		{"\"Hello, World!\" != \"Hello, World!\"", false, false},
		{"\"Hello, World!\" == \"Hello, World\"", false, false},
		{"\"Hello, World!\" != \"Hello, World\"", true, false},
		{"'a' == 'a'", true, false},
		{"'a' != 'a'", false, false},
		{"'a' == 'b'", false, false},
		{"'a' != 'b'", true, false},
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

func Test27(t *testing.T) {
	test := []struct {
		input    string
		expected float64
		hasErr   bool
	}{
		{"5.0", 5.0, false},
		{"10.0", 10.0, false},
		{"5.5", 5.5, false},
		{"-5.0", -5.0, false},
		{"-10.0", -10.0, false},
		{"4 - 5.5", -1.5, false},
		{"5.5 - 4", 1.5, false},
		{"5.5 + 5.5 + 5.5 + 5.5 - 10.0", 12.0, false},
		{"2.5 * 2.5 * 2.5 * 2.5 * 2.5", 97.65625, false},
		{"-50.0 + 100.0 + -50.0", 0.0, false},
		{"5.5 * 2.5 + 10.0", 23.75, false},
		{"5.5 + 2.5 * 10.0", 30.5, false},
		{"20.0 + 2.5 * -10.0", -5.0, false},
		{"50.0 / 2.5 * 2.5 + 10.0", 60.0, false},
		{"2.5 * (5.5 + 10.0)", 38.75, false},
		{"3.5 * 3.5 * 3.5 + 10.0", 52.875, false},
		{"3.5 * (3.5 * 3.5) + 10.0", 52.875, false},
		{"3.5 * (3.5 + 10.0)", 47.25, false},
		{"(5.5 + 10.0 * 2.5 + 15.0 / 3.0) * 2.5 + -10.0", 78.75, false},
		{"10.0 % 2.5", 0.0, false},
		{"10.0 % 3.5", 1.0, false},
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

func Test28(t *testing.T) {
	test := []struct {
		input    string
		expected string
		hasErr   bool
	}{
		{"\"Hello, World!\"", "\"Hello, World!\"", false},
		{"\"Hello, \" + \"World!\"", "\"Hello, World!\"", false},
		{"'a' + 'b'", "\"ab\"", false},
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

func Test29(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{"'a'", "'a'"},
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

func Test30(t *testing.T) {
	test := []struct {
		input    string
		expected bool
		hasErr   bool
	}{
		{"!true", false, false},
		{"!false", true, false},
		{"!!false", false, false},
		{"!!true", true, false},
		{"!5", false, true},
		{"!5.5", false, true},
		{"!\"Hello, World!\"", false, true},
		{"!'a'", false, true},
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

func Test31(t *testing.T) {
	tests := []struct {
		input  string
		object object.Object
		hasErr bool
	}{
		{"-5", &object.Integer{Value: -5}, false},
		{"-5.5", &object.Float{Value: -5.5}, false},
		{"-\"Hello, World!\"", evaluator.NULL, true},
		{"-'a'", evaluator.NULL, true},
		{"-true", evaluator.NULL, true},
		{"-false", evaluator.NULL, true},
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

func Test32(t *testing.T) {
	postfixTest := []struct {
		input  string
		object object.Object
		hasErr bool
	}{
		{"5++", &object.Integer{Value: 6}, false},
		{"5--", &object.Integer{Value: 4}, false},
		{"-5--", &object.Integer{Value: -6}, false},
		{"-5++", &object.Integer{Value: -4}, false},
		{"10.1++", &object.Float{Value: 11.1}, false},
		{"10.1--", &object.Float{Value: 9.1}, false},
		{"-10.1--", &object.Float{Value: -11.1}, false},
		{"-10.1++", &object.Float{Value: -9.1}, false},
		{"true++", &object.Null{}, true},
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
			if err != nil && i != 8 {
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
		}
	}
}

func Test33(t *testing.T) {
	ifElseTests := []struct {
		input    string
		expected interface{}
		hasErr   bool
	}{
		{"if: (true): { 10 }", 10, false},
		{"if: (false): { 10 }", nil, false},
		{"if: (1): { 10 }", evaluator.NULL, true}, // should throw error
		{"if: (1 < 2): { 10 }", 10, false},
		{"if: (1 > 2): { 10 }", nil, false},
		{"if: (1 > 2): { 10 } else: { 20 }", 20, false},
		{"if: (1 < 2): { 10 } else: { 20 }", 10, false},
		{"if: (1 > 2): { 10 } else if: (1 == 2): { 20 } else if: (1 < 2): { 30 }", 30, false},
		{"if: (1 > 2): { 10 } else if: (1 < 2): { 20 } else if: (1 == 2): { 30 }", 20, false},
		{"if: (1 > 2): { 10 } else if: (1 == 2): { 20 } else if: (1 >= 2): { 30 } else: { 40 }", 40, false},
		{"if: (1 > 2): { 10 } else if: (1): { 20 } else if: (1 >= 2): { 30 } else: { 40 }", evaluator.NULL, true}, // should throw error
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

func testEval(input string) (object.Object, bool, error) {
	l := lexer.Tokenizer(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return evaluator.Eval(program)
}

func testNil(t *testing.T, obj object.Object) bool {
	if obj != nil {
		t.Errorf("Object is not nil. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

// func testNullObject(t *testing.T, obj object.Object) bool {
// 	if obj != evaluator.NULL {
// 		t.Errorf("Object is not NULL. got=%T (%+v)", obj, obj)
// 		return false
// 	}
// 	return true
// }

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
