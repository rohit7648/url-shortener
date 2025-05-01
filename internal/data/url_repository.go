package data

import (
	"context"
	"database/sql"
	"time"

	"url-shortener/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type urlRepository struct {
	data   *Data
	logger *log.Helper
}

func NewUrlRepository(data *Data, logger log.Logger) biz.UrlRepository {
	return &urlRepository{
		data:   data,
		logger: log.NewHelper(logger),
	}
}

func (r *urlRepository) Save(ctx context.Context, shortCode, longURL string, expiry time.Duration) error {
	// Store in database
	_, err := r.data.db.ExecContext(ctx, "INSERT INTO urls (short_code, long_url, expires_at) VALUES ($1, $2, $3)",
		shortCode, longURL, time.Now().Add(expiry))
	if err != nil {
		return err
	}

	// Store in Redis cache
	err = r.data.redis.Set(ctx, shortCode, longURL, expiry)
	if err != nil {
		r.logger.Errorf("Failed to cache URL: %v", err)
		return err
	}
	r.logger.Infof("Successfully cached key: %s", shortCode)
	return nil
}

func (r *urlRepository) Get(ctx context.Context, shortCode string) (string, error) {
	// Try Redis cache first
	longURL, err := r.data.redis.Get(ctx, shortCode)
	if err == nil {
		r.logger.Infof("Cache hit for key: %s", shortCode)
		return longURL, nil
	}
	r.logger.Infof("Cache miss for key: %s, falling back to database", shortCode)

	// If not in cache, get from database
	var url string
	var expiresAt time.Time
	err = r.data.db.QueryRowContext(ctx, "SELECT long_url, expires_at FROM urls WHERE short_code = $1", shortCode).Scan(&url, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", biz.ErrURLNotFound
		}
		return "", err
	}

	// Check if URL has expired
	if time.Now().After(expiresAt) {
		return "", biz.ErrURLNotFound
	}

	// Cache the result for future requests
	cacheErr := r.data.redis.Set(ctx, shortCode, url, time.Until(expiresAt))
	if cacheErr != nil {
		r.logger.Errorf("Failed to cache URL after database fetch: %v", cacheErr)
	}

	return url, nil
}
