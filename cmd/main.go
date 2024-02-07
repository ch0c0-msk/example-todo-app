package main

import (
	"log"

	todo "github.com/ch0c0-msk/example-todo-app"
)

func main() {
	srv := new(todo.Server)
	const port = "8000"
	if err := srv.Run(port); err != nil {
		log.Fatal(err.Error())
	}
}
