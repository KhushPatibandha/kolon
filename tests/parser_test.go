package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func Test23(t *testing.T) {
	input := "var x: int = 10;var y: int = 100;var foobar: int = 10000;const age: int = 100;const heh: string = \"hello\";"
	helper(t, []string{input})
}

func Test24(t *testing.T) {
	input := "fun: test(): (int) {return: 5;}fun: test1(): (int) {return: 100;}fun: test2(): (int) {return: 312413;}fun: test3(): (int, bool, string) {return: ((5 + 1), true, \"hello\");}"
	helper(t, []string{input})
}

func Test25(t *testing.T) {
	input := "fun: hehe(name: string, age: int): (bool, int) {var a: int = 5;return: (true, 5);}"
	helper(t, []string{input})
}

func Test26(t *testing.T) {
	input := "fun: main() {var a: int = 10;var b: int = 20;var c: int = 30;if: ((a > b)): {var d: int = 40;}else if: ((b > c)): {var e: int = 50;}else: {var f: int = 60;}}fun: add(a: int, b: int): (int) {return: (a + b);}"
	helper(t, []string{input})
}

func Test27(t *testing.T) {
	input := "fun: main() {var a: int = 10;var b: int = 20;var c: int = 30;if: ((a > b)): {var d: int = 40;}else if: ((b > c)): {var e: int = 50;}else: {var f: int = 60;}}fun: add(a: int, b: int): (int) {return: (a + b);}"
	helper(t, []string{input})
}

func Test28(t *testing.T) {
	input := "var a: int;var b:int;if: ((a < b)): {var c: int = 10;}else: {var d: int = 20;}"
	expected := "var a: int = 0;var b: int = 0;if: ((a < b)): {var c: int = 10;}else: {var d: int = 20;}"
	helper1(t, []string{input}, []string{expected})
}

func Test29(t *testing.T) {
	input := "var a: int = 0;var b: int = 0;a = b;"
	helper(t, []string{input})
}

func Test30(t *testing.T) {
	input := "var a: int = 0;a = 10;"
	helper(t, []string{input})
}

func Test31(t *testing.T) {
	input := "var a: bool = true;var b: int = 0;a = (!true);b = (-1);"
	helper(t, []string{input})
}

func Test32(t *testing.T) {
	input := "a = (5 + 5);b = (5 - 5);c = (5 * 5);d = (5 / 5);e = (5 % 5);f = (5 > 5);" +
		"g = (5 < 5);h = (5 >= 5);i = (5 <= 5);j = (5 == 5);k = (5 != 5);l = (5 && 5);m = (5 || 5);" +
		"n = (5 & 5);o = (5 | 5);p = (true == true);q = (false != false);r = (true && false);s = (true || false);"
	helper(t, []string{input})
}

func Test33(t *testing.T) {
	input := []string{
		"l = a + add(b * c) + d;",
		"l = add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8));",
		"l = add(a + b + c * d / f + g);",
		"l = 1 + (2 + 3) + 4;",
		"l = (5 + 5) * 2;",
		"l = 2 / (5 + 5);",
		"l = -(5 + 5);",
		"l = !(true == true);",
		"l = true;",
		"l = false;",
		"l = 3 > 5 == false;",
		"l = 3 < 5 == true;",
		"l = -a * b;",
		"l = !-a;",
		"l = a + b + c;",
		"l = a + b - c;",
		"l = a * b * c;",
		"l = a * b / c;",
		"l = a + b / c;",
		"l = a + b * c + d / e - f;",
		"l = 3 + 4;l = -5 * 5;",
		"l = 5 > 4 == 3 < 4;",
		"l = 5 < 4 != 3 > 4;",
		"l = 3 + 4 * 5 == 3 * 1 + 4 * 5;",
		"l = -5--;",
		"l = 1 == 1 && 2 == 2 && 3 == 3;",
	}
	expected := []string{
		"l = ((a + add((b * c))) + d);",
		"l = add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)));",
		"l = add((((a + b) + ((c * d) / f)) + g));",
		"l = ((1 + (2 + 3)) + 4);",
		"l = ((5 + 5) * 2);",
		"l = (2 / (5 + 5));",
		"l = (-(5 + 5));",
		"l = (!(true == true));",
		"l = true;",
		"l = false;",
		"l = ((3 > 5) == false);",
		"l = ((3 < 5) == true);",
		"l = ((-a) * b);",
		"l = (!(-a));",
		"l = ((a + b) + c);",
		"l = ((a + b) - c);",
		"l = ((a * b) * c);",
		"l = ((a * b) / c);",
		"l = (a + (b / c));",
		"l = (((a + (b * c)) + (d / e)) - f);",
		"l = (3 + 4);l = ((-5) * 5);",
		"l = ((5 > 4) == (3 < 4));",
		"l = ((5 < 4) != (3 > 4));",
		"l = ((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",
		"l = (-(5--));",
		"l = (((1 == 1) && (2 == 2)) && (3 == 3));",
	}
	helper1(t, input, expected)
}

func Test34(t *testing.T) {
	input := []string{
		"var t: int = 1;t = 5++;",
		"var t: int = 1;t = 5--;",
	}
	expected := []string{
		"var t: int = 1;t = (5++);",
		"var t: int = 1;t = (5--);",
	}
	helper1(t, input, expected)
}

func Test35(t *testing.T) {
	input := []string{
		"var a: bool = true;",
		"var b: bool = false;",
		"fun: add(a: int, b: int) {add(a, b);}",
		"var a: int = 1;(a++);",
		"var a: int = 1;(a--);",
	}
	helper(t, input)
}

func helper(t *testing.T, input []string) {
	for _, str := range input {
		tokens := lexer.Tokenizer(str)
		parser := parser.New(tokens, true)
		program, err := parser.ParseProgram()
		if err != nil {
			t.Fatalf("ParseProgram() returned error: %s", err)
		}
		assert.Equal(t, str, program.String())
	}
}

func helper1(t *testing.T, input []string, expected []string) {
	for i, str := range input {
		tokens := lexer.Tokenizer(str)
		parser := parser.New(tokens, true)
		program, err := parser.ParseProgram()
		if err != nil {
			t.Fatalf("ParseProgram() returned error: %s", err)
		}
		assert.Equal(t, expected[i], program.String())
	}
}
