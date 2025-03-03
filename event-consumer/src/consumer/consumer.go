package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"event-consumer/logger"
	"event-consumer/redis"
	"event-consumer/schema"
)

const retryInterval = 5

type ConsumerService struct {
	logger      *logger.Log
	cancel      context.CancelFunc
	redisClient *redis.RedisClient
}

func NewConsumerService() *ConsumerService {
	log := logger.NewLogger("ConsumerService")

	redisClient := redis.NewRedisClient()

	return &ConsumerService{
		logger:      log,
		redisClient: redisClient,
	}
}

func (c *ConsumerService) Stop() error {
	c.logger.Info("Stopping Consumer Service")

	c.cancel()

	err := c.redisClient.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *ConsumerService) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	go c.startConsumerWithRetry(ctx)
}

func (c *ConsumerService) startConsumerWithRetry(ctx context.Context) {
	for {
		err := c.listenStream(ctx)
		if err != nil {
			c.logger.Error("Error consumer: %v", err)
			if errors.Is(err, context.Canceled) {
				return
			}
		}

		select {
		case <-time.After(retryInterval * time.Second):
			c.logger.Error("Restarting Consumer")
			continue
		case <-ctx.Done():
			c.logger.Error("Error consumer: %v", err)
			return
		}
	}
}

func (c *ConsumerService) listenStream(ctx context.Context) error {
	c.logger.Info("Starting Consumer Service")

	client := c.redisClient.NewClient()
	c.redisClient.Client = client

	err := c.redisClient.Ping(ctx)
	if err != nil {
		return err
	}

	err = c.redisClient.GroupCreate(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
			event, err := c.redisClient.ReadGroup(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					continue
				}
				return err
			}
			c.logger.Debug("Received event from consumer group: %v", event)

			if len(event) > 0 {
				for _, stream := range event[0].Messages {
					c.processEvent(stream.Values)

					err := c.redisClient.Acknowledge(ctx, stream.ID)
					if err != nil {
						c.logger.Error("Error acknowledge: %v", err)
					}
				}
			}
		}
	}
}

//var Data []schema.NotificationMessage

func (c *ConsumerService) processEvent(data map[string]interface{}) {
	if _, ok := data["CreateStream"]; ok {
		c.logger.Info("CreateStream: %v", data["CreateStream"])
		return
	}

	if value, ok := data["EventStream"]; ok {
		j := []byte(value.(string))
		var event schema.NotificationMessage
		err := json.Unmarshal(j, &event)
		if err != nil {
			fmt.Println("Error unmarshalling", err)
			return
		}
		c.logger.Info("Event: %v", event)

	}
}
