package repository

func (r *repository) URL(id int) (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if id >= len(r.storage) {
		return "", ErrURLNotFound
	}

	url := r.storage[id]

	return url, nil
}
