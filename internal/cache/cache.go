package cache

import (
	"context"
	"github.com/alserok/url_shortener/internal/cache/redis"
	"github.com/alserok/url_shortener/internal/config"
)

type Cache interface {
	Set(ctx context.Context, key string, val string) error
	Get(ctx context.Context, key string) (string, error)

	Close() error
}

const (
	Redis = iota
)

func New(t uint, cfg config.Cache) Cache {
	switch t {
	case Redis:
		cl := redis.MustConnect(cfg.RedisDSN())
		return redis.NewCache(cl)
	default:
		panic("invalid cache type")
	}
}
