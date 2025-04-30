package biz

import "context"

type UrlRepository interface {
	Save(ctx context.Context, url *Url) error
	Get(ctx context.Context, shortCode string) (*Url, error)
	Exists(ctx context.Context, shortCode string) (bool, error)
}
