package main

import (
	"os"

	"github.com/KhushPatibandha/Kolon/src/lexer"
)

func main() {
	bytes, _ := os.ReadFile("./tests/dotKolFiles/all.kol")
	tokens := lexer.Tokenizer(string(bytes))
	for _, token := range tokens {
		token.Help()
	}
}
