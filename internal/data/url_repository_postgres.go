package data

import (
	"context"
	"database/sql"
	"net/url"
	"time"
	"url-shortener/internal/biz"
	"url-shortener/internal/pkg/encoding"

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

func (r *urlRepo) Save(ctx context.Context, bizUrl *biz.Url) error {
	// Validate URL
	if _, err := url.Parse(bizUrl.LongURL); err != nil {
		return biz.ErrInvalidUrl
	}

	// Start transaction
	tx, err := r.data.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	bizUrl.CreatedAt = now
	bizUrl.UpdatedAt = now

	// Insert URL
	query := `INSERT INTO urls (long_url, created_at, updated_at, expires_at) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		bizUrl.LongURL,
		bizUrl.CreatedAt,
		bizUrl.UpdatedAt,
		bizUrl.ExpiresAt,
	).Scan(&bizUrl.ID)

	if err != nil {
		return err
	}

	// Generate and update short code
	shortCode := encoding.EncodeBase62(bizUrl.ID)
	bizUrl.ShortCode = shortCode

	updateQuery := `UPDATE urls SET short_code = $1, updated_at = $2 WHERE id = $3`
	_, err = tx.ExecContext(ctx, updateQuery, shortCode, now, bizUrl.ID)
	if err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// Cache the URL in Redis
	if bizUrl.ExpiresAt != nil {
		ttl := time.Until(*bizUrl.ExpiresAt)
		if ttl > 0 {
			err = r.data.redis.Set(ctx, shortCode, bizUrl.LongURL, ttl)
			if err != nil {
				r.log.Errorf("Failed to cache URL in Redis: %v", err)
			}
		}
	} else {
		err = r.data.redis.Set(ctx, shortCode, bizUrl.LongURL, 24*time.Hour)
		if err != nil {
			r.log.Errorf("Failed to cache URL in Redis: %v", err)
		}
	}

	return nil
}

func (r *urlRepo) Get(ctx context.Context, shortCode string) (*biz.Url, error) {
	// Try to get from Redis cache first
	longURL, err := r.data.redis.Get(ctx, shortCode)
	if err == nil {
		// Found in cache
		return &biz.Url{
			ShortCode: shortCode,
			LongURL:   longURL,
		}, nil
	}

	// If not in cache, get from database
	query := `SELECT id, long_url, short_code, created_at, updated_at, expires_at 
	          FROM urls 
	          WHERE short_code = $1 
	          AND (expires_at IS NULL OR expires_at > NOW())`

	url := &biz.Url{}
	var expiresAt sql.NullTime

	err = r.data.db.QueryRowContext(ctx, query, shortCode).Scan(
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

	// Cache the URL in Redis
	if url.ExpiresAt != nil {
		ttl := time.Until(*url.ExpiresAt)
		if ttl > 0 {
			err = r.data.redis.Set(ctx, shortCode, url.LongURL, ttl)
			if err != nil {
				r.log.Errorf("Failed to cache URL in Redis: %v", err)
			}
		}
	} else {
		err = r.data.redis.Set(ctx, shortCode, url.LongURL, 24*time.Hour)
		if err != nil {
			r.log.Errorf("Failed to cache URL in Redis: %v", err)
		}
	}

	return url, nil
}

func (r *urlRepo) Exists(ctx context.Context, shortCode string) (bool, error) {
	// Try to get from Redis cache first
	_, err := r.data.redis.Get(ctx, shortCode)
	if err == nil {
		return true, nil
	}

	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`
	var exists bool
	err = r.data.db.QueryRowContext(ctx, query, shortCode).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
