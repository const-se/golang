package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/const-se/golang/internal/app/shortener/handler"
	"github.com/const-se/golang/internal/app/shortener/repository"
	"log"
	"net/http"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func main() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repository.NewRepository(), cfg.BaseURL)

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, h))
}
