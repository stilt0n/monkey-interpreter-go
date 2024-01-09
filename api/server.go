package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"monkey-pl/evaluator"
	"monkey-pl/lexer"
	"monkey-pl/object"
	"monkey-pl/parser"
	"net/http"
	"os"
	"strings"
)

type EvalRequestBody struct {
	Code string `json:"code"`
}

type EvalResponse struct {
	Result string `json:"result"`
	// Errors are a valid thing for eval to send back
	// so I don't want to treat them as server errors
	// but giving that info will allow me to color them
	// differently on the front-end
	IsError bool `json:"isError"`
}

func Serve() {
	log.Printf("Running server on port :%d\n", 5150)
	http.HandleFunc("/eval", handleEvaluate)
	err := http.ListenAndServe(":5150", nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("The server is shutting down...")
		return
	}
	if err != nil {
		log.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func handleEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	parsedBody := EvalRequestBody{}
	err := fromJson(r.Body, &parsedBody)
	if err != nil {
		sendErr(w, err, 400)
		return
	}
	code := parsedBody.Code
	env := object.NewEnvironment()
	lex := lexer.New(code)
	prs := parser.New(lex)
	program := prs.ParseProgram()

	response := EvalResponse{}
	if len(prs.Errors()) != 0 {
		response.Result = strings.Join(prs.Errors(), "\n")
		response.IsError = true
		sendJson(w, func() (interface{}, error) {
			return response, nil
		})
		return
	}
	evaluated := evaluator.Eval(program, env)
	// TODO: Perhaps this should actually return a NULL object.Object
	if evaluated == nil {
		response.Result = "NULL"
	} else {
		response.Result = evaluated.Inspect()
	}
	response.IsError = false
	sendJson(w, func() (interface{}, error) {
		return response, nil
	})
}

func fromJson[T any](body io.Reader, target T) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	return json.Unmarshal(buf.Bytes(), &target)
}

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func sendJson[T any](w http.ResponseWriter, withData func() (T, error)) {
	setJsonHeader(w)

	data, serverErr := withData()
	if serverErr != nil {
		w.WriteHeader(500)
		serverErrJson, err := json.Marshal(&serverErr)
		if err != nil {
			log.Println(err)
			return
		}
		w.Write(serverErrJson)
		return
	}

	dataJson, err := json.Marshal(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(dataJson)
}

func sendErr(w http.ResponseWriter, err error, code int) {
	sendJson(w, func() (interface{}, error) {
		errorMessage := struct {
			Err string
		}{
			Err: err.Error(),
		}
		w.WriteHeader(code)
		return errorMessage, nil
	})
}
