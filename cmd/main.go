package main

import (
	"fmt"
	"os"

	"github.com/sanity-io/litter"

	"github.com/KhushPatibandha/Kolon/src/interpreter/evaluator"
	"github.com/KhushPatibandha/Kolon/src/lexer"
	"github.com/KhushPatibandha/Kolon/src/parser"
)

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("Kolon v2.0.0")
		return
	} else if len(os.Args) == 2 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println(`Usage:
    kolon [flags]
    kolon [command] ([flags] <optional>)`)

		fmt.Println()

		fmt.Println(`Available Commands:
    'run: <file.kol>'                             Run a kolon file
    'debug: <file.kol> [--tokens | --ast]'        Debug a kolon file`)

		fmt.Println()

		fmt.Println(`Flags:
    -h, --help        help for kolon
    -v, --version     show version information
    --tokens          print tokens of the file [Command: 'debug:']
    --ast             print ast of the file [Command: 'debug:']`)

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

		p := parser.New(tokens, false)
		program, err := p.ParseProgram()
		if err != nil {
			fmt.Println("Error parsing program:", err)
			return
		}
		e := evaluator.New(false)
		_, err = e.Evaluate(program)
		if err != nil {
			fmt.Println("Error evaluating program:", err)
			return
		}
		return
	} else if len(os.Args) == 4 && os.Args[1] == "debug:" && (os.Args[3] == "--tokens" || os.Args[3] == "--ast") {
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
		if os.Args[3] == "--tokens" {
			for _, token := range tokens {
				token.Help()
			}
		} else {
			p := parser.New(tokens, false)
			program, err := p.ParseProgram()
			if err != nil {
				fmt.Println("Error parsing program:", err)
				return
			}
			litter.Dump(program)
		}
		return
	} else {
		fmt.Println("Not a valid command, use `--help` or `-h` for more information")
		return
	}
}
