package db

import (
	"context"
	"github.com/alserok/url_shortener/internal/config"
	"github.com/alserok/url_shortener/internal/db/in_memory"
	"github.com/alserok/url_shortener/internal/db/postgresql"
)

type Repository interface {
	SaveURL(ctx context.Context, url, shortened string) error
	GetURL(ctx context.Context, shortened string) (string, error)

	Close() error
}

const (
	PostgreSQL = iota
	InMemory

	postgresMigrationsDir = "./internal/db/postgresql/migrations"
)

func New(t uint, cfg config.DB) Repository {
	var repo Repository

	switch t {
	case PostgreSQL:
		conn := postgresql.MustConnect(cfg.PostgresDSN(), postgresMigrationsDir)
		defer func() {
			_ = conn.Close()
		}()

		repo = postgresql.NewRepository(conn)
	case InMemory:
		return in_memory.NewRepository()
	default:
		panic("invalid repository type")
	}

	return repo
}
