package repository

import (
	"log"
	"sync"
)

type Repository interface {
	URL(id int) (string, error)
	SaveURL(url string) (int, error)
}

type repository struct {
	storagePath string
	cache       []string
	mutex       sync.Mutex
}

func NewRepository(storagePath string) Repository {
	cache, err := load(storagePath)
	if err != nil {
		log.Fatal(err)
	}

	return &repository{
		storagePath: storagePath,
		cache:       cache,
	}
}
