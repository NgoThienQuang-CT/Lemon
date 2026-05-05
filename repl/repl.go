// Package repl implements Read-Eval-Print loop for the shell
package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"lemon/evaluator"
	"lemon/parser"
	"lemon/value"
)

const (
	DEFAULT = ">>> "
	WAIT    = "... "
)

func Start(in io.Reader, out io.Writer) {
	var prompt string
	scanner := bufio.NewScanner(in)
	level := 0
	buffer := ""

	scope := value.NewScope()

	for {
		if level > 0 {
			prompt = WAIT
		} else {
			prompt = DEFAULT
		}

		_, promptErr := fmt.Fprint(out, prompt)
		if promptErr != nil {
			fmt.Println(promptErr)
			return
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}
		buffer += (line + "\n")

		level += (strings.Count(line, "(") - strings.Count(line, ")"))
		level += (strings.Count(line, "{") - strings.Count(line, "}"))

		if level > 0 {
			continue
		}

		program, errors := parser.ParseProgram(buffer)
		if len(errors) != 0 {
			printParserErrors(out, errors)
			buffer = ""
			continue
		}

		result := evaluator.Eval(program, scope)
		if result == nil {
			continue
		}

		_, err := io.WriteString(out, result.Inspect()+"\n")
		if err != nil {
			fmt.Println(err)
			continue
		}

		buffer = ""
	}
}

func printParserErrors(out io.Writer, errors []string) {
	_, err := io.WriteString(out, "at <repl>:\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, msg := range errors {
		_, err := io.WriteString(out, msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
