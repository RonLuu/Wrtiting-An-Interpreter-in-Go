package repl

import (
	"Chapter_3/evaluator"
	"Chapter_3/lexer"
	"Chapter_3/parser"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">> "

const MONKEY_FACE = ` 
           _,_
 .--.   . -" "- .   .--.
 / .. \/ .-. .-. \/ .. \
| |  '| /   Y   \ |'  | |
| \   \ \ 0 | 0 / /   / |
 \ '- ,\.-"""""-./, -' /
   ''-' /_ ^ ^ _\ '-''
       | \._ _./ |
       \ \ '~' / /
       '._'-=-'_.'
          '---'
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	fmt.Print(MONKEY_FACE)
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		lexer := lexer.NewLexer(line)
		parser := parser.NewParser(lexer)
		program := parser.ParseProgram()
		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
