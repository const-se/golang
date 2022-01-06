package main

import (
	"github.com/const-se/golang/internal/app/handlers"
	"log"
	"net/http"
)

func main() {
	handler := new(handlers.Shortener)
	http.HandleFunc("/", handler.Handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
