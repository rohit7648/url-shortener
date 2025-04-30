package data

import (
	"context"
	"database/sql"
	"time"
	"url-shortener/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type urlRepo struct {
	data *Data
	log  *log.Helper
}

func NewUrlRepository(data *Data, logger log.Logger) biz.UrlRepository {
	return &urlRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *urlRepo) Save(ctx context.Context, url *biz.Url) error {
	now := time.Now()
	url.CreatedAt = now
	url.UpdatedAt = now

	query := `INSERT INTO urls (long_url, short_code, created_at, updated_at, expires_at) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id`

	err := r.data.db.QueryRowContext(ctx, query,
		url.LongURL,
		url.ShortCode,
		url.CreatedAt,
		url.UpdatedAt,
		url.ExpiresAt,
	).Scan(&url.ID)

	if err != nil {
		return err
	}
	return nil
}

func (r *urlRepo) Get(ctx context.Context, shortCode string) (*biz.Url, error) {
	query := `SELECT id, long_url, short_code, created_at, updated_at, expires_at 
	          FROM urls 
	          WHERE short_code = $1 
	          AND (expires_at IS NULL OR expires_at > NOW())`

	url := &biz.Url{}
	var expiresAt sql.NullTime

	err := r.data.db.QueryRowContext(ctx, query, shortCode).Scan(
		&url.ID,
		&url.LongURL,
		&url.ShortCode,
		&url.CreatedAt,
		&url.UpdatedAt,
		&expiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, biz.ErrUrlNotFound
		}
		return nil, err
	}

	if expiresAt.Valid {
		url.ExpiresAt = &expiresAt.Time
	}

	return url, nil
}

func (r *urlRepo) Exists(ctx context.Context, shortCode string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`
	var exists bool
	err := r.data.db.QueryRowContext(ctx, query, shortCode).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
