package redis

import (
	"be/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCommand interface {
	Get(key string) (any, error)
	Set(key string, value any, expiration time.Duration) error
	Delete(key string) error
	Subscribe(key string) (*redis.PubSub, error)
	Publish(key string, message any) error
}

type RedisCache struct {
	Client *redis.Client
}

func NewCache(config *config.Config) (*RedisCache, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Redis.Timeout * time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB: config.Redis.DB,
		PoolSize: config.Redis.MaxConnections,
		DialTimeout: config.Redis.Timeout,
		ReadTimeout: config.Redis.Timeout,
		WriteTimeout: config.Redis.Timeout,
	})
	
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to redis:", err)
	}

	return &RedisCache{Client: client}, nil;
}

func (r *RedisCache) Get(key string) (any, error) {
	return r.Client.Get(context.Background(), key).Result()
}

func (r *RedisCache) Set(key string, value any, expiration time.Duration) error{
	return r.Client.Set(context.Background(), key, value, expiration).Err()
}

func (r *RedisCache) Delete(key string) error {
	return r.Client.Del(context.Background(), key).Err()
}

func (r *RedisCache) Subscribe(channel string) (*redis.PubSub, error) {
	return r.Client.Subscribe(context.Background(), channel), nil
}

func (r *RedisCache) Publish(channel string, message any) error {
	return r.Client.Publish(context.Background(),  channel, message).Err()
}

func (r *RedisCache) Close() error {
	return r.Client.Close()
}