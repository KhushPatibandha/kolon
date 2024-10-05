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
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err {
			t.Errorf("Error evaluating %s", tt.input)
		}
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func Test26(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err {
			t.Errorf("Error evaluating %s", tt.input)
		}
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func Test27(t *testing.T) {
	test := []struct {
		input    string
		expected float64
	}{
		{"5.0", 5.0},
		{"10.0", 10.0},
		{"5.5", 5.5},
		{"-5.0", -5.0},
		{"-10.0", -10.0},
	}

	for _, tt := range test {
		evaluated, err := testEval(tt.input)
		if err {
			t.Errorf("Error evaluating %s", tt.input)
		}
		testFloatObject(t, evaluated, tt.expected)
	}
}

func Test28(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{"\"Hello, World!\"", "\"Hello, World!\""},
	}

	for _, tt := range test {
		evaluated, err := testEval(tt.input)
		if err {
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
		evaluated, err := testEval(tt.input)
		if err {
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
	}{
		{"!true", false},
		{"!false", true},
		{"!!false", false},
		{"!!true", true},
		{"!5", false},
		{"!5.5", false},
		{"!\"Hello, World!\"", false},
		{"!'a'", false},
	}

	for i, tt := range test {
		evaluated, err := testEval(tt.input)
		// fmt.Printf("evaluated: %v\n", evaluated)
		// fmt.Printf("err: %v\n", err)
		if i == 0 || i == 1 || i == 2 || i == 3 {
			if !testBooleanObject(t, evaluated, tt.expected) {
				t.Errorf("Error evaluating %s", tt.input)
			}
		}
		if i == 4 || i == 5 || i == 6 || i == 7 {
			if !err {
				t.Errorf("Error not raised for %s", tt.input)
			}
			continue
		}
		if err {
			t.Errorf("Error evaluating %s", tt.input)
		}
	}
}

func Test31(t *testing.T) {
	tests := []struct {
		input    string
		object   object.Object
		expected bool
	}{
		{"-5", &object.Integer{Value: -5}, false},
		{"-5.5", &object.Float{Value: -5.5}, false},
		{"-\"Hello, World!\"", evaluator.NULL, true},
		{"-'a'", evaluator.NULL, true},
		{"-true", evaluator.NULL, true},
		{"-false", evaluator.NULL, true},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if tt.expected {
			if !err {
				t.Errorf("Error not raised for %s", tt.input)
			}
			if tt.object != evaluator.NULL && evaluated != tt.object {
				t.Errorf("Error not raised for %s", tt.input)
			}
			continue
		}
		if err {
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

func testEval(input string) (object.Object, bool) {
	l := lexer.Tokenizer(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return evaluator.Eval(program)
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
