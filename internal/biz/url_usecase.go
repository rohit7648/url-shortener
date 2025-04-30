package biz

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrInvalidUrl  = errors.New("invalid url")
)

// UrlUsecase is a URL usecase.
type UrlUsecase struct {
	repo UrlRepository
}

// NewUrlUsecase new a URL usecase.
func NewUrlUsecase(repo UrlRepository) *UrlUsecase {
	return &UrlUsecase{repo: repo}
}

// Get retrieves a URL by its short code.
func (uc *UrlUsecase) Get(ctx context.Context, shortCode string) (*Url, error) {
	return uc.repo.Get(ctx, shortCode)
}

// Shorten shortens a URL.
func (uc *UrlUsecase) Shorten(ctx context.Context, longURL string, expiresInSeconds int64) (string, error) {
	// Generate a short code
	shortCode := generateShortCode()

	// Check if the short code already exists
	exists, err := uc.repo.Exists(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("short code already exists")
	}

	// Create a new URL
	url := &Url{
		LongURL:   longURL,
		ShortCode: shortCode,
	}

	// Set expiry time if provided
	if expiresInSeconds > 0 {
		expiresAt := time.Now().Add(time.Duration(expiresInSeconds) * time.Second)
		url.ExpiresAt = &expiresAt
		fmt.Printf("Setting expiry time to: %v\n", expiresAt)
	}

	// Save the URL
	err = uc.repo.Save(ctx, url)
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

// Resolve resolves a short code to a long URL.
func (uc *UrlUsecase) Resolve(ctx context.Context, shortCode string) (string, error) {
	url, err := uc.repo.Get(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if url == nil {
		return "", errors.New("URL not found")
	}
	return url.LongURL, nil
}

// generateShortCode generates a random short code.
func generateShortCode() string {
	// This is a simple implementation. You might want to use a more sophisticated one.
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}
