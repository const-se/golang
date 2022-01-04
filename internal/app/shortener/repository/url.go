package repository

func (r *repository) URL(id int) (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if id >= len(r.cache) {
		return "", ErrURLNotFound
	}

	url := r.cache[id]

	return url, nil
}
