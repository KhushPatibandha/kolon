package main

import (
	"fmt"
	"os"

	"github.com/KhushPatibandha/Kolon/src/interpreter/evaluator"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/object"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println("Kolon v1.0.2")
		return
	} else if len(os.Args) == 3 && os.Args[1] == "run:" {
		filePath := os.Args[2]
		if filePath[len(filePath)-4:] != ".kol" {
			fmt.Println("Error: File should have .kol extension")
			return
		}
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		tokens := lexer.Tokenizer(string(bytes))
		// for _, token := range tokens {
		// 	token.Help()
		// }

		p := parser.New(tokens, false)
		program, err := p.ParseProgram()
		if err != nil {
			fmt.Println("Error parsing program:", err)
			return
		}
		typeCheckerEnv := parser.NewEnvironment()
		err = parser.TypeCheckProgram(program, typeCheckerEnv, false)
		if err != nil {
			fmt.Println("Error type checking program:", err)
			return
		}
		env := object.NewEnvironment()
		_, _, err = evaluator.Eval(program, env, false)
		if err != nil {
			fmt.Println("Error evaluating program:", err)
			return
		}
	} else {
		fmt.Println("Usage: kolon run: <path-to-kolon-file>")
		return
	}
}
