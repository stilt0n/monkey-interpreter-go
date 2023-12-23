package main

import (
	"fmt"
	"monkey-pl/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()

	if err != nil {
		fmt.Println("🙊 Oh no! There was an error! See below: 🙊")
		panic(err)
	}

	fmt.Printf("🐵 Hello %s! Welcome to the Monkey programming language 🐵\n", user.Username)
	fmt.Printf("🐵🍌 Try out some commands! 🍌🐵\n\n")
	repl.Start(os.Stdin, os.Stdout)
}
