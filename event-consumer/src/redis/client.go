package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"event-consumer/logger"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	logger *logger.Log
	Client *redis.Client
	cfg    *config
}

func NewRedisClient() *RedisClient {
	log := logger.NewLogger("RedisClient")

	var cfg config
	newConfig(&cfg)

	return &RedisClient{
		logger: log,
		cfg:    &cfg,
	}
}

func (rs *RedisClient) NewClient() *redis.Client {
	return redis.NewClient(newUniversalOptions(rs.cfg).Simple())
}

func newUniversalOptions(cfg *config) *redis.UniversalOptions {
	return &redis.UniversalOptions{
		Addrs:       []string{fmt.Sprintf("%s:%d", cfg.host, cfg.port)},
		Password:    cfg.password,
		ReadTimeout: time.Duration(cfg.readTimeout) * time.Second,
	}
}

func (rs *RedisClient) ReadGroup(ctx context.Context) ([]redis.XStream, error) {
	return rs.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Streams:  []string{rs.cfg.stream, ">"},
		Group:    rs.cfg.consumerGroup,
		Consumer: rs.cfg.consumer,
		Count:    1,
		Block:    0,
		NoAck:    false,
	}).Result()

}

func (rs *RedisClient) Ping(ctx context.Context) error {
	return rs.Client.Ping(ctx).Err()
}

func (rs *RedisClient) Close() error {
	return rs.Client.Close()
}

func (rs *RedisClient) GroupCreate(ctx context.Context) error {
	alreadyExists := errors.New("BUSYGROUP Consumer Group name already exists")
	err := rs.Client.XGroupCreate(ctx, rs.cfg.stream, rs.cfg.consumerGroup, "0").Err()
	if err != nil {
		if err.Error() != alreadyExists.Error() {
			return err
		}
		rs.logger.Warn("Consumer Group already exists: %v", err)
		return nil
	}
	rs.logger.Warn("Consumer Group created")
	return nil
}

func (rs *RedisClient) Acknowledge(ctx context.Context, id string) error {
	return rs.Client.XAck(ctx, rs.cfg.stream, rs.cfg.consumerGroup, id).Err()
}

func (rs *RedisClient) Claim(ctx context.Context, messages []string) ([]redis.XMessage, error) {
	return rs.Client.XClaim(ctx, &redis.XClaimArgs{
		Stream:   rs.cfg.stream,
		Group:    rs.cfg.consumerGroup,
		Consumer: rs.cfg.consumer,
		MinIdle:  10 & time.Second,
		Messages: messages,
	}).Result()
}
