package lexer_test

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/lexer"
)

func Test3(t *testing.T) {
	input := `
    // No accurate representation of Kolon
    // package main
    
    // import {
    //     println from stdio
    // }
    // // or you can do it like this -> import println from stdio

    // struct Person {
    //     var name: string
    //     var age: int
    // }

    fun main() {
        const a: int = 100
        const b: int = 100
        var bo: bool = true
        println(!bo)
        println(a + b)
        println(a - b)
        println(a * b)
        println(a / b)
        println(a % b)
        println(a != b)
        println(a == b)

        a++
        b--
        a += 10
        b -= 10
        a *= 10
        b /= 10
        a %= 10;

        ? . & |
 
        for: (var i: int = 0, i <= 10, i++) {
            println(i)
        }
        // will go from 0 to 9 (inclusive)
        for(var i: int in 0..10) {
            println(i)
        }

        a = 10
        b = 20
        if: (a == 10 && b == 20) {
            println("a is 10 and b is 20")
        }
        if: (a == 10 || b == 20) {
            println("a is 10 or b is 20")
        }
    }
    `

	tests := []struct {
		expectedType    lexer.TokenKind
		expectedLiteral string
	}{
		{lexer.FUN, "fun"},
		{lexer.IDENTIFIER, "main"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.OPEN_CURLY_BRACKET, "{"},
		{lexer.CONST, "const"},
		{lexer.IDENTIFIER, "a"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "int"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.INT, "100"},
		{lexer.CONST, "const"},
		{lexer.IDENTIFIER, "b"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "int"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.INT, "100"},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "bo"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "bool"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.BOOL, "true"},
		{lexer.PRINTLN, "!bo"},
		{lexer.PRINTLN, "a + b"},
		{lexer.PRINTLN, "a - b"},
		{lexer.PRINTLN, "a * b"},
		{lexer.PRINTLN, "a / b"},
		{lexer.PRINTLN, "a % b"},
		{lexer.PRINTLN, "a != b"},
		{lexer.PRINTLN, "a == b"},
		{lexer.IDENTIFIER, "a"},
		{lexer.PLUS_PLUS, "++"},
		{lexer.IDENTIFIER, "b"},
		{lexer.MINUS_MINUS, "--"},
		{lexer.IDENTIFIER, "a"},
		{lexer.PLUS_EQUAL, "+="},
		{lexer.INT, "10"},
		{lexer.IDENTIFIER, "b"},
		{lexer.MINUS_EQUAL, "-="},
		{lexer.INT, "10"},
		{lexer.IDENTIFIER, "a"},
		{lexer.STAR_EQUAL, "*="},
		{lexer.INT, "10"},
		{lexer.IDENTIFIER, "b"},
		{lexer.SLASH_EQUAL, "/="},
		{lexer.INT, "10"},
		{lexer.IDENTIFIER, "a"},
		{lexer.PERCENT_EQUAL, "%="},
		{lexer.INT, "10"},
		{lexer.SEMI_COLON, ";"},
		{lexer.QUESTION_MARK, "?"},
		{lexer.DOT, "."},
		{lexer.AND, "&"},
		{lexer.OR, "|"},
		{lexer.FOR, "for"},
		{lexer.COLON, ":"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "i"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "int"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.INT, "0"},
		{lexer.COMMA, ","},
		{lexer.IDENTIFIER, "i"},
		{lexer.LESS_THAN_EQUAL, "<="},
		{lexer.INT, "10"},
		{lexer.COMMA, ","},
		{lexer.IDENTIFIER, "i"},
		{lexer.PLUS_PLUS, "++"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.OPEN_CURLY_BRACKET, "{"},
		{lexer.PRINTLN, "i"},
		{lexer.CLOSE_CURLY_BRACKET, "}"},
		{lexer.FOR, "for"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "i"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "int"},
		{lexer.IN, "in"},
		{lexer.INT, "0"},
		{lexer.DOT_DOT, ".."},
		{lexer.INT, "10"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.OPEN_CURLY_BRACKET, "{"},
		{lexer.PRINTLN, "i"},
		{lexer.CLOSE_CURLY_BRACKET, "}"},
		{lexer.IDENTIFIER, "a"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.INT, "10"},
		{lexer.IDENTIFIER, "b"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.INT, "20"},
		{lexer.IF, "if"},
		{lexer.COLON, ":"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "a"},
		{lexer.DOUBLE_EQUAL, "=="},
		{lexer.INT, "10"},
		{lexer.AND_AND, "&&"},
		{lexer.IDENTIFIER, "b"},
		{lexer.DOUBLE_EQUAL, "=="},
		{lexer.INT, "20"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.OPEN_CURLY_BRACKET, "{"},
		{lexer.PRINTLN, "\"a is 10 and b is 20\""},
		{lexer.CLOSE_CURLY_BRACKET, "}"},
		{lexer.IF, "if"},
		{lexer.COLON, ":"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "a"},
		{lexer.DOUBLE_EQUAL, "=="},
		{lexer.INT, "10"},
		{lexer.OR_OR, "||"},
		{lexer.IDENTIFIER, "b"},
		{lexer.DOUBLE_EQUAL, "=="},
		{lexer.INT, "20"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.OPEN_CURLY_BRACKET, "{"},
		{lexer.PRINTLN, "\"a is 10 or b is 20\""},
		{lexer.CLOSE_CURLY_BRACKET, "}"},
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
