package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey-pl/lexer"
	"monkey-pl/parser"
)

const PROMPT = "🐒 >> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if line == ":exit" {
			io.WriteString(out, "\n🐵 See you next time!!! 🐵\n")
			return
		}
		lex := lexer.New(line)
		pars := parser.New(lex)
		program := pars.ParseProgram()

		if len(pars.Errors()) != 0 {
			printParserErrors(out, pars.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "\n🙊 Oh No! You typed something Monkey can't handle! 🙊\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
	io.WriteString(out, "\n")
}
