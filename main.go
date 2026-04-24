package main

import (
	"fmt"
	"lemon/compiler/lexer"
	"lemon/compiler/token"
	"os"
)

func main() {
	fmt.Println("Hello, Lemon!")

	filename := os.Args[1]
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	lemonLexer := lexer.New(filename, string(source))

	for {
		tok := lemonLexer.LexToken()
		fmt.Println(tok.String())

		if tok.Type == token.EOF {
			return
		}
	}
}
