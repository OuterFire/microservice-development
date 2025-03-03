package api

import (
	"rest-server/utils"
	"time"
)

const pingCheckInterval = 5 * time.Second

type config struct {
	port         int
	writeTimeout int
}

func newConfig(config *config) {
	config.port = utils.GetEnvInt("REST_SERVER_PORT")
	config.writeTimeout = utils.GetEnvInt("REST_SERVER_WRITE_TIMEOUT")
}
