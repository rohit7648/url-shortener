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
	db *sql.DB
}

// NewData .
func NewData(cfg *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)

	db, err := sql.Open("postgres", cfg.Database.Source)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.Info("closing the data resources")
		db.Close()
	}

	return &Data{db: db}, cleanup, nil
}
