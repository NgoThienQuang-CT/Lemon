package main

import (
	"fmt"
	"lemon/lexer"
	"lemon/token"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lemon [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		repl()
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Unable to read file. Error: ", err)
		return
	}

	var lemonLexer lexer.Lexer
	lemonLexer.Init(string(content))
	for tok := lemonLexer.NextToken(); tok.Type != token.EOF; tok = lemonLexer.NextToken() {
		fmt.Println(tok.String())
	}
}

func repl() {
}
