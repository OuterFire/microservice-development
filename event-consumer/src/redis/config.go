package redis

import (
	"event-consumer/utils"
)

type config struct {
	port          int
	host          string
	password      string
	stream        string
	consumerGroup string
	consumer      string
	readTimeout   int
}

func newConfig(config *config) {
	config.port = utils.GetEnvInt("REDIS_PORT")
	config.host = utils.GetEnvString("REDIS_HOST")
	config.password = utils.GetEnvString("REDIS_PASSWORD")
	config.stream = utils.GetEnvString("REDIS_STREAM")
	config.consumerGroup = utils.GetEnvString("REDIS_CONSUMER_GROUP")
	config.consumer = utils.GetEnvString("REDIS_CONSUMER")
	config.readTimeout = utils.GetEnvInt("REDIS_READ_TIMEOUT")
}
