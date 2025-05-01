package biz

import (
	"context"
	"errors"
	"time"
)

var (
	ErrURLNotFound = errors.New("URL not found or has expired")
	ErrInvalidUrl  = errors.New("invalid url")
)

// UrlRepository defines the interface for URL storage operations
type UrlRepository interface {
	Save(ctx context.Context, shortCode, longURL string, expiry time.Duration) error
	Get(ctx context.Context, shortCode string) (string, error)
}

// UrlUsecase is a URL usecase.
type UrlUsecase struct {
	repo UrlRepository
}

// NewUrlUsecase new a URL usecase.
func NewUrlUsecase(repo UrlRepository) *UrlUsecase {
	return &UrlUsecase{repo: repo}
}

// Shorten shortens a URL.
func (uc *UrlUsecase) Shorten(ctx context.Context, longURL string, expiry int64) (string, error) {
	shortCode := generateShortCode()
	err := uc.repo.Save(ctx, shortCode, longURL, time.Duration(expiry)*time.Second)
	if err != nil {
		return "", err
	}
	return shortCode, nil
}

// Resolve resolves a short code to a long URL.
func (uc *UrlUsecase) Resolve(ctx context.Context, shortCode string) (string, error) {
	return uc.repo.Get(ctx, shortCode)
}

// generateShortCode generates a random short code.
func generateShortCode() string {
	// Simple implementation - in production, use a more robust method
	return "abc123" // This is just a placeholder
}
