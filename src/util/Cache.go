package util

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheDatabase struct {
	client *redis.Client
}

var Cache *CacheDatabase

func InitCacheDatabase() {
	logger.Info("Initial cache database...")
	redisIns := redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Host,
		Password: Config.Redis.Password,
		DB:       Config.Redis.DB,
	})
	ctx := context.Background()
	_, err := redisIns.Ping(ctx).Result()
	if err != nil {
		errorMsg := "Fail connect cache database. Quitting..."
		logger.Error(errorMsg)
		println(errorMsg)
		os.Exit(-1)
	}
	Cache = &CacheDatabase{
		client: redisIns,
	}
	logger.Info("Success initial cache database!")
}

func (c *CacheDatabase) Set(key string, value interface{}) error {
	return c.client.Set(context.Background(), key, value, time.Duration(Config.Cache.ExpireTime)*time.Minute).Err()
}

func (c *CacheDatabase) Del(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

func (c *CacheDatabase) Exists(key string) (bool, error) {
	result, err := c.client.Exists(context.Background(), key).Result()
	return result > 0, err
}

func (c *CacheDatabase) Close() {
	c.client.Close()
}

func (c *CacheDatabase) Get(key string) (string, error) {
	return c.client.Get(context.Background(), key).Result()
}

func (c *CacheDatabase) SetHTTPResponseCache(key string, resp HTTPResponseCache) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal HTTPResponseCache: %w", err)
	}
	return c.Set(key, data)
}

func (c *CacheDatabase) GetHTTPResponseCache(key string) (*HTTPResponseCache, error) {
	data, err := c.Get(key)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get key from Redis: %w", err)
	}

	var resp HTTPResponseCache
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal HTTPResponseIdentifier: %w", err)
	}

	return &resp, nil
}
