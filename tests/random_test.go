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
	run(t, "./testKolFiles/fac.kol", "120")
	run(t, "./testKolFiles/test.kol", "1\nhello!! 1\n2\nhello!! 1.1 hehe!! true\n1\n1.1\ntrue\nc")
	run(t, "./testKolFiles/test1.kol", "100")
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
	run(t, "./testKolFiles/test21.kol", "Error type checking program: variable `a` is undefined/not found")
	run(t, "./testKolFiles/test22.kol", "Error type checking program: missing `return` statement for function: callMe")
	run(t, "./testKolFiles/test23.kol", "Error parsing program: everything must be inside a function")
	run(t, "./testKolFiles/test24.kol", "Error type checking program: variable `b` is a constant, can't re-declare const variables")
	run(t, "./testKolFiles/test25.kol", "int\nfloat\nstring\nchar\nbool\nint[]\nstring[int]")
	run(t, "./testKolFiles/test26.kol", "[2, 3, 4]\n[2, 4]\nhus\nhs")
	run(t, "./testKolFiles/test27.kol", "Error type checking program: call expression must return only 1 value for `var` statement, got: 0")
	run(t, "./testKolFiles/test28.kol", "102\nhello\nhello\n123\n10")
	run(t, "./testKolFiles/test29.kol", "{}\n{\"khush\": 1}")
	run(t, "./testKolFiles/test30.kol", "hello!!\n10\ntrue\nHello\nw\n1.1\nhello!!")
	run(t, "./testKolFiles/test31.kol", "0\n1\nhere 2\n0\n1\nhere2 2\n0\n1\nhere3 2")
	run(t, "./testKolFiles/test32.kol", "1\n2\n3\n4")
	run(t, "./testKolFiles/test33.kol", "[0, 1, 2, 4, 5, 6, 7, 8, 0, 1, 2, 4, 5, 6, 7, 8, 0, 1, 2, 3, 4, 5, 6, 0, 1, 2, 3, 4, 5, 6]")
	run(t, "./testKolFiles/test34.kol", "[0, 2, 4, 6, 8, 10, 0, 2, 4, 6, 8, 10, 0, 2, 4, 6, 8, 10]")
	run(t, "./testKolFiles/test35.kol", "0\n2\n4\n6\n8\n10\n12")
	run(t, "./testKolFiles/test36.kol", "0\n1\n2\n3\n4\n5\n100")
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
