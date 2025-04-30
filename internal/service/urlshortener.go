package service

import (
	"context"
	"fmt"
	v1 "url-shortener/api/urlshortener/v1"
	"url-shortener/internal/biz"
	"url-shortener/internal/conf"
)

type UrlShortenerService struct {
	v1.UnimplementedUrlShortenerServer
	uc                 *biz.UrlUsecase
	domain             string
	defaultExpiryHours int64
}

func NewUrlShortenerService(uc *biz.UrlUsecase, c *conf.Server) *UrlShortenerService {
	return &UrlShortenerService{
		uc:                 uc,
		domain:             c.Domain,
		defaultExpiryHours: 24, // Default to 24 hours if not configured
	}
}

func (s *UrlShortenerService) Shorten(ctx context.Context, req *v1.ShortenRequest) (*v1.ShortenReply, error) {
	// Always use the configured default expiry time
	expiresInSeconds := s.defaultExpiryHours * 3600 // Convert hours to seconds

	shortCode, err := s.uc.Shorten(ctx, req.LongUrl, expiresInSeconds)
	if err != nil {
		return nil, err
	}

	// Get the URL to get the expiry time
	url, err := s.uc.Get(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	var expiresAt int64
	if url.ExpiresAt != nil {
		expiresAt = url.ExpiresAt.Unix()
	}

	return &v1.ShortenReply{
		ShortUrl:  fmt.Sprintf("%s/%s", s.domain, shortCode),
		ExpiresAt: expiresAt,
	}, nil
}

func (s *UrlShortenerService) Resolve(ctx context.Context, req *v1.ResolveRequest) (*v1.ResolveReply, error) {
	url, err := s.uc.Get(ctx, req.ShortCode)
	if err != nil {
		if err == biz.ErrUrlNotFound {
			return nil, fmt.Errorf("URL not found or has expired")
		}
		return nil, err
	}

	var expiresAt int64
	if url.ExpiresAt != nil {
		expiresAt = url.ExpiresAt.Unix()
	}

	return &v1.ResolveReply{
		LongUrl:   url.LongURL,
		ExpiresAt: expiresAt,
	}, nil
}
