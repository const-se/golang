package repository

import (
	"sync"
)

type Repository interface {
	URL(id int) (string, error)
	SaveURL(url string) (int, error)
}

type repository struct {
	storage []string
	mutex   sync.Mutex
}

func NewRepository() Repository {
	return new(repository)
}
