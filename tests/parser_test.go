package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func Test23(t *testing.T) {
	input := "var x: int = 10;var y: int = 100;var foobar: int = 10000;const age: int = 100;const heh: string = \"hello\";"
	helper(t, input)
}

func Test24(t *testing.T) {
	input := "return: 5;return: 100;return: 312413;return: ((5 + 1), true, \"hello\");"
	helper(t, input)
}

func Test25(t *testing.T) {
	input := "fun: hehe(name: string, age: int): (bool, int) {var a: int = 5;return: (true, 5);}"
	helper(t, input)
}

func Test26(t *testing.T) {
	input := "fun: main() {var a: int = 10;var b: int = 20;var c: int = 30;if: ((a > b)): {var d: int = 40;}else if: ((b > c)): {var e: int = 50;}else: {var f: int = 60;}}fun: add(a: int, b: int): (int) {return: (a + b);}"
	helper(t, input)
}

func helper(t *testing.T, input string) {
	tokens := lexer.Tokenizer(input)
	parser := parser.New(tokens, true)
	program, err := parser.ParseProgram()
	if err != nil {
		t.Fatalf("ParseProgram() returned error: %s", err)
	}
	assert.Equal(t, input, program.String())
}
