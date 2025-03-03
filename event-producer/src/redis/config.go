package redis

import (
	"rest-server/utils"
)

type redisConfig struct {
	port         int
	host         string
	password     string
	stream       string
	streamMaxLen int64
	writeTimeout int
}

func newRedisConfig(config *redisConfig) {
	config.port = utils.GetEnvInt("REDIS_PORT")
	config.host = utils.GetEnvString("REDIS_HOST")
	config.password = utils.GetEnvString("REDIS_PASSWORD")
	config.stream = utils.GetEnvString("REDIS_STREAM")
	config.streamMaxLen = int64(utils.GetEnvInt("REDIS_STREAM_MAX_LEN"))
	config.writeTimeout = utils.GetEnvInt("REDIS_WRITE_TIMEOUT")
}
