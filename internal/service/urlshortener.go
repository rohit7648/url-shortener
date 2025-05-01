package service

import (
	"context"
	"time"

	v1 "url-shortener/api/urlshortener/v1"
	"url-shortener/internal/biz"
	"url-shortener/internal/data"

	"github.com/go-kratos/kratos/v2/log"
)

type UrlShortenerService struct {
	v1.UnimplementedUrlShortenerServer

	uc  *biz.UrlUsecase
	log *log.Helper
}

func NewUrlShortenerService(data *data.Data, logger log.Logger) *UrlShortenerService {
	repo := data.NewUrlRepository(logger)
	uc := biz.NewUrlUsecase(repo)
	return &UrlShortenerService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *UrlShortenerService) Shorten(ctx context.Context, req *v1.ShortenRequest) (*v1.ShortenReply, error) {
	shortCode, err := s.uc.Shorten(ctx, req.LongUrl, 24*3600) // Default 24 hours expiry
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	return &v1.ShortenReply{
		ShortUrl:  shortCode,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *UrlShortenerService) Resolve(ctx context.Context, req *v1.ResolveRequest) (*v1.ResolveReply, error) {
	longURL, err := s.uc.Resolve(ctx, req.ShortCode)
	if err != nil {
		return nil, err
	}

	return &v1.ResolveReply{
		LongUrl: longURL,
	}, nil
}
