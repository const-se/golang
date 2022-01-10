package repository

import (
	"bufio"
	"os"
)

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
