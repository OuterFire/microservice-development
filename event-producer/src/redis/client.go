package redis

import (
	"context"
	"fmt"
	"rest-server/logger"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	logger *logger.Log
	cfg    *redisConfig
	mu     *sync.Mutex
}

func NewRedisClient() *RedisClient {
	log := logger.NewLogger("RedisClient")
	return &RedisClient{
		logger: log,
		mu:     &sync.Mutex{},
	}
}

func newClient(cfg *redisConfig) *redis.Client {
	return redis.NewClient(newUniversalOptions(cfg).Simple())
}

func newUniversalOptions(cfg *redisConfig) *redis.UniversalOptions {
	return &redis.UniversalOptions{
		Addrs:        []string{fmt.Sprintf("%s:%d", cfg.host, cfg.port)},
		Password:     cfg.password,
		WriteTimeout: time.Duration(cfg.writeTimeout) * time.Second,
	}
}

func (rs *RedisClient) Connect() {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	var cfg redisConfig
	newRedisConfig(&cfg)
	rs.cfg = &cfg

	client := newClient(rs.cfg)
	rs.client = client
}

func (rs *RedisClient) Ping(ctx context.Context) error {
	return rs.client.Ping(ctx).Err()
}

func (rs *RedisClient) WriteEntry(ctx context.Context, key string, data string) error {
	return rs.client.XAdd(ctx, &redis.XAddArgs{
		Stream: rs.cfg.stream,
		MaxLen: rs.cfg.streamMaxLen,
		Values: []string{key, data},
	}).Err()
}
