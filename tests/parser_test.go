package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func Test23(t *testing.T) {
	input := map[string]bool{
		"var x: int = 10;var y: int = 100;var foobar: int = 10000;const age: int = 100;const heh: string = \"hello\";": true,

		"fun: test(): (int) {return: 5;}fun: test1(): (int) {return: 100;}fun: test2(): (int) " +
			"{return: 312413;}fun: test3(): (int, bool, string) {return: ((5 + 1), true, \"hello\");}": true,

		"fun: hehe(name: string, age: int): (bool, int) {var a: int = 5;return: (true, 5);}": true,

		"fun: main() {var a: int = 10;var b: int = 20;var c: int = 30;if: ((a > b)): {var d: int = 40;}else " +
			"if: ((b > c)): {var e: int = 50;}else: {var f: int = 60;}}fun: add(a: int, b: int): (int) {return: (a + b);}": true,

		"var a: int = 0;var b: int = 0;a = b;": true,

		"var a: int = 0;a = 10;": true,

		"var a: bool = true;var b: int = 0;a = (!true);b = (-1);": true,

		"var a: bool = true;": true,

		"var b: bool = false;": true,

		"fun: add(a: int, b: int) {add(a, b);}": true,

		"fun: add(c: int, d: int) {var a: int = c;var b: int = d;add(a, b);}": true,

		"fun: main() {var a: int = 1;var b: int = a;}": true,

		"var a: int = 1;(a++);": true,

		"var a: int = 1;(a--);": true,

		"var a: int[string] = {};": true,

		"var a: int[] = [1];":   true,
		"var a: int[] = [1.1];": false,

		"var a: int[float] = {1: 1.1};":   true,
		"var a: int[float] = {1.1: 1.1};": false,

		"fun: add(a: int, b: int): (int);fun: main() {add(1, 2);}": false,
		"fun: add(a: int, b: int);fun: main() {add(1, 2);}":        false,
	}
	helper(t, []map[string]bool{input})
}

func Test24(t *testing.T) {
	test := map[string]string{
		"var a: int;var b:int;if: ((a < b)): {var c: int = 10;}else: " +
			"{var d: int = 20;}": "var a: int = 0;var b: int = 0;if: ((a < b)): {var c: int = 10;}else: {var d: int = 20;}",

		"var t: int = 1;t = 5++;": "var t: int = 1;t = (5++);",

		"var t: int = 1;t = 5--;": "var t: int = 1;t = (5--);",
	}
	helper1(t, []map[string]string{test})
}

func Test25(t *testing.T) {
	test := map[string]string{
		"fun: add(a: int): (int) {var b: int, var c: int, var d: int, var l: int = 0, 0, 0, 0;l = a + add(b * c) + d;}": "fun: " +
			"add(a: int): (int) {var b: int = 0;var c: int = 0;var d: int = 0;var l: int = 0;l = ((a + add((b * c))) + d);}",

		"fun: add(a: int, b: int, c: int, d: int, e: int, f: int): (int) {var l: int = add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8, 3, 4, 5, 6));}": "fun: " +
			"add(a: int, b: int, c: int, d: int, e: int, f: int): (int) {var l: int = add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8), 3, 4, 5, 6));}",

		"fun: add(a: int): (int) {return: a;}fun: sub(a: int, b: int, c: int, d: int, f: int, g: int) {var l: int = add(a + b + c * d / f + g);}": "fun: add(a: int): (int) {return: a;}fun: " +
			"sub(a: int, b: int, c: int, d: int, f: int, g: int) {var l: int = add((((a + b) + ((c * d) / f)) + g));}",

		"var l: int = 1 + (2 + 3) + 4;": "var l: int = ((1 + (2 + 3)) + 4);",

		"var l: int = (5 + 5) * 2;": "var l: int = ((5 + 5) * 2);",

		"var l: int = 2 / (5 + 5);": "var l: int = (2 / (5 + 5));",

		"var l: int = -(5 + 5);": "var l: int = (-(5 + 5));",

		"var l: bool = !(true == true);": "var l: bool = (!(true == true));",

		"var l: bool = true;": "var l: bool = true;",

		"var l: bool = false;": "var l: bool = false;",

		"var l: bool = 3 > 5 == false;": "var l: bool = ((3 > 5) == false);",

		"var l: bool = 3 < 5 == true;": "var l: bool = ((3 < 5) == true);",

		"var a: int;var b: int;var l: int = -a * b;": "var a: int = 0;var b: int = 0;var l: int = ((-a) * b);",

		"var a: int;var b: int;var c: int;var l: int = a + b + c;": "var a: int = 0;var b: int = 0;var c: int = 0;var l: int = ((a + b) + c);",

		"var a: int;var b: int;var c: int;var l: int = a + b - c;": "var a: int = 0;var b: int = 0;var c: int = 0;var l: int = ((a + b) - c);",

		"var a: int;var b: int;var c: int;var l: int = a * b * c;": "var a: int = 0;var b: int = 0;var c: int = 0;var l: int = ((a * b) * c);",

		"var a: int;var b: int;var c: int;var l: int = a * b / c;": "var a: int = 0;var b: int = 0;var c: int = 0;var l: int = ((a * b) / c);",

		"var a: int;var b: int;var c: int;var l: int = a + b / c;": "var a: int = 0;var b: int = 0;var c: int = 0;var l: int = (a + (b / c));",

		"var a: int;var b: int;var c: int;var d: int;var e: int;var f: int;var l: int = a + b * c + d / e - f;": "var a: int = 0;var b: int = 0;" +
			"var c: int = 0;var d: int = 0;var e: int = 0;var f: int = 0;var l: int = (((a + (b * c)) + (d / e)) - f);",

		"var l: int = 3 + 4;l = -5 * 5;": "var l: int = (3 + 4);l = ((-5) * 5);",

		"var l: bool = 5 > 4 == 3 < 4;": "var l: bool = ((5 > 4) == (3 < 4));",

		"var l: bool = 5 < 4 != 3 > 4;": "var l: bool = ((5 < 4) != (3 > 4));",

		"var l: bool = 3 + 4 * 5 == 3 * 1 + 4 * 5;": "var l: bool = ((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",

		"var l: int = -5--;": "var l: int = (-(5--));",

		"var l: bool = 1 == 1 && 2 == 2 && 3 == 3;": "var l: bool = (((1 == 1) && (2 == 2)) && (3 == 3));",
	}
	helper1(t, []map[string]string{test})
}

func Test26(t *testing.T) {
	input := map[string]bool{
		"var a: int = (5 + 5);":           true,
		"var b: int = (5 - 5);":           true,
		"var c: int = (5 * 5);":           true,
		"var d: int = (5 / 5);":           true,
		"var e: int = (5 % 5);":           true,
		"var n: int = (5 & 5);":           true,
		"var o: int = (5 | 5);":           true,
		"var f: bool = (5 > 5);":          true,
		"var g: bool = (5 < 5);":          true,
		"var h: bool = (5 >= 5);":         true,
		"var i: bool = (5 <= 5);":         true,
		"var j: bool = (5 == 5);":         true,
		"var k: bool = (5 != 5);":         true,
		"var l: bool = (5 && 5);":         false,
		"var m: bool = (5 || 5);":         false,
		"var p: bool = (true == true);":   true,
		"var q: bool = (false != false);": true,
		"var r: bool = (true && false);":  true,
		"var s: bool = (true || false);":  true,
	}
	helper(t, []map[string]bool{input})
}

func Test27(t *testing.T) {
	input := map[string]string{
		"fun: add(a: int, b: int): (int);fun: main() {add(1, 2);}fun: add(a: int, b: int): (int) {return: (a + b);}": "fun: add(a: int, b: int): (int) {return: (a + b);}fun: main() {add(1, 2);}",
		"fun: add(a: int, b: int);fun: main() {add(1, 2);}fun: add(a: int, b: int) {}":                               "fun: add(a: int, b: int) {}fun: main() {add(1, 2);}",
	}
	helper1(t, []map[string]string{input})
}

func helper(t *testing.T, input []map[string]bool) {
	for _, test := range input {
		for key, val := range test {
			tokens := lexer.Tokenizer(key)
			parser := parser.New(tokens, true)
			program, err := parser.ParseProgram()
			if val {
				if err != nil {
					t.Fatalf("ParseProgram() returned error: %s", err)
				}
				assert.Equal(t, key, program.String())
			} else {
				assert.Error(t, err)
			}
		}
	}
}

func helper1(t *testing.T, test []map[string]string) {
	for _, pair := range test {
		for input, expected := range pair {
			tokens := lexer.Tokenizer(input)
			parser := parser.New(tokens, true)
			program, err := parser.ParseProgram()
			if err != nil {
				t.Fatalf("ParseProgram() returned error: %s", err)
			}
			assert.Equal(t, expected, program.String())
		}
	}
}
