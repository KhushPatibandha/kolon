package lexer_test

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/lexer"
)

func Test1(t *testing.T) {
	input := `
    fun main() {
        // this should not print
        var someString: string = "hey"
        var someChar: char = 'c'
        var someInt: int = 2147483647
        var someLong: long = 9223372036854775807
        var someOtherLong: long = 123l
        var someFloat: float = 1.1
        var someOtherFloat: float = 3.14
        var someDouble: double = 3.141592653589793
        var someOtherDouble: double = 3.14d
        var someBool: bool = true
        var someOtherBol: bool = false
    }`

	tests := []struct {
		expectedType    lexer.TokenKind
		expectedLiteral string
	}{
		{lexer.FUN, "fun"},
		{lexer.IDENTIFIER, "main"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.OPEN_CURLY_BRACKET, "{"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "someString"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "string"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.STRING, "\"hey\""},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "someChar"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "char"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.CHAR, "'c'"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "someInt"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "int"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.INT, "2147483647"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "someLong"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "long"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.LONG, "9223372036854775807"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "someOtherLong"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "long"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.LONG, "123l"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "someFloat"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "float"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.FLOAT, "1.1"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "someOtherFloat"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "float"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.FLOAT, "3.14"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "someDouble"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "double"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.DOUBLE, "3.141592653589793"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "someOtherDouble"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "double"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.DOUBLE, "3.14d"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "someBool"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "bool"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.BOOL, "true"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "someOtherBol"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "bool"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.BOOL, "false"},
		{lexer.CLOSE_CURLY_BRACKET, "}"},
		{lexer.EOF, "EOF"},
	}

	tokens := lexer.Tokenizer(input)
	for i, tt := range tests {
		if tokens[i].Kind != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedType, tokens[i].Kind)
		}
		if tokens[i].Value != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tokens[i].Value)
		}
	}
}
