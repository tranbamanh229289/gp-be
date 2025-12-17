package redis

import (
	"be/config"
	"be/pkg/logger"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisCommand interface {
	Get(key string) (any, error)
	Set(key string, value any, expiration time.Duration) error
	Delete(key string) error
	Subscribe(key string) (*redis.PubSub, error)
	Publish(key string, message any) error
}

type RedisCache struct {
	client *redis.Client
	logger *logger.ZapLogger
}

func NewCache(config *config.Config, logger *logger.ZapLogger) (*RedisCache, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Redis.Timeout)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     config.Redis.MaxConnections,
		DialTimeout:  config.Redis.Timeout,
		ReadTimeout:  config.Redis.Timeout,
		WriteTimeout: config.Redis.Timeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to redis: %s", err))
		return nil, err
	}
	logger.Info("Successfully connected to Redis",
		zap.String("addresses", fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)))

	return &RedisCache{client: client, logger: logger}, nil
}

func (r *RedisCache) Get(key string) (any, error) {
	return r.client.Get(context.Background(), key).Result()
}

func (r *RedisCache) Set(key string, value any, expiration time.Duration) error {
	return r.client.Set(context.Background(), key, value, expiration).Err()
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *RedisCache) Subscribe(channel string) (*redis.PubSub, error) {
	return r.client.Subscribe(context.Background(), channel), nil
}

func (r *RedisCache) Publish(channel string, message any) error {
	return r.client.Publish(context.Background(), channel, message).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
