package in_memory

import (
	"context"
	"github.com/alserok/url_shortener/internal/utils"
	"sync"
)

func NewRepository() *repository {
	return &repository{
		db: make(map[string]string),
	}
}

type repository struct {
	mu sync.RWMutex

	db map[string]string
}

func (r *repository) SaveURL(ctx context.Context, url, shortened string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.db[shortened] = url

	return nil
}

func (r *repository) GetURL(ctx context.Context, shortened string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.db[shortened]
	if !ok {
		return "", utils.NewError("url not found", utils.NotFoundErr)
	}

	return url, nil
}

func (r *repository) Close() error {
	return nil
}
