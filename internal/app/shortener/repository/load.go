package repository

import (
	"bufio"
	"os"
)

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
