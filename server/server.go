package server

import (
	"errors"
	"fmt"
	"io"
	"monkey-pl/evaluator"
	"monkey-pl/lexer"
	"monkey-pl/object"
	"monkey-pl/parser"
	"net/http"
	"os"
)

// This should maybe eventually be the real main?
func Main() {
	fmt.Printf("Running server on port :%d\n", 5150)
	// This is a sanity check since I'm not really very
	// familiar with Go's http api yet
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/eval", evaluateMonkeyCode)
	err := http.ListenAndServe(":5150", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("The server is shutting down...")
		return
	}
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the root route!")
	io.WriteString(w, "Test")
}

// TODO: Make this function actually read and send data
func evaluateMonkeyCode(w http.ResponseWriter, r *http.Request) {
	code := "let standIn = 10; standIn;"
	env := object.NewEnvironment()
	lex := lexer.New(code)
	prs := parser.New(lex)
	program := prs.ParseProgram()

	if len(prs.Errors()) != 0 {
		fmt.Println("There was an error that should be sent by response writer")
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
		return
	}

	fmt.Println("NULL")
}
