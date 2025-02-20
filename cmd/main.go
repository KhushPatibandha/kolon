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
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("Kolon v1.0.2")
		return
	} else if len(os.Args) == 2 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println(`Usage:
        kolon [flags]
        kolon [command] ([flags] <optional>)`)

		fmt.Println()

		fmt.Println(`Available Commands:
        'run: <file.kol> [optional flags]'        Run a kolon file`)

		fmt.Println()

		fmt.Println(`Flags:
        -h, --help        help for kolon
        -v, --version     show version information
        --print-tokens    print tokens of the file [optional: 'run:']`)

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
	} else if len(os.Args) == 4 && os.Args[1] == "run:" && os.Args[3] == "--print-tokens" {
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
		for _, token := range tokens {
			token.Help()
		}
	} else {
		fmt.Println("Not a valid command, use `--help` or `-h` for more information")
		return
	}
}
