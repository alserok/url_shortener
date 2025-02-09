package in_memory

import (
	"context"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/alserok/url_shortener/pkg/logger"
	"sync"
)

const (
	maxStorageElements       = 10_000
	duplicateKeyErrorMessage = "entity already exists"
)

func NewRepository() *repository {
	return &repository{
		db:   make(map[string]*node[string]),
		size: maxStorageElements,
	}
}

type repository struct {
	mu sync.RWMutex

	size int

	head *node[string]
	tail *node[string]

	db map[string]*node[string]
}

func (r *repository) SaveURL(ctx context.Context, url, shortened string) error {
	log := logger.ExtractLogger(ctx)

	log.Debug("started SaveURL repo", logger.WithArg("url", url))

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.db[shortened]; ok {
		// no need to return an error, entity already exists in db
		return nil
	}

	n := &node[string]{
		next: r.head,
		key:  shortened,
		val:  url,
	}

	if r.head == nil {
		r.tail = n
	} else {
		r.head.prev = n
	}

	r.db[shortened] = n
	r.head = n

	if len(r.db) > r.size {
		delete(r.db, r.tail.key)
		r.tail = r.tail.prev

		if r.tail != nil {
			r.tail.next = nil
		}
	}

	log.Debug("successfully finished SaveURL repo", logger.WithArg("url", url))

	return nil
}

func (r *repository) GetURL(ctx context.Context, shortened string) (string, error) {
	log := logger.ExtractLogger(ctx)

	log.Debug("started GetURL repo", logger.WithArg("shortened_url", shortened))

	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.db[shortened]
	if !ok {
		return "", utils.NewError("entity not found", utils.NotFoundErr)
	}

	if value.prev != nil {
		value.prev.next = value.next
	}

	if value.next != nil {
		value.next.prev = value.prev
	} else {
		r.tail = value.prev
	}

	if value != r.head {
		value.next = r.head
		r.head.prev = value
		r.head = value
		value.prev = nil
	}

	log.Debug("successfully finished GetURL repo", logger.WithArg("shortened_url", shortened))

	return value.val, nil
}

func (r *repository) Close() error {
	return nil
}

type node[T any] struct {
	prev *node[T]
	next *node[T]
	val  T
	key  string
}
