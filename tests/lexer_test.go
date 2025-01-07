package tests

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/lexer"
)

func Test1(t *testing.T) {
	input := `
    fun main() {
        // No accurate representation of Kolon
        // this should not print
        var someString: string = "hey"
        var someChar: char = 'c'
        var someInt: int = 2147483647
        var someOtherInt: int = 9223372036854775807
        var someFloat: float = 1.1
        var someOtherFloat: float = 3.14
        var someAnotherFloat: float= 3.141592653589793
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
		{lexer.IDENTIFIER, "someOtherInt"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "int"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.INT, "9223372036854775807"},
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
		{lexer.IDENTIFIER, "someAnotherFloat"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "float"},
		{lexer.EQUAL_ASSIGN, "="},
		{lexer.FLOAT, "3.141592653589793"},
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

func Test2(t *testing.T) {
	input := `
    // No accurate representation of Kolon
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
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "ok"},
		{lexer.CLOSE_BRACKET, ")"},
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
		{lexer.IDENTIFIER, "print"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.STRING, "\"You can vote\""},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.RETURN, "return"},
		{lexer.BOOL, "true"},
		{lexer.CLOSE_CURLY_BRACKET, "}"},
		{lexer.ELSE, "else"},
		{lexer.OPEN_CURLY_BRACKET, "{"},
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.STRING, "\"You cannot vote kid!\""},
		{lexer.CLOSE_BRACKET, ")"},
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

        & |
 
        for: (var i: int = 0, i <= 10, i++) {
            println(i)
        }
        // will go from 0 to 9 (inclusive)
        for(var i: int 0 10) {
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
        continue
        break
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
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.NOT, "!"},
		{lexer.IDENTIFIER, "bo"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "a"},
		{lexer.PLUS, "+"},
		{lexer.IDENTIFIER, "b"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "a"},
		{lexer.DASH, "-"},
		{lexer.IDENTIFIER, "b"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "a"},
		{lexer.STAR, "*"},
		{lexer.IDENTIFIER, "b"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "a"},
		{lexer.SLASH, "/"},
		{lexer.IDENTIFIER, "b"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "a"},
		{lexer.PERCENT, "%"},
		{lexer.IDENTIFIER, "b"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "a"},
		{lexer.NOT_EQUAL, "!="},
		{lexer.IDENTIFIER, "b"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "a"},
		{lexer.DOUBLE_EQUAL, "=="},
		{lexer.IDENTIFIER, "b"},
		{lexer.CLOSE_BRACKET, ")"},
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
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "i"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.CLOSE_CURLY_BRACKET, "}"},
		{lexer.FOR, "for"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.VAR, "var"},
		{lexer.IDENTIFIER, "i"},
		{lexer.COLON, ":"},
		{lexer.TYPE, "int"},
		{lexer.INT, "0"},
		{lexer.INT, "10"},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.OPEN_CURLY_BRACKET, "{"},
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.IDENTIFIER, "i"},
		{lexer.CLOSE_BRACKET, ")"},
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
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.STRING, "\"a is 10 and b is 20\""},
		{lexer.CLOSE_BRACKET, ")"},
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
		{lexer.IDENTIFIER, "println"},
		{lexer.OPEN_BRACKET, "("},
		{lexer.STRING, "\"a is 10 or b is 20\""},
		{lexer.CLOSE_BRACKET, ")"},
		{lexer.CLOSE_CURLY_BRACKET, "}"},
		{lexer.CONTINUE, "continue"},
		{lexer.BREAK, "break"},
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
