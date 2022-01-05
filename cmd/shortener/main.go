package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/const-se/golang/internal/app/shortener/handler"
	"github.com/const-se/golang/internal/app/shortener/repository"
	"log"
	"net/http"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

func main() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server Address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "File Storage Path")
	flag.Parse()

	r := repository.NewRepository(cfg.FileStoragePath)
	h := handler.NewHandler(r, cfg.BaseURL)

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, h))
}
