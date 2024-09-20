package lexer_test

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/lexer"
)

func Test2(t *testing.T) {
	input := `
    fun main() {
        var age: int = 10
        var ok: bool = canVote(10)
        println(ok)
    }
    fun canVote(age: int): (bool) {
        if: (age >= 18) {
            print("You can vote")
            println()
            return true
        } else {
            println("You cannot vote kid!")
        }
        return false
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
		{lexer.IDENTIFIER, "age"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "int"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.INT, "10"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "ok"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "bool"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.IDENTIFIER, "canVote"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.INT, "10"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.PRINTLN, "ok"},
		{lexer.CLOSE_CURLY_BRACKET, "}"},
		{lexer.FUN, "fun"},
		{lexer.IDENTIFIER, "canVote"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "age"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "int"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.COLON, ":"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.TYPE, "bool"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.OPEN_CURLY_BRACKET, "{"},
		{lexer.IF, "if"},
		{lexer.COLON, ":"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "age"},
		{lexer.GREATER_THAN_EQUAL, ">="},
		{lexer.INT, "18"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.OPEN_CURLY_BRACKET, "{"},
		{lexer.PRINT, "\"You can vote\""},
		{lexer.PRINTLN, ""},
		{lexer.RETURN, "return"},
		{lexer.BOOL, "true"},
		{lexer.CLOSE_CURLY_BRACKET, "}"},
		{lexer.ELSE, "else"},
		{lexer.OPEN_CURLY_BRACKET, "{"},
		{lexer.PRINTLN, "\"You cannot vote kid!\""},
		{lexer.CLOSE_CURLY_BRACKET, "}"},
		{lexer.RETURN, "return"},
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
