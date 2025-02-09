package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/alserok/url_shortener/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

const (
	duplicateKeyErrorCode    = "23505"
	duplicateKeyErrorMessage = "entity already exists"
)

type repository struct {
	db *sqlx.DB
}

func (r *repository) Close() error {
	return r.db.Close()
}

func (r *repository) SaveURL(ctx context.Context, url, shortened string) error {
	log := logger.ExtractLogger(ctx)

	log.Debug("started SaveURL repo", logger.WithArg("url", url))

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
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			switch pqError.Code {
			case duplicateKeyErrorCode:
				return utils.NewError(duplicateKeyErrorMessage, utils.BadRequestErr)
			}
		}

		return utils.NewError(err.Error(), utils.InternalErr)
	}

	if err = tx.Commit(); err != nil {
		return utils.NewError(err.Error(), utils.InternalErr)
	}

	log.Debug("successfully finished SaveURL repo", logger.WithArg("url", url))

	return nil
}

func (r *repository) GetURL(ctx context.Context, shortened string) (string, error) {
	log := logger.ExtractLogger(ctx)

	log.Debug("started GetURL repo", logger.WithArg("shortened_url", shortened))

	q := `SELECT url FROM urls WHERE shortened_url = $1`

	var url string
	if err := r.db.QueryRowxContext(ctx, q, shortened).Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", utils.NewError("entity not found", utils.NotFoundErr)
		}
		return "", utils.NewError(err.Error(), utils.InternalErr)
	}

	log.Debug("successfully finished GetURL repo", logger.WithArg("shortened_url", shortened))

	return url, nil
}
