package repository

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

var ErrURLNotFound = fmt.Errorf("URL not found")

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

func (r *repository) URL(id int) (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if id >= len(r.cache) {
		return "", ErrURLNotFound
	}

	url := r.cache[id]

	return url, nil
}

func (r *repository) SaveURL(url string) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.cache = append(r.cache, url)
	id := len(r.cache) - 1

	if r.storagePath != "" {
		file, err := os.OpenFile(r.storagePath, os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			return 0, err
		}

		defer func() {
			_ = file.Close()
		}()

		writer := bufio.NewWriter(file)
		defer func() {
			_ = writer.Flush()
		}()

		if _, err = writer.WriteString(url + "\n"); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func load(path string) (cache []string, err error) {
	if path == "" {
		return
	}

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if scanner.Err() != nil {
			return
		}

		url := scanner.Text()
		if len(url) > 0 {
			cache = append(cache, url)
		}
	}

	return
}
