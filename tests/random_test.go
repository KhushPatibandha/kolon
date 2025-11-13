package tests

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestRandom(t *testing.T) {
	run(t, "./testKolFiles/fibo.kol", "0\n1\n1\n2\n3")
	run(t, "./testKolFiles/fiboRec.kol", "55")
	run(t, "./testKolFiles/fac.kol", "120")
	run(t, "./testKolFiles/test.kol", "1\nhello!! 1\n2\nhello!! 1.1 hehe!! true\n1\n1.1\ntrue\nc")
	run(t, "./testKolFiles/test1.kol", "Error parsing program: can't override a built-in function, function `len` already exists")
	run(t, "./testKolFiles/test2.kol", "true")
	run(t, "./testKolFiles/test3.kol", "true")
	run(t, "./testKolFiles/test4.kol", "true\n[\"khush\", \"hehe\"]")
	run(t, "./testKolFiles/test5.kol", "[1, 2, 3]\n4\n[5, 6, 7]\n[1, 2, 3, 10]\n[1, 2, 3]\n[2, 3]\n[2, 10, 3]\n[2, 3]")
	run(t, "./testKolFiles/test6.kol", "3")
	// run(t, "./testKolFiles/test7.kol", "{\"khush\": 1, \"heeh\": 2}\n{\"yo\": 101, \"hello\": 100}\n{\"hello\": 1, \"yo\": 101}\n{\"hello\": 1, \"yo\": 101, \"hehe\": 1}\n{\"yo\": 101, \"hehe\": 1}")
	run(t, "./testKolFiles/test8.kol", "[\"khush\", \"heeh\"]\n[\"hello\", \"yo\"]")
	run(t, "./testKolFiles/test9.kol", "[\"khush\", \"heeh\"]\n[\"hello\", \"yo\"]\n10\ntrue")
	// run(t, "./testKolFiles/test10.kol", "{\"heeh\": 2, \"khush\": 1}\n1")
	run(t, "./testKolFiles/test11.kol", "[0, 1, 2, 3, 4, 6, 7, 8, 9, 0, 1, 2, 3, 4, 6, 7, 8, 9, 0, 1, 2, 3, 4, 6, 7, 8, 9]")
	run(t, "./testKolFiles/test12.kol", "310")
	run(t, "./testKolFiles/test13.kol", "310")
	run(t, "./testKolFiles/test14.kol", "1\n0")
	run(t, "./testKolFiles/test15.kol", "2")
	run(t, "./testKolFiles/test16.kol", "1")
	run(t, "./testKolFiles/test17.kol", "2")
	run(t, "./testKolFiles/test18.kol", "-1")
	run(t, "./testKolFiles/test19.kol", "0\n-2")
	run(t, "./testKolFiles/test20.kol", "110")
	run(t, "./testKolFiles/test21.kol", "Error parsing program: variable `a` is undefined/not found")
	run(t, "./testKolFiles/test22.kol", "Error parsing program: function `callMe` must have a `return` statement at the end of all branches")
	run(t, "./testKolFiles/test23.kol", "Error parsing program: everything must be inside a function")
	run(t, "./testKolFiles/test24.kol", "Error parsing program: variable `b` is a constant, can't re-declare const variables")
	run(t, "./testKolFiles/test25.kol", "int\nfloat\nstring\nchar\nbool\nint[]\nstring[int]")
	run(t, "./testKolFiles/test26.kol", "[2, 3, 4]\n[2, 4]\nhus\nhs")
	run(t, "./testKolFiles/test27.kol", "Error parsing program: variable (`var`) and constant (`const`) declarations must be assigned a single value, got: 0. in case of call expression, it must return a single value")
	run(t, "./testKolFiles/test28.kol", "102\nhello\nhello\n123\n10")
	run(t, "./testKolFiles/test29.kol", "{}\n{\"khush\": 1}")
	run(t, "./testKolFiles/test30.kol", "hello!!\n10\ntrue\nHello\nw\n1.1\nhello!!")
	run(t, "./testKolFiles/test31.kol", "0\n1\nhere 2\n0\n1\nhere2 2\n0\n1\nhere3 2")
	run(t, "./testKolFiles/test32.kol", "1\n2\n3\n4")
	run(t, "./testKolFiles/test33.kol", "[0, 1, 2, 4, 5, 6, 7, 8, 0, 1, 2, 4, 5, 6, 7, 8, 0, 1, 2, 3, 4, 5, 6, 0, 1, 2, 3, 4, 5, 6]")
	run(t, "./testKolFiles/test34.kol", "[0, 2, 4, 6, 8, 10, 0, 2, 4, 6, 8, 10, 0, 2, 4, 6, 8, 10]")
	run(t, "./testKolFiles/test35.kol", "0\n2\n4\n6\n8\n10\n12")
	run(t, "./testKolFiles/test36.kol", "0\n1\n2\n3\n4\n5\n100")
	run(t, "./testKolFiles/test37.kol", "10.0\n10.1111\nfloat")
	run(t, "./testKolFiles/test38.kol", "11\n10\n11\nint\n1\n10\n11\nint\n65\n99\nint\nError evaluating program: Error converting string to int, can't convert: 10.1")
	run(t, "./testKolFiles/test41.kol", "Error parsing program: `main` function must not take in any parameters and must not return anything, since it is the starting point of the program")
	run(t, "./testKolFiles/test42.kol", "")
	run(t, "./testKolFiles/test43.kol", "int[]\nint[]\nint[]\nint[int[string][]]")
	run(t, "./testKolFiles/test44.kol", "[1, 2, 3]\n[1, 2, 3]\n[1, 2, 3, 4]\n[1, 2, 3, 4]\ntrue\ntrue\ntrue\nfalse\ntrue")
	run(t, "./testKolFiles/test45.kol", "[1, 2, 3]\n[1, 2, 3]\n[1, 2, 3, 4]\n[1, 2, 3, 4]\ntrue\ntrue\nfalse\ntrue")
	run(t, "./testKolFiles/test46.kol", "510-5-101032020250603037375001236420")
	run(t, "./testKolFiles/test47.kol", "truefalsefalsetruetruefalsetruetruefalsefalsefalsefalsetruefalsefalsetruetruetruefalsetruetruefalsetruetruefalsefalsefalsetruetruetruefalsefalsefalsetruetruetruetruefalsefalsetruefalsefalsetruefalsetruefalsetruefalsetruefalsetruefalsetruetruetruetruefalsetruefalsetruefalsefalsetruetruefalsefalsetrue")
	run(t, "./testKolFiles/test48.kol", "5.010.05.5-5.0-10.0-1.51.512.097.656250.023.7530.5-5.060.038.7552.87552.87547.2578.752.50.5")
	run(t, "./testKolFiles/test49.kol", "Hello, World!, Hello, World!, ab")
	run(t, "./testKolFiles/test50.kol", "64-4-6-411.19.1-9.1-11.1")
	run(t, "./testKolFiles/test51.kol", "10102010302040")
	run(t, "./testKolFiles/test52.kol", "truetruefalsetrue101020truefalse24.1helloa1010trueHellow1.1")
	run(t, "./testKolFiles/test53.kol", "4.0\n4.0\n3.0\n3.0\n3.0\n3.0\n-1.0\n-2.0")
	run(t, "./testKolFiles/test54.kol", "1.0\n2.0\n2.0\n1.1116\n1.11155579001\n1.1")
}

func run(t *testing.T, filePath string, expectedOutput string) {
	cmd := exec.Command("../kolon", "run:", filePath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		return
	}

	output := stdout.String()
	output = strings.TrimSpace(output)
	fmt.Println(output)
	if output != expectedOutput {
		t.Errorf("unexpected output for %s: got %q, want %q", filePath, output, expectedOutput)
	}
}
