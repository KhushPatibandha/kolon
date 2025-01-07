package main

import (
	"fmt"
	"os"

	"github.com/KhushPatibandha/Kolon/src/interpreter/evaluator"
	"github.com/KhushPatibandha/Kolon/src/interpreter/object"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path-to-kolon-file>")
		return
	}
	filePath := os.Args[1]
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	tokens := lexer.Tokenizer(string(bytes))
	// for _, token := range tokens {
	// 	token.Help()
	// }

	parser := parser.New(tokens)
	program := parser.ParseProgram()
	env := object.NewEnvironment()
	_, _, err = evaluator.Eval(program, env)
	if err != nil {
		fmt.Println("Error evaluating program:", err)
	}
}
