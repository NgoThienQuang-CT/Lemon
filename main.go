package main

import (
	"fmt"
	"os"
	"os/user"

	"lemon/parser"
	"lemon/repl"
)

const LEMON = `
                                                 :                    
                     ,;                         t#,     L.            
            i      f#i                         ;##W.    EW:        ,ft
           LE    .E#t            ..       :   :#L:WE    E##;       t#E
          L#E   i#W,            ,W,     .Et  .KG  ,#D   E###t      t#E
         G#W.  L#D.            t##,    ,W#t  EE    ;#f  E#fE#f     t#E
        D#K. :K#Wfff;         L###,   j###t f#.     t#i E#t D#G    t#E
       E#K.  i##WLLLLt      .E#j##,  G#fE#t :#G     GK  E#t  f#E.  t#E
     .E#E.    .E#L         ;WW; ##,:K#i E#t  ;#L   LW.  E#t   t#K: t#E
    .K#E        f#E:      j#E.  ##f#W,  E#t   t#f f#:   E#t    ;#W,t#E
   .K#D          ,WW;   .D#L    ###K:   E#t    f#D#;    E#t     :K#D#E
  .W#G            .D#; :K#t     ##D.    E#t     G#t     E#t      .E##E
 :W##########Wt     tt ...      #G      ..       t      ..         G#E
 :,,,,,,,,,,,,,.                j                                   fE
                                                                     ,
`

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lemon [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runRepl()
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Unable to read file. Error: ", err)
		return
	}

	program, errors := parser.ParseProgram(string(content))
	if printErrors(filename, errors) {
		return
	}

	fmt.Println(program.String())
}

func runRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Lemon programing language!",
		user.Username)
	fmt.Printf("\033[32m%s\033[0m", LEMON)
	fmt.Println("Type 'exit' to exit the repl.")
	repl.Start(os.Stdin, os.Stdout)
}

func printErrors(filename string, errors []string) bool {
	if len(errors) == 0 {
		return false
	}

	fmt.Println("at <" + filename + ">:")
	for _, msg := range errors {
		fmt.Print(msg)
	}

	return true
}
