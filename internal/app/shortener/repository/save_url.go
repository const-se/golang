package repository

func (r *repository) SaveURL(url string) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.storage = append(r.storage, url)
	id := len(r.storage) - 1

	return id, nil
}
