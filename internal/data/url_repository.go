package data

import (
	"context"
	"url-shortener/internal/biz"
)

type UrlRepository interface {
	Save(ctx context.Context, url *biz.Url) error
	Get(ctx context.Context, shortCode string) (*biz.Url, error)
	Exists(ctx context.Context, shortCode string) (bool, error)
}
