package postgresql

import (
	"context"
	"database/sql"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/jmoiron/sqlx"
)

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) Close() error {
	return r.db.Close()
}

func (r *repository) SaveURL(ctx context.Context, url, shortened string) error {
	q := `INSERT INTO urls (url,shortened_url) VALUES ($1,$2)`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return utils.NewError(err.Error(), utils.InternalErr)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err = tx.QueryxContext(ctx, q, url, shortened); err != nil {
		return utils.NewError(err.Error(), utils.InternalErr)
	}

	if err = tx.Commit(); err != nil {
		return utils.NewError(err.Error(), utils.InternalErr)
	}

	return nil
}

func (r *repository) GetURL(ctx context.Context, shortened string) (string, error) {
	q := `SELECT url FROM urls WHERE shortened_url = $1`

	var url string
	if err := r.db.QueryRowxContext(ctx, q, shortened).Scan(&url); err != nil {
		return "", utils.NewError(err.Error(), utils.InternalErr)
	}

	return url, nil
}
