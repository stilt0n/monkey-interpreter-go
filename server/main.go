package main

import (
	"log"
	"monkey-pl/api"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// This should prevent the server from going down from many types of bugs
			log.Println("Recovering from panic...")
			api.Serve()
		}
	}()
	api.Serve()
}
