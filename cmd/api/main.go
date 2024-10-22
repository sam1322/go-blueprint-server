package main

import (
	"fmt"
	"new_project/internal/authenticate"
	"new_project/internal/server"
)

func main() {
	fmt.Println("Starting the server")
	authenticate.NewAuth()
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
