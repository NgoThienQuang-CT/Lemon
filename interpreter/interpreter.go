// Package interpreter provides an interpreter that execute file and
// start a Read - Eval - Print loop interaction for Lemon programming language
package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"lemon/internal/evaluator"
	"lemon/internal/parser"
	"lemon/internal/value"
)

const (
	DefaultPrompt = ">>> "
	WaitPrompt    = "... "
)

func RunFile(filename string) {
	path, err := filepath.Abs(filename)
	if err != nil {
		fmt.Printf("Error: Invalid path '%s'\n", os.Args[1])
		os.Exit(1)
	}

	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: file '%s' not found\n", filename)
		os.Exit(74)
	}

	if info.Size() > 10*1024*1024 /* 10 MB */ {
		fmt.Fprintln(os.Stderr, "Error: file size is too large (>10MB)")
		os.Exit(74)
	}

	input, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(74)
	}

	program, error := parser.ParseProgram(string(input))
	if len(error) != 0 {
		fmt.Fprintf(os.Stderr, "At '%s':\n", filename)
		for _, msg := range error {
			fmt.Fprint(os.Stderr, msg)
		}
		os.Exit(65)
	}

	result := evaluator.Eval(program, value.NewScope())
	if result != nil {
		_, err := fmt.Fprintln(os.Stdout, result.Inspect())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error printing result: %v\n", err)
		}

		if result.Type() == value.ErrorValType {
			os.Exit(70)
		}
	}
}

func StartRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	lemonlogo, err := filepath.Abs("interpreter/lemon.txt")
	if err != nil {
		panic(err)
	}
	LEMON, err := os.ReadFile(lemonlogo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This the Lemon programming language!\n", user.Username)
	fmt.Printf("\033[32m%s\033[0m", string(LEMON))
	fmt.Println("Type '/exit' to exit the REPL.")

	var prompt string
	buffer := ""
	level := 0
	scanner := bufio.NewScanner(os.Stdin)
	scope := value.NewScope()

	for {
		if level > 0 {
			prompt = WaitPrompt
		} else {
			prompt = DefaultPrompt
		}

		_, err := fmt.Fprint(os.Stdin, prompt)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}

		scanned := scanner.Scan()
		if !scanned {
			fmt.Fprintf(os.Stderr, "Error: cannot read line\n")
		}

		line := scanner.Text()
		if line == "/exit" {
			return
		}

		buffer += (line + "\n")
		level += (strings.Count(line, "(") - strings.Count(line, ")"))
		level += (strings.Count(line, "{") - strings.Count(line, "}"))

		if level > 0 {
			continue
		}

		program, error := parser.ParseProgram(buffer)
		if len(error) != 0 {
			fmt.Fprintln(os.Stderr, "At <repl>:")
			for _, msg := range error {
				fmt.Fprint(os.Stderr, msg)
			}
			continue
		}

		result := evaluator.Eval(program, scope)
		if result != nil {
			_, err := fmt.Fprintln(os.Stdout, result.Inspect())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error printing result: %v\n", err)
			}
		}

		buffer = ""
	}
}
