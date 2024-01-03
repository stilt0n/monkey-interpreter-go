package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey-pl/evaluator"
	"monkey-pl/lexer"
	"monkey-pl/object"
	"monkey-pl/parser"
	"runtime"
)

func getEvalOutputColor() []string {
	// windows terminal won't support this
	// so we just don't use a color
	if runtime.GOOS == "windows" {
		return []string{"", ""}
	}
	return []string{"\033[33m", "\033[0m"}

}

const PROMPT = "🐒 >> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	evalColorCodes := getEvalOutputColor()
	yellow := func(str string) string {
		return fmt.Sprintf("%s%s%s", evalColorCodes[0], str, evalColorCodes[1])
	}
	// this will allow let definitions to continue to be remembered
	env := object.NewEnvironment()
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

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, yellow(evaluated.Inspect()))
			io.WriteString(out, "\n")
		}
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
