package data

import (
	"database/sql"
	"url-shortener/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/lib/pq"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUrlRepository)

// Data .
type Data struct {
	db    *sql.DB
	redis *RedisClient
	log   *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)

	db, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, nil, err
	}

	redisClient, err := NewRedisClient(c.Redis, logger)
	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		db:    db,
		redis: redisClient,
		log:   log,
	}

	cleanup := func() {
		if err := d.db.Close(); err != nil {
			log.Error(err)
		}
		if err := d.redis.Close(); err != nil {
			log.Error(err)
		}
	}

	return d, cleanup, nil
}
