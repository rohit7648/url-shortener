package data

import (
	"context"
	"time"

	"url-shortener/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
	log    *log.Helper
}

func NewRedisClient(conf *conf.Data_Redis, logger log.Logger) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Network:      conf.Network,
		Addr:         conf.Addr,
		Password:     conf.Password,
		DB:           int(conf.Db),
		DialTimeout:  conf.DialTimeout.AsDuration(),
		ReadTimeout:  conf.ReadTimeout.AsDuration(),
		WriteTimeout: conf.WriteTimeout.AsDuration(),
	})

	log := log.NewHelper(logger)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Errorf("Failed to connect to Redis: %v", err)
		return nil, err
	}

	log.Info("Successfully connected to Redis")

	return &RedisClient{
		client: client,
		log:    log,
	}, nil
}

func (r *RedisClient) Close() error {
	r.log.Info("Closing Redis connection")
	return r.client.Close()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	start := time.Now()
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.log.Debugf("Cache miss for key: %s", key)
		} else {
			r.log.Errorf("Failed to get from Redis: %v", err)
		}
	} else {
		r.log.Debugf("Cache hit for key: %s (took %v)", key, time.Since(start))
	}
	return value, err
}

func (r *RedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	start := time.Now()
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		r.log.Errorf("Failed to set in Redis: %v", err)
	} else {
		r.log.Debugf("Successfully cached key: %s (took %v)", key, time.Since(start))
	}
	return err
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.log.Errorf("Failed to delete from Redis: %v", err)
	} else {
		r.log.Debugf("Successfully deleted key: %s", key)
	}
	return err
}

func (r *RedisClient) Ping(ctx context.Context) error {
	start := time.Now()
	err := r.client.Ping(ctx).Err()
	if err != nil {
		r.log.Errorf("Redis ping failed: %v", err)
	} else {
		r.log.Debugf("Redis ping successful (took %v)", time.Since(start))
	}
	return err
}
