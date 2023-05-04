package main

import (
	"log"

	"github.com/kei-gnu/golang_http/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}