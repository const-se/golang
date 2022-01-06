package main

import (
	"github.com/const-se/golang/internal/app/shortener/handler"
	"github.com/const-se/golang/internal/app/shortener/repository"
	"log"
	"net/http"
)

func main() {
	h := handler.NewHandler(repository.NewRepository())

	log.Fatal(http.ListenAndServe(":8080", h))
}
